package producer

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"rate-service/internal/models"
)

type MailProducer struct {
	channel *amqp.Channel
}

func NewMailProducer(channel *amqp.Channel) *MailProducer {
	return &MailProducer{
		channel: channel,
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
		"mail",
		false,
		false,
		amqp.Publishing{
			ContentType: "json",
			Body:        emailEventJSON,
		})
	if err != nil {
		log.Fatalf("failed to publish a message: %s", err.Error())
	}
}
