package service

import (
	"errors"
	clients "github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/client/rate"
	"github.com/avast/retry-go/v4"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	RetriesAmount = 3
)

type RateService struct {
	client rate.Rate
}

func NewRateService(clients *clients.Client) *RateService {
	r := &RateService{}
	r.client = clients.NbuRate
	r.client.SetNext(clients.PrivatBankRate).
		SetNext(clients.FawazahmedRate)
	return r
}

func (r *RateService) GetRate() (float64, error) {
	var rateResponse rate.Response
	var err error
	currentClient := r.client
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
		retry.Delay(time.Second),
	}
}
