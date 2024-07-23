package main

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"net/http"
	"os"
	"os/signal"
	"scheduler-service/internal/configs"
	"scheduler-service/internal/handler"
	"scheduler-service/internal/models"
	"scheduler-service/internal/scheduler"
	"syscall"
	"time"
)

const (
	ConfigPath              = "configs/config.yml"
	RabbitMQConnStrTemplate = "amqp://%s:%s@%s:%s/"
	ShutdownTimeout         = 5 * time.Second
)

func main() {
	config, err := configs.NewConfig(ConfigPath)
	if err != nil {
		log.Fatalf("failed to read config: %s", err.Error())
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
		log.Fatalf("failed to open a channel: %s", err.Error()) // nolint: gocritic
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(
		config.RabbitMQ.Queue.RateNotificationCron,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %s", err.Error())
	}
	rateNotificationScheduler := scheduler.NewRateNotificationScheduler(&config.Crons, &config.RabbitMQ, channel)
	rateNotificationScheduler.StartJob()

	h := handler.NewHandler()
	s := models.NewServer(h.InitRoutes(), &config.Server)
	if err := s.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %s", err.Error())
		}
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Errorf("error occurred while shutting down server: %s", err.Error())
	}
	if err := conn.Close(); err != nil {
		log.Errorf("error occurred while closing RabbitMQ connection: %s", err.Error())
	}
}
