package handler_test

import (
	"errors"
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
		name            string
		authHeader      string
		mockHandler     test.MockContract
		tokenStorageUse test.MockContract
		jwtService      test.MockContract

		expectedCode int
	}{
		{
			name:            "ok, 200",
			authHeader:      "Bearer token_123",
			mockHandler:     test.MockContract{Calls: 1},
			tokenStorageUse: test.MockContract{Param1: "app_token", Value1: nil, Calls: 1},
			jwtService:      test.MockContract{Param1: "token_123", Value1: map[string]interface{}{"is_verify": true, "token": "app_token"}, Calls: 1},

			expectedCode: http.StatusOK,
		},
		{
			name:       "error, can't verify token, 401",
			authHeader: "Bearer token_123",
			jwtService: test.MockContract{Param1: "token_123", Value2: errors.New("error"), Calls: 1},

			expectedCode: http.StatusUnauthorized,
		},
		{
			name:       "error, jwt token doesn't contain is_verify 401",
			authHeader: "Bearer token_123",
			jwtService: test.MockContract{Param1: "token_123", Value1: map[string]interface{}{"token": "app_token"}, Calls: 1},

			expectedCode: http.StatusUnauthorized,
		},
		{
			name:            "error, empty auth header, 401",
			tokenStorageUse: test.MockContract{Value1: errors.New("error")},

			expectedCode: http.StatusUnauthorized,
		},
		{
			name:            "error, can't use token, 401",
			authHeader:      "Bearer token_123",
			tokenStorageUse: test.MockContract{Param1: "app_token", Value1: errors.New("error"), Calls: 1},
			jwtService:      test.MockContract{Param1: "token_123", Value1: map[string]interface{}{"is_verify": true, "token": "app_token"}, Calls: 1},

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
				Use(gomock.Any(), tc.tokenStorageUse.Param1).Return(tc.tokenStorageUse.Value1).Times(tc.tokenStorageUse.Calls)

			jwtService := NewMockJWTService(ctrl)
			jwtService.EXPECT().Verify(tc.jwtService.Param1).Return(tc.jwtService.Value1, tc.jwtService.Value2).Times(tc.jwtService.Calls)

			handler.AuthMW(mockHandler, tokenStorage, jwtService).ServeHTTP(w, req)

			require.Equal(t, tc.expectedCode, w.Code)
		})
	}
}
