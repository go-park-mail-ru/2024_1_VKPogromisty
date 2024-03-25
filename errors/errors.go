package errors

import "errors"

type HTTPError struct {
	Error string `json:"error"`
}

var (
	ErrMissingFields        = errors.New("provided data missing required fields")
	ErrInvalidData          = errors.New("invalid data")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrInvalidLoginData     = errors.New("invalid login data")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrNotMatchingPasswords = errors.New("password and repeated password are not equal")
	ErrPasswordMinLength    = errors.New("password should contain at least 6 characters")
	ErrEmailsDuplicate      = errors.New("user with such email already exists")
	ErrInvalidDate          = errors.New("invalid date provided")
	ErrJSONUnmarshalling    = errors.New("unable to unmarshal json")
	ErrJSONMarshalling      = errors.New("unable to return json reponse")
	ErrInvalidFilePathGen   = errors.New("unable to open file with generated filepath")
	ErrInvalidBody          = errors.New("invalid request body provided")
	ErrForbidden            = errors.New("forbidden")
	ErrNotFound             = errors.New("not found")
	ErrInternal             = errors.New("internal server error")
	ErrInvalidFileName      = errors.New("invalid file name")
	ErrInvalidSlug          = errors.New("invalid slug parameters")
	ErrRowsAffected         = errors.New("wrong number of rows affected")
)
