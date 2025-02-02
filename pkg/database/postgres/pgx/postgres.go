package pgx

import (
	"context"
	"errors"
	"github.com/adanyl0v/pocket-ideas/pkg/database/postgres"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/adanyl0v/pocket-ideas/pkg/proxerr"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
	"time"
)

type Row struct {
	row pgx.Row
}

func newRow(row pgx.Row) *Row {
	return &Row{row: row}
}

func (r *Row) Scan(dest ...any) error {
	err := r.row.Scan(dest...)
	if errors.Is(err, pgx.ErrNoRows) {
		return proxerr.New(postgres.ErrNoRows, err.Error())
	} else {
		err = handleCommonExecErrors(err)
	}

	return err
}

type Rows struct {
	rows pgx.Rows
}

func newRows(rows pgx.Rows) *Rows {
	return &Rows{rows: rows}
}

func (r *Rows) Scan(dest ...any) error {
	return handleCommonExecErrors(r.rows.Scan(dest...))
}

func (r *Rows) Err() error {
	err := r.rows.Err()
	if errors.Is(err, pgx.ErrNoRows) {
		return proxerr.New(postgres.ErrNoRows, err.Error())
	} else {
		err = handleCommonExecErrors(err)
	}

	return err
}

func (r *Rows) Next() bool {
	return r.rows.Next()
}

func (r *Rows) Close() {
	r.rows.Close()
}

type pgxExecutor interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

func exec(executor pgxExecutor, logger log.Logger, ctx context.Context, query string, args ...any) error {
	logger = logger.With(log.Fields{"sql": inlineQuery(query)})

	startedAt := time.Now()
	_, err := executor.Exec(ctx, query, args...)
	endedAt := time.Now()

	if err != nil {
		err = handleCommonExecErrors(err)
		logger.WithError(err).Error("failed")
		return err
	}

	logger.With(log.Fields{"duration": endedAt.Sub(startedAt)}).Debug("done")
	return nil
}

type pgxRowQuerier interface {
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

func queryRow(querier pgxRowQuerier, logger log.Logger, ctx context.Context, query string, args ...any) *Row {
	logger = logger.With(log.Fields{"sql": inlineQuery(query)})

	startedAt := time.Now()
	row := querier.QueryRow(ctx, query, args...)
	endedAt := time.Now()

	logger.With(log.Fields{"duration": endedAt.Sub(startedAt)}).Debug("done")
	return newRow(row)
}

type pgxQuerier interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

func queryRows(querier pgxQuerier, logger log.Logger, ctx context.Context, query string, args ...any) (*Rows, error) {
	logger = logger.With(log.Fields{"sql": inlineQuery(query)})

	startedAt := time.Now()
	rows, err := querier.Query(ctx, query, args...)
	endedAt := time.Now()

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = proxerr.New(postgres.ErrNoRows, err.Error())
		} else {
			err = handleCommonExecErrors(err)
		}

		logger.WithError(err).Error("failed")
		return nil, err
	}

	logger.With(log.Fields{"duration": endedAt.Sub(startedAt)}).Debug("done")
	return newRows(rows), nil
}

type pgxBeginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

func begin(beginner pgxBeginner, logger log.Logger, ctx context.Context) (*Tx, error) {
	tx, err := beginner.Begin(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to start a transaction")
		return nil, err
	}

	logger.Debug("begun a new transaction")
	return newTx(tx, logger), nil
}

func handleCommonExecErrors(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		var commonErr error

		switch pgErr.Code {
		case pgerrcode.CheckViolation:
			commonErr = postgres.ErrCheckViolation
		case pgerrcode.UniqueViolation:
			commonErr = postgres.ErrUniqueViolation
		case pgerrcode.NotNullViolation:
			commonErr = postgres.ErrNotNullViolation
		case pgerrcode.ForeignKeyViolation:
			commonErr = postgres.ErrForeignKeyViolation
		}

		err = proxerr.New(commonErr, err.Error())
	}

	return err
}

func inlineQuery(query string) string {
	return strings.TrimSpace(strings.ReplaceAll(query, "\n", " "))
}
