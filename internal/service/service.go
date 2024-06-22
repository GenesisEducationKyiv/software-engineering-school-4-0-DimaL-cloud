package service

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/configs"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/models"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/repository"
)

type Subscription interface {
	GetAllSubscriptions() ([]models.Subscription, error)
	CreateSubscription(email string) error
	DeleteSubscription(email string) error
}

type Mail interface {
	SendEmails(subject string, body string, to []string) error
}

type Rate interface {
	GetRate() (float64, error)
}

type Service struct {
	Subscription
	Mail
	Rate
}

func NewService(repositories *repository.Repository, config *configs.Config, clients *client.Client) *Service {
	return &Service{
		Subscription: NewSubscriptionService(repositories.Subscription),
		Mail:         NewMailService(config.Mail),
		Rate:         NewRateService(clients),
	}
}
