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
	level, exists := slogLevels[l]
	if !exists {
		return level
	}

	return slog.Level(l)
}

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

var slogLevels = map[Level]slog.Level{
	DebugLevel: slog.LevelDebug,
	InfoLevel:  slog.LevelInfo,
	WarnLevel:  slog.LevelWarn,
	ErrorLevel: slog.LevelError,
}
