package service

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/pkg/configs"
	"net/smtp"
)

type MailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type MailService struct {
	config configs.Mail
}

func NewMailService(config configs.Mail) *MailService {
	return &MailService{config: config}
}

func (m MailService) SendEmails(subject string, body string, to []string) error {
	auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)
	msg := []byte("Subject: " + subject + "\r\n\r\n" + body)
	err := smtp.SendMail(
		m.config.Host+":"+m.config.Port,
		auth,
		m.config.Username,
		to,
		msg,
	)
	if err != nil {
		return err
	}
	return nil
}
