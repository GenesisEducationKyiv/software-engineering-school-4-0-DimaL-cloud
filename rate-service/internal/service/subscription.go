package service

import (
	log "github.com/sirupsen/logrus"
	"rate-service/internal/models"
	"rate-service/internal/repository"
)

type Subscription interface {
	GetAllSubscriptions() ([]models.Subscription, error)
	CreateSubscription(email string) error
	DeleteSubscription(email string) error
}

type SubscriptionService struct {
	repository repository.Subscription
}

func NewSubscriptionService(repository repository.Subscription) *SubscriptionService {
	return &SubscriptionService{repository: repository}
}

func (s *SubscriptionService) GetAllSubscriptions() ([]models.Subscription, error) {
	return s.repository.GetAllSubscriptions()
}

func (s *SubscriptionService) CreateSubscription(email string) error {
	err := s.repository.CreateSubscription(email)
	if err != nil {
		log.Error("failed to create subscription", err)
	} else {
		log.Info("subscription created for email: ", email)
	}
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
