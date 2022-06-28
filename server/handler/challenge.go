package handler

import (
	"net/http"

	"word_of_wisdom/server/token"
)

type challengeHandler struct {
	tStorage TokenStorage
}

func NewChallengeHandler(tokenStorage TokenStorage) *challengeHandler {
	return &challengeHandler{
		tStorage: tokenStorage,
	}
}

func (h *challengeHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			h.challengeRequest(w, req)
		case http.MethodPost:
			h.challengeVerify(w, req)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

func (h *challengeHandler) challengeRequest(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *challengeHandler) challengeVerify(w http.ResponseWriter, req *http.Request) {
	t := token.New()
	h.tStorage.Put(t, uint(1))

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(t))
}
