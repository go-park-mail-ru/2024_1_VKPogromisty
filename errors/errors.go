package errors

import "errors"

var (
	ErrMissingFields        = errors.New("provided data missing required fields")
	ErrInvalidData          = errors.New("invalid data")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrInvalidLoginData     = errors.New("invalid login data")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrNotMatchingPasswords = errors.New("password and repeated password are not equal")
	ErrPasswordMinLength    = errors.New("password should contain at least 6 characters")
	ErrEmailsDuplicate      = errors.New("uses with such email already exists")
	ErrInvalidDate          = errors.New("invalid date provided")
	ErrJSONUnmarshalling    = errors.New("unable to unmarshal json")
	ErrJSONMarshalling      = errors.New("unable to return json reponse")
	ErrInvalidFilePathGen   = errors.New("unable to open file with generated filepath")
	ErrInvalidBody          = errors.New("invalid request body provided")
	ErrInternal             = errors.New("internal server error")
	ErrInvalidFileName      = errors.New("invalid file name")
)
