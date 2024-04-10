package json_test

import (
	"context"
	goErr "errors"
	"net/http/httptest"
	"socio/errors"
	"socio/pkg/json"
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
			json.ServeJSONBody(context.Background(), rr, tt.body)

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
