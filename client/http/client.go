package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	getChallengeRespBody struct {
		Timestamp  int64  `json:"timestamp"`
		Token      string `json:"token"`
		TargetBits uint   `json:"target_bits"`
		JWT        string `json:"jwt"`
	}

	postChallengeRespBody struct {
		JWT string `json:"jwt"`
	}

	postChallengeReqBody struct {
		Nonce int `json:"nonce"`
	}

	client struct {
		host      string
		transport *http.Client
		pow       PoW

		authHeaderValue string
	}
)

func NewClient(host string, transport *http.Client, pow PoW) *client {
	return &client{
		host:      host,
		transport: transport,
		pow:       pow,
	}
}

func (c *client) GetQuote(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.host+"/quote", nil)
	if err != nil {
		return "", err
	}

	resp, err := c.doAuthRequest(ctx, req, false)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("can't get quote, status code: %d", resp.StatusCode)
	}

	bb, err := ioutil.ReadAll(resp.Body)
	return string(bb), err
}

func (c *client) doAuthRequest(ctx context.Context, req *http.Request, retry bool) (*http.Response, error) {
	resp, err := c.transport.Do(req)
	if resp != nil && resp.StatusCode == http.StatusUnauthorized && !retry {
		if err = c.auth(ctx); err == nil {
			addBearerHeader(req.Header, c.authHeaderValue)
			return c.doAuthRequest(ctx, req, true)
		}
	}
	return resp, err
}
