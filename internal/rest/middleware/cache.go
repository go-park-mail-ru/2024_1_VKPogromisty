package middleware

import "net/http"

func DisableCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		h.ServeHTTP(w, r)
	})
}
