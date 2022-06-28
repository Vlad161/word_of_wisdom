package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"word_of_wisdom/server/handler"
	"word_of_wisdom/server/token"
)

const (
	port              = 8080
	authTokenLifetime = 10 * time.Second
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	tokenStorage := token.NewStorage(ctx, authTokenLifetime)
	challengeHandler := handler.NewChallengeHandler(tokenStorage)

	mux := http.NewServeMux()
	mux.HandleFunc("/quote", handler.AuthMW(handler.QuoteHandlerFunc(), tokenStorage))
	mux.HandleFunc("/challenge", challengeHandler.Handler())

	server := http.Server{
		Addr:    fmt.Sprint(":", port),
		Handler: mux,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("server listen and server error:", err)
		}
	}()

	// Waiting OS signals or context cancellation
	wait(ctx)

	ctxShutdown, cancelCtxShutdown := context.WithTimeout(ctx, 5*time.Second)
	defer cancelCtxShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		fmt.Println("shutdown error:", err)
	}
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
