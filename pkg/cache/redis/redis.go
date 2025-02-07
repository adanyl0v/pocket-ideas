package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
	"time"
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
	rejson rejson.ReJSON
	logger log.Logger
}

func New(client *redis.Client, rejson rejson.ReJSON, logger log.Logger) *Client {
	return &Client{
		client: client,
		rejson: rejson,
		logger: logger,
	}
}

func (c *Client) Close() error {
	return c.client.Close()
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
	client := redis.NewClient(&opts)
	logger.With(log.Fields{
		"addr":     opts.Addr,
		"database": opts.DB,
	}).Debug("created a new redis client")

	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.WithError(err).Error("failed to ping the redis connection")
		return nil, err
	}

	logger.Debug("pinged the redis connection")

	rh := rejson.NewReJSONHandler()
	rh.SetGoRedisClientWithContext(ctx, client)
	return New(client, rh, logger), nil
}
