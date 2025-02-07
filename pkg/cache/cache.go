package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (Scannable, error)
	GetJSON(ctx context.Context, key, path string) (Scannable, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	SetJSON(ctx context.Context, key, path string, value any, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	DeleteJSON(ctx context.Context, key, path string) error
}

type Scannable interface {
	Scan(dest any) error
}
