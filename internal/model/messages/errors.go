package messages

import "github.com/pkg/errors"

var (
	errCategoryNotFound      = errors.New("category not found")
	errInvalidCategoryNumber = errors.New("invalid category number")
)
