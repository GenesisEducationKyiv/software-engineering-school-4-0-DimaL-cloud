package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go.uber.org/fx"
	"net/http"
	"rate-service/internal/client"
	"rate-service/internal/client/rate"
	"rate-service/internal/configs"
	"rate-service/internal/handler"
	"rate-service/internal/messaging/consumer"
	"rate-service/internal/messaging/producer"
	"rate-service/internal/models"
	"rate-service/internal/repository"
	"rate-service/internal/service"
)

const (
	MigrationsPath          = "file://migrations"
	ConfigPath              = "configs/config.yml"
	RabbitMQConnStrTemplate = "amqp://%s:%s@%s:%s/"
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
		fx.Provide(func(config *configs.Config) *configs.Server {
			return &config.Server
		}),
		fx.Provide(func(config *configs.Config) *configs.RabbitMQ {
			return &config.RabbitMQ
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
		fx.Provide(func(handler *handler.Handler, config *configs.Server) *http.Server {
			return models.NewServer(handler.InitRoutes(), config)
		}),
		fx.Provide(NewMigrateInstance),
		fx.Provide(NewRabbitMQConnection),
		fx.Provide(NewRabbitMQChannel),
		fx.Provide(consumer.NewCustomerEventConsumer),
		fx.Provide(consumer.NewRateNotificationCronConsumer),
		fx.Provide(
			fx.Annotate(
				producer.NewMessageProducer,
				fx.As(new(producer.Producer)),
			),
		),
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

func NewRabbitMQConnection(config *configs.RabbitMQ) (*amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf(RabbitMQConnStrTemplate,
		config.Username,
		config.Password,
		config.Host,
		config.Port,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	return conn, nil
}

func NewRabbitMQChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}
	return channel, nil
}

func run(
	lc fx.Lifecycle,
	server *http.Server,
	db *sqlx.DB,
	migrate *migrate.Migrate,
	conn *amqp.Connection,
	rateNotificationCronConsumer *consumer.RateNotificationCronConsumer,
	customerEventConsumer *consumer.CustomerEventConsumer,
	channel *amqp.Channel,
	rabbitMQConfig *configs.RabbitMQ) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error { //nolint:revive
			applyMigrations(migrate)
			createQueue(channel, rabbitMQConfig.Queue.Mail)
			createQueue(channel, rabbitMQConfig.Queue.CustomerCommand)
			go rateNotificationCronConsumer.StartConsuming()
			go customerEventConsumer.StartConsuming()
			go startServer(server)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := server.Shutdown(ctx); err != nil {
				log.Errorf("error occurred while shutting down server: %s", err.Error())
			}
			if err := db.Close(); err != nil {
				log.Errorf("error occurred while closing db connection: %s", err.Error())
			}
			if err := conn.Close(); err != nil {
				log.Errorf("error occurred while closing RabbitMQ connection: %s", err.Error())
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

func createQueue(channel *amqp.Channel, queue string) {
	_, err := channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %s", err.Error())
	}
}
