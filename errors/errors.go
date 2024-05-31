package errors

import "errors"

//easyjson:json
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
	MissingFieldsMsg        = "missing fields"
	InvalidDataMsg          = "invalid data"
	InvalidEmailMsg         = "invalid email"
	InvalidLoginDataMsg     = "invalid login data"
	UnauthorizedMsg         = "unauthorized"
	NotMatchingPasswordsMsg = "password and repeated password are not equal"
	PasswordMinLengthMsg    = "password should contain at least 6 characters"
	EmailsDuplicateMsg      = "user with such email already exists"
	InvalidDateMsg          = "invalid date provided"
	JSONUnmarshallingMsg    = "unable to unmarshal json"
	JSONMarshallingMsg      = "unable to return json reponse"
	InvalidJWTMsg           = "invalid JWT provided"
	NoCookieMsg             = "no cookie provided"
	NoRowsMsg               = "no rows in result set"
	InvalidFilePathGenMsg   = "unable to open file with generated filepath"
	InvalidBodyMsg          = "invalid request body provided"
	ForbiddenMsg            = "forbidden"
	NotFoundMsg             = "not found"
	InternalMsg             = "internal server error"
	InvalidFileNameMsg      = "invalid file name"
	InvalidSlugMsg          = "invalid slug parameters"
	RowsAffectedMsg         = "wrong number of rows affected"
)

var (
	ErrMissingFields        = NewCustomError(errors.New(MissingFieldsMsg))
	ErrInvalidData          = NewCustomError(errors.New(InvalidDataMsg))
	ErrInvalidEmail         = NewCustomError(errors.New(InvalidEmailMsg))
	ErrInvalidLoginData     = NewCustomError(errors.New(InvalidLoginDataMsg))
	ErrUnauthorized         = NewCustomError(errors.New(UnauthorizedMsg))
	ErrNotMatchingPasswords = NewCustomError(errors.New(NotMatchingPasswordsMsg))
	ErrPasswordMinLength    = NewCustomError(errors.New(PasswordMinLengthMsg))
	ErrEmailsDuplicate      = NewCustomError(errors.New(EmailsDuplicateMsg))
	ErrInvalidDate          = NewCustomError(errors.New(InvalidDateMsg))
	ErrJSONUnmarshalling    = NewCustomError(errors.New(JSONUnmarshallingMsg))
	ErrJSONMarshalling      = NewCustomError(errors.New(JSONMarshallingMsg))
	ErrInvalidJWT           = NewCustomError(errors.New(InvalidJWTMsg))
	ErrNoCookie             = NewCustomError(errors.New(NoCookieMsg))
	ErrNoRows               = NewCustomError(errors.New(NoRowsMsg))
	ErrInvalidFilePathGen   = NewCustomError(errors.New(InvalidFilePathGenMsg))
	ErrInvalidBody          = NewCustomError(errors.New(InvalidBodyMsg))
	ErrForbidden            = NewCustomError(errors.New(ForbiddenMsg))
	ErrNotFound             = NewCustomError(errors.New(NotFoundMsg))
	ErrInternal             = NewCustomError(errors.New(InternalMsg))
	ErrInvalidFileName      = NewCustomError(errors.New(InvalidFileNameMsg))
	ErrInvalidSlug          = NewCustomError(errors.New(InvalidSlugMsg))
	ErrRowsAffected         = NewCustomError(errors.New(RowsAffectedMsg))
)
