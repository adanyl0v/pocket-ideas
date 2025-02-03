package log

import (
	"context"
)

type Level int

type Leveler interface {
	Level() Level
}

type Fields map[string]any

type Logger interface {
	Log(leveler Leveler, msg string)
	Trace(msg string)
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	With(fields Fields) Logger
	WithError(err error) Logger
	WithContext(ctx context.Context) Logger
	WithCallerSkip(skip int) Logger
}
