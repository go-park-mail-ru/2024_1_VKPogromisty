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
)

var (
	userColumns = []string{"id", "first_name", "last_name", "email", "avatar", "date_of_birth", "created_at", "updated_at"}
)

func TestGetSubscriptions(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	timeProv := customtime.MockTimeProvider{}

	rows := pgxpoolmock.NewRows(userColumns).AddRow(
		uint(1),
		"first_name",
		"last_name",
		"email",
		"avatar",
		timeProv.Now(),
		timeProv.Now(),
		timeProv.Now(),
	).ToPgxRows()

	pool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(rows, nil)

	repo := repository.NewSubscriptions(pool, customtime.MockTimeProvider{})

	subscriptions, err := repo.GetSubscriptions(context.Background(), 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if subscriptions[0].ID != 1 {
		t.Errorf("unexpected subscription id: %d", subscriptions[0].ID)
	}
}

func TestGetFriends(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	timeProv := customtime.MockTimeProvider{}

	rows := pgxpoolmock.NewRows(userColumns).AddRow(
		uint(1),
		"first_name",
		"last_name",
		"email",
		"avatar",
		timeProv.Now(),
		timeProv.Now(),
		timeProv.Now(),
	).ToPgxRows()

	pool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(rows, nil)

	repo := repository.NewSubscriptions(pool, customtime.MockTimeProvider{})

	friends, err := repo.GetFriends(context.Background(), 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if friends[0].ID != 1 {
		t.Errorf("unexpected friend id: %d", friends[0].ID)
	}
}

func TestGetFriendsErrNoRows(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	pool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, pgx.ErrNoRows)

	repo := repository.NewSubscriptions(pool, customtime.MockTimeProvider{})

	_, err := repo.GetFriends(context.Background(), 1)
	if err != pgx.ErrNoRows {
		t.Errorf("unexpected error: %v", err)
	}

}

func TestGetSubscriptionsNoRows(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	rows := pgxpoolmock.NewRows(userColumns).ToPgxRows()

	pool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(rows, nil)

	repo := repository.NewSubscriptions(pool, customtime.MockTimeProvider{})

	subscriptions, err := repo.GetSubscriptions(context.Background(), 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if subscriptions != nil {
		t.Errorf("unexpected subscriptions: %v", subscriptions)
	}
}

func TestGetSubscriptionsErrNoRows(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	pool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, pgx.ErrNoRows)

	repo := repository.NewSubscriptions(pool, customtime.MockTimeProvider{})

	_, err := repo.GetSubscriptions(context.Background(), 1)
	if err != pgx.ErrNoRows {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetSubscibers(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	timeProv := customtime.MockTimeProvider{}

	rows := pgxpoolmock.NewRows(userColumns).AddRow(
		uint(1),
		"first_name",
		"last_name",
		"email",
		"avatar",
		timeProv.Now(),
		timeProv.Now(),
		timeProv.Now(),
	).ToPgxRows()

	pool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(rows, nil)

	repo := repository.NewSubscriptions(pool, customtime.MockTimeProvider{})

	subscriptions, err := repo.GetSubscribers(context.Background(), 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if subscriptions[0].ID != 1 {
		t.Errorf("unexpected subscription id: %d", subscriptions[0].ID)
	}
}

func TestGetSubscribersErrNoRows(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	pool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, pgx.ErrNoRows)

	repo := repository.NewSubscriptions(pool, customtime.MockTimeProvider{})

	_, err := repo.GetSubscribers(context.Background(), 1)
	if err != pgx.ErrNoRows {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestStore(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	timeProv := customtime.MockTimeProvider{}

	row := pgxpoolmock.NewRow(
		uint(1),
		uint(1),
		uint(2),
		timeProv.Now(),
		timeProv.Now(),
	)

	pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(row)

	repo := repository.NewSubscriptions(pool, customtime.MockTimeProvider{})

	subscription, err := repo.Store(context.Background(), &domain.Subscription{
		SubscriberID:   1,
		SubscribedToID: 2,
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if subscription.ID != 1 {
		t.Errorf("unexpected subscription id: %d", subscription.ID)
	}

	if subscription.SubscriberID != 1 {
		t.Errorf("unexpected subscriber id: %d", subscription.SubscriberID)
	}

	if subscription.SubscribedToID != 2 {
		t.Errorf("unexpected subscribed to id: %d", subscription.SubscribedToID)
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	tag := pgconn.CommandTag("DELETE 1")

	pool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(tag, nil)

	repo := repository.NewSubscriptions(pool, customtime.MockTimeProvider{})

	err := repo.Delete(context.Background(), 1, 2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDeleteErrNoRowsAffected(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	pool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

	repo := repository.NewSubscriptions(pool, customtime.MockTimeProvider{})

	err := repo.Delete(context.Background(), 1, 2)
	if err != errors.ErrInvalidBody {
		t.Errorf("unexpected error value: %v", err)
	}
}

func TestGetBySubscriberAndSubscribedToID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	timeProv := customtime.MockTimeProvider{}

	row := pgxpoolmock.NewRow(
		uint(1),
		uint(1),
		uint(2),
		timeProv.Now(),
		timeProv.Now(),
	)

	pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(row)

	repo := repository.NewSubscriptions(pool, customtime.MockTimeProvider{})

	subscription, err := repo.GetBySubscriberAndSubscribedToID(context.Background(), 1, 2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if subscription.ID != 1 {
		t.Errorf("unexpected subscription id: %d", subscription.ID)
	}

	if subscription.SubscriberID != 1 {
		t.Errorf("unexpected subscriber id: %d", subscription.SubscriberID)
	}

	if subscription.SubscribedToID != 2 {
		t.Errorf("unexpected subscribed to id: %d", subscription.SubscribedToID)
	}
}
