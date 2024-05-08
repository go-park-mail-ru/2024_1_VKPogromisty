package repository_test

import (
	"context"
	"reflect"
	"socio/domain"
	"socio/errors"
	repository "socio/internal/repository/postgres"
	customtime "socio/pkg/time"
	"testing"

	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
)

func TestGetUserByID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeProv := customtime.MockTimeProvider{}

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	row := pgxpoolmock.NewRow(
		uint(1),
		"first_name",
		"last_name",
		"email",
		"hashed_password",
		"salt",
		"avatar",
		timeProv.Now(),
		timeProv.Now(),
		timeProv.Now(),
	)

	pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(row)

	repo := repository.NewUsers(pool, timeProv)

	user, err := repo.GetUserByID(context.Background(), 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if user.ID != 1 {
		t.Errorf("unexpected user id: %d", user.ID)
	}
}

func TestGetUserByIDWithSubsInfo(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeProv := customtime.MockTimeProvider{}

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	row := pgxpoolmock.NewRow(
		uint(1),
		"first_name",
		"last_name",
		"email",
		"hashed_password",
		"salt",
		"avatar",
		timeProv.Now(),
		timeProv.Now(),
		timeProv.Now(),
		true,
		true,
	)

	pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(row)

	repo := repository.NewUsers(pool, timeProv)

	user, isSubscribedTo, isSubscriber, err := repo.GetUserByIDWithSubsInfo(context.Background(), 1, 2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if user.ID != 1 {
		t.Errorf("unexpected user id: %d", user.ID)
	}

	if !isSubscribedTo {
		t.Errorf("unexpected isSubscribedTo: %v", isSubscribedTo)
	}

	if !isSubscriber {
		t.Errorf("unexpected isSubscriber: %v", isSubscriber)
	}
}

func TestGetUserByEmail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeProv := customtime.MockTimeProvider{}

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	row := pgxpoolmock.NewRow(
		uint(1),
		"first_name",
		"last_name",
		"email",
		"hashed_password",
		"salt",
		"avatar",
		timeProv.Now(),
		timeProv.Now(),
		timeProv.Now(),
	)

	pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(row)

	repo := repository.NewUsers(pool, timeProv)

	user, err := repo.GetUserByEmail(context.Background(), "email")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if user.ID != 1 {
		t.Errorf("unexpected user id: %d", user.ID)
	}
}

func TestStoreUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeProv := customtime.MockTimeProvider{}

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	row := pgxpoolmock.NewRow(
		uint(1),
		"first_name",
		"last_name",
		"email",
		"avatar",
		timeProv.Now(),
		timeProv.Now(),
		timeProv.Now(),
	)

	pool.EXPECT().QueryRow(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(row)

	repo := repository.NewUsers(pool, timeProv)

	err := repo.StoreUser(context.Background(), &domain.User{
		FirstName: "first_name",
		LastName:  "last_name",
		Email:     "email",
		Avatar:    "avatar",
		Password:  "password",
		Salt:      "salt",
		DateOfBirth: customtime.CustomTime{
			Time: timeProv.Now(),
		},
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeProv := customtime.MockTimeProvider{}

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	row := pgxpoolmock.NewRow(
		uint(1),
		"first_name",
		"last_name",
		"email",
		"avatar",
		timeProv.Now(),
		timeProv.Now(),
		timeProv.Now(),
	)

	pool.EXPECT().QueryRow(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(row)

	repo := repository.NewUsers(pool, timeProv)

	user, err := repo.UpdateUser(context.Background(), &domain.User{
		ID:        1,
		FirstName: "first_name",
		LastName:  "last_name",
		Email:     "email",
		Avatar:    "avatar",
		Password:  "password",
		Salt:      "salt",
		DateOfBirth: customtime.CustomTime{
			Time: timeProv.Now(),
		},
	}, "prev_pass")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if user.ID != 1 {
		t.Errorf("unexpected user id: %d", user.ID)
	}
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	tag := pgconn.CommandTag("DELETE 1")

	pool.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(pool, nil)
	pool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(tag, nil).AnyTimes()
	pool.EXPECT().Rollback(gomock.Any()).Return(nil)
	pool.EXPECT().Commit(gomock.Any()).Return(nil)

	repo := repository.NewUsers(pool, customtime.MockTimeProvider{})

	err := repo.DeleteUser(context.Background(), 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDeleteUserErr(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	pool.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(pool, nil)
	pool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound).AnyTimes()

	repo := repository.NewUsers(pool, customtime.MockTimeProvider{})

	err := repo.DeleteUser(context.Background(), 1)
	if err != errors.ErrNotFound {
		t.Errorf("expected error")
	}
}

func TestSearchByName(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name     string
		query    string
		mock     func(pool *pgxpoolmock.MockPgxIface, query string)
		expected []*domain.User
		err      error
	}{
		{
			name:  "Test OK",
			query: "John",
			mock: func(pool *pgxpoolmock.MockPgxIface, query string) {
				// Mock the getUsersByNameQuery
				rows := pgxpoolmock.NewRows([]string{"id", "first_name", "last_name", "email", "password", "salt", "avatar", "date_of_birth", "created_at", "updated_at"})
				rows.AddRow(uint(1), "John", "Doe", "john.doe@example.com", "password", "salt", "avatar", tp.Now(), tp.Now(), tp.Now())
				pool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(rows.ToPgxRows(), nil)
			},
			expected: []*domain.User{
				{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
					Password:  "password",
					Salt:      "salt",
					Avatar:    "avatar",
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
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			users := repository.NewUsers(pool, tp)

			tt.mock(pool, tt.query)

			_, err := users.SearchByName(context.Background(), tt.query)

			if err != tt.err {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetSubscriptionIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		userID   uint
		mock     func(pool *pgxpoolmock.MockPgxIface, userID uint)
		expected []uint
		err      error
	}{
		{
			name:   "Test OK",
			userID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint) {
				// Mock the getSubscriptionIDsQuery
				rows := pgxpoolmock.NewRows([]string{"subscribed_to_id"})
				rows.AddRow(uint(2))
				rows.AddRow(uint(3))
				pool.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(rows.ToPgxRows(), nil)
			},
			expected: []uint{2, 3},
			err:      nil,
		},
		{
			name:   "Test ErrInternal",
			userID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint) {
				// Mock the getSubscriptionIDsQuery
				rows := pgxpoolmock.NewRows([]string{"subscribed_to_id"})
				rows.AddRow(uint(2))
				rows.AddRow(uint(3))
				pool.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
			expected: nil,
			err:      errors.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			users := repository.NewUsers(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.userID)

			result, err := users.GetSubscriptionIDs(context.Background(), tt.userID)

			if err != tt.err {
				t.Errorf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("unexpected result: got %v, want %v", result, tt.expected)
			}
		})
	}
}
