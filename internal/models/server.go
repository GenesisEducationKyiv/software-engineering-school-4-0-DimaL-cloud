package models

import (
	"context"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

const (
	MaxHeaderBytesSize = 1 << 20
	TimeoutDuration    = 10 * time.Second
)

func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + viper.GetString("server.port"),
		Handler:        handler,
		MaxHeaderBytes: MaxHeaderBytesSize,
		ReadTimeout:    TimeoutDuration,
		WriteTimeout:   TimeoutDuration,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
