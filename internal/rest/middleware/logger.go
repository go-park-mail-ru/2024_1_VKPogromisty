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

func NewLogger(logger *zap.SugaredLogger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := uuid.New().String()

		ctx := context.WithValue(r.Context(), requestcontext.RequestIDKey, requestID)

		h.ServeHTTP(w, r.WithContext(ctx))

		l.logger.Info(r.URL.Path,
			zap.String("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("url", r.URL.Path),
			zap.Duration("work_time", time.Since(start)),
		)
	})
}
