// Package logger is a package to wrap a log driver
package logger

import (
	"io"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Close()
	Info(msg string)
	Warn(msg string)
	Debug(msg string)
	Error(err error)
	WithField(key, value string) Logger
}

type logger struct {
	logger *zap.Logger
}

// New create a instance of Logger
func New(writer io.Writer, level Level) (Logger, error) {
	if writer == nil {
		return nil, ErrNilWriter
	}
	cfg := zap.NewProductionConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(writer),
		zapcore.Level(level),
	)

	log := &logger{
		logger: zap.New(core),
	}

	return log, nil
}

// Close method call Sync (kind of flush buffer).
func (l *logger) Close() {
	if err := l.logger.Sync(); err != nil {
		if strings.Contains(err.Error(), "bad file descriptor") || strings.Contains(err.Error(), "invalid argument") {
			// ignore because the stderr should not sync.
			return
		}
		l.logger.Error(err.Error())
	}

}

// Info level message.
func (l *logger) Info(msg string) {
	l.logger.Info(msg)
}

// Warn level message.
func (l *logger) Warn(msg string) {
	l.logger.Warn(msg)
}

// Debug level message.
func (l *logger) Debug(msg string) {
	l.logger.Debug(msg)
}

// Error level message.
func (l *logger) Error(err error) {
	l.logger.Error(err.Error())
}

// WithField to add specific keys - values
// Example usage: logger.WithField("key-1", "value-1").With.Field("key-2", "value-2").Info("message with context").
func (l *logger) WithField(key, value string) Logger {
	return &logger{logger: l.logger.With(zap.Any(key, value))}
}
