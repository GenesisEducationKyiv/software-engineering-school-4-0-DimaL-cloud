package service

import (
	"errors"
	"github.com/avast/retry-go/v4"
	log "github.com/sirupsen/logrus"
	"rate-service/internal/client/rate"
)

const (
	RetriesAmount = 3
	Delay         = 100
)

type Rate interface {
	GetRate() (float64, error)
}

type RateService struct {
	initialRateClient rate.Rate
}

func NewRateService(initialRateClient rate.Rate) *RateService {
	return &RateService{
		initialRateClient: initialRateClient,
	}
}

func (r *RateService) GetRate() (float64, error) {
	var rateResponse rate.Response
	var err error
	currentClient := r.initialRateClient
	for currentClient != nil {
		err = retry.Do(
			func() error {
				rateResponse, err = currentClient.GetRate()
				return err
			},
			r.getRetryOptions()...,
		)
		if err == nil {
			log.Infof("API response from %s: %s", rateResponse.APIName, rateResponse.APIResponse)
			return rateResponse.RateValue, nil
		}
		currentClient = currentClient.GetNext()
	}
	return 0, errors.New("failed to get exchange rate from all APIs")
}

func (r *RateService) getRetryOptions() []retry.Option {
	return []retry.Option{
		retry.Attempts(uint(RetriesAmount)),
		retry.OnRetry(func(n uint, err error) {
			log.Infof("Retry request %d to and get error: %v", n+1, err)
		}),
		retry.Delay(Delay),
	}
}
