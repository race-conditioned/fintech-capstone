package zap_adapter

import (
	"fintech-capstone/m/v2/internal/api_gateway/ports"

	"go.uber.org/zap"
)

type ZapLogger struct {
	z *zap.Logger
}

func New(z *zap.Logger) *ZapLogger { return &ZapLogger{z} }

func (l *ZapLogger) With(fields ...ports.Field) ports.Logger {
	zFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zFields[i] = zap.Any(f.Key, f.Value)
	}
	return &ZapLogger{l.z.With(zFields...)}
}

func (l *ZapLogger) Debug(msg string, fields ...ports.Field) { l.z.Debug(msg, toZapFields(fields)...) }
func (l *ZapLogger) Info(msg string, fields ...ports.Field)  { l.z.Info(msg, toZapFields(fields)...) }
func (l *ZapLogger) Warn(msg string, fields ...ports.Field)  { l.z.Warn(msg, toZapFields(fields)...) }
func (l *ZapLogger) Error(err error, fields ...ports.Field) {
	if err == nil {
		l.z.Error("<nil> error:", toZapFields(fields)...)
	} else {
		l.z.Error(err.Error(), toZapFields(fields)...)
	}
}

func (l *ZapLogger) Fatal(err error, fields ...ports.Field) {
	if err == nil {
		l.z.Fatal("<nil> error", toZapFields(fields)...)
	} else {
		l.z.Fatal(err.Error(), toZapFields(fields)...)
	}
}

func toZapFields(fields []ports.Field) []zap.Field {
	zFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zFields[i] = zap.Any(f.Key, f.Value)
	}
	return zFields
}
