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
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())

	ticker := time.NewTicker(1 * time.Second)

	cl := http.NewClient("http://localhost:8080", &gohttp.Client{
		Timeout: 3 * time.Second,
	})

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println(cl.GetQuote(ctx))
			}
		}
	}()

	// Waiting OS signals or context cancellation
	wait(ctx)
	cancelCtx()
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
