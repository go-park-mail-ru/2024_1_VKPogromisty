package middleware

import (
	"net/http"
	"socio/pkg/json"
	"socio/usecase/auth"
)

func CreateCheckIsAuthorizedMiddleware(sessionStorage auth.SessionStorage) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := r.Cookie("session_id")
			if err == http.ErrNoCookie {
				json.ServeJSONError(w, err)
				return
			}

			if _, err := sessionStorage.GetUserIDBySession(session.Value); err != nil {
				json.ServeJSONError(w, err)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
