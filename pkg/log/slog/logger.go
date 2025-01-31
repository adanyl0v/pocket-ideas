package slog

import (
	"context"
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"log/slog"
	"runtime"
	"slices"
	"time"
)

var ErrorFieldKey = "error"

type Logger struct {
	l     *slog.Logger
	ctx   context.Context
	attrs []slog.Attr
	skip  int
}

func NewLogger(logger *slog.Logger) *Logger {
	return &Logger{
		l:     logger,
		ctx:   context.Background(),
		attrs: make([]slog.Attr, 4),
		skip:  1,
	}
}

func (l *Logger) Log(leveler log.Leveler, msg string) {
	l.log(leveler, msg)
}

func (l *Logger) Debug(msg string) {
	l.log(DebugLevel, msg)
}

func (l *Logger) Info(msg string) {
	l.log(InfoLevel, msg)
}

func (l *Logger) Warn(msg string) {
	l.log(WarnLevel, msg)
}

func (l *Logger) Error(msg string) {
	l.log(ErrorLevel, msg)
}

func (l *Logger) With(fields log.Fields) log.Logger {
	clone := l.clone()

	clone.attrs = slices.Clone(l.attrs)
	slices.Grow(clone.attrs, len(fields)/2)

	computeFields(&clone.attrs, fields)
	return clone
}

func (l *Logger) WithError(err error) log.Logger {
	clone := l.clone()

	clone.attrs = slices.Clone(l.attrs)
	clone.attrs = append(clone.attrs, slog.Attr{
		Key:   ErrorFieldKey,
		Value: slog.AnyValue(err),
	})

	return clone
}

func (l *Logger) WithCallerSkip(skip int) log.Logger {
	// Don't have to copy attrs
	clone := l.clone()
	clone.skip += skip
	return clone
}

func (l *Logger) WithContext(ctx context.Context) log.Logger {
	// Don't have to copy attrs
	clone := l.clone()
	clone.ctx = ctx
	return clone
}

func (l *Logger) log(leveler log.Leveler, msg string) {
	level := Level(leveler.Level()).SlogLevel()
	if !l.l.Enabled(l.ctx, level) {
		return
	}

	pc, _, _, _ := runtime.Caller(l.skip + 1)
	record := slog.NewRecord(time.Now(), level, msg, pc)

	record.AddAttrs(l.attrs...)
	_ = l.l.Handler().Handle(l.ctx, record)
}

func (l *Logger) clone() *Logger {
	clone := new(Logger)
	*clone = *l
	return clone
}

func computeFields(dst *[]slog.Attr, fields log.Fields) {
	for k, v := range fields {
		attr := slog.Attr{Key: k}

		switch v := v.(type) {
		case log.Fields:
			attrs := make([]slog.Attr, len(v))
			computeFields(&attrs, v)

			attr.Value = slog.GroupValue(attrs...)
		case string:
			attr.Value = slog.StringValue(v)
		case int, int8, int16, int32:
			attr.Value = slog.IntValue(v.(int))
		case uint, uint8, uint16, uint32:
			attr.Value = slog.Int64Value(v.(int64))
		case int64:
			attr.Value = slog.Int64Value(v)
		case uint64:
			attr.Value = slog.Uint64Value(v)
		case float32, float64:
			attr.Value = slog.Float64Value(v.(float64))
		case bool:
			attr.Value = slog.BoolValue(v)
		case time.Time:
			attr.Value = slog.TimeValue(v)
		case time.Duration:
			attr.Value = slog.DurationValue(v)
		default:
			attr.Value = slog.AnyValue(v)
		}

		*dst = append(*dst, attr)
	}
}
