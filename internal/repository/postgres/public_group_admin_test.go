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

func TestStorePublicGroupAdmin(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name     string
		input    *domain.PublicGroupAdmin
		mock     func(pool *pgxpoolmock.MockPgxIface, input *domain.PublicGroupAdmin)
		expected *domain.PublicGroupAdmin
		err      bool
	}{
		{
			name: "Test OK",
			input: &domain.PublicGroupAdmin{
				PublicGroupID: 1,
				UserID:        1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, input *domain.PublicGroupAdmin) {
				// Mock the storePublicGroupAdminQuery
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), input.PublicGroupID, input.UserID, tp.Now(), tp.Now()))
			},
			expected: &domain.PublicGroupAdmin{
				ID:            1,
				PublicGroupID: 1,
				UserID:        1,
				CreatedAt: customtime.CustomTime{
					Time: tp.Now(),
				},
				UpdatedAt: customtime.CustomTime{
					Time: tp.Now(),
				},
			},
			err: false,
		},
		{
			name: "Test err",
			input: &domain.PublicGroupAdmin{
				PublicGroupID: 1,
				UserID:        1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, input *domain.PublicGroupAdmin) {
				// Mock the storePublicGroupAdminQuery
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrRow{})
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

			users := repository.NewUsers(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.input)

			result, err := users.StorePublicGroupAdmin(context.Background(), tt.input)

			if (err != nil) != tt.err {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.err {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDeletePublicGroupAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *domain.PublicGroupAdmin
		mock  func(pool *pgxpoolmock.MockPgxIface, input *domain.PublicGroupAdmin)
		err   error
	}{
		{
			name: "Test OK",
			input: &domain.PublicGroupAdmin{
				PublicGroupID: 1,
				UserID:        1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, input *domain.PublicGroupAdmin) {
				// Mock the deletePublicGroupAdminQuery
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgconn.CommandTag("DELETE 1"),
					nil,
				)
			},
			err: nil,
		},
		{
			name: "Test not found",
			input: &domain.PublicGroupAdmin{
				PublicGroupID: 1,
				UserID:        1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, input *domain.PublicGroupAdmin) {
				// Mock the deletePublicGroupAdminQuery
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgconn.CommandTag("DELETE 0"),
					nil,
				)
			},
			err: nil,
		},
		{
			name: "Test ErrInternal",
			input: &domain.PublicGroupAdmin{
				PublicGroupID: 1,
				UserID:        1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, input *domain.PublicGroupAdmin) {
				// Mock the deletePublicGroupAdminQuery
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil,
					errors.ErrInternal,
				)
			},
			err: errors.ErrInternal,
		},
		{
			name: "Test ErrRowsAffected",
			input: &domain.PublicGroupAdmin{
				PublicGroupID: 1,
				UserID:        1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, input *domain.PublicGroupAdmin) {
				// Mock the deletePublicGroupAdminQuery
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgconn.CommandTag("DELETE 2"),
					nil,
				)
			},
			err: errors.ErrRowsAffected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			users := repository.NewUsers(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.input)

			err := users.DeletePublicGroupAdmin(context.Background(), tt.input)

			if err != tt.err {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetAdminsByPublicGroupID(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name     string
		groupID  uint
		mock     func(pool *pgxpoolmock.MockPgxIface, groupID uint)
		expected []*domain.User
		err      bool
	}{
		{
			name:    "Test OK",
			groupID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupID uint) {
				// Mock the getAdminsByPublicGroupIDQuery
				rows := pgxpoolmock.NewRows([]string{"id", "first_name", "last_name", "email", "avatar", "password", "salt", "date_of_birth", "created_at", "updated_at"})
				rows.AddRow(uint(1), "John", "Doe", "john.doe@example.com", "avatar", "password", "salt", tp.Now(), tp.Now(), tp.Now())
				pool.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(rows.ToPgxRows(), nil)
			},
			expected: []*domain.User{
				{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
					Avatar:    "avatar",
					Password:  "password",
					Salt:      "salt",
					DateOfBirth: customtime.CustomTime{
						Time: tp.Now(),
					},
					CreatedAt: customtime.CustomTime{
						Time: tp.Now(),
					},
					UpdatedAt: customtime.CustomTime{
						Time: tp.Now(),
					},
				},
			},
			err: false,
		},
		{
			name:    "Test 2",
			groupID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupID uint) {
				// Mock the getAdminsByPublicGroupIDQuery
				rows := pgxpoolmock.NewRows([]string{"err"})
				rows.AddRow(ErrRow{})
				pool.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(rows.ToPgxRows(), nil)
			},
			expected: nil,
			err:      true,
		},
		{
			name:    "Test 3",
			groupID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupID uint) {
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

			users := repository.NewUsers(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.groupID)

			result, err := users.GetAdminsByPublicGroupID(context.Background(), tt.groupID)

			if (err != nil) != tt.err {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.err {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
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
		mock          func(pool *pgxpoolmock.MockPgxIface, publicGroupID, userID uint)
		expected      bool
		err           error
	}{
		{
			name:          "Test OK",
			publicGroupID: 1,
			userID:        1,
			mock: func(pool *pgxpoolmock.MockPgxIface, publicGroupID, userID uint) {
				// Mock the checkIfUserIsAdminQuery
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(true),
				)
			},
			expected: true,
			err:      nil,
		},
		{
			name:          "Test ErrNoRows",
			publicGroupID: 1,
			userID:        1,
			mock: func(pool *pgxpoolmock.MockPgxIface, publicGroupID, userID uint) {
				// Mock the checkIfUserIsAdminQuery
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					ErrRow{},
				)
			},
			expected: false,
			err:      pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			users := repository.NewUsers(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.publicGroupID, tt.userID)

			result, err := users.CheckIfUserIsAdmin(context.Background(), tt.publicGroupID, tt.userID)

			if err != tt.err {
				t.Errorf("unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("unexpected result: got %v, want %v", result, tt.expected)
			}
		})
	}
}
