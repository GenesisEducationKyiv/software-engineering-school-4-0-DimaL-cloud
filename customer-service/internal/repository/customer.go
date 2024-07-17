package repository

import (
	"customer-service/internal/models"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
)

const (
	CustomerTable          = "customer"
	UniqueViolationErrCode = "23505"
)

type Customer interface {
	CreateCustomer(customer models.Customer) (int, error)
	DeleteCustomerByEmail(email string) error
}

type CustomerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (s *CustomerRepository) CreateCustomer(customer models.Customer) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (email, subscription_id) VALUES ($1, $2) RETURNING id", CustomerTable)
	var id int
	err := s.db.QueryRow(query, customer.Email, customer.SubscriptionID).Scan(&id)
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) {
			if pgError.Code == UniqueViolationErrCode && strings.Contains(pgError.Message, "customer_email_key") {
				return 0, models.ErrDuplicateEmail
			}
		}
		return 0, err
	}
	return id, nil
}

func (s *CustomerRepository) DeleteCustomer(email string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE email = $1", CustomerTable)
	row := s.db.QueryRow(query, email)
	return row.Err()
}
