package client

type ExchangeRate interface {
	GetCurrentExchangeRate() (ExchangeRateResponse, error)
}

type Client struct {
	ExchangeRate
}

func NewClient() *Client {
	return &Client{
		ExchangeRate: NewExchangeRateClient(),
	}
}
