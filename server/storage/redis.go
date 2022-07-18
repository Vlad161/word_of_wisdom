package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
)

type redisTtl struct {
	cl RedisClient

	ttl time.Duration
}

func NewRedis(cl *redis.Client, ttl time.Duration) *redisTtl {
	return &redisTtl{cl: cl, ttl: ttl}
}

func (r *redisTtl) Get(ctx context.Context, k string) (empty []byte, _ error) {
	return r.cl.Get(ctx, k).Bytes()
}

func (r *redisTtl) Put(ctx context.Context, k string, v []byte) error {
	return r.cl.SetEx(ctx, k, v, r.ttl).Err()
}

func (r *redisTtl) Delete(ctx context.Context, k string) error {
	return r.cl.Del(ctx, k).Err()
}
