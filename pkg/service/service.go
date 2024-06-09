package service

import (
	exchangeratenotifierapi "exchange-rate-notifier-api/pkg/models"
	"exchange-rate-notifier-api/pkg/repository"
)

type Subscription interface {
	GetAllSubscriptions() ([]exchangeratenotifierapi.Subscription, error)
	CreateSubscription(email string) error
	DeleteSubscription(email string) error
}

type Mail interface {
	SendEmails(subject string, body string, to []string) error
}

type Service struct {
	Subscription
	Mail
}

func NewService(repositories *repository.Repository, mailConfig MailConfig) *Service {
	return &Service{
		Subscription: NewSubscriptionService(repositories.Subscription),
		Mail:         NewMailService(mailConfig),
	}
}
