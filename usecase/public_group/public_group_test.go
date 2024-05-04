package publicgroup_test

import (
	"context"
	"reflect"
	"socio/domain"
	"socio/errors"
	mock_publicgroup "socio/mocks/usecase/public_group"
	"socio/pkg/static"
	publicgroup "socio/usecase/public_group"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		groupID   uint
		userID    uint
		mock      func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, groupID, userID uint)
		wantGroup *publicgroup.PublicGroupWithInfo
		wantErr   bool
	}{
		{
			name:    "Test OK",
			groupID: 1,
			userID:  1,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, groupID, userID uint) {
				publicGroupStorage.EXPECT().GetPublicGroupByID(gomock.Any(), groupID, userID).Return(&publicgroup.PublicGroupWithInfo{PublicGroup: &domain.PublicGroup{ID: 1}}, nil)
			},
			wantGroup: &publicgroup.PublicGroupWithInfo{PublicGroup: &domain.PublicGroup{ID: 1}},
			wantErr:   false,
		},
		{
			name:    "Test Error",
			groupID: 1,
			userID:  1,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, groupID, userID uint) {
				publicGroupStorage.EXPECT().GetPublicGroupByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
			wantGroup: nil,
			wantErr:   true,
		},
		{
			name:    "Test err invalid data",
			groupID: 0,
			userID:  1,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, groupID, userID uint) {
			},
			wantGroup: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			publicGroupStorage := mock_publicgroup.NewMockPublicGroupStorage(ctrl)

			s := publicgroup.NewService(publicGroupStorage, nil)

			tt.mock(publicGroupStorage, tt.groupID, tt.userID)

			gotGroup, err := s.GetByID(context.Background(), tt.groupID, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotGroup, tt.wantGroup) {
				t.Errorf("GetByID() gotGroup = %v, want %v", gotGroup, tt.wantGroup)
			}
		})
	}
}

func TestSearchByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		query      string
		userID     uint
		mock       func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, query string, userID uint)
		wantGroups []*publicgroup.PublicGroupWithInfo
		wantErr    bool
	}{
		{
			name:   "Test OK",
			query:  "test",
			userID: 1,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, query string, userID uint) {
				publicGroupStorage.EXPECT().SearchPublicGroupsByNameWithInfo(gomock.Any(), query, userID).Return([]*publicgroup.PublicGroupWithInfo{{PublicGroup: &domain.PublicGroup{ID: 1}}}, nil)
			},
			wantGroups: []*publicgroup.PublicGroupWithInfo{{PublicGroup: &domain.PublicGroup{ID: 1}}},
			wantErr:    false,
		},
		{
			name:   "Test Error",
			query:  "test",
			userID: 1,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, query string, userID uint) {
				publicGroupStorage.EXPECT().SearchPublicGroupsByNameWithInfo(gomock.Any(), query, userID).Return(nil, errors.ErrInternal)
			},
			wantGroups: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			publicGroupStorage := mock_publicgroup.NewMockPublicGroupStorage(ctrl)

			s := publicgroup.NewService(publicGroupStorage, nil)

			tt.mock(publicGroupStorage, tt.query, tt.userID)

			gotGroups, err := s.SearchByName(context.Background(), tt.query, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("SearchByName() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotGroups, tt.wantGroups) {
				t.Errorf("SearchByName() gotGroups = %v, want %v", gotGroups, tt.wantGroups)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		publicGroup *domain.PublicGroup
		mock        func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, publicGroup *domain.PublicGroup)
		wantGroup   *domain.PublicGroup
		wantErr     bool
	}{
		{
			name:        "Test OK",
			publicGroup: &domain.PublicGroup{Name: "test", Avatar: ""},
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, publicGroup *domain.PublicGroup) {
				publicGroupStorage.EXPECT().StorePublicGroup(gomock.Any(), publicGroup).Return(&domain.PublicGroup{ID: 1, Name: "test", Avatar: static.DefaultGroupAvatarFileName}, nil)
			},
			wantGroup: &domain.PublicGroup{ID: 1, Name: "test", Avatar: static.DefaultGroupAvatarFileName},
			wantErr:   false,
		},
		{
			name:        "Test Error",
			publicGroup: &domain.PublicGroup{Name: "", Avatar: "avatar"},
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, publicGroup *domain.PublicGroup) {
			},
			wantGroup: nil,
			wantErr:   true,
		},
		{
			name:        "Test err internal",
			publicGroup: &domain.PublicGroup{Name: "test", Avatar: ""},
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, publicGroup *domain.PublicGroup) {
				publicGroupStorage.EXPECT().StorePublicGroup(gomock.Any(), publicGroup).Return(nil, errors.ErrInternal)
			},
			wantGroup: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			publicGroupStorage := mock_publicgroup.NewMockPublicGroupStorage(ctrl)

			s := publicgroup.NewService(publicGroupStorage, nil)

			tt.mock(publicGroupStorage, tt.publicGroup)

			gotGroup, err := s.Create(context.Background(), tt.publicGroup)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotGroup, tt.wantGroup) {
				t.Errorf("Create() gotGroup = %v, want %v", gotGroup, tt.wantGroup)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		publicGroup *domain.PublicGroup
		mock        func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, avatarStorage *mock_publicgroup.MockAvatarStorage, publicGroup *domain.PublicGroup)
		wantGroup   *domain.PublicGroup
		wantErr     bool
	}{
		{
			name:        "Test OK",
			publicGroup: &domain.PublicGroup{ID: 1, Name: "test", Avatar: "avatar"},
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, avatarStorage *mock_publicgroup.MockAvatarStorage, publicGroup *domain.PublicGroup) {
				publicGroupStorage.EXPECT().GetPublicGroupByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&publicgroup.PublicGroupWithInfo{
					PublicGroup: publicGroup,
				}, nil)
				publicGroupStorage.EXPECT().UpdatePublicGroup(gomock.Any(), publicGroup).Return(publicGroup, nil)
				avatarStorage.EXPECT().Delete(publicGroup.Avatar).Return(nil)
			},
			wantGroup: &domain.PublicGroup{ID: 1, Name: "test", Avatar: "avatar"},
			wantErr:   false,
		},
		{
			name:        "Test Error",
			publicGroup: &domain.PublicGroup{ID: 0, Name: "test", Avatar: "avatar"},
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, avatarStorage *mock_publicgroup.MockAvatarStorage, publicGroup *domain.PublicGroup) {
			},
			wantGroup: nil,
			wantErr:   true,
		},
		{
			name:        "Test err not found",
			publicGroup: &domain.PublicGroup{ID: 1, Name: "test", Avatar: "avatar"},
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, avatarStorage *mock_publicgroup.MockAvatarStorage, publicGroup *domain.PublicGroup) {
				publicGroupStorage.EXPECT().GetPublicGroupByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil,
					errors.ErrNotFound,
				)
			},
			wantGroup: nil,
			wantErr:   true,
		},
		{
			name:        "Test err internal",
			publicGroup: &domain.PublicGroup{ID: 1, Name: "test", Avatar: "avatar"},
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, avatarStorage *mock_publicgroup.MockAvatarStorage, publicGroup *domain.PublicGroup) {
				publicGroupStorage.EXPECT().GetPublicGroupByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&publicgroup.PublicGroupWithInfo{
					PublicGroup: publicGroup,
				}, nil)
				avatarStorage.EXPECT().Delete(publicGroup.Avatar).Return(nil)
				publicGroupStorage.EXPECT().UpdatePublicGroup(gomock.Any(), publicGroup).Return(nil, errors.ErrInternal)
			},
			wantGroup: nil,
			wantErr:   true,
		},
		{
			name:        "Test err delete avatar",
			publicGroup: &domain.PublicGroup{ID: 1, Name: "test", Avatar: "avatar"},
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, avatarStorage *mock_publicgroup.MockAvatarStorage, publicGroup *domain.PublicGroup) {
				publicGroupStorage.EXPECT().GetPublicGroupByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&publicgroup.PublicGroupWithInfo{
					PublicGroup: publicGroup,
				}, nil)
				avatarStorage.EXPECT().Delete(publicGroup.Avatar).Return(errors.ErrInternal)
			},
			wantGroup: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			publicGroupStorage := mock_publicgroup.NewMockPublicGroupStorage(ctrl)
			avatarStorage := mock_publicgroup.NewMockAvatarStorage(ctrl)

			s := publicgroup.NewService(publicGroupStorage, avatarStorage)

			tt.mock(publicGroupStorage, avatarStorage, tt.publicGroup)

			gotGroup, err := s.Update(context.Background(), tt.publicGroup)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotGroup, tt.wantGroup) {
				t.Errorf("Update() gotGroup = %v, want %v", gotGroup, tt.wantGroup)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		publicGroupID uint
		mock          func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, publicGroupID uint)
		wantErr       bool
	}{
		{
			name:          "Test OK",
			publicGroupID: 1,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, publicGroupID uint) {
				publicGroupStorage.EXPECT().DeletePublicGroup(gomock.Any(), publicGroupID).Return(nil)
			},
			wantErr: false,
		},
		{
			name:          "Test Error",
			publicGroupID: 0,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, publicGroupID uint) {
			},
			wantErr: true,
		},
		{
			name:          "Test err internal",
			publicGroupID: 1,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, publicGroupID uint) {
				publicGroupStorage.EXPECT().DeletePublicGroup(gomock.Any(), publicGroupID).Return(errors.ErrInternal)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			publicGroupStorage := mock_publicgroup.NewMockPublicGroupStorage(ctrl)

			s := publicgroup.NewService(publicGroupStorage, nil)

			tt.mock(publicGroupStorage, tt.publicGroupID)

			err := s.Delete(context.Background(), tt.publicGroupID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetSubscriptionByPublicGroupIDAndSubscriberID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		publicGroupID    uint
		subscriberID     uint
		mock             func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, publicGroupID uint, subscriberID uint)
		wantSubscription *domain.PublicGroupSubscription
		wantErr          bool
	}{
		{
			name:          "Test OK",
			publicGroupID: 1,
			subscriberID:  1,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, publicGroupID uint, subscriberID uint) {
				publicGroupStorage.EXPECT().GetSubscriptionByPublicGroupIDAndSubscriberID(gomock.Any(), publicGroupID, subscriberID).Return(&domain.PublicGroupSubscription{PublicGroupID: 1, SubscriberID: 1}, nil)
			},
			wantSubscription: &domain.PublicGroupSubscription{PublicGroupID: 1, SubscriberID: 1},
			wantErr:          false,
		},
		{
			name:          "Test Error",
			publicGroupID: 0,
			subscriberID:  0,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, publicGroupID uint, subscriberID uint) {
			},
			wantSubscription: nil,
			wantErr:          true,
		},
		{
			name:          "Test OK",
			publicGroupID: 1,
			subscriberID:  1,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, publicGroupID uint, subscriberID uint) {
				publicGroupStorage.EXPECT().GetSubscriptionByPublicGroupIDAndSubscriberID(gomock.Any(), publicGroupID, subscriberID).Return(
					nil, errors.ErrNotFound,
				)
			},
			wantSubscription: nil,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			publicGroupStorage := mock_publicgroup.NewMockPublicGroupStorage(ctrl)

			s := publicgroup.NewService(publicGroupStorage, nil)

			tt.mock(publicGroupStorage, tt.publicGroupID, tt.subscriberID)

			gotSubscription, err := s.GetSubscriptionByPublicGroupIDAndSubscriberID(context.Background(), tt.publicGroupID, tt.subscriberID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubscriptionByPublicGroupIDAndSubscriberID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotSubscription, tt.wantSubscription) {
				t.Errorf("GetSubscriptionByPublicGroupIDAndSubscriberID() gotSubscription = %v, want %v", gotSubscription, tt.wantSubscription)
			}
		})
	}
}

