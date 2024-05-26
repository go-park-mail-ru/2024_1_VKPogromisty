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
	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/assert"
)

var (
	//	attachmentsArr = pgtype.TextArray{
	//		Elements: []pgtype.Text{
	//			{String: "image.jpg", Status: pgtype.Present},
	//		},
	//	}
	//
	// columns = []string{"id", "sender_id", "receiver_id", "content", "created_at", "updated_at", "sticker_id", "attachments"}
	//
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
		"sticker_id",
		"attachments",
	}
)

func TestGetMessageByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name    string
		msgID   uint
		want    *domain.PersonalMessage
		wantErr bool
		setup   func()
	}{
		{
			name:  "test case 1",
			msgID: 1,
			want: &domain.PersonalMessage{
				ID:          1,
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Test content",
				CreatedAt:   customtime.CustomTime{Time: tp.Now()},
				UpdatedAt:   customtime.CustomTime{Time: tp.Now()},
				Attachments: []string{"attachment1", "attachment2"},
				Sticker: &domain.Sticker{
					ID:        1,
					AuthorID:  1,
					Name:      "Test sticker",
					FileName:  "sticker.jpg",
					CreatedAt: customtime.CustomTime{Time: tp.Now()},
					UpdatedAt: customtime.CustomTime{Time: tp.Now()},
				},
			},
			wantErr: false,
			setup: func() {
				row := pgxpoolmock.NewRow(uint(1), uint(1), uint(2), "Test content", tp.Now(), tp.Now(), uint(1), pgtype.TextArray{Elements: []pgtype.Text{{String: "attachment1", Status: pgtype.Present}, {String: "attachment2", Status: pgtype.Present}}, Status: pgtype.Present})
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(row)
				stickerRow := pgxpoolmock.NewRow(uint(1), uint(1), "Test sticker", "sticker.jpg", tp.Now(), tp.Now())
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(stickerRow)
			},
		},
		{
			name:    "test case 2",
			msgID:   1,
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
		{
			name:    "test case 3",
			msgID:   1,
			want:    nil,
			wantErr: true,
			setup: func() {
				row := pgxpoolmock.NewRow(uint(1), uint(1), uint(2), "Test content", tp.Now(), tp.Now(), uint(1), pgtype.TextArray{Elements: []pgtype.Text{{String: "attachment1", Status: pgtype.Present}, {String: "attachment2", Status: pgtype.Present}}, Status: pgtype.Present})
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(row)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			pm := repository.NewPersonalMessages(mockDB, tp)

			got, err := pm.GetMessageByID(context.Background(), tt.msgID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetLastMessageID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name       string
		senderID   uint
		receiverID uint
		want       uint
		wantErr    bool
		setup      func()
	}{
		{
			name:       "test case 1",
			senderID:   1,
			receiverID: 2,
			want:       1,
			wantErr:    false,
			setup: func() {
				row := pgxpoolmock.NewRow(uint(1))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(row)
			},
		},
		{
			name:       "test case 2",
			senderID:   1,
			receiverID: 2,
			want:       1,
			wantErr:    true,
			setup: func() {
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			pm := repository.NewPersonalMessages(mockDB, tp)

			got, err := pm.GetLastMessageID(context.Background(), tt.senderID, tt.receiverID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetMessagesByDialog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name           string
		senderID       uint
		receiverID     uint
		lastMessageID  uint
		messagesAmount uint
		want           []*domain.PersonalMessage
		wantErr        bool
		setup          func()
	}{
		{
			name:           "test case 1",
			senderID:       1,
			receiverID:     2,
			lastMessageID:  0,
			messagesAmount: 10,
			want: []*domain.PersonalMessage{
				{
					ID:          1,
					SenderID:    1,
					ReceiverID:  2,
					Content:     "Test content",
					CreatedAt:   customtime.CustomTime{Time: tp.Now()},
					UpdatedAt:   customtime.CustomTime{Time: tp.Now()},
					Attachments: []string{"attachment1", "attachment2"},
					Sticker: &domain.Sticker{
						ID:        1,
						AuthorID:  1,
						Name:      "Test sticker",
						FileName:  "sticker.jpg",
						CreatedAt: customtime.CustomTime{Time: tp.Now()},
						UpdatedAt: customtime.CustomTime{Time: tp.Now()},
					},
				},
			},
			wantErr: false,
			setup: func() {
				rows := pgxpoolmock.NewRows([]string{"id", "sender_id", "receiver_id", "content", "created_at", "updated_at", "sticker_id", "attachments"}).
					AddRow(uint(1), uint(1), uint(2), "Test content", tp.Now(), tp.Now(), uint(1), pgtype.TextArray{Elements: []pgtype.Text{{String: "attachment1", Status: pgtype.Present}, {String: "attachment2", Status: pgtype.Present}}, Status: pgtype.Present}).ToPgxRows()
				mockDB.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(rows, nil)
				stickerRow := pgxpoolmock.NewRow(uint(1), uint(1), "Test sticker", "sticker.jpg", tp.Now(), tp.Now())
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(stickerRow)
			},
		},
		{
			name:           "test case 2",
			senderID:       1,
			receiverID:     2,
			lastMessageID:  0,
			messagesAmount: 10,
			want:           nil,
			wantErr:        true,
			setup: func() {
				rows := pgxpoolmock.NewRows([]string{"id", "sender_id", "receiver_id", "content", "created_at", "updated_at", "sticker_id", "attachments"}).
					AddRow(uint(1), uint(1), uint(2), "Test content", tp.Now(), tp.Now(), uint(1), pgtype.TextArray{Elements: []pgtype.Text{{String: "attachment1", Status: pgtype.Present}, {String: "attachment2", Status: pgtype.Present}}, Status: pgtype.Present}).ToPgxRows()
				mockDB.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(rows, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
		{
			name:           "test case 3",
			senderID:       1,
			receiverID:     2,
			lastMessageID:  0,
			messagesAmount: 10,
			want:           nil,
			wantErr:        true,
			setup: func() {
				mockDB.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
		{
			name:           "test case 4",
			senderID:       1,
			receiverID:     2,
			lastMessageID:  0,
			messagesAmount: 10,
			want:           nil,
			wantErr:        true,
			setup: func() {
				rows := pgxpoolmock.NewRows([]string{"err"}).AddRow(ErrRow{}).ToPgxRows()
				mockDB.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(rows, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			pm := repository.NewPersonalMessages(mockDB, tp)

			got, err := pm.GetMessagesByDialog(context.Background(), tt.senderID, tt.receiverID, tt.lastMessageID, tt.messagesAmount)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetDialogsByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name    string
		userID  uint
		want    []*domain.Dialog
		wantErr bool
		setup   func()
	}{
		{
			name:   "test case 1",
			userID: 1,
			want: []*domain.Dialog{
				{
					User1: &domain.User{
						ID:          1,
						FirstName:   "User1FirstName",
						LastName:    "User1LastName",
						Email:       "asd@asd.asd",
						Avatar:      "avatar.png",
						DateOfBirth: customtime.CustomTime{Time: tp.Now()},
						CreatedAt:   customtime.CustomTime{Time: tp.Now()},
						UpdatedAt:   customtime.CustomTime{Time: tp.Now()},
					},
					User2: &domain.User{
						ID:          2,
						FirstName:   "User2FirstName",
						LastName:    "User2LastName",
						Email:       "asd@asd.asd",
						Avatar:      "avatar.png",
						DateOfBirth: customtime.CustomTime{Time: tp.Now()},
						CreatedAt:   customtime.CustomTime{Time: tp.Now()},
						UpdatedAt:   customtime.CustomTime{Time: tp.Now()},
					},
					LastMessage: &domain.PersonalMessage{
						ID:          1,
						SenderID:    1,
						ReceiverID:  2,
						Content:     "Test content",
						CreatedAt:   customtime.CustomTime{Time: tp.Now()},
						UpdatedAt:   customtime.CustomTime{Time: tp.Now()},
						Attachments: []string{"attachment1", "attachment2"},
						Sticker: &domain.Sticker{
							ID:        1,
							AuthorID:  1,
							Name:      "Test sticker",
							FileName:  "sticker.jpg",
							CreatedAt: customtime.CustomTime{Time: tp.Now()},
							UpdatedAt: customtime.CustomTime{Time: tp.Now()},
						},
					},
				},
			},
			wantErr: false,
			setup: func() {
				rows := pgxpoolmock.NewRows(dialogColumns).AddRow(
					uint(1),
					"User1FirstName",
					"User1LastName",
					"asd@asd.asd",
					"avatar.png",
					tp.Now(),
					tp.Now(),
					tp.Now(),
					uint(2),
					"User2FirstName",
					"User2LastName",
					"asd@asd.asd",
					"avatar.png",
					tp.Now(),
					tp.Now(),
					tp.Now(),
					uint(1),
					uint(1),
					uint(2),
					"Test content",
					tp.Now(),
					tp.Now(),
					uint(1),
					pgtype.TextArray{Elements: []pgtype.Text{{String: "attachment1", Status: pgtype.Present}, {String: "attachment2", Status: pgtype.Present}}, Status: pgtype.Present}).ToPgxRows()
				mockDB.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(rows, nil)
				stickerRow := pgxpoolmock.NewRow(uint(1), uint(1), "Test sticker", "sticker.jpg", tp.Now(), tp.Now())
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(stickerRow)
			},
		},
		{
			name:    "test case 2",
			userID:  1,
			want:    nil,
			wantErr: true,
			setup: func() {
				rows := pgxpoolmock.NewRows(dialogColumns).AddRow(
					uint(1),
					"User1FirstName",
					"User1LastName",
					"asd@asd.asd",
					"avatar.png",
					tp.Now(),
					tp.Now(),
					tp.Now(),
					uint(2),
					"User2FirstName",
					"User2LastName",
					"asd@asd.asd",
					"avatar.png",
					tp.Now(),
					tp.Now(),
					tp.Now(),
					uint(1),
					uint(1),
					uint(2),
					"Test content",
					tp.Now(),
					tp.Now(),
					uint(1),
					pgtype.TextArray{Elements: []pgtype.Text{{String: "attachment1", Status: pgtype.Present}, {String: "attachment2", Status: pgtype.Present}}, Status: pgtype.Present}).ToPgxRows()
				mockDB.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(rows, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
		{
			name:    "test case 3",
			userID:  1,
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
		{
			name:    "test case 4",
			userID:  1,
			want:    nil,
			wantErr: true,
			setup: func() {
				rows := pgxpoolmock.NewRows([]string{"err"}).AddRow(ErrRow{}).ToPgxRows()
				mockDB.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(rows, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			pm := repository.NewPersonalMessages(mockDB, tp)

			got, err := pm.GetDialogsByUserID(context.Background(), tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestStoreMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name    string
		msg     *domain.PersonalMessage
		want    *domain.PersonalMessage
		wantErr bool
		setup   func()
	}{
		{
			name: "test case 1",
			msg: &domain.PersonalMessage{
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Test content",
				Attachments: []string{"attachment1", "attachment2"},
			},
			want: &domain.PersonalMessage{
				ID:          1,
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Test content",
				CreatedAt:   customtime.CustomTime{Time: tp.Now()},
				UpdatedAt:   customtime.CustomTime{Time: tp.Now()},
				Attachments: []string{"attachment1", "attachment2"},
			},
			wantErr: false,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), uint(2), "Test content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment1"))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment2"))
				mockDB.EXPECT().Commit(context.Background()).Return(nil)
				mockDB.EXPECT().Rollback(context.Background()).Return(nil)
			},
		},
		{
			name: "test case 2",
			msg: &domain.PersonalMessage{
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Test content",
				Attachments: []string{"attachment1", "attachment2"},
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), uint(2), "Test content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment1"))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment2"))
				mockDB.EXPECT().Commit(context.Background()).Return(nil)
				mockDB.EXPECT().Rollback(context.Background()).Return(errors.ErrInternal)
			},
		},
		{
			name: "test case 3",
			msg: &domain.PersonalMessage{
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Test content",
				Attachments: []string{"attachment1", "attachment2"},
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), uint(2), "Test content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment1"))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment2"))
				mockDB.EXPECT().Commit(context.Background()).Return(errors.ErrInternal)
			},
		},
		{
			name: "test case 4",
			msg: &domain.PersonalMessage{
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Test content",
				Attachments: []string{"attachment1"},
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), uint(2), "Test content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
		{
			name: "test case 5",
			msg: &domain.PersonalMessage{
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Test content",
				Attachments: []string{"attachment1", "attachment2"},
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
		{
			name: "test case 6",
			msg: &domain.PersonalMessage{
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Test content",
				Attachments: []string{"attachment1", "attachment2"},
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			pm := repository.NewPersonalMessages(mockDB, tp)

			got, err := pm.StoreMessage(context.Background(), tt.msg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUpdateMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name                string
		msg                 *domain.PersonalMessage
		attachmentsToDelete []string
		want                *domain.PersonalMessage
		wantErr             bool
		setup               func()
	}{
		{
			name: "test case 1",
			msg: &domain.PersonalMessage{
				ID:          1,
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Test content",
				Attachments: []string{"attachment2"},
			},
			attachmentsToDelete: []string{"attachment1"},
			want: &domain.PersonalMessage{
				ID:          1,
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Updated content",
				CreatedAt:   customtime.CustomTime{Time: tp.Now()},
				UpdatedAt:   customtime.CustomTime{Time: tp.Now()},
				Attachments: []string{"attachment2"},
			},
			wantErr: false,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), uint(2), "Updated content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment2"))
				mockDB.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag("DELETE 1"), nil)
				mockDB.EXPECT().Commit(context.Background()).Return(nil)
				mockDB.EXPECT().Rollback(context.Background()).Return(nil)
			},
		},
		{
			name: "test case 2",
			msg: &domain.PersonalMessage{
				ID:          1,
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Test content",
				Attachments: []string{"attachment2"},
			},
			attachmentsToDelete: []string{"attachment1"},
			want:                nil,
			wantErr:             true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), uint(2), "Updated content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment2"))
				mockDB.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag("DELETE 1"), nil)
				mockDB.EXPECT().Commit(context.Background()).Return(nil)
				mockDB.EXPECT().Rollback(context.Background()).Return(errors.ErrInternal)
			},
		},
		{
			name: "test case 3",
			msg: &domain.PersonalMessage{
				ID:          1,
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Test content",
				Attachments: []string{"attachment2"},
			},
			attachmentsToDelete: []string{"attachment1"},
			want:                nil,
			wantErr:             true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), uint(2), "Updated content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment2"))
				mockDB.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag("DELETE 1"), nil)
				mockDB.EXPECT().Commit(context.Background()).Return(errors.ErrInternal)
			},
		},
		{
			name: "test case 4",
			msg: &domain.PersonalMessage{
				ID:          1,
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Test content",
				Attachments: []string{"attachment2"},
			},
			attachmentsToDelete: []string{"attachment1"},
			want:                nil,
			wantErr:             true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), uint(2), "Updated content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment2"))
				mockDB.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
		{
			name: "test case 5",
			msg: &domain.PersonalMessage{
				ID:          1,
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Test content",
				Attachments: []string{"attachment2"},
			},
			attachmentsToDelete: []string{"attachment1"},
			want:                nil,
			wantErr:             true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), uint(2), "Updated content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
		{
			name: "test case 6",
			msg: &domain.PersonalMessage{
				ID:          1,
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Test content",
				Attachments: []string{"attachment2"},
			},
			attachmentsToDelete: []string{"attachment1"},
			want:                nil,
			wantErr:             true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
		{
			name: "test case 7",
			msg: &domain.PersonalMessage{
				ID:          1,
				SenderID:    1,
				ReceiverID:  2,
				Content:     "Test content",
				Attachments: []string{"attachment2"},
			},
			attachmentsToDelete: []string{"attachment1"},
			want:                nil,
			wantErr:             true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			pm := repository.NewPersonalMessages(mockDB, tp)

			got, err := pm.UpdateMessage(context.Background(), tt.msg, tt.attachmentsToDelete)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDeleteMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name      string
		messageID uint
		wantErr   bool
		setup     func()
	}{
		{
			name:      "test case 1",
			messageID: 1,
			wantErr:   false,
			setup: func() {
				tag := pgconn.CommandTag("DELETE 1")
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(tag, nil)
			},
		},
		{
			name:      "test case 2",
			messageID: 1,
			wantErr:   true,
			setup: func() {
				tag := pgconn.CommandTag("DELETE 0")
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(tag, nil)
			},
		},
		{
			name:      "test case 3",
			messageID: 1,
			wantErr:   true,
			setup: func() {
				tag := pgconn.CommandTag("DELETE 2")
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(tag, nil)
			},
		},
		{
			name:      "test case 4",
			messageID: 1,
			wantErr:   true,
			setup: func() {
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			pm := repository.NewPersonalMessages(pool, tp)

			err := pm.DeleteMessage(context.Background(), tt.messageID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
