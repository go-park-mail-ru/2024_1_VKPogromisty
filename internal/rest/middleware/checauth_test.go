package middleware_test

import (
	"net/http"
	"net/http/httptest"
	repository "socio/internal/repository/map"
	"socio/internal/rest/middleware"
	"sync"
	"testing"
)

func TestCheckIsAuthorizedMiddleware(t *testing.T) {
	var sessionStorage = repository.NewSessions(&sync.Map{})

	sessionID := sessionStorage.CreateSession(0)

	CheckIsAuthorizedMiddleware := middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage)

	tests := []struct {
		name     string
		cookie   *http.Cookie
		wantCode int
	}{
		{
			name:     "Valid session",
			cookie:   &http.Cookie{Name: "session_id", Value: sessionID},
			wantCode: http.StatusOK,
		},
		{
			name:     "Invalid session",
			cookie:   &http.Cookie{Name: "session_id", Value: "invalidSessionValue"},
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "No session",
			cookie:   nil,
			wantCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			rr := httptest.NewRecorder()

			handler := CheckIsAuthorizedMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantCode)
			}
		})
	}
}
