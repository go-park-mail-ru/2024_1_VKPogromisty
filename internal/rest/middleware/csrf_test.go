package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"socio/pkg/requestcontext"
	customtime "socio/pkg/time"
	"socio/usecase/csrf"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestCSRFMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tp := customtime.MockTimeProvider{}

	os.Setenv("CSRF_SECRET", "testSecret")

	mockCSRFService := csrf.NewCSRFService(tp)

	validCsrfToken, _ := mockCSRFService.Create("testSessionID", uint(1), tp.Now().Add(300000*time.Hour).Unix())

	validContext := context.WithValue(context.Background(), requestcontext.SessionIDKey, "testSessionID")
	validContext = context.WithValue(validContext, requestcontext.UserIDKey, uint(1))

	contextNoSessionID := context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1))

	contextNoUserID := context.WithValue(context.Background(), requestcontext.SessionIDKey, "testSessionID")

	testCases := []struct {
		name           string
		token          string
		sessionID      string
		userID         uint
		expectedStatus int
		ctx            context.Context
	}{
		{
			name:           "valid CSRF token",
			token:          validCsrfToken,
			sessionID:      "testSessionID",
			userID:         1,
			expectedStatus: http.StatusOK,
			ctx:            validContext,
		},
		{
			name:           "invalid CSRF token",
			token:          "invalidToken",
			sessionID:      "testSessionID",
			userID:         1,
			expectedStatus: http.StatusBadRequest,
			ctx:            validContext,
		},
		{
			name:           "invalid session ID",
			token:          validCsrfToken,
			sessionID:      "invalidSessionID",
			userID:         1,
			expectedStatus: http.StatusBadRequest,
			ctx:            contextNoSessionID,
		},
		{
			name:           "invalid user ID",
			token:          validCsrfToken,
			sessionID:      "testSessionID",
			userID:         0,
			expectedStatus: http.StatusBadRequest,
			ctx:            contextNoUserID,
		},
		{
			name:           "no CSRF token",
			token:          "",
			sessionID:      "testSessionID",
			userID:         1,
			expectedStatus: http.StatusForbidden,
			ctx:            validContext,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set(CSRFHeader, tc.token)

			req = req.WithContext(tc.ctx)

			rr := httptest.NewRecorder()
			CSRFMiddleware := CreateCSRFMiddleware(mockCSRFService)
			handlerFunc := CSRFMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

			handlerFunc.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}
		})
	}
}
