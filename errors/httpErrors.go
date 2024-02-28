package errors

import (
	"fmt"
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

func ServeHttpError(w *http.ResponseWriter, err error) {
	status, ok := HTTPErrors[err]
	if !ok {
		status = 500
	}

	fmt.Printf("handled error: %s\n", err.Error())

	http.Error(*w, err.Error(), status)
}
