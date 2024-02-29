package errors

import (
	"net/http"
)

var HTTPErrors = map[error]int{
	ErrInvalidLoginData:     http.StatusUnauthorized,
	http.ErrNoCookie:        http.StatusUnauthorized,
	ErrMissingFields:        http.StatusBadRequest,
	ErrInvalidData:          http.StatusBadRequest,
	ErrInvalidEmail:         http.StatusBadRequest,
	ErrNotMatchingPasswords: http.StatusBadRequest,
	ErrPasswordMinLength:    http.StatusBadRequest,
	ErrEmailsDuplicate:      http.StatusBadRequest,
	ErrInvalidDate:          http.StatusBadRequest,
	ErrJSONUnmarshalling:    http.StatusBadRequest,
	ErrInvalidFilePathGen:   http.StatusBadRequest,
	ErrInvalidFileName:      http.StatusBadRequest,
	ErrInvalidBody:          http.StatusBadRequest,
	ErrInternal:             http.StatusInternalServerError,
}

func ServeHttpError(err error) (msg string, status int) {
	status, ok := HTTPErrors[err]
	if !ok {
		status = 500
	}

	msg = err.Error()

	return
}
