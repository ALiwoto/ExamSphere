package emailUtils

import "errors"

var (
	ErrInvalidEmailConfig   = errors.New("invalid email configuration")
	ErrEmailClientNotLoaded = errors.New("email client not loaded")
)
