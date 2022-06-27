package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"word_of_wisdom/server/handler"
)

func TestQuoteHandlerFunc(t *testing.T) {
	tests := []struct {
		name   string
		method string

		expectedCode int
		expectedBody bool
	}{
		{
			name:   "ok, 200",
			method: http.MethodGet,

			expectedCode: http.StatusOK,
			expectedBody: true,
		},
		{
			name:   "error, 405",
			method: http.MethodPost,

			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/quote", nil)
			w := httptest.NewRecorder()

			handler.QuoteHandlerFunc().ServeHTTP(w, req)

			require.Equal(t, tc.expectedCode, w.Code)
			if tc.expectedBody {
				assert.True(t, w.Body.Len() > 0)
			}
		})
	}
}
