package client

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/client/rate"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/configs"
	"net/http"
)

type Client struct {
	NbuRate        rate.Rate
	PrivatBankRate rate.Rate
	FawazahmedRate rate.Rate
}

func NewClient(client *http.Client, configs *configs.Config) *Client {
	return &Client{
		NbuRate:        rate.NewNbuRateClient(client, configs.Rate.APIUrls.Nbu),
		PrivatBankRate: rate.NewPrivatBankRateClient(client, configs.Rate.APIUrls.PrivatBank),
		FawazahmedRate: rate.NewFawazahmedRateClient(client, configs.Rate.APIUrls.Fawazahmed),
	}
}
