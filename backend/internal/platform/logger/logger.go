package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// New initializes and returns a new JSON slog.Logger.
func New(levelStr string, w io.Writer) *slog.Logger {
	if w == nil {
		w = os.Stdout
	}

	var level slog.Level
	switch strings.ToLower(strings.TrimSpace(levelStr)) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewJSONHandler(w, opts)
	return slog.New(handler)
}
