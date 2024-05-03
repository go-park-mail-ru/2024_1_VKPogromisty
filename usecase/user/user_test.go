package user_test

import (
	"context"
	"reflect"
	"socio/domain"
	"socio/errors"
	customtime "socio/pkg/time"
	"socio/usecase/user"
	"testing"
	"time"

	mock_user "socio/mocks/usecase/user"

	"github.com/golang/mock/gomock"
)

func TestGetUserByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		userID   uint
		mock     func(userStorage *mock_user.MockUserStorage, userID uint)
		wantUser *domain.User
		wantErr  bool
	}{
		{
			name:   "Test OK",
			userID: 1,
			mock: func(userStorage *mock_user.MockUserStorage, userID uint) {
				userStorage.EXPECT().GetUserByID(gomock.Any(), userID).Return(&domain.User{ID: 1}, nil)
			},
			wantUser: &domain.User{ID: 1},
			wantErr:  false,
		},
		{
			name:   "Test Error",
			userID: 0,
			mock: func(userStorage *mock_user.MockUserStorage, userID uint) {
				userStorage.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, errors.ErrInternal)
			},
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userStorage := mock_user.NewMockUserStorage(ctrl)
			avatarStorage := mock_user.NewMockAvatarStorage(ctrl)

			s := user.NewUserService(userStorage, avatarStorage)

			tt.mock(userStorage, tt.userID)

			gotUser, err := s.GetUserByID(context.Background(), tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("GetUserByID() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		email    string
		mock     func(userStorage *mock_user.MockUserStorage, email string)
		wantUser *domain.User
		wantErr  bool
	}{
		{
			name:  "Test OK",
			email: "test@example.com",
			mock: func(userStorage *mock_user.MockUserStorage, email string) {
				userStorage.EXPECT().GetUserByEmail(gomock.Any(), email).Return(&domain.User{ID: 1, Email: "test@example.com"}, nil)
			},
			wantUser: &domain.User{ID: 1, Email: "test@example.com"},
			wantErr:  false,
		},
		{
			name:  "Test Error",
			email: "error@example.com",
			mock: func(userStorage *mock_user.MockUserStorage, email string) {
				userStorage.EXPECT().GetUserByEmail(gomock.Any(), email).Return(nil, errors.ErrInternal)
			},
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userStorage := mock_user.NewMockUserStorage(ctrl)
			avatarStorage := mock_user.NewMockAvatarStorage(ctrl)

			s := user.NewUserService(userStorage, avatarStorage)

			tt.mock(userStorage, tt.email)

			gotUser, err := s.GetUserByEmail(context.Background(), tt.email)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("GetUserByEmail() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestUploadAvatar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		fileName    string
		filePath    string
		contentType string
		mock        func(avatarStorage *mock_user.MockAvatarStorage, fileName, filePath, contentType string)
		wantErr     bool
	}{
		{
			name:        "Test OK",
			fileName:    "avatar.jpg",
			filePath:    "/path/to/avatar.jpg",
			contentType: "image/jpeg",
			mock: func(avatarStorage *mock_user.MockAvatarStorage, fileName, filePath, contentType string) {
				avatarStorage.EXPECT().Store(fileName, filePath, contentType).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "Test Error",
			fileName:    "avatar.jpg",
			filePath:    "/path/to/avatar.jpg",
			contentType: "image/jpeg",
			mock: func(avatarStorage *mock_user.MockAvatarStorage, fileName, filePath, contentType string) {
				avatarStorage.EXPECT().Store(fileName, filePath, contentType).Return(errors.ErrInternal)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userStorage := mock_user.NewMockUserStorage(ctrl)
			avatarStorage := mock_user.NewMockAvatarStorage(ctrl)

			s := user.NewUserService(userStorage, avatarStorage)

			tt.mock(avatarStorage, tt.fileName, tt.filePath, tt.contentType)

			err := s.UploadAvatar(tt.fileName, tt.filePath, tt.contentType)

			if (err != nil) != tt.wantErr {
				t.Errorf("UploadAvatar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		userInput user.CreateUserInput
		mock      func(userStorage *mock_user.MockUserStorage, userInput user.CreateUserInput)
		wantUser  *domain.User
		wantErr   bool
	}{
		{
			name: "Test OK",
			userInput: user.CreateUserInput{
				FirstName:      "John",
				LastName:       "Doe",
				Password:       "password",
				RepeatPassword: "password",
				Email:          "john.doe@example.com",
				Avatar:         "avatar.jpg",
				DateOfBirth:    "2000-01-01",
			},
			mock: func(userStorage *mock_user.MockUserStorage, userInput user.CreateUserInput) {
				dateOfBirth, _ := time.Parse(customtime.DateFormat, userInput.DateOfBirth)
				user := &domain.User{
					FirstName: userInput.FirstName,
					LastName:  userInput.LastName,
					Password:  userInput.Password,
					Email:     userInput.Email,
					Avatar:    userInput.Avatar,
					DateOfBirth: customtime.CustomTime{
						Time: dateOfBirth,
					},
				}
				userStorage.EXPECT().GetUserByEmail(gomock.Any(), userInput.Email).Return(nil, errors.ErrNotFound)
				userStorage.EXPECT().StoreUser(gomock.Any(), user).Return(nil)
			},
			wantUser: &domain.User{
				FirstName: "John",
				LastName:  "Doe",
				Password:  "password",
				Email:     "john.doe@example.com",
				Avatar:    "avatar.jpg",
				DateOfBirth: customtime.CustomTime{
					Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name: "Test Error",
			userInput: user.CreateUserInput{
				FirstName:   "John",
				LastName:    "Doe",
				Password:    "password",
				Email:       "john.doe@example.com",
				Avatar:      "avatar.jpg",
				DateOfBirth: "2000-01-01",
			},
			mock: func(userStorage *mock_user.MockUserStorage, userInput user.CreateUserInput) {
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name: "Test Error not found",
			userInput: user.CreateUserInput{
				FirstName:      "John",
				LastName:       "Doe",
				Password:       "password",
				RepeatPassword: "password",
				Email:          "john.doe@example.com",
				Avatar:         "avatar.jpg",
				DateOfBirth:    "2000-01-01",
			},
			mock: func(userStorage *mock_user.MockUserStorage, userInput user.CreateUserInput) {
				dateOfBirth, _ := time.Parse(customtime.DateFormat, userInput.DateOfBirth)
				user := &domain.User{
					FirstName: userInput.FirstName,
					LastName:  userInput.LastName,
					Password:  userInput.Password,
					Email:     userInput.Email,
					Avatar:    userInput.Avatar,
					DateOfBirth: customtime.CustomTime{
						Time: dateOfBirth,
					},
				}
				userStorage.EXPECT().GetUserByEmail(gomock.Any(), userInput.Email).Return(nil, errors.ErrNotFound)
				userStorage.EXPECT().StoreUser(gomock.Any(), user).Return(errors.ErrInternal)
			},
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userStorage := mock_user.NewMockUserStorage(ctrl)
			avatarStorage := mock_user.NewMockAvatarStorage(ctrl)

			s := user.NewUserService(userStorage, avatarStorage)

			tt.mock(userStorage, tt.userInput)

			gotUser, err := s.CreateUser(context.Background(), tt.userInput)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("CreateUser() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestGetUserByIDWithSubsInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		userID           uint
		authorizedUserID uint
		mock             func(userStorage *mock_user.MockUserStorage, userID uint, authorizedUserID uint)
		wantUserWithInfo user.UserWithSubsInfo
		wantErr          bool
	}{
		{
			name:             "Test OK",
			userID:           1,
			authorizedUserID: 2,
			mock: func(userStorage *mock_user.MockUserStorage, userID uint, authorizedUserID uint) {
				user := &domain.User{ID: 1}
				userStorage.EXPECT().GetUserByIDWithSubsInfo(gomock.Any(), userID, authorizedUserID).Return(user, false, true, nil)
			},
			wantUserWithInfo: user.UserWithSubsInfo{
				User:           &domain.User{ID: 1},
				IsSubscriber:   true,
				IsSubscribedTo: false,
			},
			wantErr: false,
		},
		{
			name:             "Test Error",
			userID:           0,
			authorizedUserID: 0,
			mock: func(userStorage *mock_user.MockUserStorage, userID uint, authorizedUserID uint) {
				userStorage.EXPECT().GetUserByIDWithSubsInfo(gomock.Any(), userID, authorizedUserID).Return(nil, false, false, errors.ErrInternal)
			},
			wantUserWithInfo: user.UserWithSubsInfo{},
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userStorage := mock_user.NewMockUserStorage(ctrl)

			s := user.NewUserService(userStorage, nil)

			tt.mock(userStorage, tt.userID, tt.authorizedUserID)

			gotUserWithInfo, err := s.GetUserByIDWithSubsInfo(context.Background(), tt.userID, tt.authorizedUserID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByIDWithSubsInfo() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotUserWithInfo, tt.wantUserWithInfo) {
				t.Errorf("GetUserByIDWithSubsInfo() gotUserWithInfo = %v, want %v", gotUserWithInfo, tt.wantUserWithInfo)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    user.UpdateUserInput
		mock     func(userStorage *mock_user.MockUserStorage, avatarStorage *mock_user.MockAvatarStorage, input user.UpdateUserInput)
		wantUser *domain.User
		wantErr  bool
	}{
		{
			name: "Test OK",
			input: user.UpdateUserInput{
				ID:             1,
				FirstName:      "John",
				LastName:       "Doe",
				Password:       "password",
				RepeatPassword: "password",
				Email:          "john.doe@example.com",
				Avatar:         "avatar.jpg",
				DateOfBirth:    "2000-01-01",
			},
			mock: func(userStorage *mock_user.MockUserStorage, avatarStorage *mock_user.MockAvatarStorage, input user.UpdateUserInput) {
				oldUser := &domain.User{ID: 1, FirstName: "Old", LastName: "User", Password: "oldpassword", Email: "old.user@example.com", Avatar: "oldavatar.jpg"}
				userStorage.EXPECT().GetUserByID(gomock.Any(), input.ID).Return(oldUser, nil)
				userStorage.EXPECT().GetUserByEmail(gomock.Any(), input.Email).Return(nil, errors.ErrNotFound)
				avatarStorage.EXPECT().Delete(oldUser.Avatar).Return(nil)
				userStorage.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), oldUser.Password).Return(&domain.User{ID: 1, FirstName: "John", LastName: "Doe", Password: "password", Email: "john.doe@example.com", Avatar: "avatar.jpg"}, nil)
			},
			wantUser: &domain.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Password:  "password",
				Email:     "john.doe@example.com",
				Avatar:    "avatar.jpg",
			},
			wantErr: false,
		},
		{
			name: "Test Error",
			input: user.UpdateUserInput{
				ID:          1,
				FirstName:   "John",
				LastName:    "Doe",
				Password:    "password",
				Email:       "john.doe@example.com",
				Avatar:      "avatar.jpg",
				DateOfBirth: "2000-01-01",
			},
			mock: func(userStorage *mock_user.MockUserStorage, avatarStorage *mock_user.MockAvatarStorage, input user.UpdateUserInput) {
				userStorage.EXPECT().GetUserByID(gomock.Any(), input.ID).Return(nil, errors.ErrInternal)
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name: "Test err not found",
			input: user.UpdateUserInput{
				ID:             1,
				FirstName:      "John",
				LastName:       "Doe",
				Password:       "password",
				RepeatPassword: "password",
				Email:          "john.doe@example.com",
				Avatar:         "avatar.jpg",
				DateOfBirth:    "2000-01-01",
			},
			mock: func(userStorage *mock_user.MockUserStorage, avatarStorage *mock_user.MockAvatarStorage, input user.UpdateUserInput) {
				userStorage.EXPECT().GetUserByID(gomock.Any(), input.ID).Return(nil, errors.ErrNotFound)
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name: "Test err invalid input",
			input: user.UpdateUserInput{
				ID:          1,
				FirstName:   "John",
				LastName:    "Doe",
				Password:    "password",
				Email:       "john.doe@example.com",
				Avatar:      "avatar.jpg",
				DateOfBirth: "2000-01-01",
			},
			mock: func(userStorage *mock_user.MockUserStorage, avatarStorage *mock_user.MockAvatarStorage, input user.UpdateUserInput) {
				oldUser := &domain.User{ID: 1, FirstName: "Old", LastName: "User", Password: "oldpassword", Email: "old.user@example.com", Avatar: "oldavatar.jpg"}
				userStorage.EXPECT().GetUserByID(gomock.Any(), input.ID).Return(oldUser, nil)
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name: "Test err deleting old avatar",
			input: user.UpdateUserInput{
				ID:             1,
				FirstName:      "John",
				LastName:       "Doe",
				Password:       "password",
				RepeatPassword: "password",
				Email:          "john.doe@example.com",
				Avatar:         "avatar.jpg",
				DateOfBirth:    "2000-01-01",
			},
			mock: func(userStorage *mock_user.MockUserStorage, avatarStorage *mock_user.MockAvatarStorage, input user.UpdateUserInput) {
				oldUser := &domain.User{ID: 1, FirstName: "Old", LastName: "User", Password: "oldpassword", Email: "old.user@example.com", Avatar: "oldavatar.jpg"}
				userStorage.EXPECT().GetUserByID(gomock.Any(), input.ID).Return(oldUser, nil)
				userStorage.EXPECT().GetUserByEmail(gomock.Any(), input.Email).Return(nil, errors.ErrNotFound)
				avatarStorage.EXPECT().Delete(oldUser.Avatar).Return(errors.ErrInternal)
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name: "Test err internal",
			input: user.UpdateUserInput{
				ID:             1,
				FirstName:      "John",
				LastName:       "Doe",
				Password:       "password",
				RepeatPassword: "password",
				Email:          "john.doe@example.com",
				Avatar:         "avatar.jpg",
				DateOfBirth:    "2000-01-01",
			},
			mock: func(userStorage *mock_user.MockUserStorage, avatarStorage *mock_user.MockAvatarStorage, input user.UpdateUserInput) {
				oldUser := &domain.User{ID: 1, FirstName: "Old", LastName: "User", Password: "oldpassword", Email: "old.user@example.com", Avatar: "oldavatar.jpg"}
				userStorage.EXPECT().GetUserByID(gomock.Any(), input.ID).Return(oldUser, nil)
				userStorage.EXPECT().GetUserByEmail(gomock.Any(), input.Email).Return(nil, errors.ErrNotFound)
				avatarStorage.EXPECT().Delete(oldUser.Avatar).Return(nil)
				userStorage.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), oldUser.Password).Return(
					nil, errors.ErrInternal,
				)
			},
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userStorage := mock_user.NewMockUserStorage(ctrl)
			avatarStorage := mock_user.NewMockAvatarStorage(ctrl)

			s := user.NewUserService(userStorage, avatarStorage)

			tt.mock(userStorage, avatarStorage, tt.input)

			gotUser, err := s.UpdateUser(context.Background(), tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("UpdateUser() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userID  uint
		mock    func(userStorage *mock_user.MockUserStorage, avatarStorage *mock_user.MockAvatarStorage, userID uint)
		wantErr bool
	}{
		{
			name:   "Test OK",
			userID: 1,
			mock: func(userStorage *mock_user.MockUserStorage, avatarStorage *mock_user.MockAvatarStorage, userID uint) {
				user := &domain.User{ID: 1, Avatar: "avatar.jpg"}
				userStorage.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
				avatarStorage.EXPECT().Delete(user.Avatar).Return(nil)
				userStorage.EXPECT().DeleteUser(gomock.Any(), userID).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Test Error",
			userID: 1,
			mock: func(userStorage *mock_user.MockUserStorage, avatarStorage *mock_user.MockAvatarStorage, userID uint) {
				userStorage.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, errors.ErrInternal)
			},
			wantErr: true,
		},
		{
			name:   "Test err deleting avatar",
			userID: 1,
			mock: func(userStorage *mock_user.MockUserStorage, avatarStorage *mock_user.MockAvatarStorage, userID uint) {
				user := &domain.User{ID: 1, Avatar: "avatar.jpg"}
				userStorage.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
				avatarStorage.EXPECT().Delete(user.Avatar).Return(errors.ErrInternal)
			},
			wantErr: true,
		},
		{
			name:   "Test err internal",
			userID: 1,
			mock: func(userStorage *mock_user.MockUserStorage, avatarStorage *mock_user.MockAvatarStorage, userID uint) {
				user := &domain.User{ID: 1, Avatar: "avatar.jpg"}
				userStorage.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
				avatarStorage.EXPECT().Delete(user.Avatar).Return(nil)
				userStorage.EXPECT().DeleteUser(gomock.Any(), userID).Return(errors.ErrInternal)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userStorage := mock_user.NewMockUserStorage(ctrl)
			avatarStorage := mock_user.NewMockAvatarStorage(ctrl)

			s := user.NewUserService(userStorage, avatarStorage)

			tt.mock(userStorage, avatarStorage, tt.userID)

			err := s.DeleteUser(context.Background(), tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearchByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		query     string
		mock      func(userStorage *mock_user.MockUserStorage, query string)
		wantUsers []*domain.User
		wantErr   bool
	}{
		{
			name:  "Test OK",
			query: "John",
			mock: func(userStorage *mock_user.MockUserStorage, query string) {
				users := []*domain.User{
					{ID: 1},
					{ID: 2},
				}
				userStorage.EXPECT().SearchByName(gomock.Any(), query).Return(users, nil)
			},
			wantUsers: []*domain.User{
				{ID: 1},
				{ID: 2},
			},
			wantErr: false,
		},
		{
			name:  "Test Error",
			query: "John",
			mock: func(userStorage *mock_user.MockUserStorage, query string) {
				userStorage.EXPECT().SearchByName(gomock.Any(), query).Return(nil, errors.ErrInternal)
			},
			wantUsers: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userStorage := mock_user.NewMockUserStorage(ctrl)

			s := user.NewUserService(userStorage, nil)

			tt.mock(userStorage, tt.query)

			gotUsers, err := s.SearchByName(context.Background(), tt.query)

			if (err != nil) != tt.wantErr {
				t.Errorf("SearchByName() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotUsers, tt.wantUsers) {
				t.Errorf("SearchByName() gotUsers = %v, want %v", gotUsers, tt.wantUsers)
			}
		})
	}
}

func TestGetSubscriptionIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		userID     uint
		mock       func(userStorage *mock_user.MockUserStorage, userID uint)
		wantSubIDs []uint
		wantErr    bool
	}{
		{
			name:   "Test OK",
			userID: 1,
			mock: func(userStorage *mock_user.MockUserStorage, userID uint) {
				user := &domain.User{ID: 1}
				userStorage.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
				subIDs := []uint{2, 3, 4}
				userStorage.EXPECT().GetSubscriptionIDs(gomock.Any(), userID).Return(subIDs, nil)
			},
			wantSubIDs: []uint{2, 3, 4},
			wantErr:    false,
		},
		{
			name:   "Test Error",
			userID: 1,
			mock: func(userStorage *mock_user.MockUserStorage, userID uint) {
				userStorage.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, errors.ErrInternal)
			},
			wantSubIDs: nil,
			wantErr:    true,
		},
		{
			name:   "Test OK",
			userID: 1,
			mock: func(userStorage *mock_user.MockUserStorage, userID uint) {
				user := &domain.User{ID: 1}
				userStorage.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
				userStorage.EXPECT().GetSubscriptionIDs(gomock.Any(), userID).Return(
					nil, errors.ErrInternal,
				)
			},
			wantSubIDs: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userStorage := mock_user.NewMockUserStorage(ctrl)

			s := user.NewUserService(userStorage, nil)

			tt.mock(userStorage, tt.userID)

			gotSubIDs, err := s.GetSubscriptionIDs(context.Background(), tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubscriptionIDs() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotSubIDs, tt.wantSubIDs) {
				t.Errorf("GetSubscriptionIDs() gotSubIDs = %v, want %v", gotSubIDs, tt.wantSubIDs)
			}
		})
	}
}
