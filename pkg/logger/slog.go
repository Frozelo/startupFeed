package logger

import (
	slogmulti "github.com/samber/slog-multi"
	"log/slog"
	"os"
)

type Slog struct {
	logger *slog.Logger
}

func New(lever slog.Level, filePath string) (*Slog, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	fileHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(
		slogmulti.Fanout(fileHandler, consoleHandler))
	return &Slog{logger: logger}, nil
}

func (l *Slog) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *Slog) Error(msg string, keysAndValues ...any) {
	l.logger.Error(msg, keysAndValues...)
}

func (l *Slog) Debug(msg string, keysAndValues ...any) {
	l.logger.Debug(msg, keysAndValues...)
}

func (l *Slog) With(args ...any) *Slog {
	return &Slog{
		logger: l.logger.With(args...),
	}
}
