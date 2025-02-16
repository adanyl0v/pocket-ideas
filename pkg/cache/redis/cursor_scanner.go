package redis

import (
	"context"
	"github.com/adanyl0v/pocket-ideas/pkg/cache"
	"github.com/redis/go-redis/v9"
)

type (
	CursorScannerFn func(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd

	CursorScanner struct {
		fn     CursorScannerFn
		cursor uint64
		match  string
		count  int64
	}
)

func NewCursorScanner(fn CursorScannerFn) *CursorScanner {
	return &CursorScanner{fn: fn}
}

func (s *CursorScanner) Scan(ctx context.Context) cache.ScanIterator {
	return s.fn(ctx, s.cursor, s.match, s.count).Iterator()
}

func (s *CursorScanner) WithArgs(cursor uint64, match string, count int64) cache.CursorScanner {
	return &CursorScanner{
		fn:     s.fn,
		cursor: cursor,
		match:  match,
		count:  count,
	}
}
