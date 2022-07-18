package token_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"word_of_wisdom/server/token"
)

func TestNewStorageBytesAdapter(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	t.Run("ok", func(t *testing.T) {
		as := token.NewStorageBytesAdapter(newSimpleBytesStorage())

		expected := token.NewTestPrivateValue(1, true)
		_ = as.Put(ctx, "key", expected)

		v, _ := as.Get(ctx, "key")
		require.Equal(t, expected, v)
	})
}

type simpleBytesStorage struct {
	m map[string][]byte
}

func newSimpleBytesStorage() *simpleBytesStorage {
	return &simpleBytesStorage{m: make(map[string][]byte)}
}

func (s *simpleBytesStorage) Get(_ context.Context, k string) ([]byte, error) {
	return s.m[k], nil
}

func (s *simpleBytesStorage) Put(_ context.Context, k string, v []byte) error {
	s.m[k] = v
	return nil
}

func (s *simpleBytesStorage) Delete(_ context.Context, k string) error {
	delete(s.m, k)
	return nil
}
