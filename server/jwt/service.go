package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	hs256KeyPart2 = "stLwlC49P6WiTrl7epHE"

	AlgHS256 alg = "HS256"
)

var (
	base64NoPadding = base64.RawURLEncoding

	encodedHeaders = map[alg]string{
		AlgHS256: base64NoPadding.EncodeToString([]byte(`{"typ":"JWT","alg":"HS256"}`)),
	}
)

type (
	alg string

	service struct {
		hs256Key []byte
	}
)

func New(hs256KeyPart1 string) *service {
	return &service{
		hs256Key: []byte(hs256KeyPart1 + hs256KeyPart2),
	}
}

func (s *service) CreateToken(data map[string]interface{}, exp time.Time, alg alg) (string, error) {
	payload := make(map[string]interface{})
	if data != nil {
		payload = data
	}
	payload["exp"] = exp.Unix()

	bb, err := json.Marshal(&payload)
	if err != nil {
		return "", err
	}

	encodedHeader := encodedHeaders[alg]
	encodedPayload := base64NoPadding.EncodeToString(bb)

	switch alg {
	case AlgHS256:
		return hs256Signature(s.hs256Key, encodedHeader, encodedPayload), nil
	default:
		return "", errors.New("unknown alg")
	}
}

func (s *service) Verify(token string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("invalid token")
	}

	decodedPayload, err := base64NoPadding.DecodeString(parts[1])
	if err != nil {
		return err
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(decodedPayload, &payload); err != nil {
		return err
	}

	if exp, ok := payload["exp"].(float64); !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		return errors.New("token is expired")
	}

	decodedHeader, err := base64NoPadding.DecodeString(parts[0])
	if err != nil {
		return err
	}
	var header map[string]interface{}
	if err := json.Unmarshal(decodedHeader, &header); err != nil {
		return err
	}
	switch alg(header["alg"].(string)) {
	case AlgHS256:
		return hs256Verify(s.hs256Key, parts[0], parts[1], token)
	default:
		return errors.New("unknown alg")
	}
}

func hs256Signature(key []byte, header, payload string) string {
	unsigned := fmt.Sprintf("%s.%s", header, payload)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(unsigned))

	return fmt.Sprintf("%s.%s", unsigned, base64NoPadding.EncodeToString(h.Sum(nil)))
}

func hs256Verify(key []byte, header, payload, token string) error {
	nToken := hs256Signature(key, header, payload)
	if nToken != token {
		return errors.New("can't verify token")
	}
	return nil
}
