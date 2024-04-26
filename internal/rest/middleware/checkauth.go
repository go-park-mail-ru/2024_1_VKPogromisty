package middleware

import (
	"context"
	"net/http"
	"socio/errors"
	authpb "socio/internal/grpc/auth/proto"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
)

func CreateCheckIsAuthorizedMiddleware(authManager authpb.AuthClient) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := r.Cookie("session_id")
			if err == http.ErrNoCookie {
				json.ServeJSONError(r.Context(), w, errors.ErrUnauthorized)
				return
			}

			res, err := authManager.ValidateSession(r.Context(), &authpb.ValidateSessionRequest{SessionId: session.Value})
			if err != nil {
				json.ServeGRPCStatus(r.Context(), w, err)
				return
			}

			ctx := context.WithValue(r.Context(), requestcontext.UserIDKey, uint(res.UserId))
			ctx = context.WithValue(ctx, requestcontext.SessionIDKey, session.Value)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
