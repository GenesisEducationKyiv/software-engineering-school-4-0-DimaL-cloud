package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"rate-service/internal/models"
	"strings"
)

const (
	SubscriptionTable      = "subscription"
	UniqueViolationErrCode = "23505"
)

type Subscription interface {
	GetAllSubscriptions() ([]models.Subscription, error)
	CreateCustomer(email string) (int, error)
	DeleteCustomerByEmail(email string) error
}

type SubscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (s *SubscriptionRepository) GetAllSubscriptions() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	query := "SELECT * FROM " + SubscriptionTable
	if err := s.db.Select(&subscriptions, query); err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (s *SubscriptionRepository) CreateCustomer(email string) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (email) VALUES ($1) RETURNING id", SubscriptionTable)
	var id int
	err := s.db.QueryRow(query, email).Scan(&id)
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) {
			if pgError.Code == UniqueViolationErrCode && strings.Contains(pgError.Message, "subscription_email_key") {
				return 0, models.ErrDuplicateEmail
			}
		}
		return 0, err
	}
	return id, nil
}

func (s *SubscriptionRepository) DeleteCustomerByEmail(email string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE email = $1", SubscriptionTable)
	row := s.db.QueryRow(query, email)
	return row.Err()
}
