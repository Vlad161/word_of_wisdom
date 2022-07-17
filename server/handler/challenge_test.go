package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"word_of_wisdom/server/handler"
	"word_of_wisdom/test"
)

func TestChallengeHandler_GET(t *testing.T) {
	var (
		ctrl       = gomock.NewController(t)
		targetBits = uint(14)
	)

	t.Run("ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/challenge", nil)
		w := httptest.NewRecorder()

		tokenStorage := NewMockTokenStorage(ctrl)
		tokenStorage.EXPECT().Put(gomock.Any(), targetBits).Times(1)

		handler.NewChallengeHandler(targetBits, tokenStorage, nil).Handler().ServeHTTP(w, req)

		var body map[string]interface{}
		require.NoError(t, json.NewDecoder(w.Body).Decode(&body))

		assert.True(t, body["timestamp"].(float64) > 0)
		assert.True(t, len(body["token"].(string)) > 0)
		assert.True(t, body["target_bits"].(float64) == float64(targetBits))
	})
}

func TestChallengeHandler_Post(t *testing.T) {
	var (
		ctrl       = gomock.NewController(t)
		targetBits = uint(14)
	)

	tests := []struct {
		name          string
		reqBody       []byte
		storageGet    test.MockContract
		storageVerify test.MockContract
		powVerify     test.MockContract

		expectedCode int
	}{
		{
			name:          "200, ok",
			reqBody:       []byte(`{"timestamp": 1234, "token": "token_123", "target_bits": 14, "nonce": 3456}`),
			storageGet:    test.MockContract{Param1: "token_123", Value1: uint(14), Calls: 1},
			storageVerify: test.MockContract{Param1: "token_123", Value1: nil, Calls: 1},
			powVerify:     test.MockContract{Param1: []byte("token_123"), Param2: int64(1234), Param3: targetBits, Param4: 3456, Value1: true, Calls: 1},

			expectedCode: http.StatusOK,
		},
		{
			name:          "500, storage target bits error",
			reqBody:       []byte(`{"timestamp": 1234, "token": "token_123", "target_bits": 14, "nonce": 3456}`),
			storageGet:    test.MockContract{Param1: "token_123", Value1: uint(14), Value2: errors.New("error"), Calls: 1},
			storageVerify: test.MockContract{Value1: nil},
			powVerify:     test.MockContract{Value1: true},

			expectedCode: http.StatusInternalServerError,
		},
		{
			name:          "500, pow verify error",
			reqBody:       []byte(`{"timestamp": 1234, "token": "token_123", "target_bits": 14, "nonce": 3456}`),
			storageGet:    test.MockContract{Param1: "token_123", Value1: uint(14), Calls: 1},
			storageVerify: test.MockContract{Value1: nil},
			powVerify:     test.MockContract{Param1: []byte("token_123"), Param2: int64(1234), Param3: targetBits, Param4: 3456, Value1: false, Calls: 1},

			expectedCode: http.StatusInternalServerError,
		},
		{
			name:          "500, storage verify error",
			reqBody:       []byte(`{"timestamp": 1234, "token": "token_123", "target_bits": 14, "nonce": 3456}`),
			storageGet:    test.MockContract{Param1: "token_123", Value1: uint(14), Calls: 1},
			storageVerify: test.MockContract{Param1: "token_123", Value1: errors.New("error"), Calls: 1},
			powVerify:     test.MockContract{Param1: []byte("token_123"), Param2: int64(1234), Param3: targetBits, Param4: 3456, Value1: true, Calls: 1},

			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/challenge", bytes.NewBuffer(tc.reqBody))
			w := httptest.NewRecorder()

			tokenStorage := NewMockTokenStorage(ctrl)
			tokenStorage.EXPECT().
				Get(tc.storageGet.Param1).
				Return(tc.storageGet.Value1, tc.storageGet.Value2).Times(tc.storageGet.Calls)
			tokenStorage.EXPECT().
				Verify(tc.storageVerify.Param1).Return(tc.storageVerify.Value1).Times(tc.storageVerify.Calls)

			powAlg := NewMockPoW(ctrl)
			powAlg.EXPECT().Verify(tc.powVerify.Param1, tc.powVerify.Param2, tc.powVerify.Param3, tc.powVerify.Param4).
				Return(tc.powVerify.Value1).Times(tc.powVerify.Calls)

			handler.NewChallengeHandler(targetBits, tokenStorage, powAlg).Handler().ServeHTTP(w, req)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}
