package main

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	_ "rate-service/docs"
	"syscall"
	"time"
)

const (
	ShutdownTimeout = 5 * time.Second
)

// @title Exchange rate notifier API
// @version 1.0
// @description API server for notifying exchange rate
func main() {
	app := NewApp()

	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("failed to start application: %s", err.Error())
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Exchange rate notifier shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		log.Errorf("error occurred while shutting down server: %s", err.Error())
	}
}
