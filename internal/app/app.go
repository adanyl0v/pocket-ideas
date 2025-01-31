package app

import (
	"fmt"
	"github.com/adanyl0v/pocket-ideas/internal/config"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"github.com/adanyl0v/pocket-ideas/pkg/log/slog"
	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	stdslog "log/slog"
	"os"
	"path/filepath"
)

var logger log.Logger

func Run() {
	cfg := config.MustReadFile(config.DefaultFilePath())
	logger = mustSetupLogger(cfg.Env, &cfg.Log)
	logger.With(log.Fields{"env": cfg.Env}).Info("read config")
}

func mustSetupLogger(env string, cfg *config.LogConfig) log.Logger {
	var (
		zapLevel  zapcore.Level
		slogLevel stdslog.Level
	)
	switch cfg.Level {
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
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
			EncodeDuration: zapcore.SecondsDurationEncoder,
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
	l.With(log.Fields{"level": cfg.Level}).Info("initialized logger")
	return l
}
