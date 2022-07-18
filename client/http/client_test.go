package http_test

import (
	"context"
	gohttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"word_of_wisdom/client/http"
)

func TestClient_GetQuote(t *testing.T) {
	var (
		ctrl           = gomock.NewController(t)
		ctx, cancelCtx = context.WithCancel(context.Background())
	)
	defer cancelCtx()

	t.Run("ok", func(t *testing.T) {
		respData := "response_data"

		server := httptest.NewServer(gohttp.HandlerFunc(func(w gohttp.ResponseWriter, req *gohttp.Request) {
			require.Equal(t, "/quote", req.URL.Path)

			w.WriteHeader(gohttp.StatusOK)
			_, _ = w.Write([]byte(respData))
		}))
		defer server.Close()

		data, err := http.NewClient(server.URL, server.Client(), NewMockPoW(ctrl)).GetQuote(ctx)
		require.NoError(t, err)
		require.Equal(t, respData, data)
	})

	t.Run("error, 500", func(t *testing.T) {
		respData := gohttp.StatusText(gohttp.StatusInternalServerError)

		server := httptest.NewServer(gohttp.HandlerFunc(func(w gohttp.ResponseWriter, req *gohttp.Request) {
			require.Equal(t, "/quote", req.URL.Path)

			w.WriteHeader(gohttp.StatusInternalServerError)
			_, _ = w.Write([]byte(respData))
		}))
		defer server.Close()

		_, err := http.NewClient(server.URL, server.Client(), NewMockPoW(ctrl)).GetQuote(ctx)
		require.Error(t, err)
	})
}
