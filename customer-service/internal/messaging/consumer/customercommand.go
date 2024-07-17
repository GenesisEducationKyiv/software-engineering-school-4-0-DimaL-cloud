package consumer

import (
	"customer-service/internal/configs"
	"customer-service/internal/models"
	"customer-service/internal/service"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	ServiceName           = "customer-service"
	CreateCustomerCommand = "CreateCustomerCommand"
)

type CustomerCommandConsumer struct {
	channel         *amqp.Channel
	customerService service.Customer
	config          *configs.RabbitMQ
}

func NewCustomerConsumer(
	channel *amqp.Channel,
	customerService service.Customer,
	config *configs.RabbitMQ,
) *CustomerCommandConsumer {
	return &CustomerCommandConsumer{
		channel:         channel,
		customerService: customerService,
		config:          config,
	}
}

func (c *CustomerCommandConsumer) StartConsuming() {
	msgs, err := c.channel.Consume(
		c.config.Queue.CustomerCommand,
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
			c.handleMessage(msg)
			if err := msg.Ack(false); err != nil {
				log.Errorf("failed to acknowledge message: %s", err.Error())
			}
		}
	}()
}

func (c *CustomerCommandConsumer) handleMessage(msg amqp.Delivery) {
	var command models.Command
	err := json.Unmarshal(msg.Body, &command)
	if err != nil {
		log.Errorf("failed to deserialize message: %s", err.Error())
		return
	}

	switch command.Type {
	case CreateCustomerCommand:
		var createCustomerCommand models.CreateCustomerCommand
		err := json.Unmarshal(msg.Body, &createCustomerCommand)
		if err != nil {
			log.Errorf("failed to deserialize message: %s", err.Error())
			return
		}
		customer := models.Customer{
			Email:          createCustomerCommand.Email,
			SubscriptionID: createCustomerCommand.SubscriptionID,
		}
		_, _ = c.customerService.CreateCustomer(customer)
	default:
		log.Errorf("unknown command type: %s", command.Type)
	}
}
