package token

import (
	"context"
	"sync"
	"time"
)

type storage struct {
	sync.RWMutex
	tokens    map[string]struct{}
	oldTokens map[string]struct{}
}

// NewStorage is a single use tokens storage. After successful verify token delete from storage.
// Token guaranteed storage at least tokenLifeTime duration.
func NewStorage(ctx context.Context, tokenLifeTime time.Duration) *storage {
	s := &storage{
		tokens:    make(map[string]struct{}),
		oldTokens: make(map[string]struct{}),
	}

	s.cleanupWorker(ctx, tokenLifeTime)
	return s
}

func (s *storage) Put(v string) {
	s.Lock()
	s.tokens[v] = struct{}{}
	s.Unlock()
}

func (s *storage) Verify(v string) bool {
	s.RLock()
	if isContain(s.tokens, v) || isContain(s.oldTokens, v) {
		s.RUnlock()
		s.Lock()
		delete(s.tokens, v)
		delete(s.oldTokens, v)
		s.Unlock()
		return true
	}

	s.RUnlock()
	return false
}

func isContain(m map[string]struct{}, key string) bool {
	_, ok := m[key]
	return ok
}

func (s *storage) cleanupWorker(ctx context.Context, tokenLifeTime time.Duration) {
	go func() {
		ticker := time.NewTicker(tokenLifeTime)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.Lock()
				s.oldTokens = s.tokens
				s.tokens = make(map[string]struct{})
				s.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()
}
