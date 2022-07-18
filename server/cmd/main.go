package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v9"

	"word_of_wisdom/env"
	"word_of_wisdom/logger"
	"word_of_wisdom/pow"
	"word_of_wisdom/server/handler"
	"word_of_wisdom/server/storage"
	"word_of_wisdom/server/token"
)

var (
	port                = env.GetInt("PORT", 8081)
	redisHost           = env.GetString("REDIS_HOST", "localhost:6379")
	authTokenLifetime   = env.GetDuration("AUTH_TOKEN_LIFETIME", 10*time.Second)
	authTokenTargetBits = env.GetInt("AUTH_TOKEN_TARGET_BITS", 14)
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	log := logger.New()
	defer log.Sync()

	rdb := redis.NewClient(&redis.Options{Addr: redisHost})

	//localStorage := storage.NewLocalTemporary(ctx, authTokenLifetime)
	redisStorage := token.NewStorageBytesAdapter(storage.NewRedis(rdb, authTokenLifetime))

	tokenStorage := token.NewOnetimeStorage(redisStorage)
	powAlg := pow.NewHashCash()
	challengeHandler := handler.NewChallengeHandler(log, uint(authTokenTargetBits), tokenStorage, powAlg)

	mux := http.NewServeMux()
	mux.HandleFunc("/quote", handler.AuthMW(handler.QuoteHandlerFunc(), tokenStorage))
	mux.HandleFunc("/challenge", challengeHandler.Handler())

	server := http.Server{
		Addr:    fmt.Sprint(":", port),
		Handler: mux,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("server listen and server error:", err)
		}
	}()

	// Waiting OS signals or context cancellation
	wait(ctx, log)

	ctxShutdown, cancelCtxShutdown := context.WithTimeout(ctx, 5*time.Second)
	defer cancelCtxShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Error("shutdown error:", err)
	}
}

func wait(ctx context.Context, log logger.Logger) {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-osSignals:
	case <-ctx.Done():
		log.Error("main context was canceled:", ctx.Err())
	}

	log.Info("termination signal received")
}
