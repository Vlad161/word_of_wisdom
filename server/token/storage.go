package token

import (
	"sync"
)

// Single use tokens storage. After successful verify token delete from storage.
type storage struct {
	sync.RWMutex
	tokens map[string]struct{} //TODO tokens are infinitive. Create job to clear old values
}

func NewStorage() *storage {
	return &storage{
		tokens: make(map[string]struct{}),
	}
}

func (s *storage) Put(v string) {
	s.Lock()
	s.tokens[v] = struct{}{}
	s.Unlock()
}

func (s *storage) Verify(v string) bool {
	s.RLock()
	if _, ok := s.tokens[v]; ok {
		s.RUnlock()
		s.Lock()
		delete(s.tokens, v)
		s.Unlock()
		return true
	}

	s.RUnlock()
	return false
}
