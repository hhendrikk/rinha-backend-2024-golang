package validation

import "errors"

var (
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("not found")
)
