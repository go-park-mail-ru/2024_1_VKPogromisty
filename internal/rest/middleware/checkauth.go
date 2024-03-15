package middleware

import (
	"net/http"
	repository "socio/internal/repository/map"
	"socio/utils"
)

func CreateCheckIsAuthorizedMiddleware(sessionStorage *repository.Sessions) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := r.Cookie("session_id")
			if err == http.ErrNoCookie {
				utils.ServeJSONError(w, err)
				return
			}

			if _, err := sessionStorage.GetUserIDBySession(session.Value); err != nil {
				utils.ServeJSONError(w, err)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
