package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

var (
	// slogNOPHandler is a discarding slog.Handler
	slogNOPHandler = slog.NewTextHandler(io.Discard, nil)
	// slogDefaultHandler is default slog.Handler
	slogDefaultHandler = NewSlogHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
)

const (
	// LevelTrace is custom trace slog level
	LevelTrace = slog.Level(-8)
	// LevelFatal is custom fatal slog level
	LevelFatal = slog.Level(12)
)

// NewSlogAdapter returns SlogAdapter with specified slog.Handler
func NewSlogAdapter(h slog.Handler) SlogAdapter {
	l := slog.New(h)
	return SlogAdapter{l}
}

// SlogAdapter is wrapper for slog.Logger
type SlogAdapter struct {
	*slog.Logger
}

// Trace logs at trace level
// The attribute arguments are processed as follows:
// If an argument is a string and this is not the last argument, the following argument is treated as the value and the two are combined into an key - value pair.
// Otherwise, the argument is treated as a value with key "!BADKEY".
func (l SlogAdapter) Trace(msg string, args ...any) {
	l.Logger.Log(context.Background(), LevelTrace, msg, args...)
}

// Debug logs at debug level,
// The attribute arguments are processed as follows:
// If an argument is a string and this is not the last argument, the following argument is treated as the value and the two are combined into an key - value pair.
// Otherwise, the argument is treated as a value with key "!BADKEY".
func (l SlogAdapter) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, logArgs(args)...)
}

// Info logs at info level
// The attribute arguments are processed as follows:
// If an argument is a string and this is not the last argument, the following argument is treated as the value and the two are combined into an key - value pair.
// Otherwise, the argument is treated as a value with key "!BADKEY".
func (l SlogAdapter) Info(msg string, args ...any) {
	l.Logger.Info(msg, logArgs(args)...)
}

// Warn logs at warn level
// The attribute arguments are processed as follows:
// If an argument is a string and this is not the last argument, the following argument is treated as the value and the two are combined into an key - value pair.
// Otherwise, the argument is treated as a value with key "!BADKEY".
func (l SlogAdapter) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, logArgs(args)...)
}

// Error logs at error level
// The attribute arguments are processed as follows:
// If an argument is a string and this is not the last argument, the following argument is treated as the value and the two are combined into an key - value pair.
// Otherwise, the argument is treated as a value with key "!BADKEY".
func (l SlogAdapter) Error(msg string, args ...any) {
	l.Logger.Error(msg, logArgs(args)...)
}

// Fatal logs at fatal level, followed by an exit.
// The attribute arguments are processed as follows:
// If an argument is a string and this is not the last argument, the following argument is treated as the value and the two are combined into an key - value pair.
// Otherwise, the argument is treated as a value with key "!BADKEY".
func (l SlogAdapter) Fatal(msg string, args ...any) {
	l.Logger.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}

// Tracef formats according to a format specifier and logs the resulting string at trace level
func (l SlogAdapter) Tracef(msg string, args ...any) {
	m := fmt.Sprintf(msg, args...)
	l.Logger.Log(context.Background(), LevelTrace, m)
}

// Debugf formats according to a format specifier and logs the resulting string at debug level
func (l SlogAdapter) Debugf(msg string, args ...any) {
	m := fmt.Sprintf(msg, args...)
	l.Logger.Debug(m)
}

// Infof formats according to a format specifier and logs the resulting string at info level
func (l SlogAdapter) Infof(msg string, args ...any) {
	m := fmt.Sprintf(msg, args...)
	l.Logger.Info(m)
}

// Warnf formats according to a format specifier and logs the resulting string at warn level
func (l SlogAdapter) Warnf(msg string, args ...any) {
	m := fmt.Sprintf(msg, args...)
	l.Logger.Warn(m)
}

// Errorf formats according to a format specifier and logs the resulting string at error level
func (l SlogAdapter) Errorf(msg string, args ...any) {
	m := fmt.Sprintf(msg, args...)
	l.Logger.Error(m)
}

// Fatalf formats according to a format specifier and logs the resulting string at fatal level, followed by an exit.
func (l SlogAdapter) Fatalf(msg string, args ...any) {
	m := fmt.Sprintf(msg, args...)
	l.Logger.Log(context.Background(), LevelFatal, m)
	os.Exit(1)
}

// With returns new Interface that includes given fields with each operation
// returned logger will log each field on every logging operation it will be used
func (l SlogAdapter) With(key string, fields Fields) Interface {
	newLogger := l.Logger.With(key, fields)
	return SlogAdapter{newLogger}
}

func logArgs(args []any) []any {
	result := []any{}
	for _, arg := range args {
		switch v := arg.(type) {
		case Fields:
			result = append(result, v.Fields().Get()...)
		default:
			result = append(result, arg)
		}
	}
	return result
}
