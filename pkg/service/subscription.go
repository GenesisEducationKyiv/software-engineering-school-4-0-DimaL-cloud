package service

import (
	exchangeratenotifierapi "exchange-rate-notifier-api/pkg/models"
	"exchange-rate-notifier-api/pkg/repository"
	log "github.com/sirupsen/logrus"
)

type SubscriptionService struct {
	repository repository.Subscription
}

func NewSubscriptionService(repository repository.Subscription) *SubscriptionService {
	return &SubscriptionService{repository: repository}
}

func (s *SubscriptionService) GetAllSubscriptions() ([]exchangeratenotifierapi.Subscription, error) {
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
