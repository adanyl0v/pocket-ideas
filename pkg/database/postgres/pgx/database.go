package pgx

import (
	"context"
	"fmt"
	"github.com/adanyl0v/pocket-ideas/pkg/database/postgres"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxuuid "github.com/vgarvardt/pgx-google-uuid/v5"

	"time"
)

const DriverName = "pgx"

type Config struct {
	Host              string
	Port              int
	User              string
	Password          string
	Database          string
	SSLMode           string
	MaxConns          int
	MinConns          int
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

func (c *Config) URL() string {
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode)
}

type DB struct {
	pool   *pgxpool.Pool
	logger log.Logger
}

func New(pool *pgxpool.Pool, logger log.Logger) *DB {
	return &DB{
		pool:   pool,
		logger: logger,
	}
}

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

func Connect(ctx context.Context, config *Config, logger log.Logger) (*DB, error) {
	c, err := pgxpool.ParseConfig(config.URL())
	if err != nil {
		logger.WithError(err).Error("failed to parse the connection config")
		return nil, err
	}
	logger.With(log.Fields{
		"driver":   DriverName,
		"host":     config.Host,
		"port":     config.Port,
		"database": config.Database,
		"sslmode":  config.SSLMode,
	}).Debug("parsed the connection config")

	c.MaxConns = int32(config.MaxConns)
	c.MinConns = int32(config.MinConns)
	c.MaxConnLifetime = config.MaxConnLifetime
	c.MaxConnIdleTime = config.MaxConnIdleTime
	c.HealthCheckPeriod = config.HealthCheckPeriod

	// Ensure that pgx supports google UUIDs
	c.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		logger.WithError(err).Error("failed to create a new connection pool")
		return nil, err
	}
	logger.With(log.Fields{
		"max_conns":           c.MaxConns,
		"min_conns":           c.MinConns,
		"max_conn_lifetime":   c.MaxConnLifetime,
		"max_conn_idle_time":  c.MaxConnIdleTime,
		"health_check_period": c.HealthCheckPeriod,
	}).Debug("created a new connection pool")

	if err = pool.Ping(ctx); err != nil {
		logger.WithError(err).Error("failed to ping the connection")
		return nil, err
	}

	stat := pool.Stat()
	logger.With(log.Fields{
		"acquired_total": stat.AcquireCount(),
		"acquired_now":   stat.AcquiredConns(),
		"acquired_empty": stat.EmptyAcquireCount(),
	}).Debug("pinged the connection")

	return New(pool, logger.WithCallerSkip(1)), nil
}

func (db *DB) Exec(ctx context.Context, query string, args ...any) error {
	return exec(db.pool, db.logger, ctx, query, args...)
}

func (db *DB) Query(ctx context.Context, query string, args ...any) (postgres.Rows, error) {
	return queryRows(db.pool, db.logger, ctx, query, args...)
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...any) postgres.Row {
	return queryRow(db.pool, db.logger, ctx, query, args...)
}

func (db *DB) Begin(ctx context.Context) (postgres.Tx, error) {
	return begin(db.pool, db.logger, ctx)
}
