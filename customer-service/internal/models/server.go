package models

import (
	"customer-service/internal/configs"
	"net/http"
	"time"
)

const (
	MaxHeaderBytesSize = 1 << 20
	TimeoutDuration    = 10 * time.Second
)

func NewServer(handler http.Handler, config *configs.Server) *http.Server {
	return &http.Server{
		Addr:           ":" + config.Port,
		Handler:        handler,
		MaxHeaderBytes: MaxHeaderBytesSize,
		ReadTimeout:    TimeoutDuration,
		WriteTimeout:   TimeoutDuration,
	}
}
