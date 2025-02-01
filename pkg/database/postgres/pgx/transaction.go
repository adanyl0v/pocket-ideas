package pgx

import (
	"context"
	"github.com/adanyl0v/pocket-ideas/pkg/database/postgres"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/jackc/pgx/v5"
)

type Tx struct {
	tx     pgx.Tx
	logger log.Logger
}

func newTx(tx pgx.Tx, logger log.Logger) *Tx {
	return &Tx{
		tx:     tx,
		logger: logger,
	}
}

func (t *Tx) Exec(ctx context.Context, query string, args ...any) error {
	return exec(t.tx, t.logger, ctx, query, args...)
}

func (t *Tx) Query(ctx context.Context, query string, args ...any) (postgres.Rows, error) {
	return queryRows(t.tx, t.logger, ctx, query, args...)
}

func (t *Tx) QueryRow(ctx context.Context, query string, args ...any) postgres.Row {
	return queryRow(t.tx, t.logger, ctx, query, args...)
}

func (t *Tx) Begin(ctx context.Context) (postgres.Tx, error) {
	return begin(t.tx, t.logger, ctx)
}

// [Tx.Commit] and [Tx.Rollback] method calls must occur in helper functions for proper logging.
// Since the database object stores a logger with a pre-increased number of call frame skips,
// if you implement these methods directly, the wrong path to the caller will be logged.

func (t *Tx) Commit(ctx context.Context) error {
	return t.commit(ctx)
}

func (t *Tx) Rollback(ctx context.Context) error {
	return t.rollback(ctx)
}

func (t *Tx) commit(ctx context.Context) error {
	if err := t.tx.Commit(ctx); err != nil {
		t.logger.WithError(err).Error("failed to commit the transaction")
		return err
	}

	t.logger.Debug("committed the transaction")
	return nil
}

func (t *Tx) rollback(ctx context.Context) error {
	if err := t.tx.Rollback(ctx); err != nil {
		t.logger.WithError(err).Error("failed to rollback the transaction")
		return err
	}

	t.logger.Debug("rolled back the transaction")
	return nil
}
