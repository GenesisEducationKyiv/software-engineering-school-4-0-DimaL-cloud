package producer

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"rate-service/internal/configs"
	"rate-service/internal/models"
)

const (
	ContentType = "application/json"
)

type MailProducer struct {
	channel *amqp.Channel
	config  *configs.RabbitMQ
}

func NewMailProducer(channel *amqp.Channel, config *configs.RabbitMQ) *MailProducer {
	return &MailProducer{
		channel: channel,
		config:  config,
	}
}

func (p *MailProducer) PublishMail(sendEmailCommand models.SendEmailCommand) {
	emailEventJSON, err := json.Marshal(sendEmailCommand)
	if err != nil {
		log.Errorf("failed to serialize SendEmailCommand to JSON: %s", err.Error())
		return
	}
	err = p.channel.Publish(
		"",
		p.config.Queue.Mail,
		false,
		false,
		amqp.Publishing{
			ContentType: ContentType,
			Body:        emailEventJSON,
		})
	if err != nil {
		log.Fatalf("failed to publish a message: %s", err.Error())
	}
}
