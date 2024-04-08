package csrf_test

import (
	"os"
	"socio/usecase/csrf"
	"testing"
	"time"

	customtime "socio/pkg/time"
)

func TestCSRFService_Create(t *testing.T) {
	type args struct {
		sessionID    string
		userID       uint
		tokenExpTime int64
	}
	tests := []struct {
		name      string
		c         *csrf.CSRFService
		args      args
		wantToken bool
		wantErr   bool
		prepare   func()
	}{
		{
			name: "Test 1",
			c:    csrf.NewCSRFService(customtime.MockTimeProvider{}),
			args: args{
				sessionID:    "sessionID",
				userID:       1,
				tokenExpTime: 1,
			},
			wantToken: true,
			wantErr:   false,
			prepare: func() {
				os.Setenv("CSRF_SECRET", "secret")
			},
		},
		{
			name: "Test 1",
			c:    csrf.NewCSRFService(customtime.MockTimeProvider{}),
			args: args{
				sessionID:    "sessionID",
				userID:       1,
				tokenExpTime: 1,
			},
			wantToken: true,
			wantErr:   false,
			prepare: func() {
				os.Unsetenv("CSRF_SECRET")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := tt.c.Create(tt.args.sessionID, tt.args.userID, tt.args.tokenExpTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("CSRFService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (gotToken != "") != tt.wantToken {
				t.Errorf("CSRFService.Create() = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}

func TestCSRFService_Check(t *testing.T) {
	validToken, _ := csrf.NewCSRFService(customtime.MockTimeProvider{}).Create("sessionID", 1, time.Now().Add(time.Hour).Unix())

	type args struct {
		sessionID  string
		userID     uint
		inputToken string
	}
	tests := []struct {
		name    string
		c       *csrf.CSRFService
		args    args
		wantErr bool
	}{
		{
			name: "unhashed token",
			c:    csrf.NewCSRFService(customtime.MockTimeProvider{}),
			args: args{
				sessionID:  "sessionID",
				userID:     1,
				inputToken: "inputToken",
			},
			wantErr: true,
		},
		{
			name: "invalid sessionID",
			c:    csrf.NewCSRFService(customtime.MockTimeProvider{}),
			args: args{
				sessionID:  "session",
				userID:     1,
				inputToken: validToken,
			},
			wantErr: true,
		},
		{
			name: "Test 2",
			c:    csrf.NewCSRFService(customtime.MockTimeProvider{}),
			args: args{
				sessionID:  "sessionID",
				userID:     1,
				inputToken: validToken,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Check(tt.args.sessionID, tt.args.userID, tt.args.inputToken); (err != nil) != tt.wantErr {
				t.Errorf("CSRFService.Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
