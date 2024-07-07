package repository

import (
	"github.com/jmoiron/sqlx"
	"scheduler-service/internal/models"
)

type Event interface {
	SaveEvent(event models.Event) error
}

type EventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (e *EventRepository) SaveEvent(event models.Event) error {
	query := "INSERT INTO events (type, timestamp, body) VALUES ($1, $2, $3)"
	_, err := e.db.Exec(query, event.Type, event.Timestamp, event.Body)
	return err
}
