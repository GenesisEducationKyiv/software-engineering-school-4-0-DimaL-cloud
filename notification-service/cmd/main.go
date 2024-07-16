package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"notification-service/internal/configs"
	"notification-service/internal/messaging/consumer"
	"notification-service/internal/service"
	"os"
	"os/signal"
	"syscall"
)

const (
	ConfigPath = "configs/config.yml"
)

func main() {
	config, err := configs.NewConfig(ConfigPath)
	if err != nil {
		log.Fatalf("failed to read config: %s", err.Error())
	}

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
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

	mailService := service.NewMailService(config.Mail)

	mailConsumer := consumer.NewMailConsumer(channel, mailService, &config.RabbitMQ)
	mailConsumer.StartConsuming()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := conn.Close(); err != nil {
		log.Errorf("error occurred while closing RabbitMQ connection: %s", err.Error())
	}
}
