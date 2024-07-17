package consumer

import (
	"customer-service/internal/configs"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	ServiceName = "customer-service"
)

type CustomerConsumer struct {
	channel         *amqp.Channel
	customerService service.Customer
	config          *configs.RabbitMQ
}

func NewCustomerConsumer(channel *amqp.Channel, customerService service.Customer, config *configs.RabbitMQ) *CustomerConsumer {
	return &CustomerConsumer{
		channel:     channel,
		mailService: customerService,
		config:      config,
	}
}

func (c *CustomerConsumer) StartConsuming() {
	msgs, err := c.channel.Consume(
		c.config.Queue.Customer,
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

func (c *CustomerConsumer) handleMessage(msg amqp.Delivery) {
	var sendEmailCommand models.SendEmailCommand
	err := json.Unmarshal(msg.Body, &sendEmailCommand)
	if err != nil {
		log.Errorf("failed to deserialize message: %s", err.Error())
		return
	}
	c.mailService.SendEmail(sendEmailCommand.Subject, sendEmailCommand.Body, sendEmailCommand.To)
}
