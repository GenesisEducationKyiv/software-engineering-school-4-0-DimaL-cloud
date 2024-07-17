package producer

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	ApplicationJson = "application/json"
)

type Producer interface {
	PublishMessage(message any, queue string)
}

type MessageProducer struct {
	channel *amqp.Channel
}

func NewMessageProducer(channel *amqp.Channel) *MessageProducer {
	return &MessageProducer{
		channel: channel,
	}
}

func (p *MessageProducer) PublishMessage(message any, queue string) {
	serializedMessage, err := json.Marshal(message)
	if err != nil {
		log.Errorf("failed to serialize message to JSON: %s", err.Error())
		return
	}
	err = p.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: ApplicationJson,
			Body:        serializedMessage,
		})
	if err != nil {
		log.Errorf("failed to publish a message: %s", err.Error())
	} else {
		log.Infof("message %s published to queue: %s", serializedMessage, queue)
	}
}
