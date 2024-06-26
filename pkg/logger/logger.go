package logger

import (
	"context"
	"net/http"
	"socio/pkg/requestcontext"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	DefaultOutputPaths = []string{
		"/var/log/socio/socio.log",
		"stderr",
	}
)

type Logger struct {
	logger *zap.Logger
}

func NewZapLogger(outputPaths []string) (logger *zap.Logger, err error) {
	if outputPaths == nil {
		outputPaths = DefaultOutputPaths
	}

	cfg := zap.NewProductionConfig()

	cfg.Sampling = nil

	cfg.OutputPaths = outputPaths

	logger, err = cfg.Build()
	if err != nil {
		return
	}

	return
}

func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := uuid.New().String()

		currLogger := l.logger.With(
			zap.String("requestID", requestID),
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("url", r.URL.Path),
		)

		ctx := context.WithValue(r.Context(), requestcontext.LoggerKey, currLogger)

		h.ServeHTTP(w, r.WithContext(ctx))

		currLogger.Info(
			"Working time: ",
			zap.Duration("work_time", time.Since(start)),
		)
	})
}

func (l *Logger) UnaryLoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	requestID, err := requestcontext.GetRequestID(ctx)
	if err != nil {
		return nil, err
	}

	currLogger := l.logger.With(
		zap.String("requestID", requestID),
		zap.String("method", info.FullMethod),
	)

	newCtx := context.WithValue(ctx, requestcontext.LoggerKey, currLogger)

	resp, err = handler(newCtx, req)
	if err != nil {
		currLogger.Error(err.Error())
	}

	currLogger.Info(
		"Working time: ",
		zap.Duration("work_time", time.Since(start)),
	)

	return resp, err
}
