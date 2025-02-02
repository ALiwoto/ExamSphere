package logging

const (
	MaxLogTitleLen = 32
)

const (
	// LogTypeDebug is a debug log.
	LogTypeDebug LogType = "debug"

	// LogTypeInfo is an info log.
	LogTypeInfo LogType = "info"

	// LogTypeWarning is a warning log.
	LogTypeWarning LogType = "warning"

	// LogTypeError is an error log.
	LogTypeError LogType = "error"
)
