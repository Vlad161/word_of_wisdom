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
}
