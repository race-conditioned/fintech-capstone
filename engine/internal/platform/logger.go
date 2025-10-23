package platform

// Logger defines the logging interface.
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(err error, fields ...Field)
	With(fields ...Field) Logger
	Fatal(err error, fields ...Field)
}

// Field represents a key-value pair for structured logging.
type Field struct {
	Key   string
	Value any
}
