package middleware

import (
	"context"
	"net/http"
	"socio/pkg/requestcontext"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger() (sugar *zap.SugaredLogger) {
	logger, _ := zap.NewProduction()

	sugar = logger.Sugar()

	return
}

func NewLogger(logger *zap.SugaredLogger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := uuid.New().String()

		currLogger := l.logger.With(
			"requestID", requestID,
			zap.Duration("work_time", time.Since(start)),
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("url", r.URL.Path),
		)

		ctx := context.WithValue(r.Context(), requestcontext.LoggerKey, currLogger)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
