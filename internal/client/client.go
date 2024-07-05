package client

import (
	"net/http"
	"time"
)

const (
	TimeoutDuration = 30 * time.Second
)

func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: TimeoutDuration,
	}
}
