package logger

import (
	"log/slog"
	"os"
	"strings"
)

func New(levelStr, formatStr string) *slog.Logger {
	level := parseLevel(levelStr)
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch strings.ToLower(formatStr) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

func parseLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
