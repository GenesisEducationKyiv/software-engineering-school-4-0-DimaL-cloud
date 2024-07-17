package models

type SendEmailCommand struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	To      string `json:"to"`
}

type CreateCustomerCommand struct {
	Email          string `json:"email"`
	SubscriptionID int    `json:"subscription_id"`
}
