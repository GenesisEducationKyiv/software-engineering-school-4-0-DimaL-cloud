package models

import "time"

type Event struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Body      string    `json:"body"`
}
