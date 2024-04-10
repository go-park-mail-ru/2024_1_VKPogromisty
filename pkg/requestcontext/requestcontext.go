package requestcontext

import (
	"context"
	"socio/errors"

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

func GetLogger(ctx context.Context) (logger *zap.SugaredLogger, err error) {
	logger, ok := ctx.Value(LoggerKey).(*zap.SugaredLogger)
	if !ok {
		err = errors.ErrInvalidData
		return
	}

	return
}
