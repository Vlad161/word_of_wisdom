package token

import (
	"bytes"
	"context"
	"encoding/gob"
)

type storageBytesAdapter struct {
	storage StorageBytes
}

func init() {
	gob.Register(value{})
}

func NewStorageBytesAdapter(storage StorageBytes) *storageBytesAdapter {
	return &storageBytesAdapter{storage: storage}
}

func (a *storageBytesAdapter) Get(ctx context.Context, k string) (interface{}, error) {
	bb, err := a.storage.Get(ctx, k)
	if err != nil {
		return nil, err
	}

	var v value
	err = gob.NewDecoder(bytes.NewBuffer(bb)).Decode(&v)
	return v, err
}

func (a *storageBytesAdapter) Put(ctx context.Context, k string, v interface{}) error {
	bb := bytes.Buffer{}
	if err := gob.NewEncoder(&bb).Encode(v); err != nil {
		return err
	}

	return a.storage.Put(ctx, k, bb.Bytes())
}

func (a *storageBytesAdapter) Delete(ctx context.Context, k string) error {
	return a.storage.Delete(ctx, k)
}
