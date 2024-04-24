package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var GRPCErrors = map[CustomError]codes.Code{
	ErrUnauthorized:         codes.Unauthenticated,
	ErrInvalidLoginData:     codes.Unauthenticated,
	ErrMissingFields:        codes.InvalidArgument,
	ErrInvalidData:          codes.InvalidArgument,
	ErrInvalidEmail:         codes.InvalidArgument,
	ErrInvalidSlug:          codes.InvalidArgument,
	ErrInvalidJWT:           codes.InvalidArgument,
	ErrNotMatchingPasswords: codes.InvalidArgument,
	ErrPasswordMinLength:    codes.InvalidArgument,
	ErrEmailsDuplicate:      codes.InvalidArgument,
	ErrInvalidDate:          codes.InvalidArgument,
	ErrJSONUnmarshalling:    codes.InvalidArgument,
	ErrInvalidFilePathGen:   codes.InvalidArgument,
	ErrInvalidFileName:      codes.InvalidArgument,
	ErrInvalidBody:          codes.InvalidArgument,
	ErrRowsAffected:         codes.InvalidArgument,
	ErrForbidden:            codes.PermissionDenied,
	ErrNotFound:             codes.NotFound,
	ErrNoRows:               codes.NotFound,
	ErrJSONMarshalling:      codes.Internal,
	ErrInternal:             codes.Internal,
}

var GRPCStatuses = map[codes.Code]int{
	codes.Unauthenticated:  http.StatusUnauthorized,
	codes.InvalidArgument:  http.StatusBadRequest,
	codes.PermissionDenied: http.StatusForbidden,
	codes.NotFound:         http.StatusNotFound,
	codes.Internal:         http.StatusInternalServerError,
}

func (e CustomError) GRPCStatus() (grpcStatus *status.Status) {
	code, ok := GRPCErrors[e]
	if !ok {
		code = codes.Internal
	}

	return status.New(code, e.Error())
}

func ParseGRPCError(err CustomError) (msg string, code int) {
	if err.Error() != "" {
		st, ok := status.FromError(err)
		if ok {
			msg = st.Message()
			code = GRPCStatuses[st.Code()]
		} else {
			msg = err.Error()
			code = http.StatusInternalServerError
		}
	}

	return
}
