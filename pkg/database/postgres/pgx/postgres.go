package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/adanyl0v/pocket-ideas/pkg/database"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/adanyl0v/pocket-ideas/pkg/proxerr"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxuuid "github.com/vgarvardt/pgx-google-uuid/v5"
	"time"
)

const DriverName = "pgx"

var ErrNotTransaction = errors.New("the connection is not a transaction")

type DriverConn interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Begin(context.Context) (pgx.Tx, error)
}

type DriverTx interface {
	DriverConn
	Commit(context.Context) error
	Rollback(context.Context) error
}

type Conn struct {
	conn   DriverConn
	logger log.Logger
}

func newConn(conn DriverConn, logger log.Logger) Conn {
	return Conn{
		conn:   conn,
		logger: logger,
	}
}

func (c *Conn) Execute(ctx context.Context, query string, args ...any) (database.Result, error) {
	logger := c.logger.With(log.Fields{"query": query})

	now := time.Now()
	tag, err := c.conn.Exec(ctx, query, args...)
	t := time.Since(now)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.CheckViolation:
				err = proxerr.New(database.ErrCheckViolation, pgErr.Error())
			case pgerrcode.UniqueViolation:
				err = proxerr.New(database.ErrUniqueViolation, pgErr.Error())
			case pgerrcode.NotNullViolation:
				err = proxerr.New(database.ErrNotNullViolation, pgErr.Error())
			case pgerrcode.ForeignKeyViolation:
				err = proxerr.New(database.ErrForeignKeyViolation, pgErr.Error())
			}

			logger = logger.With(log.Fields{"driverError": *pgErr})
		}

		logger.WithError(err).Error("failed sql query execution")
		return nil, err
	}

	logger.With(log.Fields{"duration": t}).Debug("executed sql")
	return tag, err
}

func (c *Conn) Query(ctx context.Context, query string, args ...any) (database.Rows, error) {
	logger := c.logger.With(log.Fields{"query": query})

	now := time.Now()
	rows, err := c.conn.Query(ctx, query, args...)
	t := time.Since(now)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && errors.Is(pgErr, pgx.ErrNoRows) {
			err = proxerr.New(database.ErrNoRows, pgErr.Error())
			logger = logger.With(log.Fields{"driverError": *pgErr})
		}

		logger.WithError(err).Error("failed sql query execution")
		return nil, err
	}

	logger.With(log.Fields{"duration": t}).Debug("executed sql")
	return rows, err
}

func (c *Conn) QueryRow(ctx context.Context, query string, args ...any) database.Row {
	logger := c.logger.With(log.Fields{"query": query})

	now := time.Now()
	row := c.conn.QueryRow(ctx, query, args...)
	t := time.Since(now)

	logger.With(log.Fields{"duration": t}).Debug("executed sql")
	return row
}

func (c *Conn) Begin(ctx context.Context) (database.Tx, error) {
	tx, err := c.conn.Begin(ctx)
	if err != nil {
		c.logger.WithError(err).Error("failed to begin an sql transaction")
		return nil, err
	}

	c.logger.Debug("begun an sql transaction")
	return newTx(newConn(tx, c.logger)), nil
}

type Config struct {
	Host              string
	Port              int
	User              string
	Password          string
	Database          string
	SSLMode           string
	MaxConns          int
	MinConns          int
	MaxConnIdleTime   time.Duration
	MaxConnLifetime   time.Duration
	HealthCheckPeriod time.Duration
}

func (c *Config) URL() string {
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode)
}

type Client struct {
	Conn
	pool *pgxpool.Pool
}

func (c *Client) Pool() *pgxpool.Pool {
	return c.pool
}

func (c *Client) Close() {
	c.pool.Close()
	c.logger.Info("closed the postgres connection")
}

func Connect(ctx context.Context, logger log.Logger, config *Config) (*Client, error) {
	pc, err := pgxpool.ParseConfig(config.URL())
	if err != nil {
		logger.WithError(err).Error("failed to parse the postgres connection config")
		return nil, err
	}
	logger.With(log.Fields{
		"driver":   DriverName,
		"host":     config.Host,
		"port":     config.Port,
		"user":     config.User,
		"database": config.Database,
		"sslmode":  config.SSLMode,
	}).Debug("parsed the postgres connection config")

	// Ensure that there is no preparation and the entire request goes through in
	// a single network call [https://habr.com/ru/companies/avito/articles/461935]
	pc.ConnConfig.RuntimeParams = map[string]string{"standard_conforming_strings": "on"}
	pc.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pc.MaxConns = int32(config.MaxConns)
	pc.MinConns = int32(config.MinConns)
	pc.MaxConnIdleTime = config.MaxConnIdleTime
	pc.MaxConnLifetime = config.MaxConnLifetime
	pc.HealthCheckPeriod = config.HealthCheckPeriod

	// Ensure that pgx supports google UUIDs
	pc.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, pc)
	if err != nil {
		logger.WithError(err).Error("failed to create a postgres connection pool")
		return nil, err
	}
	logger.With(log.Fields{
		"max_conns":           pc.MaxConns,
		"min_conns":           pc.MinConns,
		"max_conn_lifetime":   pc.MaxConnLifetime,
		"max_conn_idle_time":  pc.MaxConnIdleTime,
		"health_check_period": pc.HealthCheckPeriod,
	}).Debug("created a new postgres connection pool")

	if err = pool.Ping(ctx); err != nil {
		logger.WithError(err).Error("failed to ping the postgres database")
		return nil, err
	}

	stat := pool.Stat()
	logger.With(log.Fields{
		"acquired_total": stat.AcquireCount(),
		"acquired_now":   stat.AcquiredConns(),
		"acquired_empty": stat.EmptyAcquireCount(),
	}).Debug("pinged the postgres connection")

	return &Client{
		Conn: newConn(pool, logger),
		pool: pool,
	}, nil
}

type Tx struct {
	Conn
}

func (t *Tx) Commit(ctx context.Context) error {
	tx, ok := t.conn.(DriverTx)
	if !ok {
		return ErrNotTransaction
	}

	if err := tx.Commit(ctx); err != nil {
		t.logger.WithError(err).Error("failed to commit an sql transaction")
		return err
	}

	t.logger.Debug("committed the sql transaction")
	return nil
}

func (t *Tx) Rollback(ctx context.Context) error {
	tx, ok := t.conn.(DriverTx)
	if !ok {
		return ErrNotTransaction
	}

	if err := tx.Rollback(ctx); err != nil {
		t.logger.WithError(err).Error("failed to rollback an sql transaction")
		return err
	}

	t.logger.Debug("rolled back the sql transaction")
	return nil
}

func newTx(conn Conn) *Tx {
	return &Tx{
		Conn: conn,
	}
}
