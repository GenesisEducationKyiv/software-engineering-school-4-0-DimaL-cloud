package rate

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type PrivatBankRateClient struct {
	client *http.Client
	apiURL string
	ClientChain
}

func NewPrivatBankRateClient(client *http.Client, apiURL string) *PrivatBankRateClient {
	return &PrivatBankRateClient{
		client: client,
		apiURL: apiURL,
	}
}

type PrivatBankRateResponse struct {
	Rate string `json:"sale"`
}

func (p *PrivatBankRateClient) GetRate() (Response, error) {
	resp, err := p.client.Get(p.apiURL)
	if err != nil {
		log.Errorf("failed to fetch PrivatBank exchange rate: %s", err.Error())
		return Response{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read PrivatBank exchange rate response: %s", err.Error())
		return Response{}, err
	}
	return p.parseRateResponse(body)
}

func (p *PrivatBankRateClient) parseRateResponse(body []byte) (Response, error) {
	var exchangeRates []PrivatBankRateResponse
	if err := json.Unmarshal(body, &exchangeRates); err != nil {
		log.Errorf("failed to unmarshal exchange rate: %s", err.Error())
		return Response{}, err
	}

	rate, err := strconv.ParseFloat(exchangeRates[1].Rate, 64)
	if err != nil {
		log.Errorf("failed to convert rate to float: %s", err.Error())
		return Response{}, err
	}

	return Response{
		RateValue:   rate,
		APIResponse: string(body),
		APIName:     "PrivatBank",
	}, nil
}
