package redis

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/adanyl0v/pocket-ideas/pkg/proxerr"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	ErrNonexistentKey = errors.New("the key does not exist")
)

const (
	Root = "$"
	Keep = redis.KeepTTL
)

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

type Client struct {
	client *redis.Client
	logger log.Logger
}

func New(client *redis.Client, logger log.Logger) *Client {
	return &Client{
		client: client,
		logger: logger,
	}
}

func (c *Client) Close() error {
	if err := c.client.Close(); err != nil {
		c.logger.WithError(err).Debug("failed to close the redis connection")
		return err
	}

	c.logger.Debug("closed the redis connection")
	return nil
}

func (c *Client) Get(ctx context.Context, key string, dest any) error {
	logger := c.logger.With(log.Fields{"key": key})

	cmd := c.client.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			logger.WithError(err).Debug("key does not exist")
			return proxerr.New(ErrNonexistentKey, err.Error())
		}

		logger.WithError(err).Debug("failed to get the key")
		return err
	}

	logger.Debug("got the key")
	return cmd.Scan(dest)
}

func (c *Client) GetJSON(ctx context.Context, key, path string, dest any) error {
	logger := c.logger.With(log.Fields{
		"key":  key,
		"path": path,
	})

	cmd := c.client.JSONGet(ctx, key, path)
	if err := cmd.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			logger.WithError(err).Debug("key does not exist")
			return proxerr.New(ErrNonexistentKey, err.Error())
		}

		logger.WithError(err).Debug("failed to get the key")
		return err
	}

	logger.Debug("got the key")
	return json.Unmarshal([]byte(cmd.String()), dest)
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	logger := c.logger.With(log.Fields{
		"key": key,
		"exp": expiration,
	})

	cmd := c.client.Set(ctx, key, value, expiration)
	if err := cmd.Err(); err != nil {
		logger.WithError(err).Debug("failed to set the key")
		return err
	}

	logger.Debug("set the key")
	return nil
}

func (c *Client) SetJSON(ctx context.Context, key, path string, value any, expiration time.Duration) error {
	logger := c.logger.With(log.Fields{
		"key":  key,
		"path": path,
		"exp":  expiration,
	})

	cmd := c.client.JSONSet(ctx, key, path, value)
	if err := cmd.Err(); err != nil {
		logger.WithError(err).Debug("failed to set the key")
		return err
	}

	if err := c.client.Expire(ctx, key, expiration).Err(); err != nil {
		logger.WithError(err).Debug("failed to expire the key")
		return err
	}

	logger.Debug("set the key")
	return nil
}

func (c *Client) Delete(ctx context.Context, keys ...string) error {
	logger := c.logger.With(log.Fields{"keys": keys})

	cmd := c.client.Del(ctx, keys...)
	if err := cmd.Err(); err != nil {
		logger.WithError(err).Debug("failed to delete the keys")
		return err
	}

	logger.Debug("deleted the keys")
	return nil
}

func (c *Client) DeleteJSON(ctx context.Context, key, path string) error {
	logger := c.logger.With(log.Fields{
		"key":  key,
		"path": path,
	})

	cmd := c.client.JSONDel(ctx, key, path)
	if err := cmd.Err(); err != nil {
		logger.WithError(err).Debug("failed to delete the key")
		return err
	}

	logger.Debug("deleted the key")
	return nil
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
	return New(client, logger), nil
}
