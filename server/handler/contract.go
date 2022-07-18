package handler

import (
	"context"
	"net/http"
)

type Handler interface {
	http.Handler
}

type TokenStorage interface {
	Get(ctx context.Context, k string) (uint, error)
	Put(ctx context.Context, k string, targetBits uint) error
	Use(ctx context.Context, k string) error
	Verify(ctx context.Context, k string) error
}

type PoW interface {
	Verify(payload []byte, timestamp int64, targetBits uint, nonce int) bool
}
