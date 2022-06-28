package handler

import (
	"net/http"
)

type Handler interface {
	http.Handler
}

type TokenStorage interface {
	Put(k string, v uint)
	Use(k string) bool
	Verify(k string) bool
	TargetBits(k string) (uint, bool)
}

type PoW interface {
	Verify(payload []byte, timestamp int64, targetBits uint, nonce int) bool
}
