package token

import (
	"context"
)

type Storage interface {
	Get(ctx context.Context, k string) (interface{}, error)
	Put(ctx context.Context, k string, v interface{}) error
	Delete(ctx context.Context, k string) error
}

type StorageBytes interface {
	Get(ctx context.Context, k string) ([]byte, error)
	Put(ctx context.Context, k string, v []byte) error
	Delete(ctx context.Context, k string) error
}
