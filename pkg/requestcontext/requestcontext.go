package requestcontext

import (
	"net/http"
	"socio/errors"
)

type ContextKey string

const (
	UserIDKey    ContextKey = "userID"
	SessionIDKey ContextKey = "sessionID"
	RequestIDKey ContextKey = "requestID"
)

func GetUserID(r *http.Request) (userID uint, err error) {
	userID, ok := r.Context().Value(UserIDKey).(uint)
	if !ok {
		err = errors.ErrInvalidData
		return
	}

	return
}

func GetSessionID(r *http.Request) (sessionID string, err error) {
	sessionID, ok := r.Context().Value(SessionIDKey).(string)
	if !ok {
		err = errors.ErrInvalidData
		return
	}

	return
}
