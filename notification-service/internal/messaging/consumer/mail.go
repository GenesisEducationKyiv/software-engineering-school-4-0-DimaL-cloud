package consumer

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"notification-service/internal/configs"
	"notification-service/internal/models"
	"notification-service/internal/service"
)

const (
	ServiceName = "mail-service"
)

type MailConsumer struct {
	channel     *amqp.Channel
	mailService service.Mail
	config      *configs.RabbitMQ
}

func NewMailConsumer(channel *amqp.Channel, mailService service.Mail, config *configs.RabbitMQ) *MailConsumer {
	return &MailConsumer{
		channel:     channel,
		mailService: mailService,
		config:      config,
	}
}

func (c *MailConsumer) StartConsuming() {
	msgs, err := c.channel.Consume(
		c.config.Queue.Mail,
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
			log.Errorf("received a message: %s", msg.Body)
			c.handleMessage(msg)
			if err := msg.Ack(false); err != nil {
				log.Errorf("failed to acknowledge message: %s", err.Error())
			}
		}
	}()
}

func (c *MailConsumer) handleMessage(msg amqp.Delivery) {
	var sendEmailCommand models.SendEmailCommand
	err := json.Unmarshal(msg.Body, &sendEmailCommand)
	if err != nil {
		log.Errorf("failed to deserialize message: %s", err.Error())
		return
	}
	c.mailService.SendEmail(sendEmailCommand.Subject, sendEmailCommand.Body, sendEmailCommand.To)
}
