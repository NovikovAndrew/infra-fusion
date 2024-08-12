package logger

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const (
	LevelDebug string = "DEBUG"
	LevelInfo  string = "INFO"
	LevelWarn  string = "WARN"
	LevelError string = "ERROR"
)

var l = logger{
	opts: &slog.HandlerOptions{
		Level: slog.LevelDebug,
	},
}

type logger struct {
	opts *slog.HandlerOptions
	slog *slog.Logger
}

func (l *logger) Debug(ctx context.Context, msg string, args ...any) {
	l.slog.DebugContext(ctx, msg, args...)
}

func (l *logger) Info(ctx context.Context, msg string, args ...any) {
	l.slog.InfoContext(ctx, msg, args...)
}

func (l *logger) Warn(ctx context.Context, msg string, args ...any) {
	l.slog.WarnContext(ctx, msg, args...)
}

func (l *logger) Error(ctx context.Context, msg string, args ...any) {
	l.slog.ErrorContext(ctx, msg, args...)
}

func (l *logger) GetSlogLogger() *slog.Logger {
	return l.slog
}

func (l *logger) GetLogLogger(level string, prefix string) *log.Logger {
	slogLevel := slog.LevelDebug
	switch level {
	case LevelDebug:
		slogLevel = slog.LevelDebug
	case LevelInfo:
		slogLevel = slog.LevelInfo
	case LevelWarn:
		slogLevel = slog.LevelWarn
	case LevelError:
		slogLevel = slog.LevelError
	}

	return log.New(&handlerWriter{l.slog.Handler(), slogLevel, true}, prefix, 0)
}

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)

	GetSlogLogger() *slog.Logger
	GetLogLogger(level string, prefix string) *log.Logger
}

func InitLogger(_ context.Context, logLevel string) Logger {
	switch logLevel {
	case LevelDebug:
		l.opts.Level = slog.LevelDebug
	case LevelInfo:
		l.opts.Level = slog.LevelInfo
	case LevelWarn:
		l.opts.Level = slog.LevelWarn
	case LevelError:
		l.opts.Level = slog.LevelError
	}

	l.slog = slog.New(slog.NewJSONHandler(os.Stdout, l.opts))
	return &l
}

func InterceptorLogger(l Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			l.Debug(ctx, msg, fields...)
		case logging.LevelInfo:
			l.Info(ctx, msg, fields...)
		case logging.LevelWarn:
			l.Warn(ctx, msg, fields...)
		case logging.LevelError:
			l.Error(ctx, msg, fields...)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
