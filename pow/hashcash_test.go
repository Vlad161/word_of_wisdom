package pow_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"word_of_wisdom/pow"
)

func TestHashCash(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		var (
			ctx, cancelCtx = context.WithTimeout(context.Background(), 1*time.Second)
			timestamp      = time.Now().Unix()
			data           = uuid.New().String()
			targetBits     = uint(14)
			hc             = pow.NewHashCash(timestamp, data, targetBits)
		)
		defer cancelCtx()

		nonce, hash, ok := hc.Calculate(ctx)
		fmt.Printf("nonce: %d, hash: %x \n", nonce, hash)

		require.True(t, ok)
		require.True(t, hc.Verify(nonce))
	})

	t.Run("timeout exceeded", func(t *testing.T) {
		var (
			ctx, cancelCtx = context.WithTimeout(context.Background(), 1*time.Second)
			timestamp      = time.Now().Unix()
			data           = uuid.New().String()
			targetBits     = uint(48)
			hc             = pow.NewHashCash(timestamp, data, targetBits)
		)
		defer cancelCtx()

		nonce, hash, ok := hc.Calculate(ctx)
		fmt.Printf("nonce: %d, hash: %x \n", nonce, hash)

		require.False(t, ok)
		require.False(t, hc.Verify(nonce))
	})
}
