package config

import "time"

const (
	EnvLocal = "local"
	EnvDev   = "dev"
)

const (
	LogLevelTrace = "trace"
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

type Config struct {
	Env            string         `yaml:"env" env:"ENV" env-required:"true"`
	Log            LogConfig      `yaml:"log"`
	PostgresConfig PostgresConfig `yaml:"postgres"`
}

type LogConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"warn"`
}

type PostgresConfig struct {
	Host              string        `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port              int           `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	User              string        `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password          string        `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	Database          string        `yaml:"database" env:"POSTGRES_DATABASE" env-required:"true"`
	Schema            string        `yaml:"schema" env:"POSTGRES_SCHEMA" env-default:"public"`
	SSLMode           string        `yaml:"ssl_mode" env:"POSTGRES_SSL_MODE" env-default:"disable"`
	ConnTimout        time.Duration `yaml:"conn_timout" env:"POSTGRES_CONN_TIMOUT" env-default:"5s"`
	MaxConns          int           `yaml:"max_conns" env:"POSTGRES_MAX_CONNS" env-default:"4"`
	MinConns          int           `yaml:"min_conns" env:"POSTGRES_MIN_CONNS" env-default:"0"`
	MaxConnLifetime   time.Duration `yaml:"max_conn_lifetime" env:"POSTGRES_MAX_CONN_LIFETIME" env-default:"60m"`
	MaxConnIdleTime   time.Duration `yaml:"max_conn_idle_time" env:"POSTGRES_MAX_CONN_IDLE_TIME" env-default:"30m"`
	HealthCheckPeriod time.Duration `yaml:"health_check_period" env:"POSTGRES_HEALTH_CHECK_PERIOD" env-default:"1m"`
}
