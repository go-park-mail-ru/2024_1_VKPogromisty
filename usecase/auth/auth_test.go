package auth_test

import (
	"context"
	"reflect"
	"socio/domain"
	"socio/errors"
	mock_auth "socio/mocks/usecase/auth"
	"socio/pkg/hash"
	customtime "socio/pkg/time"
	"socio/usecase/auth"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestService_Login(t *testing.T) {
	timeProv := customtime.MockTimeProvider{}

	type fields struct {
		SessionStorage *mock_auth.MockSessionStorage
	}

	type args struct {
		ctx        context.Context
		loginInput auth.LoginInput
		user       *domain.User
	}

	tests := []struct {
		name        string
		args        args
		prepareMock func(*fields)
		wantSession string
		wantErr     bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				loginInput: auth.LoginInput{
					Email:    "john@mail.ru",
					Password: "password",
				},
				user: &domain.User{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@mail.ru",
					Password:  hash.HashPassword("password", []byte("salt")),
					Salt:      "salt",
					Avatar:    "default_avatar.png",
					DateOfBirth: customtime.CustomTime{
						Time: timeProv.Now(),
					},
				},
			},
			wantSession: "session_id",
			wantErr:     false,
			prepareMock: func(f *fields) {
				f.SessionStorage.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return("session_id", nil)
			},
		},
		{
			name: "invalid password",
			args: args{
				ctx: context.Background(),
				loginInput: auth.LoginInput{
					Email:    "john@mail.ru",
					Password: "",
				},
				user: &domain.User{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@mail.ru",
					Password:  hash.HashPassword("password", []byte("salt")),
					Salt:      "salt",
					Avatar:    "default_avatar.png",
					DateOfBirth: customtime.CustomTime{
						Time: timeProv.Now(),
					},
				},
			},
			wantSession: "",
			wantErr:     true,
			prepareMock: func(f *fields) {},
		},
		{
			name: "err internal",
			args: args{
				ctx: context.Background(),
				loginInput: auth.LoginInput{
					Email:    "john@mail.ru",
					Password: "password",
				},
				user: &domain.User{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@mail.ru",
					Password:  hash.HashPassword("password", []byte("salt")),
					Salt:      "salt",
					Avatar:    "default_avatar.png",
					DateOfBirth: customtime.CustomTime{
						Time: timeProv.Now(),
					},
				},
			},
			wantSession: "",
			wantErr:     true,
			prepareMock: func(f *fields) {
				f.SessionStorage.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return("", errors.ErrInternal)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				SessionStorage: mock_auth.NewMockSessionStorage(ctrl),
			}

			if tt.prepareMock != nil {
				tt.prepareMock(&f)
			}

			s := auth.NewService(f.SessionStorage)

			gotSession, err := s.Login(tt.args.ctx, tt.args.loginInput, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSession, tt.wantSession) {
				t.Errorf("Service.Login() gotSession = %v, want %v", gotSession, tt.wantSession)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		sessionID string
		mock      func(storage *mock_auth.MockSessionStorage, sessionID string)
		wantErr   bool
	}{
		{
			name:      "Test OK",
			sessionID: "testSessionID",
			mock: func(storage *mock_auth.MockSessionStorage, sessionID string) {
				storage.EXPECT().DeleteSession(gomock.Any(), sessionID).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Test Error",
			sessionID: "testSessionID",
			mock: func(storage *mock_auth.MockSessionStorage, sessionID string) {
				storage.EXPECT().DeleteSession(gomock.Any(), sessionID).Return(errors.ErrInternal)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mock_auth.NewMockSessionStorage(ctrl)

			a := auth.NewService(storage)

			tt.mock(storage, tt.sessionID)

			err := a.Logout(context.Background(), tt.sessionID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Logout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsAuthorized(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		sessionID  string
		mock       func(storage *mock_auth.MockSessionStorage, sessionID string)
		wantUserID uint
		wantErr    bool
	}{
		{
			name:      "Test OK",
			sessionID: "testSessionID",
			mock: func(storage *mock_auth.MockSessionStorage, sessionID string) {
				storage.EXPECT().GetUserIDBySession(gomock.Any(), sessionID).Return(uint(1), nil)
			},
			wantUserID: uint(1),
			wantErr:    false,
		},
		{
			name:      "Test Error",
			sessionID: "testSessionID",
			mock: func(storage *mock_auth.MockSessionStorage, sessionID string) {
				storage.EXPECT().GetUserIDBySession(gomock.Any(), sessionID).Return(uint(0), errors.ErrInternal)
			},
			wantUserID: uint(0),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mock_auth.NewMockSessionStorage(ctrl)

			a := auth.NewService(storage)

			tt.mock(storage, tt.sessionID)

			gotUserID, err := a.IsAuthorized(context.Background(), tt.sessionID)

			if (err != nil) != tt.wantErr {
				t.Errorf("IsAuthorized() error = %v, wantErr %v", err, tt.wantErr)
			}

			if gotUserID != tt.wantUserID {
				t.Errorf("IsAuthorized() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
