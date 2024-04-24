package errors

import "errors"

type HTTPError struct {
	Error string `json:"error"`
}

type CustomError struct {
	error
}

func NewCustomError(err error) CustomError {
	if err == nil {
		return CustomError{errors.New("")}
	}

	return CustomError{err}
}

var (
	ErrMissingFields        = CustomError{errors.New("provided data missing required fields")}
	ErrInvalidData          = CustomError{errors.New("invalid data")}
	ErrInvalidEmail         = CustomError{errors.New("invalid email")}
	ErrInvalidLoginData     = CustomError{errors.New("invalid login data")}
	ErrUnauthorized         = CustomError{errors.New("unauthorized")}
	ErrNotMatchingPasswords = CustomError{errors.New("password and repeated password are not equal")}
	ErrPasswordMinLength    = CustomError{errors.New("password should contain at least 6 characters")}
	ErrEmailsDuplicate      = CustomError{errors.New("user with such email already exists")}
	ErrInvalidDate          = CustomError{errors.New("invalid date provided")}
	ErrJSONUnmarshalling    = CustomError{errors.New("unable to unmarshal json")}
	ErrJSONMarshalling      = CustomError{errors.New("unable to return json reponse")}
	ErrInvalidJWT           = CustomError{errors.New("invalid JWT provided")}
	ErrNoCookie             = CustomError{errors.New("no cookie provided")}
	ErrNoRows               = CustomError{errors.New("no rows in result set")}
	ErrInvalidFilePathGen   = CustomError{errors.New("unable to open file with generated filepath")}
	ErrInvalidBody          = CustomError{errors.New("invalid request body provided")}
	ErrForbidden            = CustomError{errors.New("forbidden")}
	ErrNotFound             = CustomError{errors.New("not found")}
	ErrInternal             = CustomError{errors.New("internal server error")}
	ErrInvalidFileName      = CustomError{errors.New("invalid file name")}
	ErrInvalidSlug          = CustomError{errors.New("invalid slug parameters")}
	ErrRowsAffected         = CustomError{errors.New("wrong number of rows affected")}
)
