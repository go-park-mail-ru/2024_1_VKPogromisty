package utils

import "errors"

var (
	ErrMissingFields    = errors.New("provided data missing required fields")
	ErrInvalidData      = errors.New("invalid data")
	ErrInvalidLoginData = errors.New("invalid login data")
)
