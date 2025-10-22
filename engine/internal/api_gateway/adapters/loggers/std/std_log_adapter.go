package std_log_adapter

import (
	"fintech-capstone/m/v2/internal/api_gateway/ports"
	"fmt"
	"log"
	"os"
	"strings"
)

type StdLogger struct {
	prefix string
}

func New() *StdLogger { return &StdLogger{} }

func (l *StdLogger) With(fields ...ports.Field) ports.Logger {
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
func (l *StdLogger) log(level string, msg string, fields ...ports.Field) {
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

func (l *StdLogger) Debug(msg string, args ...ports.Field) { l.log("[DEBUG]", msg, args...) }
func (l *StdLogger) Info(msg string, args ...ports.Field)  { l.log("[INFO]", msg, args...) }
func (l *StdLogger) Warn(msg string, args ...ports.Field)  { l.log("[WARN]", msg, args...) }
func (l *StdLogger) Error(err error, args ...ports.Field) {
	if err == nil {
		l.log("[ERROR]", "<nil> error", args...)
	} else {
		l.log("[ERROR]", err.Error(), args...)
	}
}

func (l *StdLogger) Fatal(err error, args ...ports.Field) {
	if err == nil {
		l.log("[FATAL]", "<nil> error", args...)
	} else {
		l.log("[FATAL]", err.Error(), args...)
	}
	// Ensure consistent semantics with zap's Fatal.
	os.Exit(1)
}
