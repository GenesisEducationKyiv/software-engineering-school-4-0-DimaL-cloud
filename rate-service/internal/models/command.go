package models

type Command struct {
	Type string `json:"type"`
}

type SendEmailCommand struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	To      string `json:"to"`
}

type CreateCustomerCommand struct {
	Command
	Email          string `json:"email"`
	SubscriptionID int    `json:"subscription_id"`
}
