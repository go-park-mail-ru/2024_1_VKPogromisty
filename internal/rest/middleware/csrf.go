package middleware

import (
	"net/http"
	"socio/errors"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
	"socio/usecase/csrf"
)

const (
	CSRFHeader = "X-CSRF-Token"
)

func CSRFMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrfService := csrf.NewCSRFService()

		token := r.Header.Get(CSRFHeader)
		if token == "" {
			json.ServeJSONError(r.Context(), w, errors.ErrForbidden)
			return
		}

		sessionID, err := requestcontext.GetSessionID(r.Context())
		if err != nil {
			json.ServeJSONError(r.Context(), w, err)
			return
		}

		userID, err := requestcontext.GetUserID(r.Context())
		if err != nil {
			json.ServeJSONError(r.Context(), w, err)
			return
		}

		err = csrfService.Check(sessionID, userID, token)
		if err != nil {
			json.ServeJSONError(r.Context(), w, err)
			return
		}

		h.ServeHTTP(w, r)
	})
}
