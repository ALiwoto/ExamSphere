package database

import "errors"

var (
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrInternalDatabaseError = errors.New("internal database error")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrOperationNotAllowed   = errors.New("operation not allowed")
)
