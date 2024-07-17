package main

import (
	"customer-service/internal/configs"
	"customer-service/internal/messaging/consumer"
	"customer-service/internal/messaging/producer"
	"customer-service/internal/repository"
	"customer-service/internal/service"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"syscall"
)

const (
	MigrationsPath          = "file://migrations"
	ConfigPath              = "configs/config.yml"
	RabbitMQConnStrTemplate = "amqp://%s:%s@%s:%s/"
)

func main() {
	config, err := configs.NewConfig(ConfigPath)
	if err != nil {
		log.Fatalf("failed to read config: %s", err.Error())
	}
	db, err := repository.NewDB(&config.DB)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create postgres driver: %s", err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(MigrationsPath, config.DB.DBName, driver)
	if err != nil {
		log.Fatalf("failed to create migration instance: %s", err.Error())
	}
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("No migrations to apply")
		} else {
			log.Fatalf("failed to apply migrations: %s", err.Error())
		}
	}

	conn, err := amqp.Dial(fmt.Sprintf(RabbitMQConnStrTemplate,
		config.RabbitMQ.Username,
		config.RabbitMQ.Password,
		config.RabbitMQ.Host,
		config.RabbitMQ.Port,
	))
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %s", err.Error())
	}
	defer conn.Close()
	log.Info("Connected to RabbitMQ")

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %s", err.Error())
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(
		config.RabbitMQ.Queue.CustomerEvent,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %s", err.Error())
	}

	messageProducer := producer.NewMessageProducer(channel)
	customerRepository := repository.NewCustomerRepository(db)
	customerService := service.NewCustomerService(customerRepository, &config.RabbitMQ, messageProducer)
	customerCommandConsumer := consumer.NewCustomerConsumer(channel, customerService, &config.RabbitMQ)
	customerCommandConsumer.StartConsuming()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := conn.Close(); err != nil {
		log.Errorf("error occurred while closing RabbitMQ connection: %s", err.Error())
	}
}
