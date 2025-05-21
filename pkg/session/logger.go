package session

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/log"
	"github.com/hashicorp/go-retryablehttp"
)

// GetRetryableLogger returns wrapper retryablehttp.LeveledLogger for log.Interface
func GetRetryableLogger(log log.Interface) retryablehttp.LeveledLogger {
	return &retryableLogger{log: log}
}

// retryableLogger is wrapper for log.Interface to expose retryablehttp.LeveledLogger needed for retry logic
type retryableLogger struct {
	log log.Interface
}

// Error level formatted message.
func (l *retryableLogger) Error(msg string, keysAndValues ...interface{}) {
	kv := append([]interface{}{}, keysAndValues...)
	switch msg {
	case "request failed":
		l.log.Errorf("%s %s request failed: %v", kv[1], kv[3], kv[5])
	case "error reading response body":
		l.log.Errorf("error reading response body: %v", kv[1])
	default:
		l.log.Errorf(msg, keysAndValues...)
	}
}

// Info level formatted message.
func (l *retryableLogger) Info(msg string, keysAndValues ...interface{}) {
	l.log.Infof(msg, keysAndValues...)
}

// Debug level formatted message.
func (l *retryableLogger) Debug(msg string, keysAndValues ...interface{}) {
	kv := append([]interface{}{}, keysAndValues...)
	switch msg {
	case "retrying request":
		l.log.Debugf("%s: retrying in %s (%d left)", kv[1], kv[3], kv[5])
	case "performing request":
		l.log.Debugf("%s %s", kv[1], kv[3])
	default:
		l.log.Debugf(msg, keysAndValues...)
	}
}

// Warn level formatted message.
func (l *retryableLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.log.Warnf(msg, keysAndValues...)
}
