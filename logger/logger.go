package logger

import (
	"context"
	"log/slog"
	"os"
	"services-api/constants"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(logFile string) (*Logger, error) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(handler)

	return &Logger{logger}, nil
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	if requestID, ok := ctx.Value(constants.RequestId).(string); ok {
		return &Logger{l.Logger.With("request_id", requestID)}
	}
	return l
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.WithContext(ctx).Logger.Info(msg, args...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.WithContext(ctx).Logger.Error(msg, args...)
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.WithContext(ctx).Logger.Debug(msg, args...)
}
