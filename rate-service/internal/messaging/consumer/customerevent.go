package consumer

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"rate-service/internal/configs"
	"rate-service/internal/models"
	"rate-service/internal/service"
)

const (
	CustomerCreatedEvent        = "CustomerCreatedEvent"
	CustomerCreationFailedEvent = "CustomerCreationFailedEvent"
)

type CustomerEventConsumer struct {
	channel             *amqp.Channel
	subscriptionService service.Subscription
	config              *configs.RabbitMQ
}

func NewCustomerEventConsumer(
	channel *amqp.Channel,
	subscriptionService service.Subscription,
	config *configs.RabbitMQ) *CustomerEventConsumer {
	return &CustomerEventConsumer{
		channel:             channel,
		subscriptionService: subscriptionService,
		config:              config,
	}
}

func (c *CustomerEventConsumer) StartConsuming() {
	msgs, err := c.channel.Consume(
		c.config.Queue.CustomerEvent,
		"rate-service-customer-event-consumer",
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

func (c *CustomerEventConsumer) handleMessage(msg amqp.Delivery) {
	var event models.Event
	err := json.Unmarshal(msg.Body, &event)
	if err != nil {
		log.Errorf("failed to deserialize message: %s", err.Error())
		return
	}

	switch event.Type {
	case CustomerCreatedEvent:
		var customerCreatedEvent models.CustomerCreatedEvent
		err := json.Unmarshal(msg.Body, &customerCreatedEvent)
		if err != nil {
			log.Errorf("failed to deserialize message: %s", err.Error())
			return
		}
		log.Infof("Customer with id %d created", customerCreatedEvent.CustomerId)
	case CustomerCreationFailedEvent:
		log.Infof("Customer creation failed")
		var customerCreationFailedEvent models.CustomerCreationFailedEvent
		err := json.Unmarshal(msg.Body, &customerCreationFailedEvent)
		if err != nil {
			log.Errorf("failed to deserialize message: %s", err.Error())
			return
		}
		err = c.subscriptionService.DeleteSubscriptionByID(customerCreationFailedEvent.SubscriptionId)
		if err != nil {
			log.Errorf("failed to remove subscription: %s", err.Error())
			return
		}
	default:
		log.Errorf("unknown event type: %s", event.Type)
	}
}
