package contextlogger

import (
	"context"
	"socio/pkg/requestcontext"

	"go.uber.org/zap"
)

func LogInfo(ctx context.Context) {
	logger, ctxErr := requestcontext.GetLogger(ctx)
	if ctxErr != nil {
		return
	}

	logger.Info()
}

func LogErr(ctx context.Context, err error) {
	logger, ctxErr := requestcontext.GetLogger(ctx)
	if ctxErr != nil {
		err = ctxErr
	}

	logger.Error(zap.Error(err))
}

func LogSQL(ctx context.Context, query string, args ...interface{}) {
	logger, ctxErr := requestcontext.GetLogger(ctx)
	if ctxErr != nil {
		return
	}

	logger.Info(
		"SQL query: ",
		zap.String("query", query),
		zap.Any("args", args),
	)
}

func LogRedisAction(ctx context.Context, action string, key interface{}, value interface{}) {
	logger, ctxErr := requestcontext.GetLogger(ctx)
	if ctxErr != nil {
		return
	}

	logger.Info(
		"Redis query: ",
		zap.String("action", action),
		zap.Any("key", key),
		zap.Any("value", value),
	)
}
