package messages

import "github.com/pkg/errors"

var (
	errInvalidAmount = errors.New("amount cannot be negative or 0")
	errLimitExceeded = errors.New("limit exceeded")
	errUserNotFound  = errors.New("user not found")
)
