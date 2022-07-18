package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"word_of_wisdom/logger"
	"word_of_wisdom/server/jwt"
	"word_of_wisdom/server/token"
)

const (
	keyTimestamp  = "timestamp"
	keyToken      = "token"
	keyTargetBits = "target_bits"
	keyJWT        = "jwt"
	keyIsVerify   = "is_verify"
)

type (
	postChallengeReqBody struct {
		Nonce int `json:"nonce"`
	}

	challengeHandler struct {
		log logger.Logger
		jwt JWTService

		tokenLifeTime time.Duration
		targetBits    uint
		tStorage      TokenStorage
		pow           PoW
	}
)

func NewChallengeHandler(log logger.Logger, jwt JWTService, tokenLifeTime time.Duration, targetBits uint, tokenStorage TokenStorage, pow PoW) *challengeHandler {
	return &challengeHandler{
		log:           log,
		jwt:           jwt,
		tokenLifeTime: tokenLifeTime,
		targetBits:    targetBits,
		tStorage:      tokenStorage,
		pow:           pow,
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
	ts := time.Now().Unix()

	payload := map[string]interface{}{
		keyTimestamp:  ts,
		keyToken:      tc,
		keyTargetBits: h.targetBits,
	}
	jwtToken, err := h.jwt.CreateToken(payload, time.Now().Add(h.tokenLifeTime), jwt.AlgHS256)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	payload[keyJWT] = jwtToken
	bb, err := json.Marshal(&payload)
	if err != nil {
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

	reqJwtPayload, err := h.jwt.Verify(extractHeaderBearer(req.Header))
	if err != nil {
		h.log.Error("can't verify jwt token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var (
		tc         = reqJwtPayload[keyToken].(string)
		ts         = int64(reqJwtPayload[keyTimestamp].(float64))
		targetBits = uint(reqJwtPayload[keyTargetBits].(float64))
	)
	if !h.pow.Verify([]byte(tc), ts, targetBits, reqBody.Nonce) {
		h.log.Error("can't verify pow", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = h.tStorage.Put(req.Context(), tc); err != nil {
		h.log.Error("can't put to token storage", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jwtToken, err := h.jwt.CreateToken(map[string]interface{}{
		keyToken:    tc,
		keyIsVerify: true,
	}, time.Now().Add(h.tokenLifeTime), jwt.AlgHS256)
	if err != nil {
		h.log.Error("can't create jwt token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bb, err := json.Marshal(&map[string]interface{}{
		keyJWT: jwtToken,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bb)
	return

}
