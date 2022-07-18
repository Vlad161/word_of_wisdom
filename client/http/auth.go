package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (c *client) auth(ctx context.Context) error {
	challengeData, err := c.getChallenge(ctx)
	if err != nil {
		return err
	}

	nonce, _, ok := c.pow.Calculate(ctx, []byte(challengeData.Token), challengeData.Timestamp, challengeData.TargetBits)
	if !ok {
		return fmt.Errorf("can't calculate hashcash")
	}

	postChallengeData, err := c.postChallenge(ctx, postChallengeReqBody{Nonce: nonce}, challengeData.JWT)
	if err != nil {
		return err
	}

	c.authHeaderValue = postChallengeData.JWT
	return nil
}

func (c *client) getChallenge(ctx context.Context) (getChallengeRespBody, error) {
	var respBody getChallengeRespBody

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.host+"/challenge", nil)
	if err != nil {
		return respBody, err
	}

	resp, err := c.transport.Do(req)
	if err != nil {
		return respBody, err
	}
	if resp.StatusCode != http.StatusOK {
		return respBody, errors.New("can't get challenge")
	}

	return respBody, json.NewDecoder(resp.Body).Decode(&respBody)
}

func (c *client) postChallenge(ctx context.Context, body postChallengeReqBody, jwt string) (postChallengeRespBody, error) {
	var respBody postChallengeRespBody

	bb := &bytes.Buffer{}
	if err := json.NewEncoder(bb).Encode(&body); err != nil {
		return respBody, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.host+"/challenge", bb)
	if err != nil {
		return respBody, err
	}
	addBearerHeader(req.Header, jwt)

	resp, err := c.transport.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return respBody, fmt.Errorf("can't post challenge: %w", err)
	}

	return respBody, json.NewDecoder(resp.Body).Decode(&respBody)
}

func addBearerHeader(h http.Header, token string) {
	h.Add("Authorization", "Bearer "+token)
}
