package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"word_of_wisdom/logger"
	"word_of_wisdom/server/handler"
	"word_of_wisdom/server/jwt"
	"word_of_wisdom/test"
)

func TestChallengeHandler_GET(t *testing.T) {
	var (
		ctrl          = gomock.NewController(t)
		log           = logger.New()
		targetBits    = uint(14)
		tokenLifetime = 1 * time.Second
	)

	t.Run("200, ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/challenge", nil)
		w := httptest.NewRecorder()

		jwtService := NewMockJWTService(ctrl)
		jwtService.EXPECT().CreateToken(gomock.Any(), gomock.Any(), jwt.AlgHS256).Return("jwt_123", nil).Times(1)

		handler.NewChallengeHandler(log, jwtService, tokenLifetime, targetBits, nil, nil).Handler().ServeHTTP(w, req)

		var body map[string]interface{}
		require.NoError(t, json.NewDecoder(w.Body).Decode(&body))

		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, body["timestamp"].(float64) > 0)
		assert.True(t, len(body["token"].(string)) > 0)
		assert.True(t, body["target_bits"].(float64) == float64(targetBits))
		assert.True(t, len(body["jwt"].(string)) > 0)
	})

	t.Run("500, can't create token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/challenge", nil)
		w := httptest.NewRecorder()

		jwtService := NewMockJWTService(ctrl)
		jwtService.EXPECT().CreateToken(gomock.Any(), gomock.Any(), jwt.AlgHS256).Return("", errors.New("error")).Times(1)

		handler.NewChallengeHandler(log, jwtService, tokenLifetime, targetBits, nil, nil).Handler().ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestChallengeHandler_Post(t *testing.T) {
	var (
		ctrl          = gomock.NewController(t)
		log           = logger.New()
		timestamp     = int64(1234)
		targetBits    = uint(14)
		tokenLifetime = 1 * time.Second
	)

	tests := []struct {
		name       string
		reqBody    []byte
		bearer     string
		jwtVerify  test.MockContract
		jwtCreate  test.MockContract
		storagePut test.MockContract
		powVerify  test.MockContract

		expectedCode int
		expectedBody string
	}{
		{
			name:       "200, ok",
			reqBody:    []byte(`{"nonce": 3456}`),
			bearer:     "Bearer token_123",
			jwtVerify:  test.MockContract{Param1: "token_123", Value1: map[string]interface{}{"timestamp": float64(timestamp), "token": "app_token", "target_bits": float64(targetBits)}, Calls: 1},
			storagePut: test.MockContract{Param1: "app_token", Calls: 1},
			powVerify:  test.MockContract{Param1: []byte("app_token"), Param2: timestamp, Param3: targetBits, Param4: 3456, Value1: true, Calls: 1},
			jwtCreate:  test.MockContract{Param1: map[string]interface{}{"token": "app_token", "is_verify": true}, Param2: gomock.Any(), Param3: jwt.AlgHS256, Value1: "jwt_123", Calls: 1},

			expectedCode: http.StatusOK,
			expectedBody: `{"jwt":"jwt_123"}`,
		},
		{
			name:      "500, empty jwt payload",
			reqBody:   []byte(`{"nonce": 3456}`),
			bearer:    "Bearer token_123",
			jwtVerify: test.MockContract{Param1: "token_123", Value1: map[string]interface{}{}, Calls: 1},
			powVerify: test.MockContract{Param1: gomock.Any(), Param2: gomock.Any(), Param3: gomock.Any(), Param4: gomock.Any(), Value1: false, Calls: 1},
			jwtCreate: test.MockContract{Value1: ""},

			expectedCode: http.StatusInternalServerError,
		},
		{
			name:      "500, can't jwt verify",
			reqBody:   []byte(`{"nonce": 3456}`),
			bearer:    "Bearer token_123",
			jwtVerify: test.MockContract{Param1: "token_123", Value2: errors.New("error"), Calls: 1},
			powVerify: test.MockContract{Value1: false},
			jwtCreate: test.MockContract{Value1: ""},

			expectedCode: http.StatusInternalServerError,
		},
		{
			name:       "500, can't put to storage",
			reqBody:    []byte(`{"nonce": 3456}`),
			bearer:     "Bearer token_123",
			jwtVerify:  test.MockContract{Param1: "token_123", Value1: map[string]interface{}{"timestamp": float64(timestamp), "token": "app_token", "target_bits": float64(targetBits)}, Calls: 1},
			storagePut: test.MockContract{Param1: "app_token", Value1: errors.New("error"), Calls: 1},
			powVerify:  test.MockContract{Param1: []byte("app_token"), Param2: timestamp, Param3: targetBits, Param4: 3456, Value1: true, Calls: 1},
			jwtCreate:  test.MockContract{Value1: ""},

			expectedCode: http.StatusInternalServerError,
		},
		{
			name:       "500, create token error",
			reqBody:    []byte(`{"nonce": 3456}`),
			bearer:     "Bearer token_123",
			jwtVerify:  test.MockContract{Param1: "token_123", Value1: map[string]interface{}{"timestamp": float64(timestamp), "token": "app_token", "target_bits": float64(targetBits)}, Calls: 1},
			storagePut: test.MockContract{Param1: "app_token", Calls: 1},
			powVerify:  test.MockContract{Param1: []byte("app_token"), Param2: timestamp, Param3: targetBits, Param4: 3456, Value1: true, Calls: 1},
			jwtCreate:  test.MockContract{Param1: map[string]interface{}{"token": "app_token", "is_verify": true}, Param2: gomock.Any(), Param3: jwt.AlgHS256, Value1: "", Value2: errors.New("error"), Calls: 1},

			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/challenge", bytes.NewBuffer(tc.reqBody))
			req.Header.Add("Authorization", tc.bearer)

			w := httptest.NewRecorder()

			jwtService := NewMockJWTService(ctrl)
			jwtService.EXPECT().Verify(tc.jwtVerify.Param1).Return(tc.jwtVerify.Value1, tc.jwtVerify.Value2).Times(tc.jwtVerify.Calls)
			jwtService.EXPECT().CreateToken(tc.jwtCreate.Param1, tc.jwtCreate.Param2, tc.jwtCreate.Param3).Return(tc.jwtCreate.Value1, tc.jwtCreate.Value2).Times(tc.jwtCreate.Calls)

			tokenStorage := NewMockTokenStorage(ctrl)
			tokenStorage.EXPECT().Put(gomock.Any(), tc.storagePut.Param1).Return(tc.storagePut.Value1).Times(tc.storagePut.Calls)

			powAlg := NewMockPoW(ctrl)
			powAlg.EXPECT().Verify(tc.powVerify.Param1, tc.powVerify.Param2, tc.powVerify.Param3, tc.powVerify.Param4).
				Return(tc.powVerify.Value1).Times(tc.powVerify.Calls)

			handler.NewChallengeHandler(log, jwtService, tokenLifetime, targetBits, tokenStorage, powAlg).Handler().ServeHTTP(w, req)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}
