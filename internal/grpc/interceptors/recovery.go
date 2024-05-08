package interceptors

import (
	"context"
	"socio/pkg/contextlogger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryRecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = status.Errorf(codes.Internal, "panic: %v", r)
			}

			contextlogger.LogErr(ctx, err)
		}
	}()

	return handler(ctx, req)
}
