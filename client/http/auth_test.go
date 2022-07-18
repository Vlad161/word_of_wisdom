package http_test

import (
	"context"
	gohttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"word_of_wisdom/client/http"
	"word_of_wisdom/test"
)

func TestClient_Auth(t *testing.T) {
	var (
		ctrl           = gomock.NewController(t)
		ctx, cancelCtx = context.WithCancel(context.Background())
	)
	defer cancelCtx()

	type requestServerData struct {
		path     string
		method   string
		respCode int
		respData string
	}

	tests := []struct {
		name         string
		powCalculate test.MockContract
		requests     []requestServerData

		expected      string
		expectedError bool
	}{
		{
			name: "ok",
			requests: []requestServerData{
				{
					path:     "/quote",
					method:   gohttp.MethodGet,
					respCode: gohttp.StatusUnauthorized,
				},
				{
					path:     "/challenge",
					method:   gohttp.MethodGet,
					respCode: gohttp.StatusOK,
					respData: "{}",
				},
				{
					path:     "/challenge",
					method:   gohttp.MethodPost,
					respCode: gohttp.StatusOK,
				},
				{
					path:     "/quote",
					method:   gohttp.MethodGet,
					respCode: gohttp.StatusOK,
					respData: "response_data",
				},
			},
			powCalculate: test.MockContract{Value3: true, Calls: 1},
			expected:     "response_data",
		},
		{
			name: "error, get challenge",
			requests: []requestServerData{
				{
					path:     "/quote",
					method:   gohttp.MethodGet,
					respCode: gohttp.StatusUnauthorized,
				},
				{
					path:     "/challenge",
					method:   gohttp.MethodGet,
					respCode: gohttp.StatusInternalServerError,
				},
			},
			powCalculate:  test.MockContract{Value3: false, Calls: 0},
			expectedError: true,
		},
		{
			name: "errors, can't do PoW",
			requests: []requestServerData{
				{
					path:     "/quote",
					method:   gohttp.MethodGet,
					respCode: gohttp.StatusUnauthorized,
				},
				{
					path:     "/challenge",
					method:   gohttp.MethodGet,
					respCode: gohttp.StatusOK,
					respData: "{}",
				},
				{
					path:     "/challenge",
					method:   gohttp.MethodPost,
					respCode: gohttp.StatusOK,
				},
			},
			powCalculate:  test.MockContract{Value3: false, Calls: 1},
			expectedError: true,
		},
		{
			name: "errors, post challenge",
			requests: []requestServerData{
				{
					path:     "/quote",
					method:   gohttp.MethodGet,
					respCode: gohttp.StatusUnauthorized,
				},
				{
					path:     "/challenge",
					method:   gohttp.MethodGet,
					respCode: gohttp.StatusOK,
					respData: "{}",
				},
				{
					path:     "/challenge",
					method:   gohttp.MethodPost,
					respCode: gohttp.StatusInternalServerError,
				},
			},
			powCalculate:  test.MockContract{Value3: true, Calls: 1},
			expectedError: true,
		},
		{
			name: "error, can't do request after challenge",
			requests: []requestServerData{
				{
					path:     "/quote",
					method:   gohttp.MethodGet,
					respCode: gohttp.StatusUnauthorized,
				},
				{
					path:     "/challenge",
					method:   gohttp.MethodGet,
					respCode: gohttp.StatusOK,
					respData: "{}",
				},
				{
					path:     "/challenge",
					method:   gohttp.MethodPost,
					respCode: gohttp.StatusOK,
				},
				{
					path:     "/quote",
					method:   gohttp.MethodGet,
					respCode: gohttp.StatusUnauthorized,
				},
			},
			powCalculate:  test.MockContract{Value3: true, Calls: 1},
			expectedError: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(gohttp.HandlerFunc(func(w gohttp.ResponseWriter, req *gohttp.Request) {
				require.True(t, calls < len(tc.requests), "unexpected request")

				data := tc.requests[calls]
				calls++
				assert.Equal(t, data.method, req.Method)
				assert.Equal(t, data.path, req.URL.Path)

				w.WriteHeader(data.respCode)
				_, _ = w.Write([]byte(data.respData))
			}))
			defer server.Close()

			mockPow := NewMockPoW(ctrl)
			mockPow.EXPECT().
				Calculate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(1, []byte{}, tc.powCalculate.Value3).
				Times(tc.powCalculate.Calls)

			data, err := http.NewClient(server.URL, server.Client(), mockPow).GetQuote(ctx)
			if tc.expectedError {
				assert.Error(t, err)
				return
			}
			require.Equal(t, len(tc.requests), calls)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, data)
		})
	}
}
