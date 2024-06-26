package requestcontext

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func TestGetUserID(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		args       args
		wantUserID uint
		wantErr    bool
	}{
		{
			name: "Valid user ID",
			args: args{
				ctx: context.WithValue(context.Background(), UserIDKey, uint(1)),
			},
			wantUserID: 1,
			wantErr:    false,
		},
		{
			name: "Invalid user ID",
			args: args{
				ctx: context.Background(),
			},
			wantUserID: 0,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := GetUserID(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("GetUserID() = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetSessionID(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name          string
		args          args
		wantSessionID string
		wantErr       bool
	}{
		{
			name: "Valid session ID",
			args: args{
				ctx: context.WithValue(context.Background(), SessionIDKey, "sessionID"),
			},
			wantSessionID: "sessionID",
			wantErr:       false,
		},
		{
			name: "Invalid session ID",
			args: args{
				ctx: context.Background(),
			},
			wantSessionID: "",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSessionID, err := GetSessionID(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSessionID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSessionID != tt.wantSessionID {
				t.Errorf("GetSessionID() = %v, want %v", gotSessionID, tt.wantSessionID)
			}
		})
	}
}

func TestGetLogger(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		args       args
		wantLogger *zap.Logger
		wantErr    bool
	}{
		{
			name: "Valid logger",
			args: args{
				ctx: context.WithValue(context.Background(), LoggerKey, zap.NewNop()),
			},
			wantLogger: zap.NewNop(),
			wantErr:    false,
		},
		{
			name: "Invalid logger",
			args: args{
				ctx: context.Background(),
			},
			wantLogger: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLogger, err := GetLogger(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLogger, tt.wantLogger) {
				t.Errorf("GetLogger() = %v, want %v", gotLogger, tt.wantLogger)
			}
		})
	}
}

func TestGetRequestID(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "Context with RequestID",
			ctx:      context.WithValue(context.Background(), RequestIDKey, "test-id"),
			expected: "test-id",
		},
		{
			name:     "Context without RequestID",
			ctx:      context.Background(),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRequestID(tt.ctx)
			if err != nil {
				t.Errorf("GetRequestID() error = %v", err)
				return
			}
			if tt.name == "Context without RequestID" {
				_, err := uuid.Parse(got)
				if err != nil {
					t.Errorf("GetRequestID() = %v, want a valid UUID", got)
				}
			} else if got != tt.expected {
				t.Errorf("GetRequestID() = %v, want %v", got, tt.expected)
			}
		})
	}
}
