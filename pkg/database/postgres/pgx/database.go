package pgx

import (
	"context"
	"fmt"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
)

const DriverName = "pgx"

type Config struct {
	Host              string
	Port              int
	User              string
	Password          string
	Database          string
	Schema            string
	SSLMode           string
	MaxConns          int
	MinConns          int
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

func (c *Config) URL() string {
	if c.Schema == "" {
		c.Schema = "public"
	}
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?search_path=%s&sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.Schema, c.SSLMode)
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
	logger = logger.With(log.Fields{"pid": os.Getpid()})

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
		"schema":   config.Schema,
		"sslmode":  config.SSLMode,
	}).Debug("parsed the connection config")

	c.MaxConns = int32(config.MaxConns)
	c.MinConns = int32(config.MinConns)
	c.MaxConnLifetime = config.MaxConnLifetime
	c.MaxConnIdleTime = config.MaxConnIdleTime
	c.HealthCheckPeriod = config.HealthCheckPeriod

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

	return New(pool, logger), nil
}
