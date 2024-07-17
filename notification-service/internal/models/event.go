package models

type SendEmailCommand struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	To      string `json:"to"`
}
