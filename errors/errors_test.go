package errors_test

import (
	"errors"
	errorsCustom "socio/errors"
	"testing"
)

func TestParseHTTPError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantMsg    string
		wantStatus int
	}{
		{
			name:       "Error is nil",
			err:        nil,
			wantMsg:    "internal server error",
			wantStatus: 500,
		},
		{
			name:       "Error is in HTTPErrors",
			err:        errorsCustom.ErrUnauthorized,
			wantMsg:    "unauthorized",
			wantStatus: 401,
		},
		{
			name:       "Error is not in HTTPErrors",
			err:        errors.New("unknown error"),
			wantMsg:    "internal server error",
			wantStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMsg, gotStatus := errorsCustom.ParseHTTPError(tt.err)
			if gotMsg != tt.wantMsg {
				t.Errorf("ParseHTTPError() gotMsg = %v, want %v", gotMsg, tt.wantMsg)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("ParseHTTPError() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}
