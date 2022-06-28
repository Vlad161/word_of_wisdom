package token

import (
	"context"
	"sync"
	"time"
)

type (
	value struct {
		targetBits uint
		isVerified bool
	}

	storage struct {
		mux       sync.RWMutex
		tokens    map[string]*value
		oldTokens map[string]*value
	}
)

// NewStorage is a single use tokens storage. After successful verify token delete from storage.
// Token guaranteed storage at least tokenLifeTime duration.
func NewStorage(ctx context.Context, tokenLifeTime time.Duration) *storage {
	s := &storage{
		tokens:    make(map[string]*value),
		oldTokens: make(map[string]*value),
	}

	s.cleanupWorker(ctx, tokenLifeTime)
	return s
}

func (s *storage) Put(k string, v uint) {
	s.mux.Lock()
	s.tokens[k] = &value{targetBits: v}
	s.mux.Unlock()
}

func (s *storage) Use(k string) bool {
	if v, ok := s.get(k); ok && v.isVerified {
		s.mux.Lock()
		delete(s.tokens, k)
		delete(s.oldTokens, k)
		s.mux.Unlock()
		return true
	}
	return false
}

func (s *storage) Verify(k string) bool {
	if v, ok := s.get(k); ok {
		s.mux.Lock()
		v.isVerified = true
		s.mux.Unlock()
		return true
	}
	return false
}

func (s *storage) TargetBits(k string) (uint, bool) {
	if v, ok := s.get(k); ok {
		return v.targetBits, true
	}
	return 0, false
}

func (s *storage) get(k string) (*value, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if v, ok := s.tokens[k]; ok {
		return v, ok
	}
	v, ok := s.oldTokens[k]
	return v, ok
}

// Primitive realisation to store token at least tokenLifeTime duration. Not guarantee to store token not more tokenLifeTime duration.
func (s *storage) cleanupWorker(ctx context.Context, tokenLifeTime time.Duration) {
	go func() {
		ticker := time.NewTicker(tokenLifeTime)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.mux.Lock()
				s.oldTokens = s.tokens
				s.tokens = make(map[string]*value)
				s.mux.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()
}
