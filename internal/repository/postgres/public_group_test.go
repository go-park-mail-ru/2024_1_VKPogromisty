package repository_test

import (
	"context"
	"reflect"
	"socio/domain"
	"socio/errors"
	repository "socio/internal/repository/postgres"
	customtime "socio/pkg/time"
	publicgroup "socio/usecase/public_group"
	"testing"

	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetPublicGroupByID(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name     string
		groupID  uint
		userID   uint
		mock     func(pool *pgxpoolmock.MockPgxIface, groupID, userID uint)
		expected *publicgroup.PublicGroupWithInfo
		err      bool
	}{
		{
			name:    "Test OK",
			groupID: 1,
			userID:  1,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupID, userID uint) {
				// Mock the getPublicGroupByIDWithInfoQuery
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(1), "Group", "Description", "avatar", tp.Now(), tp.Now(), uint(10), true),
				)
			},
			expected: &publicgroup.PublicGroupWithInfo{
				PublicGroup: &domain.PublicGroup{
					ID:               1,
					Name:             "Group",
					Description:      "Description",
					Avatar:           "avatar",
					CreatedAt:        customtime.CustomTime{Time: tp.Now()},
					UpdatedAt:        customtime.CustomTime{Time: tp.Now()},
					SubscribersCount: 10,
				},
				IsSubscribed: true,
			},
			err: false,
		},
		{
			name:    "Test ErrNotFound",
			groupID: 1,
			userID:  1,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupID, userID uint) {
				// Mock the getPublicGroupByIDWithInfoQuery
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					ErrRow{},
				)
			},
			expected: &publicgroup.PublicGroupWithInfo{
				PublicGroup:  &domain.PublicGroup{},
				IsSubscribed: false,
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			publicGroup := repository.NewPublicGroup(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.groupID, tt.userID)

			result, err := publicGroup.GetPublicGroupByID(context.Background(), tt.groupID, tt.userID)

			if (err != nil) != tt.err {
				t.Errorf("unexpected error: got %v, want %v", err, tt.err)
			}

			if !tt.err {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSearchPublicGroupsByNameWithInfo(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name     string
		query    string
		userID   uint
		mock     func(pool *pgxpoolmock.MockPgxIface, query string, userID uint)
		expected []*publicgroup.PublicGroupWithInfo
		err      bool
	}{
		{
			name:   "Test OK",
			query:  "Group",
			userID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, query string, userID uint) {
				// Mock the searchPublicGroupsByNameWithInfoQuery
				rows := pgxpoolmock.NewRows([]string{"id", "name", "description", "avatar", "created_at", "updated_at", "subscribers_count", "is_subscribed"})
				rows.AddRow(uint(1), "Group", "Description", "avatar", tp.Now(), tp.Now(), uint(10), true)
				pool.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(rows.ToPgxRows(), nil)
			},
			expected: []*publicgroup.PublicGroupWithInfo{
				{
					PublicGroup: &domain.PublicGroup{
						ID:               1,
						Name:             "Group",
						Description:      "Description",
						Avatar:           "avatar",
						CreatedAt:        customtime.CustomTime{Time: tp.Now()},
						UpdatedAt:        customtime.CustomTime{Time: tp.Now()},
						SubscribersCount: 10,
					},
					IsSubscribed: true,
				},
			},
			err: false,
		},
		{
			name:   "Test Empty",
			query:  "Group",
			userID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, query string, userID uint) {
				// Mock the searchPublicGroupsByNameWithInfoQuery
				pool.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, pgx.ErrNoRows)
			},
			expected: nil,
			err:      true,
		},
		{
			name:   "Test 3",
			query:  "Group",
			userID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, query string, userID uint) {
				// Mock the searchPublicGroupsByNameWithInfoQuery
				rows := pgxpoolmock.NewRows([]string{"err"})
				rows.AddRow(ErrRow{})
				pool.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(rows.ToPgxRows(), nil)
			},
			expected: nil,
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			publicGroup := repository.NewPublicGroup(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.query, tt.userID)

			result, err := publicGroup.SearchPublicGroupsByNameWithInfo(context.Background(), tt.query, tt.userID)

			if (err != nil) != tt.err {
				t.Errorf("unexpected error: got %v, want %v", err, tt.err)
			}

			if !tt.err {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestStorePublicGroup(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name     string
		group    *domain.PublicGroup
		mock     func(pool *pgxpoolmock.MockPgxIface, group *domain.PublicGroup)
		expected *domain.PublicGroup
		err      error
	}{
		{
			name: "Test ErrNoRows",
			group: &domain.PublicGroup{
				Name:        "Group",
				Description: "Description",
				Avatar:      "avatar",
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, group *domain.PublicGroup) {
				// Mock the storePublicGroupQuery
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					ErrRow{},
				)
			},
			expected: &domain.PublicGroup{},
			err:      pgx.ErrNoRows,
		},
		{
			name: "Test OK",
			group: &domain.PublicGroup{
				Name:        "Group",
				Description: "Description",
				Avatar:      "avatar",
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, group *domain.PublicGroup) {
				// Mock the storePublicGroupQuery
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(1), group.Name, group.Description, group.Avatar, tp.Now(), tp.Now()),
				)
			},
			expected: &domain.PublicGroup{
				ID:          1,
				Name:        "Group",
				Description: "Description",
				Avatar:      "avatar",
				CreatedAt:   customtime.CustomTime{Time: tp.Now()},
				UpdatedAt:   customtime.CustomTime{Time: tp.Now()},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			publicGroup := repository.NewPublicGroup(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.group)

			result, err := publicGroup.StorePublicGroup(context.Background(), tt.group)

			if err != tt.err {
				t.Errorf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("unexpected result: got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUpdatePublicGroup(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name     string
		group    *domain.PublicGroup
		mock     func(pool *pgxpoolmock.MockPgxIface, group *domain.PublicGroup)
		expected *domain.PublicGroup
		err      error
	}{
		{
			name: "Test OK",
			group: &domain.PublicGroup{
				ID:          1,
				Name:        "Updated Group",
				Description: "Updated Description",
				Avatar:      "updated_avatar",
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, group *domain.PublicGroup) {
				// Mock the updatePublicGroupQuery
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(group.ID, group.Name, group.Description, group.Avatar, tp.Now(), tp.Now()),
				)
			},
			expected: &domain.PublicGroup{
				ID:          1,
				Name:        "Updated Group",
				Description: "Updated Description",
				Avatar:      "updated_avatar",
				CreatedAt:   customtime.CustomTime{Time: tp.Now()},
				UpdatedAt:   customtime.CustomTime{Time: tp.Now()},
			},
			err: nil,
		},
		{
			name: "Test ErrNotFound",
			group: &domain.PublicGroup{
				ID:          1,
				Name:        "Updated Group",
				Description: "Updated Description",
				Avatar:      "updated_avatar",
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, group *domain.PublicGroup) {
				// Mock the updatePublicGroupQuery
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					ErrRow{},
				)
			},
			expected: &domain.PublicGroup{},
			err:      errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			publicGroup := repository.NewPublicGroup(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.group)

			result, err := publicGroup.UpdatePublicGroup(context.Background(), tt.group)

			if err != tt.err {
				t.Errorf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("unexpected result: got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDeletePublicGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		groupID uint
		mock    func(pool *pgxpoolmock.MockPgxIface, groupID uint)
		err     error
	}{
		{
			name:    "Test OK",
			groupID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupID uint) {
				// Mock the deletePublicGroupQuery
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(
					pgconn.CommandTag("DELETE 1"),
					nil,
				)
			},
			err: nil,
		},
		{
			name:    "Test Not Found",
			groupID: 2,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupID uint) {
				// Mock the deletePublicGroupQuery
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(
					pgconn.CommandTag("DELETE 0"),
					nil,
				)
			},
			err: nil,
		},
		{
			name:    "Test Rows Affected",
			groupID: 3,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupID uint) {
				// Mock the deletePublicGroupQuery
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(
					pgconn.CommandTag("DELETE 2"),
					nil,
				)
			},
			err: errors.ErrRowsAffected,
		},
		{
			name:    "Test ErrInternal",
			groupID: 3,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupID uint) {
				// Mock the deletePublicGroupQuery
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(
					nil,
					errors.ErrInternal,
				)
			},
			err: errors.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			publicGroup := repository.NewPublicGroup(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.groupID)

			err := publicGroup.DeletePublicGroup(context.Background(), tt.groupID)

			if err != tt.err {
				t.Errorf("unexpected error: got %v, want %v", err, tt.err)
			}
		})
	}
}

