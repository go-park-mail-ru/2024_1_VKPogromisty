package json_test

import (
	"context"
	goErr "errors"
	"net/http"
	"net/http/httptest"
	"socio/errors"
	"socio/pkg/json"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServeJSONBody(t *testing.T) {
	tests := []struct {
		name string
		body interface{}
		want string
	}{
		{"valid body", "test", `{"body":"test"}`},
		{"empty body", "", `{"body":""}`},
		{"invalid body", func() {}, `{"error":"unable to return json reponse"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			json.ServeJSONBody(context.Background(), rr, tt.body, http.StatusOK)

			if rr.Body.String() != tt.want {
				t.Errorf("ServeJSONBody() = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestServeJSONError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"valid error", errors.ErrUnauthorized, `{"error":"unauthorized"}`},
		{"empty error", goErr.New(""), `{"error":"internal server error"}`},
		{"nil error", nil, `{"error":"internal server error"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			json.ServeJSONError(context.Background(), rr, tt.err)

			if rr.Body.String() != tt.want {
				t.Errorf("ServeJSONError() = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestServeGRPCStatus(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Test OK",
			err:            nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "{\"error\":\"internal server error\"}",
		},
		{
			name:           "Test NotFound",
			err:            status.Error(codes.NotFound, "not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"not found"}`,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "http://example.com/foo", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json.ServeGRPCStatus(r.Context(), w, tt.err)
			})

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
