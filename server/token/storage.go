package token

import (
	"context"
	"errors"
)

var (
	ErrCastValue        = errors.New("can't cast to value struct")
	ErrTokenNotVerified = errors.New("token is not verified")
)

type (
	value struct {
		TargetBits uint
		IsVerified bool
	}

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

func (s *onetimeStorage) Get(ctx context.Context, k string) (uint, error) {
	v, err := s.get(ctx, k)
	if err != nil {
		return 0, err
	}
	return v.TargetBits, nil
}

func (s *onetimeStorage) Put(ctx context.Context, k string, targetBits uint) error {
	return s.storage.Put(ctx, k, value{TargetBits: targetBits})
}

func (s *onetimeStorage) Use(ctx context.Context, k string) error {
	v, err := s.get(ctx, k)
	if err != nil {
		return err
	}

	if !v.IsVerified {
		return ErrTokenNotVerified
	}
	return s.storage.Delete(ctx, k)
}

func (s *onetimeStorage) Verify(ctx context.Context, k string) error {
	v, err := s.get(ctx, k)
	if err != nil {
		return err
	}

	v.IsVerified = true
	return s.storage.Put(ctx, k, v)
}

func (s *onetimeStorage) get(ctx context.Context, k string) (empty value, _ error) {
	v, err := s.storage.Get(ctx, k)
	if err != nil {
		return empty, err
	}
	tv, tOk := v.(value)
	if !tOk {
		return empty, ErrCastValue
	}
	return tv, nil
}
