package userHandlers

import "errors"

var (
	ErrTooManyPasswordChangeAttempts = errors.New("too many password change attempts")
	ErrAccountAlreadyConfirmed       = errors.New("account already confirmed")
)
