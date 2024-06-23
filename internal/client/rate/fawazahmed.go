package rate

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type FawazahmedRateClient struct {
	client *http.Client
	apiURL string
	next   Rate
}

func NewFawazahmedRateClient(client *http.Client, apiURL string) *FawazahmedRateClient {
	return &FawazahmedRateClient{
		client: client,
		apiURL: apiURL,
	}
}

type FawazahmedRateResponse struct {
	USD USD `json:"usd"`
}

type USD struct {
	UAH float64 `json:"uah"`
}

func (f *FawazahmedRateClient) GetRate() (float64, error) {
	resp, err := f.client.Get(f.apiURL)
	if err != nil {
		log.Printf("failed to fetch Fawazahmed rate API: %s", err.Error())
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read Fawazahmed exchange rate response: %s", err.Error())
		return 0, err
	}
	log.Printf("Fawazahmed rate API response: %s", string(body))
	return f.parseRateResponse(body)
}

func (f *FawazahmedRateClient) parseRateResponse(body []byte) (float64, error) {
	var exchangeRate FawazahmedRateResponse
	if err := json.Unmarshal(body, &exchangeRate); err != nil {
		log.Printf("failed to unmarshal exchange rate: %s", err.Error())
		return 0, err
	}
	return exchangeRate.USD.UAH, nil
}

func (f *FawazahmedRateClient) GetNext() Rate {
	return f.next
}

func (f *FawazahmedRateClient) SetNext(next Rate) Rate {
	f.next = next
	return next
}
