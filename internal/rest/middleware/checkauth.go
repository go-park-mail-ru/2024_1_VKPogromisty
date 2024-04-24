package middleware

import (
	"context"
	"fmt"
	"net/http"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
	"socio/usecase/auth"
)

func CreateCheckIsAuthorizedMiddleware(sessionStorage auth.SessionStorage) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := r.Cookie("session_id")
			if err == http.ErrNoCookie {
				fmt.Println("no cookie")
				json.ServeJSONError(r.Context(), w, err)
				return
			}

			userID, err := sessionStorage.GetUserIDBySession(r.Context(), session.Value)
			if err != nil {
				json.ServeJSONError(r.Context(), w, err)
				return
			}

			ctx := context.WithValue(r.Context(), requestcontext.UserIDKey, userID)
			ctx = context.WithValue(ctx, requestcontext.SessionIDKey, session.Value)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
