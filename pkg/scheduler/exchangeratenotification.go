package scheduler

import (
	"exchange-rate-notifier-api/pkg/client"
	"exchange-rate-notifier-api/pkg/repository"
	"exchange-rate-notifier-api/pkg/service"
	"fmt"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ExchangeRateNotificationScheduler struct {
	subscriptionRepository repository.Subscription
	exchangeRateClient     client.ExchangeRate
	mailService            service.Mail
}

func NewExchangeRateNotificationScheduler(subscriptionRepository repository.Subscription, exchangeRateClient client.ExchangeRate, mailService service.Mail) *ExchangeRateNotificationScheduler {
	return &ExchangeRateNotificationScheduler{
		subscriptionRepository: subscriptionRepository,
		exchangeRateClient:     exchangeRateClient,
		mailService:            mailService,
	}
}

func (e ExchangeRateNotificationScheduler) StartJob() {
	c := cron.New()
	err := c.AddFunc(viper.GetString("exchange_rate.notification_cron"), func() {
		exchangerate, err := e.exchangeRateClient.GetCurrentExchangeRate()
		if err != nil {
			return
		}
		subscriptions, err := e.subscriptionRepository.GetAllSubscriptions()
		if err != nil {
			return
		}
		var emails []string
		for _, subscription := range subscriptions {
			emails = append(emails, subscription.Email)
		}
		err = e.mailService.SendEmails("Курс НБУ", fmt.Sprintf("Курс долара НБУ станом на %s: %f грн", exchangerate.ExchangeDate, exchangerate.Rate), emails)
		if err != nil {
			return
		}
	})
	if err != nil {
		log.Fatalf("failed to schedule exchange rate notification job: %s", err.Error())
	}
	c.Start()
}
