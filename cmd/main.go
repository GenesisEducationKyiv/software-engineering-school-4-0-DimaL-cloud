package main

import (
	"context"
	"errors"
	_ "github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/docs"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/pkg/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/pkg/handler"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/pkg/models"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/pkg/repository"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/pkg/scheduler"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/pkg/service"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	MigrationsPath         = "file://migrations"
	ShutdownTimeoutSeconds = 5
)

// @title Exchange rate notifier API
// @version 1.0
// @description API server for notifying exchange rate
func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}
	dbConfig := repository.DBConfig{
		Host:       viper.GetString("db.host"),
		Port:       viper.GetString("db.port"),
		Username:   viper.GetString("db.username"),
		Password:   viper.GetString("db.password"),
		DBName:     viper.GetString("db.name"),
		SSLMode:    viper.GetString("db.ssl_mode"),
		DriverName: viper.GetString("db.driver_name"),
	}
	mailConfig := service.MailConfig{
		Host:     viper.GetString("mail.host"),
		Port:     viper.GetString("mail.port"),
		Username: viper.GetString("mail.username"),
		Password: viper.GetString("mail.password"),
	}
	db, err := repository.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create postgres driver: %s", err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		MigrationsPath,
		dbConfig.DBName, driver)
	if err != nil {
		log.Fatalf("failed to create migration instance: %s", err.Error())
	}
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Print("No migrations to apply")
		} else {
			log.Fatalf("failed to apply migrations: %s", err.Error())
		}
	}
	repositories := repository.NewRepository(db)
	services := service.NewService(repositories, mailConfig)
	clients := client.NewClient()
	handlers := handler.NewHandler(services)
	notificationScheduler := scheduler.NewExchangeRateNotificationScheduler(
		repositories.Subscription,
		clients.ExchangeRate,
		services.Mail)
	notificationScheduler.StartJob()
	server := new(models.Server)
	go func() {
		if err := server.Run(viper.GetString("server.port"), handlers.InitRoutes()); err != nil {
			log.Fatalf("failed to start server: %s", err.Error())
		}
	}()

	log.Print("Exchange rate notifier started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Exchange rate notifier shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeoutSeconds*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("error occurred while shutting down server: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		log.Errorf("error occurred while closing db connection: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
