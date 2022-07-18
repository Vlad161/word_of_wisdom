package jwt_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"word_of_wisdom/server/jwt"
)

func TestService(t *testing.T) {
	const (
		keyPart1 = "123abc"
	)

	t.Run("ok", func(t *testing.T) {
		jwtService := jwt.New(keyPart1)
		payload := map[string]interface{}{"id": "123"}
		token, err := jwtService.CreateToken(
			payload,
			time.Now().Add(10*time.Second),
			jwt.AlgHS256,
		)

		require.NoError(t, err)

		tokenPayload, err := jwtService.Verify(token)
		require.NoError(t, err)
		require.NotEmpty(t, tokenPayload["id"])
		require.NotEmpty(t, tokenPayload["exp"])
	})

	t.Run("ok, empty data", func(t *testing.T) {
		jwtService := jwt.New(keyPart1)
		token, err := jwtService.CreateToken(
			nil,
			time.Now().Add(10*time.Second),
			jwt.AlgHS256,
		)

		require.NoError(t, err)

		_, err = jwtService.Verify(token)
		require.NoError(t, err)
	})

	t.Run("error, unknown alg", func(t *testing.T) {
		jwtService := jwt.New(keyPart1)
		_, err := jwtService.CreateToken(
			nil,
			time.Now().Add(10*time.Second),
			"unknown",
		)

		require.Error(t, err)
	})

	t.Run("error, verify", func(t *testing.T) {
		jwtService := jwt.New(keyPart1)
		token, err := jwtService.CreateToken(
			nil,
			time.Now().Add(10*time.Second),
			jwt.AlgHS256,
		)

		require.NoError(t, err)

		_, err = jwtService.Verify(token + "abc")
		require.Error(t, err)
	})

	t.Run("error, expired token", func(t *testing.T) {
		jwtService := jwt.New(keyPart1)
		token, err := jwtService.CreateToken(
			nil,
			time.Now().Add(-10*time.Second),
			jwt.AlgHS256,
		)

		require.NoError(t, err)

		_, err = jwtService.Verify(token)
		require.Error(t, err)
	})

	t.Run("error, verify invalid token", func(t *testing.T) {
		jwtService := jwt.New(keyPart1)
		token, err := jwtService.CreateToken(
			nil,
			time.Now().Add(10*time.Second),
			jwt.AlgHS256,
		)

		require.NoError(t, err)

		_, err = jwtService.Verify(strings.Join(strings.Split(token, ".")[:2], "."))
		require.Error(t, err)
	})
}
