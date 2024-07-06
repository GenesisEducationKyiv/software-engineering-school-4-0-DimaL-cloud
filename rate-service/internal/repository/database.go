package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"rate-service/internal/configs"
)

const (
	subscriptionsTable = "subscription"
)

func NewDB(config *configs.DB) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode)
	db, err := sqlx.Open(config.DriverName, connectionString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
