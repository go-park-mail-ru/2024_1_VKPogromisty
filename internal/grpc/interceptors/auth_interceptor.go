package interceptors

import (
	"context"

	"socio/errors"
	authpb "socio/internal/grpc/auth/proto"
	"socio/pkg/requestcontext"

	"google.golang.org/grpc"
)

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// AuthUnaryInterceptor is a unary server interceptor for authentication and authorization
func CreateAuthUnaryInterceptor(authClient authpb.AuthClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if _, ok := PublicMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		userID, err := checkSessionID(ctx, authClient)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, requestcontext.UserIDKey, userID)

		return handler(ctx, req)
	}
}

func CreateAuthStreamInterceptor(authClient authpb.AuthClient) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if _, ok := PublicMethods[info.FullMethod]; ok {
			return handler(srv, stream)
		}

		ctx := stream.Context()

		userID, err := checkSessionID(ctx, authClient)
		if err != nil {
			return err
		}

		newCtx := context.WithValue(ctx, requestcontext.UserIDKey, userID)
		wrapped := &wrappedStream{ServerStream: stream, ctx: newCtx}

		return handler(srv, wrapped)
	}
}

// checkSessionID checks the session_id in the metadata of the context
func checkSessionID(ctx context.Context, authClient authpb.AuthClient) (userID uint, err error) {
	sessionID, err := requestcontext.GetSessionID(ctx)
	if err != nil {
		err = errors.ErrUnauthorized
		return
	}

	res, err := authClient.ValidateSession(ctx, &authpb.ValidateSessionRequest{
		SessionId: sessionID,
	})
	if err != nil {
		err = errors.ErrUnauthorized
		return
	}

	userID = uint(res.UserId)

	return
}
