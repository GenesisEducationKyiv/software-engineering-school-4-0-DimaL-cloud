package client

import (
	"encoding/json"
	"github.com/avast/retry-go/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"time"
)

const (
	RetriesAmount = 3
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
	var exchangeRate ExchangeRateResponse
	err := retry.Do(
		func() error {
			var err error
			exchangeRate, err = e.fetchExchangeRate()
			return err
		},
		e.getRetryOptions()...,
	)
	return exchangeRate, err
}

func (e *ExchangeRateClient) fetchExchangeRate() (ExchangeRateResponse, error) {
	resp, err := e.client.Get(viper.GetString("exchange_rate.api_url"))
	if err != nil {
		log.Errorf("failed to fetch exchange rate: %s", err.Error())
		return ExchangeRateResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read exchange rate response: %s", err.Error())
		return ExchangeRateResponse{}, err
	}
	return e.parseExchangeRateResponse(body)
}

func (e *ExchangeRateClient) parseExchangeRateResponse(body []byte) (ExchangeRateResponse, error) {
	var exchangeRate ExchangeRateResponse
	var exchangeRates []ExchangeRateResponse
	if err := json.Unmarshal(body, &exchangeRates); err != nil {
		log.Errorf("failed to unmarshal exchange rate: %s", err.Error())
		return exchangeRate, err
	}
	exchangeRate = exchangeRates[0]
	return exchangeRate, nil
}

func (e *ExchangeRateClient) getRetryOptions() []retry.Option {
	return []retry.Option{
		retry.Attempts(uint(RetriesAmount)),
		retry.OnRetry(func(n uint, err error) {
			log.Infof("Retry request %d to and get error: %v", n+1, err)
		}),
		retry.Delay(time.Second),
	}
}
