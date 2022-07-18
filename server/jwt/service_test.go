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
		token, err := jwtService.CreateToken(
			map[string]interface{}{"id": "123"},
			time.Now().Add(10*time.Second),
			jwt.AlgHS256,
		)

		require.NoError(t, err)
		require.NoError(t, jwtService.Verify(token))
	})

	t.Run("ok, empty data", func(t *testing.T) {
		jwtService := jwt.New(keyPart1)
		token, err := jwtService.CreateToken(
			nil,
			time.Now().Add(10*time.Second),
			jwt.AlgHS256,
		)

		require.NoError(t, err)
		require.NoError(t, jwtService.Verify(token))
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
		require.Error(t, jwtService.Verify(token+"abc"))
	})

	t.Run("error, expired token", func(t *testing.T) {
		jwtService := jwt.New(keyPart1)
		token, err := jwtService.CreateToken(
			nil,
			time.Now().Add(-10*time.Second),
			jwt.AlgHS256,
		)

		require.NoError(t, err)
		require.Error(t, jwtService.Verify(token))
	})

	t.Run("error, verify invalid token", func(t *testing.T) {
		jwtService := jwt.New(keyPart1)
		token, err := jwtService.CreateToken(
			nil,
			time.Now().Add(10*time.Second),
			jwt.AlgHS256,
		)

		require.NoError(t, err)
		require.Error(t, jwtService.Verify(strings.Join(strings.Split(token, ".")[:2], ".")))
	})
}
