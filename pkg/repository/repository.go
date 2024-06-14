package repository

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/pkg/models"
	"github.com/jmoiron/sqlx"
)

type Subscription interface {
	GetAllSubscriptions() ([]models.Subscription, error)
	CreateSubscription(email string) error
	DeleteSubscription(email string) error
}

type Repository struct {
	Subscription
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Subscription: NewSubscriptionRepository(db),
	}
}
