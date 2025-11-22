package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

func NewLogger() *slog.Logger {
	out := os.Getenv("LOG_OUTPUT")
	var w io.Writer = os.Stdout
	if out != "" && out != "stdout" {
		f, err := os.OpenFile(out, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "logging: could not open log file %s: %v\n", out, err)
		} else {
			w = f
		}
	}

	var h slog.Handler

	switch strings.ToLower(os.Getenv("APP_STATUS")) {
	case "dev":
		opts := &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}
		h = slog.NewTextHandler(w, opts)
	case "prod":
		opts := &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelError,
		}
		h = slog.NewJSONHandler(w, opts)
	default:
		h = slog.NewTextHandler(w, nil)
	}

	logger := slog.New(h)

	return logger
}
