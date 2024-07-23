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
	ClientChain
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

func (f *FawazahmedRateClient) GetRate() (Response, error) {
	resp, err := f.client.Get(f.apiURL)
	if err != nil {
		log.Printf("failed to fetch Fawazahmed rate API: %s", err.Error())
		return Response{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read Fawazahmed exchange rate response: %s", err.Error())
		return Response{}, err
	}
	return f.parseRateResponse(body)
}

func (f *FawazahmedRateClient) parseRateResponse(body []byte) (Response, error) {
	var exchangeRate FawazahmedRateResponse
	if err := json.Unmarshal(body, &exchangeRate); err != nil {
		log.Printf("failed to unmarshal exchange rate: %s", err.Error())
		return Response{}, err
	}
	return Response{
		RateValue:   exchangeRate.USD.UAH,
		APIResponse: string(body),
		APIName:     "Fawazahmed",
	}, nil
}
