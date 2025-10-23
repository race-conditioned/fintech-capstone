package std_log_adapter

import (
	"fintech-capstone/m/v2/internal/platform"
	"fmt"
	"log"
	"os"
	"strings"
)

// StdLogger adapts the standard library log package to the platform.Logger interface.
type StdLogger struct {
	prefix string
}

// New constructs a new StdLogger.
func New() *StdLogger { return &StdLogger{} }

// With enriches the logger with additional fields.
func (l *StdLogger) With(fields ...platform.Field) platform.Logger {
	if len(fields) == 0 {
		return l
	}
	// Simple prefix enrichment
	buf := l.prefix
	for _, field := range fields {
		buf += "[" + field.Key + "=" + fmt.Sprint(field.Value) + "]"
	}
	return &StdLogger{prefix: buf}
}

// log emits a structured log line using the stdlib log package.
// Since stdlog is unstructured, we stringify all fields inline.
func (l *StdLogger) log(level string, msg string, fields ...platform.Field) {
	// Render inline fields (if any)
	if len(fields) > 0 {
		var sb strings.Builder
		for _, f := range fields {
			sb.WriteString(fmt.Sprintf(" %s=%v", f.Key, f.Value))
		}
		log.Printf("%s %s%s%s", level, l.prefix, msg, sb.String())
		return
	}

	log.Printf("%s %s%s", level, l.prefix, msg)
}

// ---------------
// Logging methods
// ---------------

func (l *StdLogger) Debug(msg string, args ...platform.Field) { l.log("[DEBUG]", msg, args...) }
func (l *StdLogger) Info(msg string, args ...platform.Field)  { l.log("[INFO]", msg, args...) }
func (l *StdLogger) Warn(msg string, args ...platform.Field)  { l.log("[WARN]", msg, args...) }
func (l *StdLogger) Error(err error, args ...platform.Field) {
	if err == nil {
		l.log("[ERROR]", "<nil> error", args...)
	} else {
		l.log("[ERROR]", err.Error(), args...)
	}
}

func (l *StdLogger) Fatal(err error, args ...platform.Field) {
	if err == nil {
		l.log("[FATAL]", "<nil> error", args...)
	} else {
		l.log("[FATAL]", err.Error(), args...)
	}
	// Ensure consistent semantics with zap's Fatal.
	os.Exit(1)
}
