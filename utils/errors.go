package utils

import "errors"

var (
	ErrMissingFields        = errors.New("provided data missing required fields")
	ErrInvalidData          = errors.New("invalid data")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrInvalidLoginData     = errors.New("invalid login data")
	ErrNotMatchingPasswords = errors.New("password and repeated password are not equal")
	ErrPasswordMinLength    = errors.New("password should contain at least 6 characters")
	ErrEmailsDuplicate      = errors.New("uses with such email already exists")
	ErrInvalidBirthDate     = errors.New("invalid date of birth")
	ErrNotFound             = errors.New("not found")
	ErrBadRequest           = errors.New("bad request")
	ErrEternal              = errors.New("eternal server error")
)
