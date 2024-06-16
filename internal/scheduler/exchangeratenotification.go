package scheduler

import (
	"fmt"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/repository"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/service"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	SubjectName = "Курс НБУ"
	BodyFormat  = "Курс долара НБУ станом на %s: %f грн"
)

type ExchangeRateNotificationScheduler struct {
	subscriptionRepository repository.Subscription
	exchangeRateClient     client.ExchangeRate
	mailService            service.Mail
}

func NewExchangeRateNotificationScheduler(
	subscriptionRepository repository.Subscription,
	exchangeRateClient client.ExchangeRate,
	mailService service.Mail,
) *ExchangeRateNotificationScheduler {
	return &ExchangeRateNotificationScheduler{
		subscriptionRepository: subscriptionRepository,
		exchangeRateClient:     exchangeRateClient,
		mailService:            mailService,
	}
}

func (e *ExchangeRateNotificationScheduler) StartJob() {
	c := cron.New()
	err := c.AddFunc(viper.GetString("exchange_rate.notification_cron"), func() {
		rate, err := e.exchangeRateClient.GetCurrentExchangeRate()
		if err != nil {
			log.Errorf("failed to get current exchange rate: %s", err.Error())
			return
		}
		subscriptions, err := e.subscriptionRepository.GetAllSubscriptions()
		if err != nil {
			log.Errorf("failed to get subscriptions: %s", err.Error())
			return
		}
		var emails []string
		for _, subscription := range subscriptions {
			emails = append(emails, subscription.Email)
		}
		body := fmt.Sprintf(BodyFormat, rate.ExchangeDate, rate.Rate)
		err = e.mailService.SendEmails(
			SubjectName,
			body,
			emails,
		)
		if err != nil {
			log.Errorf("failed to send emails: %s", err.Error())
			return
		}
	})
	if err != nil {
		log.Fatalf("failed to schedule exchange rate notification job: %s", err.Error())
	}
	c.Start()
}
