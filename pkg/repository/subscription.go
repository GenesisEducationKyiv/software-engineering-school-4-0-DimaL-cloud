package repository

import (
	exchangeratenotifierapi "exchange-rate-notifier-api/pkg/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type SubscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (s *SubscriptionRepository) GetAllSubscriptions() ([]exchangeratenotifierapi.Subscription, error) {
	var subscriptions []exchangeratenotifierapi.Subscription
	query := fmt.Sprintf("SELECT * FROM %s", subscriptionsTable)
	if err := s.db.Select(&subscriptions, query); err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (s *SubscriptionRepository) CreateSubscription(email string) error {
	query := fmt.Sprintf("INSERT INTO %s (email) VALUES ($1) ON CONFLICT (email) DO NOTHING", subscriptionsTable)
	row := s.db.QueryRow(query, email)
	err := row.Err()
	return err
}

func (s *SubscriptionRepository) DeleteSubscription(email string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE email = $1", subscriptionsTable)
	row := s.db.QueryRow(query, email)
	return row.Err()
}
