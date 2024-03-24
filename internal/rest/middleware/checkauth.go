package middleware

import (
	"context"
	"net/http"
	"socio/pkg/json"
	"socio/usecase/auth"
)

type ContextKey string

const UserIDKey ContextKey = "userID"
const SessionIDKey ContextKey = "sessionID"

func CreateCheckIsAuthorizedMiddleware(sessionStorage auth.SessionStorage) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := r.Cookie("session_id")
			if err == http.ErrNoCookie {
				json.ServeJSONError(w, err)
				return
			}

			userID, err := sessionStorage.GetUserIDBySession(session.Value)

			if err != nil {
				json.ServeJSONError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, SessionIDKey, session.Value)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
