package redis

import (
	"context"
	"github.com/adanyl0v/pocket-ideas/pkg/cache"
	"github.com/redis/go-redis/v9"
)

type (
	KeyCursorScannerFn func(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd

	KeyCursorScanner struct {
		fn     KeyCursorScannerFn
		key    string
		cursor uint64
		match  string
		count  int64
	}
)

func NewKeyCursorScanner(fn KeyCursorScannerFn) *KeyCursorScanner {
	return &KeyCursorScanner{fn: fn}
}

func (s *KeyCursorScanner) Scan(ctx context.Context) cache.ScanIterator {
	return s.fn(ctx, s.key, s.cursor, s.match, s.count).Iterator()
}

func (s *KeyCursorScanner) WithArgs(key string, cursor uint64, match string, count int64) cache.KeyCursorScanner {
	return &KeyCursorScanner{
		key:    key,
		cursor: cursor,
		match:  match,
		count:  count,
	}
}
