package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string, dest any) error
	GetJSON(ctx context.Context, key, path string, dest any) error
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	SetJSON(ctx context.Context, key, path string, value any, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	DeleteJSON(ctx context.Context, key, path string) error
}
