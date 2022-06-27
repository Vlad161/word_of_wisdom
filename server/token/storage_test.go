package token_test

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"word_of_wisdom/server/token"
)

func TestStorage(t *testing.T) {
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	t.Run("multi goroutine read/write", func(t *testing.T) {
		var (
			gCount  = 10
			maxJ    = 100
			r       = rand.New(rand.NewSource(time.Now().Unix()))
			storage = token.NewStorage(ctx, 1*time.Second)
			wg      = sync.WaitGroup{}
		)

		wg.Add(gCount)
		for i := 0; i < gCount; i++ {
			go func() {
				for j := 0; j < r.Intn(maxJ); j++ {
					ts := token.New()
					storage.Put(ts)
					assert.True(t, storage.Verify(ts))
					assert.False(t, storage.Verify(ts))
				}
				wg.Done()
			}()
		}
		wg.Wait()
	})

	t.Run("token lifetime", func(t *testing.T) {
		var (
			tokenLifeTime = 50 * time.Millisecond
			storage       = token.NewStorage(ctx, tokenLifeTime)
		)

		ts := token.New()
		storage.Put(ts)
		time.Sleep(tokenLifeTime*2 + 1)

		require.False(t, storage.Verify(ts))
	})
}
