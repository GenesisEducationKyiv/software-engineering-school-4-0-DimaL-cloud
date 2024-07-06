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

func (nrc *NbuRateClient) GetRate() (Response, error) {
	resp, err := nrc.client.Get(nrc.apiURL)
	if err != nil {
		log.Errorf("failed to fetch NBU exchange rate: %s", err.Error())
		return Response{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read NBU exchange rate response: %s", err.Error())
		return Response{}, err
	}
	return nrc.parseRateResponse(body)
}

func (nrc *NbuRateClient) parseRateResponse(body []byte) (Response, error) {
	var exchangeRate NbuRateResponse
	var exchangeRates []NbuRateResponse
	if err := json.Unmarshal(body, &exchangeRates); err != nil {
		log.Errorf("failed to unmarshal exchange rate: %s", err.Error())
		return Response{}, err
	}
	exchangeRate = exchangeRates[0]
	return Response{
		RateValue:   exchangeRate.Rate,
		APIResponse: string(body),
		APIName:     "NBU",
	}, nil
}
