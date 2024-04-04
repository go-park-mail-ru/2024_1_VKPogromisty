package errors

import (
	"encoding/json"
	"net/http"
)

var HTTPErrors = map[error]int{
	ErrUnauthorized:         http.StatusUnauthorized,
	ErrInvalidLoginData:     http.StatusUnauthorized,
	http.ErrNoCookie:        http.StatusUnauthorized,
	ErrMissingFields:        http.StatusBadRequest,
	ErrInvalidData:          http.StatusBadRequest,
	ErrInvalidEmail:         http.StatusBadRequest,
	ErrInvalidSlug:          http.StatusBadRequest,
	ErrNotMatchingPasswords: http.StatusBadRequest,
	ErrPasswordMinLength:    http.StatusBadRequest,
	ErrEmailsDuplicate:      http.StatusBadRequest,
	ErrInvalidDate:          http.StatusBadRequest,
	ErrJSONUnmarshalling:    http.StatusBadRequest,
	ErrInvalidFilePathGen:   http.StatusBadRequest,
	ErrInvalidFileName:      http.StatusBadRequest,
	ErrInvalidBody:          http.StatusBadRequest,
	ErrRowsAffected:         http.StatusBadRequest,
	ErrForbidden:            http.StatusForbidden,
	ErrNotFound:             http.StatusNotFound,
	ErrJSONMarshalling:      http.StatusInternalServerError,
	ErrInternal:             http.StatusInternalServerError,
}

func ParseHTTPError(err error) (msg string, status int) {
	if err == nil {
		err = ErrInternal
	}

	status, ok := HTTPErrors[err]
	if !ok {
		status = 500
		err = ErrInternal
	}

	msg = err.Error()

	return
}

func MarshalError(err error) (data []byte) {
	data, _ = json.Marshal(map[string]string{"error": err.Error()})
	return
}
