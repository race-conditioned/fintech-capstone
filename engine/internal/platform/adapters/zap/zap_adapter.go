package zap_adapter

import (
	"fintech-capstone/m/v2/internal/platform"

	"go.uber.org/zap"
)

// ZapLogger adapts a zap.Logger to implement the platform.Logger interface.
type ZapLogger struct {
	z *zap.Logger
}

// New creates a new ZapLogger wrapping the provided zap.Logger.
func New(z *zap.Logger) *ZapLogger { return &ZapLogger{z} }

// With adds structured fields to the logger.
func (l *ZapLogger) With(fields ...platform.Field) platform.Logger {
	zFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zFields[i] = zap.Any(f.Key, f.Value)
	}
	return &ZapLogger{l.z.With(zFields...)}
}

// toZapFields converts platform.Fields to zap.Fields.
func toZapFields(fields []platform.Field) []zap.Field {
	zFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zFields[i] = zap.Any(f.Key, f.Value)
	}
	return zFields
}

// ---------------
// Logging methods
// ---------------

func (l *ZapLogger) Debug(msg string, fields ...platform.Field) {
	l.z.Debug(msg, toZapFields(fields)...)
}
func (l *ZapLogger) Info(msg string, fields ...platform.Field) { l.z.Info(msg, toZapFields(fields)...) }
func (l *ZapLogger) Warn(msg string, fields ...platform.Field) { l.z.Warn(msg, toZapFields(fields)...) }
func (l *ZapLogger) Error(err error, fields ...platform.Field) {
	if err == nil {
		l.z.Error("<nil> error:", toZapFields(fields)...)
	} else {
		l.z.Error(err.Error(), toZapFields(fields)...)
	}
}

func (l *ZapLogger) Fatal(err error, fields ...platform.Field) {
	if err == nil {
		l.z.Fatal("<nil> error", toZapFields(fields)...)
	} else {
		l.z.Fatal(err.Error(), toZapFields(fields)...)
	}
}
