package service

import (
	log "github.com/sirupsen/logrus"
	"rate-service/internal/configs"
	"rate-service/internal/messaging/producer"
	"rate-service/internal/models"
	"rate-service/internal/repository"
)

type Subscription interface {
	GetAllSubscriptions() ([]models.Subscription, error)
	CreateCustomer(email string) (int, error)
	DeleteCustomerByEmail(email string) error
}

type SubscriptionService struct {
	repository      repository.Subscription
	messageProducer producer.Producer
	rabbitMQConfig  *configs.RabbitMQ
}

func NewSubscriptionService(repository repository.Subscription) *SubscriptionService {
	return &SubscriptionService{repository: repository}
}

func (s *SubscriptionService) GetAllSubscriptions() ([]models.Subscription, error) {
	return s.repository.GetAllSubscriptions()
}

func (s *SubscriptionService) CreateSubscription(email string) error {
	id, err := s.repository.CreateCustomer(email)
	if err != nil {
		log.Error("failed to create subscription", err)
	} else {
		log.Info("subscription created for email: ", email)
	}
	createCustomerCommand := models.CreateCustomerCommand{
		Email:          email,
		SubscriptionID: id,
	}
	s.messageProducer.PublishMessage(createCustomerCommand, s.rabbitMQConfig.Queue.Customer)
	return err
}

func (s *SubscriptionService) DeleteSubscription(email string) error {
	err := s.repository.DeleteCustomerByEmail(email)
	if err != nil {
		log.Error("failed to delete subscription", err)
	} else {
		log.Info("subscription deleted for email: ", email)
	}
	return err
}