func TestGetBySubscriberID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		subscriberID uint
		mock         func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, subscriberID uint)
		wantGroups   []*domain.PublicGroup
		wantErr      bool
	}{
		{
			name:         "Test OK",
			subscriberID: 1,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, subscriberID uint) {
				groups := []*domain.PublicGroup{{ID: 1, Name: "test", Avatar: "avatar"}}
				publicGroupStorage.EXPECT().GetPublicGroupsBySubscriberID(gomock.Any(), subscriberID).Return(groups, nil)
			},
			wantGroups: []*domain.PublicGroup{{ID: 1, Name: "test", Avatar: "avatar"}},
			wantErr:    false,
		},
		{
			name:         "Test Error",
			subscriberID: 0,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, subscriberID uint) {
			},
			wantGroups: nil,
			wantErr:    true,
		},
		{
			name:         "Test err internal",
			subscriberID: 1,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, subscriberID uint) {
				publicGroupStorage.EXPECT().GetPublicGroupsBySubscriberID(gomock.Any(), subscriberID).Return(
					nil, errors.ErrInternal,
				)
			},
			wantGroups: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			publicGroupStorage := mock_publicgroup.NewMockPublicGroupStorage(ctrl)

			s := publicgroup.NewService(publicGroupStorage, nil)

			tt.mock(publicGroupStorage, tt.subscriberID)

			gotGroups, err := s.GetBySubscriberID(context.Background(), tt.subscriberID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBySubscriberID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotGroups, tt.wantGroups) {
				t.Errorf("GetBySubscriberID() gotGroups = %v, want %v", gotGroups, tt.wantGroups)
			}
		})
	}
}

