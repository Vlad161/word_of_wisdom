package handler

import (
	"net/http"
)

type Handler interface {
	http.Handler
}

type TokenStorage interface {
	Put(v string)
	Verify(v string) bool
}
