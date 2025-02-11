package database

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNoRows              = sql.ErrNoRows
	ErrCheckViolation      = errors.New("check violation")
	ErrUniqueViolation     = errors.New("unique violation")
	ErrNotNullViolation    = errors.New("not null violation")
	ErrForeignKeyViolation = errors.New("foreign key violation")
)

type Result interface {
	String() string
	Insert() bool
	Select() bool
	Update() bool
	Delete() bool
	RowsAffected() int64
}

type Row interface {
	Scan(dest ...any) error
}

type Rows interface {
	Row
	Err() error
	Next() bool
	Values() ([]any, error)
	Close()
}

type Conn interface {
	Execute(ctx context.Context, query string, args ...any) (Result, error)
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
	Begin(ctx context.Context) (Tx, error)
}

type Tx interface {
	Conn
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
