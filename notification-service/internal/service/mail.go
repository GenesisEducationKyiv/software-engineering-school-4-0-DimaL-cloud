package service

import (
	log "github.com/sirupsen/logrus"
	"net/smtp"
	"notification-service/internal/configs"
)

type Mail interface {
	SendEmail(subject string, body string, to string)
}

type MailService struct {
	config configs.Mail
	auth   smtp.Auth
}

func NewMailService(config configs.Mail) *MailService {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	return &MailService{config: config, auth: auth}
}

func (m MailService) SendEmail(subject string, body string, to string) {
	msg := []byte("Subject: " + subject + "\r\n\r\n" + body)

	err := smtp.SendMail(
		m.config.Host+":"+m.config.Port,
		m.auth,
		m.config.Username,
		[]string{to},
		msg,
	)
	if err != nil {
		log.Errorf("failed to send email to %s: %s", to, err.Error())
	}
}
