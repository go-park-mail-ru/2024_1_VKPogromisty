package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var GRPCErrors = map[string]codes.Code{
	MissingFieldsMsg:        codes.InvalidArgument,
	InvalidDataMsg:          codes.InvalidArgument,
	InvalidEmailMsg:         codes.InvalidArgument,
	InvalidLoginDataMsg:     codes.Unauthenticated,
	UnauthorizedMsg:         codes.Unauthenticated,
	NotMatchingPasswordsMsg: codes.InvalidArgument,
	PasswordMinLengthMsg:    codes.InvalidArgument,
	EmailsDuplicateMsg:      codes.InvalidArgument,
	InvalidDateMsg:          codes.InvalidArgument,
	JSONUnmarshallingMsg:    codes.InvalidArgument,
	InvalidJWTMsg:           codes.InvalidArgument,
	NoCookieMsg:             codes.Unauthenticated,
	NoRowsMsg:               codes.NotFound,
	InvalidFilePathGenMsg:   codes.InvalidArgument,
	InvalidBodyMsg:          codes.InvalidArgument,
	ForbiddenMsg:            codes.PermissionDenied,
	NotFoundMsg:             codes.NotFound,
	InternalMsg:             codes.Internal,
	InvalidFileNameMsg:      codes.InvalidArgument,
	InvalidSlugMsg:          codes.InvalidArgument,
	RowsAffectedMsg:         codes.InvalidArgument,
	JSONMarshallingMsg:      codes.Internal,
}

var GRPCStatuses = map[codes.Code]int{
	codes.Unauthenticated:  http.StatusUnauthorized,
	codes.InvalidArgument:  http.StatusBadRequest,
	codes.PermissionDenied: http.StatusForbidden,
	codes.NotFound:         http.StatusNotFound,
	codes.Internal:         http.StatusInternalServerError,
}

func (e *CustomError) GRPCStatus() (grpcStatus *status.Status) {
	code, ok := GRPCErrors[e.Error()]
	if !ok {
		code = codes.Internal
	}

	return status.New(code, e.Error())
}

func ParseGRPCError(err error) (msg string, code int) {
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
