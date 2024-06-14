package service

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/pkg/models"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/pkg/repository"
)

type Subscription interface {
	GetAllSubscriptions() ([]models.Subscription, error)
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
