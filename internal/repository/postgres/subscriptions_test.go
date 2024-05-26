package repository_test

import (
	"context"
	"socio/domain"
	"socio/errors"
	repository "socio/internal/repository/postgres"
	customtime "socio/pkg/time"
	"testing"

	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetSubscriptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name    string
		userID  uint
		want    []*domain.User
		wantErr bool
		setup   func()
	}{
		{
			name:   "test case 1",
			userID: 1,
			want: []*domain.User{
				{
					ID:          1,
					FirstName:   "User 1",
					LastName:    "User 1",
					Email:       "2@2.2",
					Avatar:      "avatar",
					DateOfBirth: customtime.CustomTime{Time: tp.Now()},
					CreatedAt:   customtime.CustomTime{Time: tp.Now()},
					UpdatedAt:   customtime.CustomTime{Time: tp.Now()},
				},
			},
			wantErr: false,
			setup: func() {
				mockDB.EXPECT().Query(context.Background(), repository.GetSubscriptionsQuery, gomock.Any()).Return(
					pgxpoolmock.NewRows([]string{
						"id", "first_name", "last_name", "email", "avatar", "date_of_birth", "created_at", "updated_at",
					}).AddRow(
						uint(1), "User 1", "User 1", "2@2.2", "avatar", tp.Now(), tp.Now(), tp.Now(),
					).ToPgxRows(), nil,
				)
			},
		},
		{
			name:    "test case 2",
			userID:  1,
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().Query(context.Background(), repository.GetSubscriptionsQuery, gomock.Any()).Return(
					pgxpoolmock.NewRows([]string{
						"err",
					}).AddRow(
						ErrRow{},
					).ToPgxRows(), nil,
				)
			},
		},
		{
			name:    "test case 3",
			userID:  1,
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().Query(context.Background(), repository.GetSubscriptionsQuery, gomock.Any()).Return(nil, pgx.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := repository.NewSubscriptions(mockDB, customtime.MockTimeProvider{})

			got, err := s.GetSubscriptions(context.Background(), tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetFriends(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name    string
		userID  uint
		want    []*domain.User
		wantErr bool
		setup   func()
	}{
		{
			name:   "test case 1 - friends retrieved successfully",
			userID: 1,
			want: []*domain.User{
				{
					ID:          1,
					FirstName:   "User 1",
					LastName:    "User 1",
					Email:       "2@2.2",
					Avatar:      "avatar",
					DateOfBirth: customtime.CustomTime{Time: tp.Now()},
					CreatedAt:   customtime.CustomTime{Time: tp.Now()},
					UpdatedAt:   customtime.CustomTime{Time: tp.Now()},
				},
			},
			wantErr: false,
			setup: func() {
				mockDB.EXPECT().Query(context.Background(), repository.GetFriendsQuery, gomock.Any()).Return(
					pgxpoolmock.NewRows([]string{
						"id", "first_name", "last_name", "email", "avatar", "date_of_birth", "created_at", "updated_at",
					}).AddRow(
						uint(1), "User 1", "User 1", "2@2.2", "avatar", tp.Now(), tp.Now(), tp.Now(),
					).ToPgxRows(), nil,
				)
			},
		},
		{
			name:    "test case 2",
			userID:  1,
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().Query(context.Background(), repository.GetFriendsQuery, gomock.Any()).Return(
					pgxpoolmock.NewRows([]string{
						"err",
					}).AddRow(
						ErrRow{},
					).ToPgxRows(), nil,
				)
			},
		},
		{
			name:    "test case 3",
			userID:  1,
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().Query(context.Background(), repository.GetFriendsQuery, gomock.Any()).Return(nil, pgx.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := repository.NewSubscriptions(mockDB, customtime.MockTimeProvider{})

			got, err := s.GetFriends(context.Background(), tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetSubscribers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name    string
		userID  uint
		want    []*domain.User
		wantErr bool
		setup   func()
	}{
		{
			name:   "test case 1 - subscribers retrieved successfully",
			userID: 1,
			want: []*domain.User{
				{
					ID:          1,
					FirstName:   "User 1",
					LastName:    "User 1",
					Email:       "2@2.2",
					Avatar:      "avatar",
					DateOfBirth: customtime.CustomTime{Time: tp.Now()},
					CreatedAt:   customtime.CustomTime{Time: tp.Now()},
					UpdatedAt:   customtime.CustomTime{Time: tp.Now()},
				},
			},
			wantErr: false,
			setup: func() {
				mockDB.EXPECT().Query(context.Background(), repository.GetSubscribersQuery, gomock.Any()).Return(
					pgxpoolmock.NewRows([]string{
						"id", "first_name", "last_name", "email", "avatar", "date_of_birth", "created_at", "updated_at",
					}).AddRow(
						uint(1), "User 1", "User 1", "2@2.2", "avatar", tp.Now(), tp.Now(), tp.Now(),
					).ToPgxRows(), nil,
				)
			},
		},
		{
			name:    "test case 2 - no subscribers found",
			userID:  1,
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().Query(context.Background(), repository.GetSubscribersQuery, gomock.Any()).Return(nil, pgx.ErrNoRows)
			},
		},
		{
			name:    "test case 3",
			userID:  1,
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().Query(context.Background(), repository.GetSubscribersQuery, gomock.Any()).Return(
					pgxpoolmock.NewRows([]string{
						"err",
					}).AddRow(
						ErrRow{},
					).ToPgxRows(), nil,
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := repository.NewSubscriptions(mockDB, tp)

			got, err := s.GetSubscribers(context.Background(), tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name     string
		sub      *domain.Subscription
		want     *domain.Subscription
		wantErr  bool
		setup    func()
		teardown func()
	}{
		{
			name: "test case 1 - subscription stored successfully",
			sub: &domain.Subscription{
				SubscriberID:   1,
				SubscribedToID: 2,
			},
			want: &domain.Subscription{
				ID:             1,
				SubscriberID:   1,
				SubscribedToID: 2,
				CreatedAt:      customtime.CustomTime{Time: tp.Now()},
				UpdatedAt:      customtime.CustomTime{Time: tp.Now()},
			},
			wantErr: false,
			setup: func() {
				mockDB.EXPECT().QueryRow(context.Background(), repository.StoreSubscriptionQuery, gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(1), uint(1), uint(2), tp.Now(), tp.Now()),
				)
			},
			teardown: func() {},
		},
		{
			name: "test case 2 - subscription not stored",
			sub: &domain.Subscription{
				SubscriberID:   1,
				SubscribedToID: 2,
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().QueryRow(context.Background(), repository.StoreSubscriptionQuery, gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := repository.NewSubscriptions(mockDB, customtime.MockTimeProvider{})

			got, err := s.Store(context.Background(), tt.sub)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tests := []struct {
		name           string
		subscriberID   uint
		subscribedToID uint
		wantErr        bool
		setup          func()
	}{
		{
			name:           "test case 1 - subscription deleted successfully",
			subscriberID:   1,
			subscribedToID: 2,
			wantErr:        false,
			setup: func() {
				mockDB.EXPECT().Exec(context.Background(), repository.DeleteSubscriptionQuery, gomock.Any(), gomock.Any()).Return(
					pgconn.CommandTag("DELETE 1"), nil,
				)
			},
		},
		{
			name:           "test case 2 - subscription not found",
			subscriberID:   1,
			subscribedToID: 2,
			wantErr:        true,
			setup: func() {
				mockDB.EXPECT().Exec(context.Background(), repository.DeleteSubscriptionQuery, gomock.Any(), gomock.Any()).Return(
					pgconn.CommandTag("DELETE 0"), nil,
				)
			},
		},
		{
			name:           "test case 3",
			subscriberID:   1,
			subscribedToID: 2,
			wantErr:        true,
			setup: func() {
				mockDB.EXPECT().Exec(context.Background(), repository.DeleteSubscriptionQuery, gomock.Any(), gomock.Any()).Return(
					pgconn.CommandTag("DELETE 2"), nil,
				)
			},
		},
		{
			name:           "test case 4",
			subscriberID:   1,
			subscribedToID: 2,
			wantErr:        true,
			setup: func() {
				mockDB.EXPECT().Exec(context.Background(), repository.DeleteSubscriptionQuery, gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal,
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := repository.NewSubscriptions(mockDB, customtime.MockTimeProvider{})

			err := s.Delete(context.Background(), tt.subscriberID, tt.subscribedToID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetBySubscriberAndSubscribedToID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name           string
		subscriberID   uint
		subscribedToID uint
		want           *domain.Subscription
		wantErr        bool
		setup          func()
	}{
		{
			name:           "test case 1 - subscription found",
			subscriberID:   1,
			subscribedToID: 2,
			want: &domain.Subscription{
				ID:             1,
				SubscriberID:   1,
				SubscribedToID: 2,
				CreatedAt:      customtime.CustomTime{Time: tp.Now()},
				UpdatedAt:      customtime.CustomTime{Time: tp.Now()},
			},
			wantErr: false,
			setup: func() {
				mockDB.EXPECT().QueryRow(context.Background(), repository.GetSubscriptionBySubscriberAndSubscribedToIDQuery, gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(1), uint(1), uint(2), tp.Now(), tp.Now()),
				)
			},
		},
		{
			name:           "test case 2 - subscription not found",
			subscriberID:   1,
			subscribedToID: 2,
			want:           nil,
			wantErr:        true,
			setup: func() {
				mockDB.EXPECT().QueryRow(context.Background(), repository.GetSubscriptionBySubscriberAndSubscribedToIDQuery, gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := repository.NewSubscriptions(mockDB, tp)

			got, err := s.GetBySubscriberAndSubscribedToID(context.Background(), tt.subscriberID, tt.subscribedToID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
