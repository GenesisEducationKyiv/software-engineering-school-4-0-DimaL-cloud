package scheduler

import (
	"fmt"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/repository"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/service"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

const (
	SubjectName = "Курс НБУ"
	BodyFormat  = "Курс долара НБУ станом на %s: %f грн"
)

type ExchangeRateNotificationScheduler struct {
	subscriptionRepository repository.Subscription
	rateService            service.Rate
	mailService            service.Mail
}

func NewExchangeRateNotificationScheduler(
	subscriptionRepository repository.Subscription,
	rateService service.Rate,
	mailService service.Mail,
) *ExchangeRateNotificationScheduler {
	return &ExchangeRateNotificationScheduler{
		subscriptionRepository: subscriptionRepository,
		rateService:            rateService,
		mailService:            mailService,
	}
}

func (e *ExchangeRateNotificationScheduler) StartJob() {
	c := cron.New()
	err := c.AddFunc(viper.GetString("rate.notification_cron"), func() {
		currentDate := time.Now().Format("02.01.2006")
		rate, err := e.rateService.GetRate()
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
		body := fmt.Sprintf(BodyFormat, currentDate, rate)
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
