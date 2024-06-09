package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	subscriptionsTable = "subscription"
)

type DBConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	DBName     string
	DriverName string
	SSLMode    string
}

func NewDB(config DBConfig) (*sqlx.DB, error) {
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
