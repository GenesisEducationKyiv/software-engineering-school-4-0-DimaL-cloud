package models

type Command struct {
	Type string `json:"type"`
}

type CreateCustomerCommand struct {
	Command
	Email          string `json:"email"`
	SubscriptionID int    `json:"subscription_id"`
}
