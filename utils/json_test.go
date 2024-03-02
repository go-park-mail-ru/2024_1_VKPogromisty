package utils_test

import (
	"errors"
	"net/http/httptest"
	"socio/utils"
	"testing"
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
			utils.ServeJSONBody(rr, tt.body)

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
		{"valid error", errors.New("test error"), `{"error":"test error"}`},
		{"empty error", errors.New(""), `{"error":""}`},
		{"nil error", nil, `{"error":"internal server error"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			utils.ServeJSONError(rr, tt.err)

			if rr.Body.String() != tt.want {
				t.Errorf("ServeJSONError() = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}
