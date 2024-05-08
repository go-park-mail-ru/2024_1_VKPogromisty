package errors_test

import (
	"encoding/json"
	"errors"
	"net/http"
	errorsCustom "socio/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		{
			name:       "Error is not in HTTPErrors",
			err:        errors.New(""),
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

func TestParseGRPCError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedMsg    string
		expectedStatus int
	}{
		{
			name:           "Parse GRPC error with status",
			err:            status.Error(codes.InvalidArgument, errorsCustom.InvalidDataMsg),
			expectedMsg:    errorsCustom.InvalidDataMsg,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Parse GRPC error without status",
			err:            errors.New("some error"),
			expectedMsg:    "some error",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Parse nil error",
			err:            nil,
			expectedMsg:    errorsCustom.InternalMsg,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, status := errorsCustom.ParseGRPCError(tt.err)
			assert.Equal(t, tt.expectedMsg, msg)
			assert.Equal(t, tt.expectedStatus, status)
		})
	}
}

func TestCustomError_GRPCStatus(t *testing.T) {
	unk := errorsCustom.NewCustomError(errors.New("asd"))
	tests := []struct {
		name         string
		err          *errorsCustom.CustomError
		expectedCode codes.Code
	}{
		{
			name:         "GRPC status with known error",
			err:          &errorsCustom.ErrInvalidData,
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "GRPC status with unknown error",
			err:          &unk,
			expectedCode: codes.Internal,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grpcStatus := tt.err.GRPCStatus()
			assert.Equal(t, tt.expectedCode, grpcStatus.Code())
			assert.Equal(t, tt.err.Error(), grpcStatus.Message())
		})
	}
}
