package main

import (
	"context"
	"fmt"
	gohttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"word_of_wisdom/client/http"
	"word_of_wisdom/pow"
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	cl := http.NewClient("http://server:8001", &gohttp.Client{
		Timeout: 3 * time.Second,
	}, pow.NewHashCash())

	go func() {
	LOOP:
		for {
			select {
			case <-ticker.C:
				fmt.Println(cl.GetQuote(ctx))
			case <-ctx.Done():
				break LOOP
			}
		}
	}()

	// Waiting OS signals or context cancellation
	wait(ctx)
}

func wait(ctx context.Context) {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-osSignals:
	case <-ctx.Done():
		fmt.Println("main context was canceled:", ctx.Err())
	}

	fmt.Println("termination signal received")
}
