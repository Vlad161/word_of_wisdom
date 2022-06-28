package http

import (
	"context"
)

type PoW interface {
	Calculate(ctx context.Context, payload []byte, timestamp int64, targetBits uint) (int, []byte, bool)
}
