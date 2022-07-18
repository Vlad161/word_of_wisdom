package handler

import (
	"context"
	"net/http"
	"time"

	"word_of_wisdom/server/jwt"
)

type Handler interface {
	http.Handler
}

type TokenStorage interface {
	Get(ctx context.Context, k string) error
	Put(ctx context.Context, k string) error
	Use(ctx context.Context, k string) error
}

type PoW interface {
	Verify(payload []byte, timestamp int64, targetBits uint, nonce int) bool
}

type JWTService interface {
	CreateToken(data map[string]interface{}, exp time.Time, alg jwt.Alg) (string, error)
	Verify(token string) (map[string]interface{}, error)
}
