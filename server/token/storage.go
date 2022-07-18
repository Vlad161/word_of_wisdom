package token

import (
	"context"
)

type (
	onetimeStorage struct {
		storage Storage
	}
)

// NewOnetimeStorage is a single use tokens storage. Once a verified token has been used, it is removed from storage.
// Token guaranteed storage at least tokenLifeTime duration.
func NewOnetimeStorage(storage Storage) *onetimeStorage {
	return &onetimeStorage{
		storage: storage,
	}
}

func (s *onetimeStorage) Get(ctx context.Context, k string) error {
	_, err := s.storage.Get(ctx, k)
	return err
}

func (s *onetimeStorage) Put(ctx context.Context, k string) error {
	return s.storage.Put(ctx, k, struct{}{})
}

func (s *onetimeStorage) Use(ctx context.Context, k string) error {
	err := s.Get(ctx, k)
	if err != nil {
		return err
	}
	return s.storage.Delete(ctx, k)
}
