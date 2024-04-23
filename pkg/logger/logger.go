package logger

import (
	"log/slog"
	"os"

	"groupie-tracker/pkg/config"
)

// Logger structure Базовый кейс который даёт полный путь и укажет строку ошибки

// Init Initializing and customisation the slog logger
func Init(cfg *config.Config) (log *slog.Logger) {
	// handler, output

	defer func() {
		slog.SetDefault(log)
	}()

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.Level(cfg.Logger.Level), // -4 = DEBUG level
		AddSource: cfg.SourceKey,
	}).WithAttrs([]slog.Attr{slog.String("service_name", cfg.ServiceName)}))
}
