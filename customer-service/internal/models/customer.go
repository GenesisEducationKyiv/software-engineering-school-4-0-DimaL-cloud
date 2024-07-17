package models

type Customer struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	SubscriptionID int    `json:"subscription_id"`
}
