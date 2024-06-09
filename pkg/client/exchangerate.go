package client

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

type ExchangeRateClient struct {
	client *http.Client
}

func NewExchangeRateClient() *ExchangeRateClient {
	return &ExchangeRateClient{
		client: &http.Client{},
	}
}

type ExchangeRateResponse struct {
	Rate         float64
	ExchangeDate string
}

func (e *ExchangeRateClient) GetCurrentExchangeRate() (ExchangeRateResponse, error) {
	resp, err := e.client.Get(viper.GetString("exchange_rate.api_url"))
	if err != nil {
		log.Errorf("failed to fetch exchange rate: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read exchange rate response: %s", err.Error())
	}
	var exchangeRates []ExchangeRateResponse
	if err := json.Unmarshal(body, &exchangeRates); err != nil {
		log.Errorf("failed to unmarshal exchange rate: %s", err.Error())
	}
	return exchangeRates[0], nil
}
