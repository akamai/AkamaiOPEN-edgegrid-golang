// Package log provides logging interface
package log

import (
	"sync"
)

var (
	protect sync.Once
	log     *logControler
)

// Interface is the logging interface
type Interface interface {
	// Trace logs at trace level
	// The attribute arguments are processed as follows:
	// If an argument is a string and this is not the last argument, the following argument is treated as the value and the two are combined into an key - value pair.
	// Otherwise, the argument is treated as a value with key "!BADKEY".
	Trace(msg string, args ...any)
	// Debug logs at debug level,
	// The attribute arguments are processed as follows:
	// If an argument is a string and this is not the last argument, the following argument is treated as the value and the two are combined into an key - value pair.
	// Otherwise, the argument is treated as a value with key "!BADKEY".
	Debug(msg string, args ...any)
	// Info logs at info level
	// The attribute arguments are processed as follows:
	// If an argument is a string and this is not the last argument, the following argument is treated as the value and the two are combined into an key - value pair.
	// Otherwise, the argument is treated as a value with key "!BADKEY".
	Info(msg string, args ...any)
	// Warn logs at warn level
	// The attribute arguments are processed as follows:
	// If an argument is a string and this is not the last argument, the following argument is treated as the value and the two are combined into an key - value pair.
	// Otherwise, the argument is treated as a value with key "!BADKEY".
	Warn(msg string, args ...any)
	// Error logs at error level
	// The attribute arguments are processed as follows:
	// If an argument is a string and this is not the last argument, the following argument is treated as the value and the two are combined into an key - value pair.
	// Otherwise, the argument is treated as a value with key "!BADKEY".
	Error(msg string, args ...any)
	// Fatal logs at fatal level, followed by an exit.
	// The attribute arguments are processed as follows:
	// If an argument is a string and this is not the last argument, the following argument is treated as the value and the two are combined into an key - value pair.
	// Otherwise, the argument is treated as a value with key "!BADKEY".
	Fatal(msg string, args ...any)

	// Tracef formats according to a format specifier and logs the resulting string at trace level
	Tracef(msg string, args ...any)
	// Debugf formats according to a format specifier and logs the resulting string at debug level
	Debugf(msg string, args ...any)
	// Infof formats according to a format specifier and logs the resulting string at info level
	Infof(msg string, args ...any)
	// Warnf formats according to a format specifier and logs the resulting string at warn level
	Warnf(msg string, args ...any)
	// Errorf formats according to a format specifier and logs the resulting string at error level
	Errorf(msg string, args ...any)
	// Fatalf formats according to a format specifier and logs the resulting string at fatal level, followed by an exit.
	Fatalf(msg string, args ...any)

	// With returns new Interface that includes given fields with each operation
	// returned logger will log each field on every logging operation it will be used
	With(key string, fields Fields) Interface
}

// logControler is a wrapper for Interface
type logControler struct {
	log Interface
}

// SetLogger sets default logger interface
func SetLogger(l Interface) {
	log = &logControler{log: l}
}

// Default returns logger interface
func Default() Interface {
	protect.Do(func() {
		if log == nil {
			SetLogger(NewSlogAdapter(slogDefaultHandler))
		}
	})
	return log.log
}

// NOPLogger returns discarding logger
func NOPLogger() Interface {
	return NewSlogAdapter(slogNOPHandler)
}
