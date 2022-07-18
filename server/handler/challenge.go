package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"word_of_wisdom/logger"
	"word_of_wisdom/server/token"
)

type (
	getChallengeRespBody struct {
		Timestamp  int64  `json:"timestamp"`
		Token      string `json:"token"`
		TargetBits uint   `json:"target_bits"`
	}

	postChallengeReqBody struct {
		Timestamp  int64  `json:"timestamp"`
		Token      string `json:"token"`
		TargetBits uint   `json:"target_bits"`
		Nonce      int    `json:"nonce"`
	}

	challengeHandler struct {
		log logger.Logger

		targetBits uint
		tStorage   TokenStorage
		pow        PoW
	}
)

func NewChallengeHandler(log logger.Logger, targetBits uint, tokenStorage TokenStorage, pow PoW) *challengeHandler {
	return &challengeHandler{
		log:        log,
		targetBits: targetBits,
		tStorage:   tokenStorage,
		pow:        pow,
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
	if err = h.tStorage.Put(req.Context(), tc, h.targetBits); err != nil {
		h.log.Error("challenge request, can't put to token storage", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bb)
}

func (h *challengeHandler) challengeVerify(w http.ResponseWriter, req *http.Request) {
	var reqBody postChallengeReqBody
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sTargetBits, err := h.tStorage.Get(req.Context(), reqBody.Token)
	if err != nil {
		h.log.Error("challenge verify, can't get from token storage", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isPowVerify := h.pow.Verify([]byte(reqBody.Token), reqBody.Timestamp, reqBody.TargetBits, reqBody.Nonce)
	if sTargetBits == reqBody.TargetBits && isPowVerify {
		if err = h.tStorage.Verify(req.Context(), reqBody.Token); err == nil {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	h.log.Error("challenge verify, can't verify", err)
	w.WriteHeader(http.StatusInternalServerError)
}
