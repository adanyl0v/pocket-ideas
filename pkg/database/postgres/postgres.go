package postgres

import (
	"context"
	"errors"
)

var (
	ErrNoRows              = errors.New("no rows in result set")
	ErrCheckViolation      = errors.New("check violation")
	ErrUniqueViolation     = errors.New("unique violation")
	ErrNotNullViolation    = errors.New("not null violation")
	ErrForeignKeyViolation = errors.New("foreign key violation")
)

type Conn interface {
	Executor
	Querier
	RowQuerier
	Beginner
}

type Tx interface {
	Conn
	Committer
	Rollbacker
}

type Row interface {
	Scanner
}

type Rows interface {
	Scanner
	Err() error
	Next() bool
	Close()
}

type Scanner interface {
	Scan(dest ...any) error
}

type Executor interface {
	Exec(ctx context.Context, query string, args ...any) error
}

type Querier interface {
	Query(ctx context.Context, query string, args ...any) (Rows, error)
}

type RowQuerier interface {
	QueryRow(ctx context.Context, query string, args ...any) Row
}

type Beginner interface {
	Begin(ctx context.Context) (Tx, error)
}

type Committer interface {
	Commit(ctx context.Context) error
}

type Rollbacker interface {
	Rollback(ctx context.Context) error
}
