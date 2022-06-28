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

	t.Run("token success use after verification", func(t *testing.T) {
		var (
			tokenLifeTime = 1 * time.Second
			storage       = token.NewStorage(ctx, tokenLifeTime)
			v             = uint(1)
		)

		ts := token.New()
		storage.Put(ts, v)
		require.True(t, storage.Verify(ts))
		require.True(t, storage.Use(ts))
	})

	t.Run("token can't use without verification", func(t *testing.T) {
		var (
			tokenLifeTime = 1 * time.Second
			storage       = token.NewStorage(ctx, tokenLifeTime)
			v             = uint(1)
		)

		ts := token.New()
		storage.Put(ts, v)
		require.False(t, storage.Use(ts))
	})

	t.Run("multi goroutine read/write", func(t *testing.T) {
		var (
			gCount  = 10
			maxJ    = 100
			r       = rand.New(rand.NewSource(time.Now().Unix()))
			storage = token.NewStorage(ctx, 1*time.Second)
			wg      = sync.WaitGroup{}
			v       = uint(1)
		)

		wg.Add(gCount)
		for i := 0; i < gCount; i++ {
			go func() {
				for j := 0; j < r.Intn(maxJ); j++ {
					ts := token.New()
					storage.Put(ts, v)
					tb, _ := storage.TargetBits(ts)

					assert.Equal(t, tb, v)
					assert.True(t, storage.Verify(ts))
					assert.True(t, storage.Use(ts))
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
			v             = uint(1)
		)

		ts := token.New()
		storage.Put(ts, v)
		require.True(t, storage.Verify(ts))
		time.Sleep(tokenLifeTime * 3)
		require.False(t, storage.Use(ts))
	})
}
