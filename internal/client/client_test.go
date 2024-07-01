package client

import (
	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetCurrentExchangeRate_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", viper.GetString("exchange_rate.api_url"),
		httpmock.NewStringResponder(200, `[{"Rate": 40.5, "ExchangeDate": "2024-16-06"}]`))

	client := NewExchangeRateClient()

	resp, err := client.GetCurrentExchangeRate()

	assert.NoError(t, err)
	assert.Equal(t, 40.5, resp.Rate)
	assert.Equal(t, "2024-16-06", resp.ExchangeDate)
}

func TestGetCurrentExchangeRate_Retries(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	counter := 0
	httpmock.RegisterResponder("GET", viper.GetString("exchange_rate.api_url"),
		func(_ *http.Request) (*http.Response, error) {
			counter++
			if counter < RetriesAmount {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return httpmock.NewStringResponse(200, `[{"Rate": 40.5, "ExchangeDate": "2024-16-06"}]`), nil
		})

	client := NewExchangeRateClient()

	resp, err := client.GetCurrentExchangeRate()

	assert.NoError(t, err)
	assert.Equal(t, 40.5, resp.Rate)
	assert.Equal(t, "2024-16-06", resp.ExchangeDate)
	assert.Equal(t, RetriesAmount, counter)
}

func TestGetCurrentExchangeRate_Error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", viper.GetString("exchange_rate.api_url"),
		httpmock.NewStringResponder(500, ""))

	client := NewExchangeRateClient()

	_, err := client.GetCurrentExchangeRate()

	assert.Error(t, err)
}
