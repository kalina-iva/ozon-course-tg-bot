package messages

import "github.com/pkg/errors"

var errInvalidAmount = errors.New("amount cannot be negative or 0")
