package messages

import "github.com/pkg/errors"

var categoryNotFoundErr = errors.New("category not found")
var invalidCategoryNumberErr = errors.New("invalid category number")
