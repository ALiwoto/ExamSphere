package userHandlers

import "time"

const (
	MinPasswordLength = 4
	MaxPasswordLength = 24
)

const (
	MaxPasswordRequestAttempts = 10
	MinPasswordAttemptWaitTime = 2 * time.Minute
)

const (
	reqFirst = "req_"
)
