package consumer

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"rate-service/internal/messaging/producer"
	"rate-service/internal/models"
	"rate-service/internal/service"
	"time"
)

const (
	EmailSubject = "Курс НБУ"
	EmailBody    = "Курс долара НБУ станом на %s: %f грн"
	EventName    = "RateNotificationCronReceived"
)

type RateNotificationCronConsumer struct {
	channel             *amqp.Channel
	mailProducer        *producer.MailProducer
	subscriptionService service.Subscription
	rateService         service.Rate
}

func NewRateNotificationCronConsumer(
	channel *amqp.Channel,
	producer *producer.MailProducer,
	subscriptionService service.Subscription,
	rateService service.Rate) *RateNotificationCronConsumer {
	return &RateNotificationCronConsumer{
		channel:             channel,
		mailProducer:        producer,
		subscriptionService: subscriptionService,
		rateService:         rateService,
	}
}

func (c *RateNotificationCronConsumer) StartConsuming() {
	msgs, err := c.channel.Consume(
		"rate-notification-cron",
		"rate-service",
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
	var emails []string
	for _, subscription := range subscriptions {
		emails = append(emails, subscription.Email)
	}
	sendEmailCommand := models.SendEmailCommand{
		Subject: EmailSubject,
		Body:    fmt.Sprintf(EmailBody, currentDate, rate),
		To:      emails,
	}
	c.mailProducer.PublishMail(sendEmailCommand)
}