func TestSubscribe(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		subscription     *domain.PublicGroupSubscription
		mock             func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, subscription *domain.PublicGroupSubscription)
		wantSubscription *domain.PublicGroupSubscription
		wantErr          bool
	}{
		{
			name:         "Test OK",
			subscription: &domain.PublicGroupSubscription{PublicGroupID: 1, SubscriberID: 1},
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, subscription *domain.PublicGroupSubscription) {
				publicGroupStorage.EXPECT().GetPublicGroupByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&publicgroup.PublicGroupWithInfo{
						PublicGroup: &domain.PublicGroup{ID: 1, Name: "test", Avatar: "avatar"},
					}, nil)
				publicGroupStorage.EXPECT().StorePublicGroupSubscription(gomock.Any(), subscription).Return(subscription, nil)
			},
			wantSubscription: &domain.PublicGroupSubscription{PublicGroupID: 1, SubscriberID: 1},
			wantErr:          false,
		},
		{
			name:         "Test Error",
			subscription: &domain.PublicGroupSubscription{PublicGroupID: 0, SubscriberID: 0},
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, subscription *domain.PublicGroupSubscription) {
			},
			wantSubscription: nil,
			wantErr:          true,
		},
		{
			name:         "Test err not found",
			subscription: &domain.PublicGroupSubscription{PublicGroupID: 1, SubscriberID: 1},
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, subscription *domain.PublicGroupSubscription) {
				publicGroupStorage.EXPECT().GetPublicGroupByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrNotFound,
				)
			},
			wantSubscription: nil,
			wantErr:          true,
		},
		{
			name:         "Test err internal",
			subscription: &domain.PublicGroupSubscription{PublicGroupID: 1, SubscriberID: 1},
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, subscription *domain.PublicGroupSubscription) {
				publicGroupStorage.EXPECT().GetPublicGroupByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					&publicgroup.PublicGroupWithInfo{
						PublicGroup: &domain.PublicGroup{ID: 1, Name: "test", Avatar: "avatar"},
					}, nil)
				publicGroupStorage.EXPECT().StorePublicGroupSubscription(gomock.Any(), subscription).Return(
					nil, errors.ErrInternal,
				)
			},
			wantSubscription: nil,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			publicGroupStorage := mock_publicgroup.NewMockPublicGroupStorage(ctrl)

			s := publicgroup.NewService(publicGroupStorage, nil)

			tt.mock(publicGroupStorage, tt.subscription)

			gotSubscription, err := s.Subscribe(context.Background(), tt.subscription)

			if (err != nil) != tt.wantErr {
				t.Errorf("Subscribe() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotSubscription, tt.wantSubscription) {
				t.Errorf("Subscribe() gotSubscription = %v, want %v", gotSubscription, tt.wantSubscription)
			}
		})
	}
}

