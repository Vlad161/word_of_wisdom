package storage

import (
	"context"
	"sync"
	"time"
)

type localTemporary struct {
	mux       sync.RWMutex
	tokens    map[string]interface{}
	oldTokens map[string]interface{}
}

func NewLocalTemporary(ctx context.Context, ttl time.Duration) *localTemporary {
	s := &localTemporary{
		tokens:    make(map[string]interface{}),
		oldTokens: make(map[string]interface{}),
	}

	s.cleanupWorker(ctx, ttl)
	return s
}

func (s *localTemporary) Get(k string) (interface{}, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if v, ok := s.tokens[k]; ok {
		return v, nil
	} else if v, ok = s.oldTokens[k]; ok {
		return v, nil
	}
	return nil, ErrNotFound
}

func (s *localTemporary) Put(k string, v interface{}) error {
	s.mux.Lock()
	s.tokens[k] = v
	s.mux.Unlock()
	return nil
}

func (s *localTemporary) Delete(k string) error {
	if _, err := s.Get(k); err != nil {
		return err
	}

	s.mux.Lock()
	delete(s.tokens, k)
	delete(s.oldTokens, k)
	s.mux.Unlock()
	return nil
}

// Primitive realisation to store value at least ttl duration. Not guarantee to store value not more ttl duration.
func (s *localTemporary) cleanupWorker(ctx context.Context, ttl time.Duration) {
	go func() {
		ticker := time.NewTicker(ttl)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.mux.Lock()
				s.oldTokens = s.tokens
				s.tokens = make(map[string]interface{})
				s.mux.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()
}
