package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Level  string // debug, info, warn, error
	Format string // json, text
}

func NewLogger(cfg Config, w io.Writer) *slog.Logger {
	if w == nil {
		w = os.Stdout
	}

	var level slog.Level
	switch strings.ToLower(cfg.Level) {

	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handlerOptions := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if strings.ToLower(cfg.Format) == "text" {
		handler = slog.NewTextHandler(w, handlerOptions)
	} else {
		handler = slog.NewJSONHandler(w, handlerOptions)
	}

	return slog.New(handler)
}
