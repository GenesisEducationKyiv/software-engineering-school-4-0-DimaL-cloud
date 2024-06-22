package service

import (
	"errors"
	clients "github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/client/rate"
)

type RateService struct {
	client rate.Rate
}

func NewRateService(clients *clients.Client) *RateService {
	r := &RateService{}
	r.client = clients.NbuRate
	r.client.SetNext(clients.PrivatBankRate)
	return r
}

func (r *RateService) GetRate() (float64, error) {
	currentClient := r.client

	for currentClient != nil {
		value, err := currentClient.GetRate()
		if err == nil {
			return value, nil
		}
		currentClient = currentClient.GetNext()
	}

	return 0, errors.New("failed to get exchange rate from all clients")
}
