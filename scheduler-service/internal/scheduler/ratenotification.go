package scheduler

import (
	"encoding/json"
	"github.com/VictoriaMetrics/metrics"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"scheduler-service/internal/configs"
	"scheduler-service/internal/models"
	"time"
)

const (
	EventType = "RateSchedulerTriggeredEvent"
)

var (
	rateNotificationJobTotal       = metrics.GetOrCreateCounter(`rate_notification_job_total`)
	rateNotificationJobErrorsTotal = metrics.GetOrCreateCounter(`rate_notification_job_errors_total`)
)

type RateNotificationScheduler struct {
	cronConfig     *configs.Crons
	rabbitMQConfig *configs.RabbitMQ
	channel        *amqp.Channel
}

func NewRateNotificationScheduler(
	cronConfig *configs.Crons,
	rabbitMQConfig *configs.RabbitMQ,
	channel *amqp.Channel,
) *RateNotificationScheduler {
	return &RateNotificationScheduler{
		cronConfig:     cronConfig,
		rabbitMQConfig: rabbitMQConfig,
		channel:        channel,
	}
}

func (e *RateNotificationScheduler) StartJob() {
	c := cron.New()
	err := c.AddFunc(e.cronConfig.RateNotification, func() {
		rateNotificationJobTotal.Inc()
		event := models.Event{
			Type:      EventType,
			Timestamp: time.Now(),
		}
		serializedEvent, err := json.Marshal(event)
		if err != nil {
			log.Fatalf("failed to serialize event: %s", err.Error())
			rateNotificationJobErrorsTotal.Inc()
		}
		err = e.channel.Publish(
			"",
			e.rabbitMQConfig.Queue.RateNotificationCron,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        serializedEvent,
			},
		)
		if err != nil {
			log.Fatalf("failed to publish event: %s", err.Error())
			rateNotificationJobErrorsTotal.Inc()
		} else {
			log.Info("event published")
		}
	})
	if err != nil {
		log.Fatalf("failed to schedule rate notification job: %s", err.Error())
	}
	c.Start()
}
