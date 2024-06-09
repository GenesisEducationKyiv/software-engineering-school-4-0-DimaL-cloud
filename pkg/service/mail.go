package service

import (
	"net/smtp"
)

type MailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type MailService struct {
	config MailConfig
}

func NewMailService(config MailConfig) *MailService {
	return &MailService{config: config}
}

func (m MailService) SendEmails(subject string, body string, to []string) error {
	err := smtp.SendMail(m.config.Host+":"+m.config.Port, smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host), m.config.Username, to, []byte("Subject: "+subject+"\r\n\r\n"+body))
	if err != nil {
		return err
	}
	return nil
}
