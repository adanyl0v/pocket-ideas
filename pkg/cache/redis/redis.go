package redis

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/adanyl0v/pocket-ideas/pkg/cache"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/adanyl0v/pocket-ideas/pkg/proxerr"
	"github.com/redis/go-redis/v9"
)

type (
	DriverConn interface {
		Get(ctx context.Context, key string) *redis.StringCmd
		Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
		Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
		SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd
		HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd
		ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd
		Del(ctx context.Context, keys ...string) *redis.IntCmd
		Exists(ctx context.Context, keys ...string) *redis.IntCmd
		TxPipeline() redis.Pipeliner
	}

	DriverPipeline interface {
		DriverConn
		Exec(ctx context.Context) ([]redis.Cmder, error)
		Discard()
	}
)

func newConn(conn DriverConn, logger log.Logger) Conn {
	return Conn{
		conn:   conn,
		logger: logger,
	}
}

type Conn struct {
	conn   DriverConn
	logger log.Logger
}

func (c *Conn) DriverConn() DriverConn {
	return c.conn
}

func (c *Conn) Get(ctx context.Context, key string, dest any) error {
	logger := c.logger.With(log.Fields{"key": key})

	if err := c.conn.Get(ctx, key).Scan(dest); err != nil {
		if errors.Is(err, redis.Nil) {
			err = proxerr.New(cache.ErrKeyDoesNotExist, err.Error())
		}

		logger.WithError(err).Error("failed to get the key")
		return err
	}

	logger.Debug("got the key")
	return nil
}

func (c *Conn) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	logger := c.logger.With(log.Fields{
		"key":        key,
		"expiration": expiration,
	})

	if err := c.conn.Set(ctx, key, value, expiration).Err(); err != nil {
		logger.WithError(err).Error("failed to set the key")
		return err
	}

	logger.Debug("set the key")
	return nil
}

func (c *Conn) Scan(ctx context.Context, scanner cache.Scanner) cache.ScanIterator {
	return scanner.Scan(ctx)
}

func (c *Conn) Delete(ctx context.Context, key string) (int64, error) {
	logger := c.logger.With(log.Fields{"key": key})

	n, err := c.conn.Del(ctx, key).Result()
	if err != nil {
		logger.WithError(err).Error("failed to delete the key")
		return 0, err
	}

	logger.Debug("deleted the key")
	return n, nil
}

func (c *Conn) Exists(ctx context.Context, keys ...string) (int64, error) {
	logger := c.logger.With(log.Fields{"keys": keys})

	n, err := c.conn.Exists(ctx, keys...).Result()
	if err != nil {
		logger.WithError(err).Error("failed to check keys existence")
		return 0, err
	}

	logger.Debug(fmt.Sprintf("%d out of %d keys exist", n, len(keys)))
	return n, nil
}

func (c *Conn) Begin(_ context.Context) cache.Tx {
	return newTx(c.conn.TxPipeline(), c.logger)
}

type Tx struct {
	Conn
}

func newTx(pipeline DriverPipeline, logger log.Logger) *Tx {
	return &Tx{
		Conn: Conn{
			conn:   pipeline,
			logger: logger,
		},
	}
}

func (t *Tx) Exec(ctx context.Context) error {
	tx, ok := t.conn.(DriverPipeline)
	if !ok {
		t.logger.Error("failed to cast the driver connection to a pipeline")
		return nil
	}

	if _, err := tx.Exec(ctx); err != nil {
		t.logger.WithError(err).Error("failed to execute commands")
		return err
	}

	t.logger.Debug("executed commands")
	return nil
}

func (t *Tx) Discard(_ context.Context) error {
	tx, ok := t.conn.(DriverPipeline)
	if !ok {
		t.logger.Error("failed to cast the driver connection to a pipeline")
		return nil
	}

	tx.Discard()

	t.logger.Debug("discarded commands")
	return nil
}

type Client struct {
	Conn
	redisClient *redis.Client
}

func (c *Client) Close() error {
	if err := c.redisClient.Close(); err != nil {
		c.logger.WithError(err).Error("failed to close the redis connection")
		return err
	}

	c.logger.Info("closed the redis connection")
	return nil
}

type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        int
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxIdleConns    int
	MinIdleConns    int
	MaxActiveConns  int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
	TLSConfig       *tls.Config
}

func Connect(ctx context.Context, logger log.Logger, config *Config) (*Client, error) {
	opts := redis.Options{
		Addr:            fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username:        config.User,
		Password:        config.Password,
		DB:              config.Database,
		DialTimeout:     config.DialTimeout,
		ReadTimeout:     config.ReadTimeout,
		WriteTimeout:    config.WriteTimeout,
		MinIdleConns:    config.MinIdleConns,
		MaxIdleConns:    config.MaxIdleConns,
		MaxActiveConns:  config.MaxActiveConns,
		ConnMaxIdleTime: config.ConnMaxIdleTime,
		ConnMaxLifetime: config.ConnMaxLifetime,
		TLSConfig:       config.TLSConfig,
	}
	logger.With(log.Fields{
		"addr":          opts.Addr,
		"user":          opts.Username,
		"database":      opts.DB,
		"dial_timeout":  opts.DialTimeout,
		"read_timeout":  opts.ReadTimeout,
		"write_timeout": opts.WriteTimeout,
	}).Debug("parsed the redis connection config")

	client := redis.NewClient(&opts)
	logger.With(log.Fields{
		"min_idle_conns":     opts.MinIdleConns,
		"max_idle_conns":     opts.MaxIdleConns,
		"max_active_conns":   opts.MaxActiveConns,
		"conn_max_idle_time": opts.ConnMaxIdleTime,
		"conn_max_lifetime":  opts.ConnMaxLifetime,
	}).Debug("created a new redis client")

	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.WithError(err).Error("failed to ping the redis connection")
		return nil, err
	}

	stats := client.PoolStats()
	logger.With(log.Fields{
		"hits":        stats.Hits,
		"misses":      stats.Misses,
		"timeouts":    stats.Timeouts,
		"total_conns": stats.TotalConns,
		"idle_conns":  stats.IdleConns,
		"stale_conns": stats.StaleConns,
	}).Debug("pinged the redis connection")
	return &Client{
		Conn: Conn{
			conn:   client,
			logger: logger,
		},
		redisClient: client,
	}, nil
}
