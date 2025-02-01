package slog

import (
	"github.com/adanyl0v/pocket-ideas/pkg/log"
	"log/slog"
)

type Level int

func (l Level) Level() log.Level {
	return log.Level(l)
}

func (l Level) SlogLevel() slog.Level {
	switch l {
	case DebugLevel:
		return slog.LevelDebug
	case InfoLevel:
		return slog.LevelInfo
	case WarnLevel:
		return slog.LevelWarn
	case ErrorLevel:
		return slog.LevelError
	default:
		return slog.Level(l)
	}
}

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)
