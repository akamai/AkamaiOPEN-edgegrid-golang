package dns

import "errors"

var (
	// ErrBadRequest is returned when a required parameter is missing
	ErrBadRequest = errors.New("missing argument")
)
