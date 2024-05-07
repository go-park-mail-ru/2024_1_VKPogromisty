package middleware

import (
	"fmt"
	"net/http"
	"socio/pkg/appmetrics"
	"time"

	"github.com/gorilla/mux"
)

func TrackDuration(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer that allows us to capture the status code
		// ww := &responseWriter{ResponseWriter: w}

		// Call the next handler
		next.ServeHTTP(w, r)

		// Calculate the duration and record it in the histogram
		duration := time.Since(start)

		route := mux.CurrentRoute(r)
		pathTemplate := ""
		var err error
		if route == nil {
			return
		}

		pathTemplate, err = route.GetPathTemplate()
		if err != nil {
			fmt.Println(err)
			return
		}

		appmetrics.AppHitDuration.WithLabelValues(r.Method, pathTemplate).Set(float64(duration.Milliseconds()))
		appmetrics.AppHits.WithLabelValues(r.Method, pathTemplate, http.StatusText(http.StatusOK)).Inc()
		appmetrics.AppTotalHits.WithLabelValues().Inc()
	})
}
