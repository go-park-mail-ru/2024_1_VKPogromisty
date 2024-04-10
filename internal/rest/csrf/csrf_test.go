package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"socio/pkg/requestcontext"
	customtime "socio/pkg/time"
	"socio/usecase/csrf"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetCSRFToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTimeProvider := customtime.MockTimeProvider{}

	handler := NewCSRFHandler(mockTimeProvider)
	handler.CSRFService = csrf.NewCSRFService(customtime.MockTimeProvider{})

	validContext := context.WithValue(context.Background(), requestcontext.SessionIDKey, "testSessionID")
	validContext = context.WithValue(validContext, requestcontext.UserIDKey, uint(1))

	contextInvalidUser := context.WithValue(context.Background(), requestcontext.SessionIDKey, "testSessionID")
	contextInvalidUser = context.WithValue(contextInvalidUser, requestcontext.UserIDKey, "invalid")

	testCases := []struct {
		name           string
		expected       string
		expectedStatus int
		ctx            context.Context
	}{
		{
			name:           "valid request",
			expected:       `{"body":{"csrfToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWQiOiJ0ZXN0U2Vzc2lvbklEIiwidWlkIjoxLCJleHAiOjE2MDk0NjEwMDAsImlhdCI6MTYwOTQ1OTIwMH0.2fL4Gq_2gwKs8yl76Ytu98B3FvBirnzhPnjd-7j6m1A"}}`,
			expectedStatus: http.StatusOK,
			ctx:            validContext,
		},
		{
			name:           "invalid request",
			expected:       `{"error":"invalid data"}`,
			expectedStatus: http.StatusBadRequest,
			ctx:            context.Background(),
		},
		{
			name:           "invalid user",
			expected:       `{"error":"invalid data"}`,
			expectedStatus: http.StatusBadRequest,
			ctx:            contextInvalidUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/csrf/", nil)
			if err != nil {
				t.Fatal(err)
			}

			req = req.WithContext(tc.ctx)

			rr := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(handler.GetCSRFToken)

			handlerFunc.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			if rr.Body.String() != tc.expected {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tc.expected)
			}
		})
	}
}
