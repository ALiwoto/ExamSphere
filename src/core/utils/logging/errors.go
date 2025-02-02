package logging

import "errors"

var (
	// ErrInvalidLogEntry is returned when a log entry is invalid.
	ErrInvalidLogEntry = errors.New("invalid log entry")

	// ErrLogNotFound is returned when a log is not found.
	ErrLogNotFound = errors.New("log not found")

	ErrLogStorageNotFound = errors.New("log storage not found")
)
