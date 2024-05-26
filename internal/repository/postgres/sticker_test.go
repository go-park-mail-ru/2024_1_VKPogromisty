package repository_test

import (
	"context"
	"socio/domain"
	repository "socio/internal/repository/postgres"
	customtime "socio/pkg/time"
	"testing"

	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetStickersByAuthorID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name     string
		authorID uint
		want     []*domain.Sticker
		wantErr  bool
		setup    func()
	}{
		{
			name:     "test case 1 - stickers found",
			authorID: 1,
			want: []*domain.Sticker{
				{
					ID:        1,
					AuthorID:  1,
					Name:      "sticker1",
					FileName:  "sticker1.png",
					CreatedAt: customtime.CustomTime{Time: tp.Now()},
					UpdatedAt: customtime.CustomTime{Time: tp.Now()},
				},
				{
					ID:        2,
					AuthorID:  1,
					Name:      "sticker2",
					FileName:  "sticker2.png",
					CreatedAt: customtime.CustomTime{Time: tp.Now()},
					UpdatedAt: customtime.CustomTime{Time: tp.Now()},
				},
			},
			wantErr: false,
			setup: func() {
				rows := pgxpoolmock.NewRows([]string{"id", "author_id", "name", "file_name", "created_at", "updated_at"}).
					AddRow(uint(1), uint(1), "sticker1", "sticker1.png", tp.Now(), tp.Now()).
					AddRow(uint(2), uint(1), "sticker2", "sticker2.png", tp.Now(), tp.Now()).ToPgxRows()
				mockDB.EXPECT().Query(context.Background(), repository.GetStickersByAuthorIDQuery, gomock.Any()).Return(rows, nil)
			},
		},
		{
			name:     "test case 2 - no stickers found",
			authorID: 2,
			want:     nil,
			wantErr:  true,
			setup: func() {
				mockDB.EXPECT().Query(context.Background(), repository.GetStickersByAuthorIDQuery, gomock.Any()).Return(nil, pgx.ErrNoRows)
			},
		},
		{
			name:     "test case 3",
			authorID: 1,
			want:     nil,
			wantErr:  true,
			setup: func() {
				rows := pgxpoolmock.NewRows([]string{"err"}).AddRow(ErrRow{}).ToPgxRows()
				mockDB.EXPECT().Query(context.Background(), repository.GetStickersByAuthorIDQuery, gomock.Any()).Return(rows, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			pm := repository.NewPersonalMessages(mockDB, customtime.MockTimeProvider{})

			got, err := pm.GetStickersByAuthorID(context.Background(), tt.authorID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetAllStickers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name    string
		want    []*domain.Sticker
		wantErr bool
		setup   func()
	}{
		{
			name: "test case 1 - stickers found",
			want: []*domain.Sticker{
				{
					ID:        1,
					AuthorID:  1,
					Name:      "sticker1",
					FileName:  "sticker1.png",
					CreatedAt: customtime.CustomTime{Time: tp.Now()},
					UpdatedAt: customtime.CustomTime{Time: tp.Now()},
				},
				{
					ID:        2,
					AuthorID:  1,
					Name:      "sticker2",
					FileName:  "sticker2.png",
					CreatedAt: customtime.CustomTime{Time: tp.Now()},
					UpdatedAt: customtime.CustomTime{Time: tp.Now()},
				},
			},
			wantErr: false,
			setup: func() {
				rows := pgxpoolmock.NewRows([]string{"id", "author_id", "name", "file_name", "created_at", "updated_at"}).
					AddRow(uint(1), uint(1), "sticker1", "sticker1.png", tp.Now(), tp.Now()).
					AddRow(uint(2), uint(1), "sticker2", "sticker2.png", tp.Now(), tp.Now()).ToPgxRows()
				mockDB.EXPECT().Query(context.Background(), repository.GetAllStickersQuery).Return(rows, nil)
			},
		},
		{
			name:    "test case 2 - no stickers found",
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().Query(context.Background(), repository.GetAllStickersQuery).Return(nil, pgx.ErrNoRows)
			},
		},
		{
			name:    "test case 3",
			want:    nil,
			wantErr: true,
			setup: func() {
				rows := pgxpoolmock.NewRows([]string{"err"}).AddRow(ErrRow{}).ToPgxRows()
				mockDB.EXPECT().Query(context.Background(), repository.GetAllStickersQuery).Return(rows, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			pm := repository.NewPersonalMessages(mockDB, customtime.MockTimeProvider{})

			got, err := pm.GetAllStickers(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestStoreSticker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name     string
		sticker  *domain.Sticker
		want     *domain.Sticker
		wantErr  bool
		setup    func()
		teardown func()
	}{
		{
			name: "test case 1 - sticker stored successfully",
			sticker: &domain.Sticker{
				AuthorID: 1,
				Name:     "sticker1",
				FileName: "sticker1.png",
			},
			want: &domain.Sticker{
				ID:        1,
				AuthorID:  1,
				Name:      "sticker1",
				FileName:  "sticker1.png",
				CreatedAt: customtime.CustomTime{Time: tp.Now()},
				UpdatedAt: customtime.CustomTime{Time: tp.Now()},
			},
			wantErr: false,
			setup: func() {
				mockDB.EXPECT().QueryRow(context.Background(), repository.StoreStickerQuery, gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(
						uint(1), uint(1), "sticker1", "sticker1.png", tp.Now(), tp.Now(),
					),
				)
			},
			teardown: func() {},
		},
		{
			name: "test case 2 - error storing sticker",
			sticker: &domain.Sticker{
				AuthorID: 1,
				Name:     "sticker2",
				FileName: "sticker2.png",
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().QueryRow(context.Background(), repository.StoreStickerQuery, gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
			teardown: func() {},
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			pm := repository.NewPersonalMessages(mockDB, customtime.MockTimeProvider{})

			got, err := pm.StoreSticker(context.Background(), tt.sticker)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDeleteSticker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
		setup   func()
	}{
		{
			name:    "test case 1 - sticker deleted successfully",
			id:      1,
			wantErr: false,
			setup: func() {
				mockDB.EXPECT().Exec(context.Background(), repository.DeleteStickerQuery, gomock.Any()).Return(pgconn.CommandTag("DELETE 1"), nil)
			},
		},
		{
			name:    "test case 2 - error deleting sticker",
			id:      2,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().Exec(context.Background(), repository.DeleteStickerQuery, gomock.Any()).Return(nil, pgx.ErrNoRows)
			},
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			pm := repository.NewPersonalMessages(mockDB, customtime.MockTimeProvider{})

			err := pm.DeleteSticker(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStoreStickerMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name       string
		senderID   uint
		receiverID uint
		stickerID  uint
		want       *domain.PersonalMessage
		wantErr    bool
		setup      func()
	}{
		{
			name:       "test case 1 - sticker message stored successfully",
			senderID:   1,
			receiverID: 2,
			stickerID:  3,
			want: &domain.PersonalMessage{
				ID:         1,
				SenderID:   1,
				ReceiverID: 2,
				Sticker: &domain.Sticker{
					ID:        3,
					AuthorID:  1,
					Name:      "sticker3",
					FileName:  "sticker3.png",
					CreatedAt: customtime.CustomTime{Time: tp.Now()},
					UpdatedAt: customtime.CustomTime{Time: tp.Now()},
				},
				CreatedAt: customtime.CustomTime{Time: tp.Now()},
				UpdatedAt: customtime.CustomTime{Time: tp.Now()},
			},
			wantErr: false,
			setup: func() {
				mockDB.EXPECT().QueryRow(context.Background(), repository.StoreStickerMessageQuery, gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(1), uint(1), uint(2), uint(3), tp.Now(), tp.Now()),
				)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(3), uint(1), "sticker3", "sticker3.png", tp.Now(), tp.Now()),
				)
			},
		},
		{
			name:       "test case 2",
			senderID:   1,
			receiverID: 2,
			stickerID:  3,
			want:       nil,
			wantErr:    true,
			setup: func() {
				mockDB.EXPECT().QueryRow(context.Background(), repository.StoreStickerMessageQuery, gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(1), uint(1), uint(2), uint(3), tp.Now(), tp.Now()),
				)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(
					ErrRow{},
				)
			},
		},
		{
			name:       "test case 3",
			senderID:   1,
			receiverID: 2,
			stickerID:  3,
			want:       nil,
			wantErr:    true,
			setup: func() {
				mockDB.EXPECT().QueryRow(context.Background(), repository.StoreStickerMessageQuery, gomock.Any(), gomock.Any(), gomock.Any()).Return(
					ErrRow{},
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			pm := repository.NewPersonalMessages(mockDB, customtime.MockTimeProvider{})

			got, err := pm.StoreStickerMessage(context.Background(), tt.senderID, tt.receiverID, tt.stickerID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
