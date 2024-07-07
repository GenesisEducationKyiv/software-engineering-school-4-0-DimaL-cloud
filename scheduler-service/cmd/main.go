package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"scheduler-service/internal/configs"
	"scheduler-service/internal/repository"
	"scheduler-service/internal/scheduler"
	"syscall"
	"time"
)

const (
	ConfigPath      = "configs/config.yml"
	ShutdownTimeout = 5 * time.Second
)

func main() {
	config, err := configs.NewConfig(ConfigPath)
	if err != nil {
		log.Fatalf("failed to read config: %s", err.Error())
	}

	conn, err := amqp.Dial("amqp://rmuser:rmpassword@localhost:5672/")
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
		"rate-notification-cron",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %s", err.Error())
	}
	db, err := repository.NewDB(&config.DB)
	eventRepository := repository.NewEventRepository(db)
	rateNotificationScheduler := scheduler.NewRateNotificationScheduler(&config.Crons, channel, eventRepository)
	rateNotificationScheduler.StartJob()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := db.Close(); err != nil {
		log.Errorf("error occurred while closing db connection: %s", err.Error())
	}
	if err := conn.Close(); err != nil {
		log.Errorf("error occurred while closing RabbitMQ connection: %s", err.Error())
	}
}
