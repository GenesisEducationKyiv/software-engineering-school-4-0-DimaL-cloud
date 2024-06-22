package client

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/client/rate"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/configs"
	"net/http"
)

type ExchangeRate interface {
	GetCurrentExchangeRate() (ExchangeRateResponse, error)
}

type Client struct {
	NbuRate
	PrivatBankRate
}

type NbuRate interface {
	GetRate() (float64, error)
	GetNext() rate.Rate
	SetNext(next rate.Rate) rate.Rate
}

type PrivatBankRate interface {
	GetRate() (float64, error)
	GetNext() rate.Rate
	SetNext(next rate.Rate) rate.Rate
}

func NewClient(client *http.Client, configs *configs.Config) *Client {
	return &Client{
		NbuRate:        rate.NewNbuRateClient(client, configs.Rate.APIUrls.Nbu),
		PrivatBankRate: rate.NewPrivatBankRateClient(client, configs.Rate.APIUrls.PrivatBank),
	}
}
