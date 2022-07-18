package token

import (
	"context"
)

type storageBytesAdapter struct {
	storage StorageBytes
}

func NewStorageBytesAdapter(storage StorageBytes) *storageBytesAdapter {
	return &storageBytesAdapter{storage: storage}
}

func (a *storageBytesAdapter) Get(ctx context.Context, k string) (interface{}, error) {
	_, err := a.storage.Get(ctx, k)
	return nil, err
}

func (a *storageBytesAdapter) Put(ctx context.Context, k string, v interface{}) error {
	return a.storage.Put(ctx, k, []byte{})
}

func (a *storageBytesAdapter) Delete(ctx context.Context, k string) error {
	return a.storage.Delete(ctx, k)
}
