package handler

import (
	"net/http"
)

type Handler interface {
	http.Handler
}

type TokenStorage interface {
	Get(k string) (uint, error)
	Put(k string, targetBits uint) error
	Use(k string) error
	Verify(k string) error
}

type PoW interface {
	Verify(payload []byte, timestamp int64, targetBits uint, nonce int) bool
}
