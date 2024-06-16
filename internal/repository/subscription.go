package repository

import (
	"errors"
	"fmt"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
)

type SubscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (s *SubscriptionRepository) GetAllSubscriptions() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	query := "SELECT * FROM " + subscriptionsTable
	if err := s.db.Select(&subscriptions, query); err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (s *SubscriptionRepository) CreateSubscription(email string) error {
	query := fmt.Sprintf("INSERT INTO %s (email) VALUES ($1)", subscriptionsTable)
	_, err := s.db.Exec(query, email)
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) {
			if pgError.Code == "23505" && strings.Contains(pgError.Message, "subscription_email_key") {
				return models.ErrEmailAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (s *SubscriptionRepository) DeleteSubscription(email string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE email = $1", subscriptionsTable)
	row := s.db.QueryRow(query, email)
	return row.Err()
}
