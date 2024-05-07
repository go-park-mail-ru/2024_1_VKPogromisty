package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestTrackDuration(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		path     string
		expected int
	}{
		{
			name:     "GET /test",
			method:   "GET",
			path:     "/test",
			expected: 200,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(100 * time.Millisecond) // simulate some processing time
				w.WriteHeader(tt.expected)
			})

			router.Use(TrackDuration)

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expected {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expected)
			}
		})
	}
}
