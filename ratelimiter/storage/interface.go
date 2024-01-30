package storage

import (
	"context"
	"time"
)

type RateLimiterStorageInterface interface {
	AddBlock(ctx context.Context, key string, blockInMilliseconds int64) (*time.Time, error)
	GetBlock(ctx context.Context, key string) (*time.Time, error)
	Increment(ctx context.Context, key string, maxAccess int64) (bool, error)
}
