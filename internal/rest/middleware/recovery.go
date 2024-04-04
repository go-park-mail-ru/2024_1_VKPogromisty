package middleware

import (
	"net/http"
	"fmt"
	"socio/pkg/json"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				json.ServeJSONError(r.Context(), w, fmt.Errorf("%v", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
