package middleware

import (
	"net/http"
	"socio/utils"
)

func SetUpCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", utils.ALLOWED_ORIGIN)
		w.Header().Set("Access-Control-Allow-Headers", utils.ALLOWED_HEADERS)
		w.Header().Set("Access-Control-Allow-Methods", utils.ALLOWED_METHODS)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
