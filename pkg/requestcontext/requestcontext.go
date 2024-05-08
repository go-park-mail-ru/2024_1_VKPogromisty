package requestcontext

import (
	"context"
	"socio/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ContextKey string

const (
	UserIDKey    ContextKey = "userID"
	SessionIDKey ContextKey = "sessionID"
	RequestIDKey ContextKey = "requestID"
	LoggerKey    ContextKey = "logger"
)

func GetUserID(ctx context.Context) (userID uint, err error) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	if !ok {
		err = errors.ErrInvalidData
		return
	}

	return
}

func GetSessionID(ctx context.Context) (sessionID string, err error) {
	sessionID, ok := ctx.Value(SessionIDKey).(string)
	if !ok {
		err = errors.ErrInvalidData
		return
	}

	return
}

func GetLogger(ctx context.Context) (logger *zap.Logger, err error) {
	logger, ok := ctx.Value(LoggerKey).(*zap.Logger)
	if !ok {
		err = errors.ErrInvalidData
		return
	}

	return
}

func GetRequestID(ctx context.Context) (requestID string, err error) {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		requestID = uuid.NewString()
		return
	}

	return
}
