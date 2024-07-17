package service

import (
	log "github.com/sirupsen/logrus"
	"rate-service/internal/configs"
	"rate-service/internal/messaging/producer"
	"rate-service/internal/models"
	"rate-service/internal/repository"
)

const (
	CreateCustomerCommand = "CreateCustomerCommand"
)

type Subscription interface {
	GetAllSubscriptions() ([]models.Subscription, error)
	CreateSubscription(email string) (int, error)
	DeleteSubscriptionByEmail(email string) error
	DeleteSubscriptionByID(id int) error
}

type SubscriptionService struct {
	repository      repository.Subscription
	messageProducer producer.Producer
	rabbitMQConfig  *configs.RabbitMQ
}

func NewSubscriptionService(
	repository repository.Subscription,
	rabbitMQConfig *configs.RabbitMQ,
	messageProducer producer.Producer,
) *SubscriptionService {
	return &SubscriptionService{
		repository:      repository,
		messageProducer: messageProducer,
		rabbitMQConfig:  rabbitMQConfig,
	}
}

func (s *SubscriptionService) GetAllSubscriptions() ([]models.Subscription, error) {
	return s.repository.GetAllSubscriptions()
}

func (s *SubscriptionService) CreateSubscription(email string) (int, error) {
	id, err := s.repository.CreateSubscription(email)
	if err != nil {
		log.Error("failed to create subscription: ", err)
		return 0, err
	}
	log.Info("subscription created for email: ", email)

	createCustomerCommand := models.CreateCustomerCommand{
		Command: models.Command{
			Type: CreateCustomerCommand,
		},
		Email:          email,
		SubscriptionID: id,
	}
	s.messageProducer.PublishMessage(createCustomerCommand, s.rabbitMQConfig.Queue.CustomerCommand)
	return id, err
}

func (s *SubscriptionService) DeleteSubscriptionByEmail(email string) error {
	err := s.repository.DeleteSubscriptionByEmail(email)
	if err != nil {
		log.Error("failed to delete subscription", err)
	} else {
		log.Info("subscription deleted for email: ", email)
	}
	return err
}

func (s *SubscriptionService) DeleteSubscriptionByID(id int) error {
	err := s.repository.DeleteSubscriptionByID(id)
	if err != nil {
		log.Error("failed to delete subscription", err)
	} else {
		log.Info("subscription deleted for id: ", id)
	}
	return err
}
