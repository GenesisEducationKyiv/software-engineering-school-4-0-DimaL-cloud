package repository

import (
	exchangeratenotifierapi "exchange-rate-notifier-api/pkg/models"
	"github.com/jmoiron/sqlx"
)

type Subscription interface {
	GetAllSubscriptions() ([]exchangeratenotifierapi.Subscription, error)
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
