package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"rate-service/internal/models"
)

const (
	eventsTable = "event"
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
	query := fmt.Sprintf("INSERT INTO %s (type, timestamp, body) VALUES ($1, $2, $3)", eventsTable)
	_, err := e.db.Exec(query, event.Type, event.Timestamp, event.Body)
	return err
}
