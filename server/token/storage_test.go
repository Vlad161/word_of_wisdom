package token_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"word_of_wisdom/server/token"
)

func TestStorage(t *testing.T) {
	var (
		gCount  = 10
		maxJ    = 100
		r       = rand.New(rand.NewSource(time.Now().Unix()))
		storage = token.NewStorage()
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
}
