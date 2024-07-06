package main

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"net/http"
	"rate-service/internal/client"
	"rate-service/internal/client/rate"
	"rate-service/internal/configs"
	"rate-service/internal/handler"
	"rate-service/internal/models"
	"rate-service/internal/repository"
	"rate-service/internal/scheduler"
	"rate-service/internal/service"
)

const (
	MigrationsPath = "file://migrations"
	ConfigPath     = "configs/config.yml"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(client.NewHTTPClient),
		fx.Provide(func() (*configs.Config, error) {
			return configs.NewConfig(ConfigPath)
		}),
		fx.Provide(func(config *configs.Config) *configs.DB {
			return &config.DB
		}),
		fx.Provide(func(config *configs.Config) *configs.Rate {
			return &config.Rate
		}),
		fx.Provide(func(config *configs.Config) *configs.Mail {
			return &config.Mail
		}),
		fx.Provide(func(config *configs.Config) *configs.Server {
			return &config.Server
		}),
		fx.Provide(repository.NewDB),
		fx.Provide(func(client *http.Client, config *configs.Config) *rate.NbuRateClient {
			return rate.NewNbuRateClient(client, config.Rate.APIUrls.Nbu)
		}),
		fx.Provide(func(client *http.Client, config *configs.Config) *rate.PrivatBankRateClient {
			return rate.NewPrivatBankRateClient(client, config.Rate.APIUrls.PrivatBank)
		}),
		fx.Provide(func(client *http.Client, config *configs.Config) *rate.FawazahmedRateClient {
			return rate.NewFawazahmedRateClient(client, config.Rate.APIUrls.Fawazahmed)
		}),
		fx.Provide(
			fx.Annotate(
				repository.NewSubscriptionRepository,
				fx.As(new(repository.Subscription)),
			),
		),
		fx.Provide(service.NewMailService),
		fx.Provide(
			fx.Annotate(
				service.NewSubscriptionService,
				fx.As(new(service.Subscription)),
			),
		),
		fx.Provide(
			fx.Annotate(
				func(
					nbuClient *rate.NbuRateClient,
					privatBankClient *rate.PrivatBankRateClient,
					fawazahmedClient *rate.FawazahmedRateClient,
				) *service.RateService {
					nbuClient.SetNext(privatBankClient)
					privatBankClient.SetNext(fawazahmedClient)
					return service.NewRateService(nbuClient)
				},
				fx.As(new(service.Rate)),
			),
		),
		fx.Provide(handler.NewHandler),
		fx.Provide(scheduler.NewRateNotificationScheduler),
		fx.Provide(func(handler *handler.Handler, config *configs.Server) *http.Server {
			return models.NewServer(handler.InitRoutes(), config)
		}),
		fx.Provide(NewMigrateInstance),
		fx.Invoke(run),
	)
}

func NewMigrateInstance(db *sqlx.DB, config *configs.DB) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create postgres driver: %s", err.Error())
		return nil, err
	}
	return migrate.NewWithDatabaseInstance(MigrationsPath, config.DBName, driver)
}

func run(lc fx.Lifecycle, s *http.Server, db *sqlx.DB, m *migrate.Migrate) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error { //nolint:revive
			applyMigrations(m)
			go startServer(s)
			log.Info("Exchange rate notifier started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := s.Shutdown(ctx); err != nil {
				log.Errorf("error occurred while shutting down server: %s", err.Error())
			}
			if err := db.Close(); err != nil {
				log.Errorf("error occurred while closing db connection: %s", err.Error())
			}
			return nil
		},
	})
}

func applyMigrations(m *migrate.Migrate) {
	err := m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("No migrations to apply")
		} else {
			log.Fatalf("failed to apply migrations: %s", err.Error())
		}
	}
}

func startServer(s *http.Server) {
	if err := s.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %s", err.Error())
		}
	}
}
