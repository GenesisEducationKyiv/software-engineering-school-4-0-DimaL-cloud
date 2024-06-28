package rate

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type NbuRateClient struct {
	client *http.Client
	apiURL string
	ClientChain
}

func NewNbuRateClient(client *http.Client, apiURL string) *NbuRateClient {
	return &NbuRateClient{
		client: client,
		apiURL: apiURL,
	}
}

type NbuRateResponse struct {
	Rate float64 `json:"rate"`
}

func (nrc *NbuRateClient) GetRate() (float64, error) {
	resp, err := nrc.client.Get(nrc.apiURL)
	if err != nil {
		log.Errorf("failed to fetch NBU exchange rate: %s", err.Error())
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read NBU exchange rate response: %s", err.Error())
		return 0, err
	}
	log.Printf("NBU rate API response: %s", string(body))
	return nrc.parseRateResponse(body)
}

func (nrc *NbuRateClient) parseRateResponse(body []byte) (float64, error) {
	var exchangeRate NbuRateResponse
	var exchangeRates []NbuRateResponse
	if err := json.Unmarshal(body, &exchangeRates); err != nil {
		log.Errorf("failed to unmarshal exchange rate: %s", err.Error())
		return 0, err
	}
	exchangeRate = exchangeRates[0]
	return exchangeRate.Rate, nil
}
