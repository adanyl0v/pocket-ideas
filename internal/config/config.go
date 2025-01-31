package config

const (
	EnvLocal = "local"
	EnvDev   = "dev"
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

type Config struct {
	Env string    `yaml:"env" env:"ENV" env-required:"true"`
	Log LogConfig `yaml:"log"`
}

type LogConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"warn"`
}
