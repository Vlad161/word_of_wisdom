package token

import (
	"errors"
)

var (
	ErrCastValue        = errors.New("can't cast to value struct")
	ErrTokenNotVerified = errors.New("token is not verified")
)

type (
	value struct {
		targetBits uint
		isVerified bool
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

func (s *onetimeStorage) Get(k string) (uint, error) {
	v, err := s.get(k)
	if err != nil {
		return 0, err
	}
	return v.targetBits, nil
}

func (s *onetimeStorage) Put(k string, targetBits uint) error {
	return s.storage.Put(k, value{targetBits: targetBits})
}

func (s *onetimeStorage) Use(k string) error {
	v, err := s.get(k)
	if err != nil {
		return err
	}

	if !v.isVerified {
		return ErrTokenNotVerified
	}
	return s.storage.Delete(k)
}

func (s *onetimeStorage) Verify(k string) error {
	v, err := s.get(k)
	if err != nil {
		return err
	}

	v.isVerified = true
	return s.storage.Put(k, v)
}

func (s *onetimeStorage) get(k string) (empty value, _ error) {
	v, err := s.storage.Get(k)
	if err != nil {
		return empty, err
	}
	tv, tOk := v.(value)
	if !tOk {
		return empty, ErrCastValue
	}
	return tv, nil
}
