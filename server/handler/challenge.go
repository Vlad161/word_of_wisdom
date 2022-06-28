package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"word_of_wisdom/pow"
	"word_of_wisdom/server/token"
)

type (
	getChallengeRespBody struct {
		Timestamp  int64  `json:"timestamp"`
		Token      string `json:"token"`
		TargetBits uint   `json:"targetBits"`
	}

	postChallengeReqBody struct {
		Timestamp  int64  `json:"timestamp"`
		Token      string `json:"token"`
		TargetBits uint   `json:"targetBits"`
		Nonce      int    `json:"nonce"`
	}

	challengeHandler struct {
		targetBits uint
		tStorage   TokenStorage
	}
)

func NewChallengeHandler(targetBits uint, tokenStorage TokenStorage) *challengeHandler {
	return &challengeHandler{
		targetBits: targetBits,
		tStorage:   tokenStorage,
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

func (h *challengeHandler) challengeRequest(w http.ResponseWriter, _ *http.Request) {
	tc := token.New()

	respBody := getChallengeRespBody{
		Timestamp:  time.Now().Unix(),
		Token:      tc,
		TargetBits: h.targetBits,
	}
	bb, err := json.Marshal(&respBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.tStorage.Put(tc, h.targetBits)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bb)
}

func (h *challengeHandler) challengeVerify(w http.ResponseWriter, req *http.Request) {
	var reqBody postChallengeReqBody
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sTargetBits, ok := h.tStorage.TargetBits(reqBody.Token)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hc := pow.NewHashCash(reqBody.Timestamp, reqBody.Token, reqBody.TargetBits)
	if sTargetBits == reqBody.TargetBits && hc.Verify(reqBody.Nonce) && h.tStorage.Verify(reqBody.Token) {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}
