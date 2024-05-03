package user_test

import (
	"context"
	"reflect"
	"socio/domain"
	"socio/errors"
	"testing"

	mock_users "socio/mocks/usecase/user"
	"socio/usecase/user"

	"github.com/golang/mock/gomock"
)

func TestCreatePublicGroupAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		publicGroupAdmin     *domain.PublicGroupAdmin
		mock                 func(userStorage *mock_users.MockUserStorage, publicGroupAdmin *domain.PublicGroupAdmin)
		wantPublicGroupAdmin *domain.PublicGroupAdmin
		wantErr              bool
	}{
		{
			name:             "Test OK",
			publicGroupAdmin: &domain.PublicGroupAdmin{UserID: 1, PublicGroupID: 1},
			mock: func(userStorage *mock_users.MockUserStorage, publicGroupAdmin *domain.PublicGroupAdmin) {
				userStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&domain.User{ID: 1}, nil)
				userStorage.EXPECT().StorePublicGroupAdmin(gomock.Any(), gomock.Any()).Return(publicGroupAdmin, nil)
			},
			wantPublicGroupAdmin: &domain.PublicGroupAdmin{UserID: 1, PublicGroupID: 1},
			wantErr:              false,
		},
		{
			name:             "Test Error",
			publicGroupAdmin: &domain.PublicGroupAdmin{UserID: 0, PublicGroupID: 0},
			mock: func(userStorage *mock_users.MockUserStorage, publicGroupAdmin *domain.PublicGroupAdmin) {
				userStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
			wantPublicGroupAdmin: nil,
			wantErr:              true,
		},
		{
			name:             "Test tyazhelo",
			publicGroupAdmin: &domain.PublicGroupAdmin{UserID: 1, PublicGroupID: 1},
			mock: func(userStorage *mock_users.MockUserStorage, publicGroupAdmin *domain.PublicGroupAdmin) {
				userStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&domain.User{ID: 1}, nil)
				userStorage.EXPECT().StorePublicGroupAdmin(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
			wantPublicGroupAdmin: nil,
			wantErr:              true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userStorage := mock_users.NewMockUserStorage(ctrl)

			s := user.NewUserService(userStorage, nil)

			tt.mock(userStorage, tt.publicGroupAdmin)

			gotPublicGroupAdmin, err := s.CreatePublicGroupAdmin(context.Background(), tt.publicGroupAdmin)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePublicGroupAdmin() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotPublicGroupAdmin, tt.wantPublicGroupAdmin) {
				t.Errorf("CreatePublicGroupAdmin() gotPublicGroupAdmin = %v, want %v", gotPublicGroupAdmin, tt.wantPublicGroupAdmin)
			}
		})
	}
}

func TestDeletePublicGroupAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		publicGroupAdmin *domain.PublicGroupAdmin
		mock             func(userStorage *mock_users.MockUserStorage, publicGroupAdmin *domain.PublicGroupAdmin)
		wantErr          bool
	}{
		{
			name:             "Test OK",
			publicGroupAdmin: &domain.PublicGroupAdmin{UserID: 1, PublicGroupID: 1},
			mock: func(userStorage *mock_users.MockUserStorage, publicGroupAdmin *domain.PublicGroupAdmin) {
				userStorage.EXPECT().DeletePublicGroupAdmin(gomock.Any(), publicGroupAdmin).Return(nil)
			},
			wantErr: false,
		},
		{
			name:             "Test Error",
			publicGroupAdmin: &domain.PublicGroupAdmin{UserID: 0, PublicGroupID: 0},
			mock: func(userStorage *mock_users.MockUserStorage, publicGroupAdmin *domain.PublicGroupAdmin) {
				userStorage.EXPECT().DeletePublicGroupAdmin(gomock.Any(), publicGroupAdmin).Return(errors.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userStorage := mock_users.NewMockUserStorage(ctrl)

			s := user.NewUserService(userStorage, nil)

			tt.mock(userStorage, tt.publicGroupAdmin)

			err := s.DeletePublicGroupAdmin(context.Background(), tt.publicGroupAdmin)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeletePublicGroupAdmin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAdminsByPublicGroupID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		publicGroupID uint
		mock          func(userStorage *mock_users.MockUserStorage, publicGroupID uint)
		wantAdmins    []*domain.User
		wantErr       bool
	}{
		{
			name:          "Test OK",
			publicGroupID: 1,
			mock: func(userStorage *mock_users.MockUserStorage, publicGroupID uint) {
				admins := []*domain.User{
					{ID: 1},
					{ID: 2},
				}
				userStorage.EXPECT().GetAdminsByPublicGroupID(gomock.Any(), publicGroupID).Return(admins, nil)
			},
			wantAdmins: []*domain.User{
				{ID: 1},
				{ID: 2},
			},
			wantErr: false,
		},
		{
			name:          "Test Error",
			publicGroupID: 0,
			mock: func(userStorage *mock_users.MockUserStorage, publicGroupID uint) {
				userStorage.EXPECT().GetAdminsByPublicGroupID(gomock.Any(), publicGroupID).Return(
					nil, errors.ErrNotFound,
				)
			},
			wantAdmins: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userStorage := mock_users.NewMockUserStorage(ctrl)

			s := user.NewUserService(userStorage, nil)

			tt.mock(userStorage, tt.publicGroupID)

			gotAdmins, err := s.GetAdminsByPublicGroupID(context.Background(), tt.publicGroupID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAdminsByPublicGroupID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotAdmins, tt.wantAdmins) {
				t.Errorf("GetAdminsByPublicGroupID() gotAdmins = %v, want %v", gotAdmins, tt.wantAdmins)
			}
		})
	}
}

func TestCheckIfUserIsAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		publicGroupID uint
		userID        uint
		mock          func(userStorage *mock_users.MockUserStorage, publicGroupID, userID uint)
		wantIsAdmin   bool
		wantErr       bool
	}{
		{
			name:          "Test OK",
			publicGroupID: 1,
			userID:        1,
			mock: func(userStorage *mock_users.MockUserStorage, publicGroupID, userID uint) {
				userStorage.EXPECT().CheckIfUserIsAdmin(gomock.Any(), publicGroupID, userID).Return(true, nil)
			},
			wantIsAdmin: true,
			wantErr:     false,
		},
		{
			name:          "Test Error",
			publicGroupID: 0,
			userID:        0,
			mock: func(userStorage *mock_users.MockUserStorage, publicGroupID, userID uint) {
				userStorage.EXPECT().CheckIfUserIsAdmin(gomock.Any(), publicGroupID, userID).Return(
					false, errors.ErrNotFound,
				)
			},
			wantIsAdmin: false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userStorage := mock_users.NewMockUserStorage(ctrl)

			s := user.NewUserService(userStorage, nil)

			tt.mock(userStorage, tt.publicGroupID, tt.userID)

			gotIsAdmin, err := s.CheckIfUserIsAdmin(context.Background(), tt.publicGroupID, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CheckIfUserIsAdmin() error = %v, wantErr %v", err, tt.wantErr)
			}

			if gotIsAdmin != tt.wantIsAdmin {
				t.Errorf("CheckIfUserIsAdmin() gotIsAdmin = %v, want %v", gotIsAdmin, tt.wantIsAdmin)
			}
		})
	}
}
