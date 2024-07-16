package service

import (
	log "github.com/sirupsen/logrus"
	"net/smtp"
	"notification-service/internal/configs"
	"sync"
)

const (
	maxWorkers = 100
)

type Mail interface {
	SendEmails(subject string, body string, to []string)
}

type MailService struct {
	config configs.Mail
	auth   smtp.Auth
}

func NewMailService(config configs.Mail) *MailService {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	return &MailService{config: config, auth: auth}
}

func (m MailService) SendEmails(subject string, body string, to []string) {
	msg := []byte("Subject: " + subject + "\r\n\r\n" + body)

	var wg sync.WaitGroup
	emailTasks := make(chan string, len(to))
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for email := range emailTasks {
				err := smtp.SendMail(
					m.config.Host+":"+m.config.Port,
					m.auth,
					m.config.Username,
					[]string{email},
					msg,
				)
				if err != nil {
					log.Errorf("failed to send email to %s: %s", email, err.Error())
				}
			}
		}()
	}
	for _, recipient := range to {
		emailTasks <- recipient
	}
	close(emailTasks)
	wg.Wait()
}
