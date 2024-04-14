package repository_test

import (
	"context"
	"socio/domain"
	"socio/errors"
	repository "socio/internal/repository/postgres"
	customtime "socio/pkg/time"
	"testing"
	"time"

	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

var (
	columns       = []string{"id", "sender_id", "receiver_id", "content", "created_at", "updated_at"}
	dialogColumns = []string{
		"id",
		"first_name",
		"last_name",
		"email",
		"avatar",
		"date_of_birth",
		"created_at",
		"updated_at",
		"id",
		"first_name",
		"last_name",
		"email",
		"avatar",
		"date_of_birth",
		"created_at",
		"updated_at",
		"id",
		"sender_id",
		"receiver_id",
		"content",
		"created_at",
		"updated_at",
	}
)

func TestGetLastMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		expectedID    uint
		expectedRow   *pgxpoolmock.Row
		expectedError error
	}{
		{
			name:          "TestGetLastMessage",
			expectedID:    1,
			expectedRow:   pgxpoolmock.NewRow(uint(1)),
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			timeProv := customtime.MockTimeProvider{}

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.expectedRow)

			repo := repository.NewPersonalMessages(pool, timeProv)

			lastMessageID, err := repo.GetLastMessageID(context.Background(), 1, 2)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if lastMessageID != tt.expectedID {
				t.Errorf("unexpected last message id: %d", lastMessageID)
				return
			}
		})
	}
}

func TestGetMessagesByDialog(t *testing.T) {
	t.Parallel()

	timeProv := customtime.MockTimeProvider{}

	tests := []struct {
		name          string
		expectedCount int
		expectedError error
		rows          pgx.Rows
	}{
		{
			name:          "TestGetMessagesByDialog",
			expectedCount: 2,
			expectedError: nil,
			rows: pgxpoolmock.NewRows(columns).AddRow(
				uint(1), uint(1), uint(2), "Hello", timeProv.Now(), timeProv.Now(),
			).AddRow(
				uint(2), uint(2), uint(1), "Hi", timeProv.Now().Add(time.Hour), timeProv.Now().Add(time.Hour),
			).ToPgxRows(),
		},
		{
			name:          "TestGetMessagesByDialogError",
			expectedCount: 0,
			expectedError: errors.ErrNotFound,
			rows:          pgxpoolmock.NewRows(columns).ToPgxRows(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			timeProv := customtime.MockTimeProvider{}

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			pool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.rows, tt.expectedError)

			repo := repository.NewPersonalMessages(pool, timeProv)

			messages, err := repo.GetMessagesByDialog(context.Background(), 1, 2, 0, 20)
			if err != tt.expectedError {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(messages) != tt.expectedCount {
				t.Errorf("unexpected messages count: %d", len(messages))
				return
			}
		})
	}
}

func TestUpdatePersonalMessage(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeProv := customtime.MockTimeProvider{}

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	row := pgxpoolmock.NewRow(uint(1), uint(1), uint(2), "Hello", timeProv.Now(), timeProv.Now())

	pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(row)

	repo := repository.NewPersonalMessages(pool, timeProv)

	message, err := repo.UpdateMessage(context.Background(), &domain.PersonalMessage{ID: 1, SenderID: 1, ReceiverID: 2, Content: "Hello"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if message.ID != 1 {
		t.Errorf("unexpected message id: %d", message.ID)
		return
	}
}

func TestDeletePersonalMessage(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeProv := customtime.MockTimeProvider{}

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	// Create a CommandTag that reports 1 row affected
	tag := pgconn.CommandTag("DELETE 1")

	pool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(tag, nil)

	repo := repository.NewPersonalMessages(pool, timeProv)

	err := repo.DeleteMessage(context.Background(), 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
}

func TestStorePersonalMessage(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeProv := customtime.MockTimeProvider{}

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	row := pgxpoolmock.NewRow(uint(1), uint(1), uint(2), "Hello", timeProv.Now(), timeProv.Now())

	pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(row)

	repo := repository.NewPersonalMessages(pool, timeProv)

	message, err := repo.StoreMessage(context.Background(), &domain.PersonalMessage{SenderID: 1, ReceiverID: 2, Content: "Hello"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if message.ID != 1 {
		t.Errorf("unexpected message id: %d", message.ID)
		return
	}
}

func TestGetDialogs(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeProv := customtime.MockTimeProvider{}

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	rows := pgxpoolmock.NewRows(dialogColumns).AddRow(
		uint(1),
		"John",
		"Smith",
		"email@email",
		"avatar",
		timeProv.Now(),
		timeProv.Now(),
		timeProv.Now(),
		uint(1),
		"John",
		"Smith",
		"email@email",
		"avatar",
		timeProv.Now(),
		timeProv.Now(),
		timeProv.Now(),
		uint(1),
		uint(1),
		uint(2),
		"Hello",
		timeProv.Now(),
		timeProv.Now(),
	).ToPgxRows()

	pool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(rows, nil)

	repo := repository.NewPersonalMessages(pool, timeProv)

	dialogs, err := repo.GetDialogsByUserID(context.Background(), 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if len(dialogs) != 1 {
		t.Errorf("unexpected dialogs count: %d", len(dialogs))
		return
	}
}
