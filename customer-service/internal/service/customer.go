package service

import (
	"customer-service/internal/configs"
	"customer-service/internal/messaging/producer"
	"customer-service/internal/models"
	"customer-service/internal/repository"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	CustomerCreatedEvent        = "CustomerCreatedEvent"
	CustomerCreationFailedEvent = "CustomerCreationFailedEvent"
)

type Customer interface {
	CreateCustomer(customer models.Customer) (int, error)
	DeleteCustomerByEmail(email string) error
}

type CustomerService struct {
	repository      repository.Customer
	messageProducer producer.Producer
	rabbitMQConfig  *configs.RabbitMQ
}

func NewCustomerService(
	repository repository.Customer,
	rabbitMQConfig *configs.RabbitMQ,
	messageProducer producer.Producer,
) *CustomerService {
	return &CustomerService{
		repository:      repository,
		rabbitMQConfig:  rabbitMQConfig,
		messageProducer: messageProducer,
	}
}

func (s *CustomerService) CreateCustomer(customer models.Customer) (int, error) {
	id, err := s.repository.CreateCustomer(customer)
	if err != nil {
		log.Error("failed to create customer: ", err)
		customerCreationFailedEvent := models.CustomerCreationFailedEvent{
			Event: models.Event{
				Type:      CustomerCreationFailedEvent,
				Timestamp: time.Now(),
			},
			SubscriptionID: customer.SubscriptionID,
		}

		s.messageProducer.PublishMessage(customerCreationFailedEvent, s.rabbitMQConfig.Queue.CustomerEvent)
	} else {
		customerCreatedEvent := models.CustomerCreatedEvent{
			Event: models.Event{
				Type:      CustomerCreatedEvent,
				Timestamp: time.Now(),
			},
			CustomerID: id,
		}
		s.messageProducer.PublishMessage(customerCreatedEvent, s.rabbitMQConfig.Queue.CustomerEvent)
	}
	return id, err
}

func (s *CustomerService) DeleteCustomerByEmail(email string) error {
	err := s.repository.DeleteCustomerByEmail(email)
	if err != nil {
		log.Error("failed to delete customer", err)
	} else {
		log.Info("customer deleted for email: ", email)
	}
	return err
}
