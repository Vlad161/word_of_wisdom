package http

import (
	"context"
	"io/ioutil"
	"net/http"
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

	bb, err := ioutil.ReadAll(resp.Body)
	return string(bb), err
}

func (c *client) doAuthRequest(ctx context.Context, req *http.Request, retry bool) (*http.Response, error) {
	resp, err := c.transport.Do(req)
	if resp != nil && resp.StatusCode == http.StatusUnauthorized && !retry {
		if err := c.auth(ctx); err == nil {
			req.Header.Add("Authorization", "Bearer "+c.authHeaderValue)
			return c.doAuthRequest(ctx, req, true)
		}
	}
	return resp, err
}
