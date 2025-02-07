package app

import (
	"context"
	"fmt"
	"github.com/adanyl0v/pocket-ideas/internal/config"
	"github.com/adanyl0v/pocket-ideas/internal/repository/postgres"
	"github.com/adanyl0v/pocket-ideas/pkg/cache/redis"
	"github.com/adanyl0v/pocket-ideas/pkg/database/postgres/pgx"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/adanyl0v/pocket-ideas/pkg/log/slog"
	googleuuidgen "github.com/adanyl0v/pocket-ideas/pkg/uuid/google"
	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	stdslog "log/slog"
	"os"
	"path/filepath"
)

func Run() {
	cfg := config.MustReadFile(config.DefaultFilePath())
	logger := mustSetupLogger(cfg.Env, &cfg.Log)
	logger.With(log.Fields{"env": cfg.Env}).Info("read config")

	db := mustConnectToPostgres(logger, &cfg.PostgresConfig)
	defer func() { _ = db.Close() }()

	redisCache := mustConnectToRedis(logger, &cfg.RedisConfig)
	defer func() { _ = redisCache.Close() }()

	userRepo := postgres.NewUserRepository(db, logger, googleuuidgen.New())
	_ = userRepo
}

func mustSetupLogger(env string, cfg *config.LogConfig) log.Logger {
	var (
		zapLevel  zapcore.Level
		slogLevel stdslog.Level
	)
	switch cfg.Level {
	case config.LogLevelTrace:
		fmt.Println("used logger doesn't support trace level, so using debug as the closest one")
		cfg.Level = config.LogLevelDebug

		zapLevel = zapcore.DebugLevel
		slogLevel = stdslog.LevelDebug
	case config.LogLevelDebug:
		zapLevel = zapcore.DebugLevel
		slogLevel = stdslog.LevelDebug
	case config.LogLevelInfo:
		zapLevel = zapcore.InfoLevel
		slogLevel = stdslog.LevelInfo
	case config.LogLevelWarn:
		zapLevel = zapcore.WarnLevel
		slogLevel = stdslog.LevelWarn
	case config.LogLevelError:
		zapLevel = zapcore.ErrorLevel
		slogLevel = stdslog.LevelError
	default:
		panic(fmt.Errorf("invalid log level: %s", cfg.Level))
	}

	var zapEncoder zapcore.Encoder
	switch env {
	case config.EnvLocal:
		zapEncoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey:     "@MESSAGE",
			LevelKey:       "@LEVEL",
			TimeKey:        "@TIMESTAMP",
			NameKey:        "@NAME",
			CallerKey:      "@FILE",
			FunctionKey:    "@FUNCTION",
			StacktraceKey:  "@STACKTRACE",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout("02/01/2006 15:04:05"),
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller: func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
				cwd, _ := os.Getwd()
				file, _ := filepath.Rel(cwd, caller.FullPath())
				encoder.AppendString(file)
			},
			EncodeName: zapcore.FullNameEncoder,
		})
	case config.EnvDev:
		zapEncoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:     "@message",
			LevelKey:       "@level",
			TimeKey:        "@timestamp",
			NameKey:        "@name",
			CallerKey:      "@file",
			StacktraceKey:  "@stacktrace",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		})
	default:
		panic(fmt.Errorf("invalid env: %s", env))
	}

	zapCore := zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLevel)
	zapLogger := zap.New(zapCore)
	defer func() { _ = zapLogger.Sync() }()

	zapHandler := slogzap.Option{
		Level:     slogLevel,
		Logger:    zapLogger,
		AddSource: true,
	}.NewZapHandler()

	slog.ErrorFieldKey = "error"
	l := slog.NewLogger(stdslog.New(zapHandler))
	l = l.With(log.Fields{"pid": os.Getpid()}).(*slog.Logger)

	l.With(log.Fields{"level": cfg.Level}).Info("initialized logger")
	return l
}

func mustConnectToPostgres(logger log.Logger, cfg *config.PostgresConfig) *pgx.DB {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnTimout)
	defer cancel()

	connConfig := pgx.Config{
		Host:              cfg.Host,
		Port:              cfg.Port,
		User:              cfg.User,
		Password:          cfg.Password,
		Database:          cfg.Database,
		SSLMode:           cfg.SSLMode,
		MaxConns:          cfg.MaxConns,
		MinConns:          cfg.MinConns,
		MaxConnLifetime:   cfg.MaxConnLifetime,
		MaxConnIdleTime:   cfg.MaxConnIdleTime,
		HealthCheckPeriod: cfg.HealthCheckPeriod,
	}

	db, err := pgx.Connect(ctx, &connConfig, logger)
	if err != nil {
		panic(err)
	}

	logger.Info("connected to postgres")
	return db
}

func mustConnectToRedis(logger log.Logger, cfg *config.RedisConfig) *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	client, err := redis.Connect(ctx, logger, &redis.Config{
		Host:        cfg.Host,
		Port:        cfg.Port,
		User:        cfg.User,
		Password:    cfg.Password,
		Database:    cfg.Database,
		DialTimeout: cfg.DialTimeout,
	})
	if err != nil {
		panic(err)
	}

	logger.Info("connected to redis")
	return client
}
