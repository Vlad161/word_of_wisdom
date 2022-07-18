package main

import (
	"context"
	gohttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"word_of_wisdom/client/http"
	"word_of_wisdom/env"
	"word_of_wisdom/logger"
	"word_of_wisdom/pow"
)

var (
	serverHost = env.GetString("SERVER_HOST", "http://localhost:8080")
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	log := logger.New()
	defer log.Sync()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	cl := http.NewClient(serverHost, &gohttp.Client{Timeout: 3 * time.Second}, pow.NewHashCash())

	go func() {
	LOOP:
		for {
			select {
			case <-ticker.C:
				log.Info(cl.GetQuote(ctx))
			case <-ctx.Done():
				break LOOP
			}
		}
	}()

	// Waiting OS signals or context cancellation
	wait(ctx, log)
}

func wait(ctx context.Context, log logger.Logger) {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-osSignals:
	case <-ctx.Done():
		log.Error("main context was canceled:", ctx.Err())
	}

	log.Error("termination signal received")
}
