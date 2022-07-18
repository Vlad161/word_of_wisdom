package storage_test

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	storage2 "word_of_wisdom/server/storage"
	"word_of_wisdom/server/token"
)

func TestNewLocalTemporaryStorage(t *testing.T) {
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	t.Run("multi goroutine read/write", func(t *testing.T) {
		var (
			gCount = 10
			r      = rand.New(rand.NewSource(time.Now().Unix()))
			s      = storage2.NewLocalTemporary(ctx, 1*time.Second)
			wg     = sync.WaitGroup{}
			v      = uint(1)
		)

		wg.Add(gCount)
		for i := 0; i < gCount; i++ {
			maxJ := r.Intn(100)
			go func() {
				for j := 0; j < maxJ; j++ {
					ts := token.New()
					assert.NoError(t, s.Put(ctx, ts, v))

					_, err := s.Get(ctx, ts)
					assert.NoError(t, err)

					assert.NoError(t, s.Delete(ctx, ts))

					_, err = s.Get(ctx, ts)
					assert.Equal(t, storage2.ErrNotFound, err)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	})

	t.Run("delete err not found", func(t *testing.T) {
		s := storage2.NewLocalTemporary(ctx, 50*time.Millisecond)
		require.Equal(t, storage2.ErrNotFound, s.Delete(ctx, token.New()))
	})

	t.Run("value lifetime", func(t *testing.T) {
		var (
			tokenLifeTime = 50 * time.Millisecond
			s             = storage2.NewLocalTemporary(ctx, tokenLifeTime)
			v             = uint(1)
		)

		ts := token.New()
		require.NoError(t, s.Put(ctx, ts, v))

		_, err := s.Get(ctx, ts)
		require.NoError(t, err)

		time.Sleep(tokenLifeTime * 3)

		_, err = s.Get(ctx, ts)
		require.Equal(t, storage2.ErrNotFound, err)
	})

	t.Run("value from old tokens", func(t *testing.T) {
		var (
			tokenLifeTime = 50 * time.Millisecond
			s             = storage2.NewLocalTemporary(ctx, tokenLifeTime)
			v             = uint(1)
		)

		ts := token.New()
		require.NoError(t, s.Put(ctx, ts, v))

		_, err := s.Get(ctx, ts)
		require.NoError(t, err)

		time.Sleep(time.Duration(float64(tokenLifeTime) * 1.5))

		nv, err := s.Get(ctx, ts)
		require.NoError(t, err)
		require.Equal(t, v, nv)
	})
}
