//go:generate mockgen -source=contract.go -package=$GOPACKAGE -destination=mock_contract_test.go

package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"word_of_wisdom/server/handler"
	"word_of_wisdom/test"
)

func TestAuthMW(t *testing.T) {
	var (
		ctrl = gomock.NewController(t)
	)

	tests := []struct {
		name               string
		authHeader         string
		mockHandler        test.MockContract
		tokenStorageVerify test.MockContract

		expectedCode int
	}{
		{
			name:               "ok, 200",
			authHeader:         "Bearer token_123",
			mockHandler:        test.MockContract{Calls: 1},
			tokenStorageVerify: test.MockContract{Param1: "token_123", Value1: true, Calls: 1},

			expectedCode: http.StatusOK,
		},
		{
			name:               "empty auth header, 401",
			tokenStorageVerify: test.MockContract{Value1: false},

			expectedCode: http.StatusUnauthorized,
		},
		{
			name:               "can't verify token, 401",
			authHeader:         "Bearer token_123",
			tokenStorageVerify: test.MockContract{Param1: "token_123", Value1: false, Calls: 1},

			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Add("Authorization", tc.authHeader)

			w := httptest.NewRecorder()

			mockHandler := NewMockHandler(ctrl)
			mockHandler.EXPECT().
				ServeHTTP(gomock.Any(), gomock.Any()).Times(tc.mockHandler.Calls)

			tokenStorage := NewMockTokenStorage(ctrl)
			tokenStorage.EXPECT().
				Use(tc.tokenStorageVerify.Param1).Return(tc.tokenStorageVerify.Value1).Times(tc.tokenStorageVerify.Calls)

			handler.AuthMW(mockHandler, tokenStorage).ServeHTTP(w, req)

			require.Equal(t, tc.expectedCode, w.Code)
		})
	}
}
