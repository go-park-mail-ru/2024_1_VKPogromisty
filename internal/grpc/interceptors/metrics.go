package interceptors

import (
	"context"
	"socio/pkg/appmetrics"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func AuthHitMetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Call the next handler
	h, err := handler(ctx, req)

	// Calculate the duration and record it in the histogram
	statusCode := status.Code(err).String()
	duration := time.Since(start)
	appmetrics.AuthHitDuration.WithLabelValues(info.FullMethod).Set(float64(duration.Milliseconds()))

	appmetrics.AuthHits.WithLabelValues(info.FullMethod, statusCode).Inc()

	appmetrics.AuthTotalHits.WithLabelValues().Inc()

	return h, err
}

func PostHitMetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Call the next handler
	h, err := handler(ctx, req)

	// Calculate the duration and record it in the histogram
	statusCode := status.Code(err).String()
	duration := time.Since(start)
	appmetrics.PostHitDuration.WithLabelValues(info.FullMethod).Set(float64(duration.Milliseconds()))

	appmetrics.PostHits.WithLabelValues(info.FullMethod, statusCode).Inc()

	appmetrics.PostTotalHits.WithLabelValues().Inc()

	return h, err
}

func PublicGroupHitMetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Call the next handler
	h, err := handler(ctx, req)

	// Calculate the duration and record it in the histogram
	statusCode := status.Code(err).String()
	duration := time.Since(start)
	appmetrics.PublicGroupHitDuration.WithLabelValues(info.FullMethod).Set(float64(duration.Milliseconds()))

	appmetrics.PublicGroupHits.WithLabelValues(info.FullMethod, statusCode).Inc()

	appmetrics.PublicGroupTotalHits.WithLabelValues().Inc()

	return h, err
}

func UserHitMetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Call the next handler
	h, err := handler(ctx, req)

	// Calculate the duration and record it in the histogram
	statusCode := status.Code(err).String()
	duration := time.Since(start)
	appmetrics.UserHitDuration.WithLabelValues(info.FullMethod).Set(float64(duration.Milliseconds()))

	appmetrics.UserHits.WithLabelValues(info.FullMethod, statusCode).Inc()

	appmetrics.UserTotalHits.WithLabelValues().Inc()

	return h, err
}
