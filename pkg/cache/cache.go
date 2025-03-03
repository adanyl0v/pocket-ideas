package cache

import (
	"context"
	"errors"
	"time"
)

var ErrKeyDoesNotExist = errors.New("key does not exist")

type (
	ScanIterator interface {
		Err() error
		Val() string
		Next(ctx context.Context) bool
	}

	Scanner interface {
		Scan(ctx context.Context) ScanIterator
	}

	CursorScanner interface {
		Scanner
		WithArgs(cursor uint64, match string, count int64) CursorScanner
	}

	KeyCursorScanner interface {
		Scanner
		WithArgs(key string, cursor uint64, match string, count int64) KeyCursorScanner
	}
)

var (
	// DefaultScanner corresponds to the Redis [SCAN] command
	//
	// [SCAN]: https://redis.io/docs/latest/commands/scan
	DefaultScanner CursorScanner

	// DefaultSetScanner corresponds to the Redis [SSCAN] command
	//
	// [SSCAN]: https://redis.io/docs/latest/commands/sscan
	DefaultSetScanner KeyCursorScanner

	// DefaultHashScanner corresponds to the Redis [HSCAN] command
	//
	// [HSCAN]: https://redis.io/docs/latest/commands/hscan
	DefaultHashScanner KeyCursorScanner

	// DefaultSortedSetScanner corresponds to the Redis [ZSCAN] command
	//
	// [ZSCAN]: https://redis.io/docs/latest/commands/zscan
	DefaultSortedSetScanner KeyCursorScanner
)

type (
	Conn interface {
		Get(ctx context.Context, key string, dest any) error
		Set(ctx context.Context, key string, value any, expiration time.Duration) error
		Scan(ctx context.Context, scanner Scanner) ScanIterator
		Delete(ctx context.Context, key string) (int64, error)
		Exists(ctx context.Context, keys ...string) (int64, error)
		Begin(ctx context.Context) Tx
	}

	Tx interface {
		Conn
		Exec(ctx context.Context) error
		Discard(ctx context.Context) error
	}
)