func TestGetSubscriptionByPublicGroupIDAndSubscriberID(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name          string
		publicGroupID uint
		subscriberID  uint
		mock          func(pool *pgxpoolmock.MockPgxIface, publicGroupID, subscriberID uint)
		expected      *domain.PublicGroupSubscription
		err           error
	}{
		{
			name:          "Test OK",
			publicGroupID: 1,
			subscriberID:  1,
			mock: func(pool *pgxpoolmock.MockPgxIface, publicGroupID, subscriberID uint) {
				// Mock the getSubscriptionByPublicGroupIDAndSubscriberIDQuery
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(1), uint(1), uint(1), tp.Now(), tp.Now()),
				)
			},
			expected: &domain.PublicGroupSubscription{
				ID:            1,
				PublicGroupID: 1,
				SubscriberID:  1,
				CreatedAt:     customtime.CustomTime{Time: tp.Now()},
				UpdatedAt:     customtime.CustomTime{Time: tp.Now()},
			},
			err: nil,
		},
		{
			name:          "Test Not Found",
			publicGroupID: 2,
			subscriberID:  2,
			mock: func(pool *pgxpoolmock.MockPgxIface, publicGroupID, subscriberID uint) {
				// Mock the getSubscriptionByPublicGroupIDAndSubscriberIDQuery
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					ErrRow{},
				)
			},
			expected: &domain.PublicGroupSubscription{},
			err:      errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			publicGroup := repository.NewPublicGroup(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.publicGroupID, tt.subscriberID)

			result, err := publicGroup.GetSubscriptionByPublicGroupIDAndSubscriberID(context.Background(), tt.publicGroupID, tt.subscriberID)

			if err != tt.err {
				t.Errorf("unexpected error: got %v, want %v", err, tt.err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("unexpected result: got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetPublicGroupsBySubscriberID(t *testing.T) {
	tp := customtime.MockTimeProvider{}
	tests := []struct {
		name         string
		subscriberID uint
		mock         func(pool *pgxpoolmock.MockPgxIface, subscriberID uint)
		expected     []*domain.PublicGroup
		err          bool
	}{
		{
			name:         "Test OK",
			subscriberID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, subscriberID uint) {
				rows := pgxpoolmock.NewRows([]string{"id", "name", "description", "avatar", "created_at", "updated_at", "subscribers_count"})
				rows.AddRow(uint(1), "Group 1", "Description 1", "Avatar 1", tp.Now(), tp.Now(), uint(10))
				pool.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(rows.ToPgxRows(), nil)
			},
			expected: []*domain.PublicGroup{
				{
					ID:               1,
					Name:             "Group 1",
					Description:      "Description 1",
					Avatar:           "Avatar 1",
					CreatedAt:        customtime.CustomTime{Time: tp.Now()},
					UpdatedAt:        customtime.CustomTime{Time: tp.Now()},
					SubscribersCount: 10,
				},
			},
			err: false,
		},
		{
			name:         "Test ErrNoRows",
			subscriberID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, subscriberID uint) {
				pool.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(nil, pgx.ErrNoRows)
			},
			expected: nil,
			err:      true,
		},
		{
			name:         "Test 3",
			subscriberID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, subscriberID uint) {
				rows := pgxpoolmock.NewRows([]string{"err"})
				rows.AddRow(ErrRow{})
				pool.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(rows.ToPgxRows(), nil)
			},
			expected: nil,
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			publicGroup := repository.NewPublicGroup(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.subscriberID)

			result, err := publicGroup.GetPublicGroupsBySubscriberID(context.Background(), tt.subscriberID)

			if (err != nil) != tt.err {
				t.Errorf("unexpected error: got %v, want %v", err, tt.err)
			}

			if !tt.err {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestStorePublicGroupSubscription(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name         string
		subscription *domain.PublicGroupSubscription
		mock         func(pool *pgxpoolmock.MockPgxIface, subscription *domain.PublicGroupSubscription)
		expected     *domain.PublicGroupSubscription
		err          error
	}{
		{
			name: "Test OK",
			subscription: &domain.PublicGroupSubscription{
				PublicGroupID: 1,
				SubscriberID:  1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, subscription *domain.PublicGroupSubscription) {
				// Mock the storePublicGroupSubscriptionQuery
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(1), subscription.PublicGroupID, subscription.SubscriberID, tp.Now(), tp.Now()),
				)
			},
			expected: &domain.PublicGroupSubscription{
				ID:            1,
				PublicGroupID: 1,
				SubscriberID:  1,
				CreatedAt:     customtime.CustomTime{Time: tp.Now()},
				UpdatedAt:     customtime.CustomTime{Time: tp.Now()},
			},
			err: nil,
		},
		{
			name: "Test ErrNoRows",
			subscription: &domain.PublicGroupSubscription{
				PublicGroupID: 1,
				SubscriberID:  1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, subscription *domain.PublicGroupSubscription) {
				// Mock the storePublicGroupSubscriptionQuery
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					ErrRow{},
				)
			},
			expected: &domain.PublicGroupSubscription{},
			err:      pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			publicGroup := repository.NewPublicGroup(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.subscription)

			result, err := publicGroup.StorePublicGroupSubscription(context.Background(), tt.subscription)

			if err != tt.err {
				t.Errorf("unexpected error: got %v, want %v", err, tt.err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("unexpected result: got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDeletePublicGroupSubscription(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		subscription *domain.PublicGroupSubscription
		mock         func(pool *pgxpoolmock.MockPgxIface, subscription *domain.PublicGroupSubscription)
		err          error
	}{
		{
			name: "Test OK",
			subscription: &domain.PublicGroupSubscription{
				PublicGroupID: 1,
				SubscriberID:  1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, subscription *domain.PublicGroupSubscription) {
				// Mock the deletePublicSubscriptionQuery
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgconn.CommandTag("DELETE 1"),
					nil,
				)
			},
			err: nil,
		},
		{
			name: "Test Empty",
			subscription: &domain.PublicGroupSubscription{
				PublicGroupID: 1,
				SubscriberID:  1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, subscription *domain.PublicGroupSubscription) {
				// Mock the deletePublicSubscriptionQuery
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgconn.CommandTag("DELETE 0"),
					nil,
				)
			},
			err: nil,
		},
		{
			name: "Test rows affected",
			subscription: &domain.PublicGroupSubscription{
				PublicGroupID: 1,
				SubscriberID:  1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, subscription *domain.PublicGroupSubscription) {
				// Mock the deletePublicSubscriptionQuery
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgconn.CommandTag("DELETE 2"),
					nil,
				)
			},
			err: errors.ErrRowsAffected,
		},
		{
			name: "Test err",
			subscription: &domain.PublicGroupSubscription{
				PublicGroupID: 1,
				SubscriberID:  1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, subscription *domain.PublicGroupSubscription) {
				// Mock the deletePublicSubscriptionQuery
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil,
					errors.ErrInternal,
				)
			},
			err: errors.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			publicGroup := repository.NewPublicGroup(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.subscription)

			err := publicGroup.DeletePublicGroupSubscription(context.Background(), tt.subscription)

			if err != tt.err {
				t.Errorf("unexpected error: got %v, want %v", err, tt.err)
			}
		})
	}
}

func TestGetPublicGroupSubscriptionIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		userID   uint
		mock     func(pool *pgxpoolmock.MockPgxIface, userID uint)
		expected []uint
		err      bool
	}{
		{
			name:   "Test OK",
			userID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint) {
				// Mock the getPublicGroupSubscriptionIDsQuery
				rows := pgxpoolmock.NewRows([]string{"id"})
				rows.AddRow(uint(1))
				rows.AddRow(uint(2))
				pool.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(rows.ToPgxRows(), nil)
			},
			expected: []uint{1, 2},
			err:      false,
		},
		{
			name:   "Test 2",
			userID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint) {
				// Mock the getPublicGroupSubscriptionIDsQuery
				rows := pgxpoolmock.NewRows([]string{"id"})
				rows.AddRow(ErrRow{})
				pool.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(rows.ToPgxRows(), nil)
			},
			expected: nil,
			err:      true,
		},
		{
			name:   "Test 3",
			userID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint) {
				pool.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
			expected: nil,
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			publicGroup := repository.NewPublicGroup(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.userID)

			result, err := publicGroup.GetPublicGroupSubscriptionIDs(context.Background(), tt.userID)

			if (err != nil) != tt.err {
				t.Errorf("unexpected error: got %v, want %v", err, tt.err)
			}

			if !tt.err {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
