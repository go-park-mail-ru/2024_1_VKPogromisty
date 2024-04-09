package profile_test

import (
	"context"
	"mime/multipart"
	"reflect"
	"socio/domain"
	"socio/errors"
	mock_profile "socio/mocks/usecase/profile"
	"socio/pkg/sanitizer"
	customtime "socio/pkg/time"
	"socio/usecase/profile"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/microcosm-cc/bluemonday"
)

type fields struct {
	UserStorage    *mock_profile.MockUserStorage
	SessionStorage *mock_profile.MockSessionStorage
	Sanitizer      *sanitizer.Sanitizer
}

var timeProv = customtime.MockTimeProvider{}

func TestService_GetUserByIDWithSubsInfo(t *testing.T) {
	type args struct {
		ctx              context.Context
		userID           uint
		authorizedUserID uint
	}

	tests := []struct {
		name             string
		args             args
		wantUserWithInfo profile.UserWithSubsInfo
		wantErr          bool
		prepareMock      func(*fields)
	}{
		{
			name: "success",
			args: args{
				ctx:              context.Background(),
				userID:           1,
				authorizedUserID: 2,
			},
			wantUserWithInfo: profile.UserWithSubsInfo{
				User: &domain.User{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "email",
					Avatar:    "avatar",
					DateOfBirth: customtime.CustomTime{
						Time: timeProv.Now(),
					},
					CreatedAt: customtime.CustomTime{
						Time: timeProv.Now(),
					},
					UpdatedAt: customtime.CustomTime{
						Time: timeProv.Now(),
					},
				},
				IsSubscriber:   true,
				IsSubscribedTo: true,
			},
			wantErr: false,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByIDWithSubsInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&domain.User{
						ID:        1,
						FirstName: "John",
						LastName:  "Doe",
						Email:     "email",
						Avatar:    "avatar",
						DateOfBirth: customtime.CustomTime{
							Time: timeProv.Now(),
						},
						CreatedAt: customtime.CustomTime{
							Time: timeProv.Now(),
						},
						UpdatedAt: customtime.CustomTime{
							Time: timeProv.Now(),
						},
					}, true, true, nil)
			},
		},
		{
			name: "error getting user",
			args: args{
				ctx:              context.Background(),
				userID:           1,
				authorizedUserID: 2,
			},
			wantUserWithInfo: profile.UserWithSubsInfo{
				User:           nil,
				IsSubscriber:   false,
				IsSubscribedTo: false,
			},
			wantErr: true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByIDWithSubsInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, false, false, errors.ErrNotFound)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				UserStorage:    mock_profile.NewMockUserStorage(ctrl),
				SessionStorage: mock_profile.NewMockSessionStorage(ctrl),
				Sanitizer:      sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
			}

			if tt.prepareMock != nil {
				tt.prepareMock(&f)
			}

			p := profile.NewProfileService(f.UserStorage, f.SessionStorage, f.Sanitizer)

			gotUserWithInfo, err := p.GetUserByIDWithSubsInfo(tt.args.ctx, tt.args.userID, tt.args.authorizedUserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetUserByIDWithSubsInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUserWithInfo, tt.wantUserWithInfo) {
				t.Errorf("Service.GetUserByIDWithSubsInfo() = %v, want %v", gotUserWithInfo, tt.wantUserWithInfo)
			}
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	type args struct {
		ctx   context.Context
		input profile.UpdateUserInput
	}

	tests := []struct {
		name            string
		args            args
		wantUpdatedUser *domain.User
		wantErr         bool
		prepareMock     func(*fields)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				input: profile.UpdateUserInput{
					ID:             1,
					FirstName:      "John",
					LastName:       "Doe",
					Email:          "email@email",
					Password:       "password",
					RepeatPassword: "password",
					DateOfBirth:    "2006-01-02",
				},
			},
			wantUpdatedUser: &domain.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "email",
			},
			wantErr: false,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(
					&domain.User{
						ID:        1,
						FirstName: "John",
						LastName:  "Doe",
						Email:     "email",
					}, nil)
				f.UserStorage.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
				f.UserStorage.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(&domain.User{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "email",
				}, nil)
			},
		},
		{
			name: "error getting user",
			args: args{
				ctx: context.Background(),
				input: profile.UpdateUserInput{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "email",
				},
			},
			wantUpdatedUser: nil,
			wantErr:         true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name: "error updating user",
			args: args{
				ctx: context.Background(),
				input: profile.UpdateUserInput{
					ID:             1,
					FirstName:      "John",
					LastName:       "Doe",
					Email:          "email@email",
					Password:       "password",
					RepeatPassword: "password",
					DateOfBirth:    "2006-01-02",
				},
			},
			wantUpdatedUser: nil,
			wantErr:         true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(
					&domain.User{
						ID:        1,
						FirstName: "John",
						LastName:  "Doe",
						Email:     "email",
					}, nil)
				f.UserStorage.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
				f.UserStorage.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
		{
			name: "error invalid image",
			args: args{
				ctx: context.Background(),
				input: profile.UpdateUserInput{
					ID:             1,
					FirstName:      "John",
					LastName:       "Doe",
					Email:          "email@email",
					Password:       "password",
					RepeatPassword: "password",
					DateOfBirth:    "2006-01-02",
					Avatar:         &multipart.FileHeader{},
				},
			},
			wantUpdatedUser: nil,
			wantErr:         true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(
					&domain.User{
						ID:        1,
						FirstName: "John",
						LastName:  "Doe",
						Email:     "email",
					}, nil)
				f.UserStorage.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				UserStorage:    mock_profile.NewMockUserStorage(ctrl),
				SessionStorage: mock_profile.NewMockSessionStorage(ctrl),
				Sanitizer:      sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
			}

			if tt.prepareMock != nil {
				tt.prepareMock(&f)
			}

			p := profile.NewProfileService(f.UserStorage, f.SessionStorage, f.Sanitizer)

			gotUpdatedUser, err := p.UpdateUser(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUpdatedUser, tt.wantUpdatedUser) {
				t.Errorf("Service.UpdateUser() = %v, want %v", gotUpdatedUser, tt.wantUpdatedUser)
			}
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	type args struct {
		ctx       context.Context
		userID    uint
		sessionID string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		prepareMock func(*fields)
	}{
		{
			name: "success",
			args: args{
				ctx:       context.Background(),
				userID:    1,
				sessionID: "session",
			},
			wantErr: false,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(nil)
				f.SessionStorage.EXPECT().DeleteSession(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "error deleting user",
			args: args{
				ctx:       context.Background(),
				userID:    1,
				sessionID: "session",
			},
			wantErr: true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(errors.ErrInternal)
			},
		},
		{
			name: "error deleting session",
			args: args{
				ctx:       context.Background(),
				userID:    1,
				sessionID: "session",
			},
			wantErr: true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(nil)
				f.SessionStorage.EXPECT().DeleteSession(gomock.Any(), gomock.Any()).Return(errors.ErrInternal)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				UserStorage:    mock_profile.NewMockUserStorage(ctrl),
				SessionStorage: mock_profile.NewMockSessionStorage(ctrl),
				Sanitizer:      sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
			}

			if tt.prepareMock != nil {
				tt.prepareMock(&f)
			}

			p := profile.NewProfileService(f.UserStorage, f.SessionStorage, f.Sanitizer)

			if err := p.DeleteUser(tt.args.ctx, tt.args.userID, tt.args.sessionID); (err != nil) != tt.wantErr {
				t.Errorf("Service.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}