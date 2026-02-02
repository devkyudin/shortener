package logger

import (
	"log/slog"
	"os"
)

type Container struct {
	Logger *slog.Logger
}

func NewLoggerContainer() *Container {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return &Container{Logger: logger}
}
