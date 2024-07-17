package service

import (
	"customer-service/internal/configs"
	"customer-service/internal/repository"
	log "github.com/sirupsen/logrus"
)

type Subscription interface {
	CreateCustomer(email string) (int, error)
	DeleteCustomerByEmail(email string) error
}

type SubscriptionService struct {
	repository     repository.Customer
	rabbitMQConfig *configs.RabbitMQ
}

func NewSubscriptionService(repository repository.Subscription) *SubscriptionService {
	return &SubscriptionService{repository: repository}
}

func (s *SubscriptionService) GetAllSubscriptions() ([]models.Subscription, error) {
	return s.repository.GetAllSubscriptions()
}

func (s *SubscriptionService) CreateSubscription(email string) error {
	id, err := s.repository.CreateSubscription(email)
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
	err := s.repository.DeleteSubscription(email)
	if err != nil {
		log.Error("failed to delete subscription", err)
	} else {
		log.Info("subscription deleted for email: ", email)
	}
	return err
}
