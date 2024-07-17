package models

import "time"

type Event struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Body      string    `json:"body"`
}

type CustomerCreatedEvent struct {
	Event
	CustomerID int `json:"customer_id"`
}

type CustomerCreationFailedEvent struct {
	Event
	SubscriptionID int `json:"subscription_id"`
}
