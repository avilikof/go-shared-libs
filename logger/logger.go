package logger

import (
	"log/slog"
	"os"
)

func SetLogger() *slog.Logger {
	// Configure JSON logger with custom settings
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}).WithAttrs([]slog.Attr{
		slog.String("service", "location-api"),
		slog.String("version", "1.0"),
	}))

	return logger
}
