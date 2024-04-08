package errors_test

import (
	"encoding/json"
	"errors"
	errorsCustom "socio/errors"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestMarshalError(t *testing.T) {
	testCases := []struct {
		name          string
		err           error
		expectedError string
	}{
		{
			name:          "Test error",
			err:           errors.New("Test error"),
			expectedError: `{"error":"Test error"}`,
		},
		{
			name:          "Another test error",
			err:           errors.New("Another test error"),
			expectedError: `{"error":"Another test error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := errorsCustom.MarshalError(tc.err)
			assert.NoError(t, err)

			var result map[string]string
			err = json.Unmarshal(data, &result)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedError, string(data))
		})
	}
}
