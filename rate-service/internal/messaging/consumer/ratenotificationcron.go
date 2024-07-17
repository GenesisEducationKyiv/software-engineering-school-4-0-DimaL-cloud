package consumer

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"rate-service/internal/configs"
	"rate-service/internal/messaging/producer"
	"rate-service/internal/models"
	"rate-service/internal/service"
	"time"
)

const (
	EmailSubject = "Курс НБУ"
	EmailBody    = "Курс долара НБУ станом на %s: %f грн"
	ServiceName  = "rate-service"
)

type RateNotificationCronConsumer struct {
	channel             *amqp.Channel
	messageProducer     producer.Producer
	subscriptionService service.Subscription
	rateService         service.Rate
	config              *configs.RabbitMQ
}

func NewRateNotificationCronConsumer(
	channel *amqp.Channel,
	producer producer.Producer,
	subscriptionService service.Subscription,
	rateService service.Rate,
	config *configs.RabbitMQ) *RateNotificationCronConsumer {
	return &RateNotificationCronConsumer{
		channel:             channel,
		messageProducer:     producer,
		subscriptionService: subscriptionService,
		rateService:         rateService,
		config:              config,
	}
}

func (c *RateNotificationCronConsumer) StartConsuming() {
	msgs, err := c.channel.Consume(
		c.config.Queue.RateNotificationCron,
		ServiceName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to start consuming messages: %s", err.Error())
		return
	}
	go func() {
		for msg := range msgs {
			log.Infof("received a message: %s", msg.Body)
			c.handleMessage()
			if err := msg.Ack(false); err != nil {
				log.Errorf("failed to acknowledge message: %s", err.Error())
			}
		}
	}()
}

func (c *RateNotificationCronConsumer) handleMessage() {
	rate, err := c.rateService.GetRate()
	if err != nil {
		log.Errorf("failed to get current exchange rate: %s", err.Error())
		return
	}
	subscriptions, err := c.subscriptionService.GetAllSubscriptions()
	if err != nil {
		log.Errorf("failed to get subscriptions: %s", err.Error())
		return
	}
	currentDate := time.Now().Format("02.01.2006")
	sendEmailCommand := models.SendEmailCommand{
		Subject: EmailSubject,
		Body:    fmt.Sprintf(EmailBody, currentDate, rate),
	}
	for _, subscription := range subscriptions {
		sendEmailCommand.To = subscription.Email
		c.messageProducer.PublishMessage(sendEmailCommand, c.config.Queue.Mail)
	}
}