func TestUnsubscribe(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		subscription *domain.PublicGroupSubscription
		mock         func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, subscription *domain.PublicGroupSubscription)
		wantErr      bool
	}{
		{
			name:         "Test OK",
			subscription: &domain.PublicGroupSubscription{PublicGroupID: 1, SubscriberID: 1},
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, subscription *domain.PublicGroupSubscription) {
				publicGroupStorage.EXPECT().DeletePublicGroupSubscription(gomock.Any(), subscription).Return(nil)
			},
			wantErr: false,
		},
		{
			name:         "Test Error",
			subscription: &domain.PublicGroupSubscription{PublicGroupID: 0, SubscriberID: 0},
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, subscription *domain.PublicGroupSubscription) {
			},
			wantErr: true,
		},
		{
			name:         "Test err internal",
			subscription: &domain.PublicGroupSubscription{PublicGroupID: 1, SubscriberID: 1},
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, subscription *domain.PublicGroupSubscription) {
				publicGroupStorage.EXPECT().DeletePublicGroupSubscription(gomock.Any(), subscription).Return(errors.ErrInternal)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			publicGroupStorage := mock_publicgroup.NewMockPublicGroupStorage(ctrl)

			s := publicgroup.NewService(publicGroupStorage, nil)

			tt.mock(publicGroupStorage, tt.subscription)

			err := s.Unsubscribe(context.Background(), tt.subscription)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unsubscribe() error = %v, wantErr %v", err, tt.wantErr)
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
		mock        func(avatarStorage *mock_publicgroup.MockAvatarStorage, fileName string, filePath string, contentType string)
		wantErr     bool
	}{
		{
			name:        "Test OK",
			fileName:    "avatar.jpg",
			filePath:    "/path/to/avatar.jpg",
			contentType: "image/jpeg",
			mock: func(avatarStorage *mock_publicgroup.MockAvatarStorage, fileName string, filePath string, contentType string) {
				avatarStorage.EXPECT().Store(fileName, filePath, contentType).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "Test Error",
			fileName:    "avatar.jpg",
			filePath:    "/path/to/avatar.jpg",
			contentType: "image/jpeg",
			mock: func(avatarStorage *mock_publicgroup.MockAvatarStorage, fileName string, filePath string, contentType string) {
				avatarStorage.EXPECT().Store(fileName, filePath, contentType).Return(errors.ErrInternal)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			avatarStorage := mock_publicgroup.NewMockAvatarStorage(ctrl)

			s := publicgroup.NewService(nil, avatarStorage)

			tt.mock(avatarStorage, tt.fileName, tt.filePath, tt.contentType)

			err := s.UploadAvatar(tt.fileName, tt.filePath, tt.contentType)

			if (err != nil) != tt.wantErr {
				t.Errorf("UploadAvatar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetSubscriptionIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		userID     uint
		mock       func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, userID uint)
		wantSubIDs []uint
		wantErr    bool
	}{
		{
			name:   "Test OK",
			userID: 1,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, userID uint) {
				subIDs := []uint{1, 2, 3}
				publicGroupStorage.EXPECT().GetPublicGroupSubscriptionIDs(gomock.Any(), userID).Return(subIDs, nil)
			},
			wantSubIDs: []uint{1, 2, 3},
			wantErr:    false,
		},
		{
			name:   "Test Error",
			userID: 0,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, userID uint) {
			},
			wantSubIDs: nil,
			wantErr:    true,
		},
		{
			name:   "Test err internal",
			userID: 1,
			mock: func(publicGroupStorage *mock_publicgroup.MockPublicGroupStorage, userID uint) {
				publicGroupStorage.EXPECT().GetPublicGroupSubscriptionIDs(gomock.Any(), userID).Return(
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

			publicGroupStorage := mock_publicgroup.NewMockPublicGroupStorage(ctrl)

			s := publicgroup.NewService(publicGroupStorage, nil)

			tt.mock(publicGroupStorage, tt.userID)

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
