package repository_test

import (
	"context"
	"socio/domain"
	"socio/errors"
	repository "socio/internal/repository/postgres"
	customtime "socio/pkg/time"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/assert"

	"github.com/chrisyxlee/pgxpoolmock"
)

func TestGetCommentsByPostID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name    string
		postID  uint
		want    []*domain.Comment
		wantErr bool
		setup   func()
	}{
		{
			name:   "test case 1",
			postID: 1,
			want: []*domain.Comment{
				{
					ID:         1,
					PostID:     1,
					AuthorID:   1,
					Content:    "Test content",
					CreatedAt:  customtime.CustomTime{Time: tp.Now()},
					UpdatedAt:  customtime.CustomTime{Time: tp.Now()},
					LikedByIDs: []uint64{1, 2, 3},
				},
			},
			wantErr: false,
			setup: func() {
				rows := pgxpoolmock.NewRows([]string{"id", "post_id", "author_id", "content", "created_at", "updated_at", "liked_by_ids"}).
					AddRow(uint(1), uint(1), uint(1), "Test content", tp.Now(), tp.Now(), pgtype.Int8Array{Elements: []pgtype.Int8{
						{Int: 1, Status: pgtype.Present},
						{Int: 2, Status: pgtype.Present},
						{Int: 3, Status: pgtype.Present},
					}, Status: pgtype.Present}).ToPgxRows()
				pool.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(rows, nil)
			},
		},
		{
			name:    "test case 2",
			postID:  1,
			want:    nil,
			wantErr: true,
			setup: func() {
				pool.EXPECT().Query(context.Background(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			p := repository.NewPosts(pool, tp)

			got, err := p.GetCommentsByPostID(context.Background(), tt.postID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetCommentByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name      string
		commentID uint
		want      *domain.Comment
		wantErr   bool
		setup     func()
	}{
		{
			name:      "test case 1",
			commentID: 1,
			want: &domain.Comment{
				ID:        1,
				PostID:    1,
				AuthorID:  1,
				Content:   "Test content",
				CreatedAt: customtime.CustomTime{Time: tp.Now()},
				UpdatedAt: customtime.CustomTime{Time: tp.Now()},
			},
			wantErr: false,
			setup: func() {
				row := pgxpoolmock.NewRow(uint(1), uint(1), uint(1), "Test content", tp.Now(), tp.Now())
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(row)
			},
		},
		{
			name:      "test case 2",
			commentID: 1,
			want:      nil,
			wantErr:   true,
			setup: func() {
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			p := repository.NewPosts(pool, tp)

			got, err := p.GetCommentByID(context.Background(), tt.commentID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestStoreComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name    string
		comment *domain.Comment
		want    *domain.Comment
		wantErr bool
		setup   func()
	}{
		{
			name:    "test case 1",
			comment: &domain.Comment{PostID: 1, AuthorID: 1, Content: "Test content"},
			want: &domain.Comment{
				ID:        1,
				PostID:    1,
				AuthorID:  1,
				Content:   "Test content",
				CreatedAt: customtime.CustomTime{Time: tp.Now()},
				UpdatedAt: customtime.CustomTime{Time: tp.Now()},
			},
			wantErr: false,
			setup: func() {
				row := pgxpoolmock.NewRow(uint(1), uint(1), uint(1), "Test content", tp.Now(), tp.Now())
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(row)
			},
		},
		{
			name:    "test case 2",
			comment: &domain.Comment{PostID: 1, AuthorID: 1, Content: "Test content"},
			want:    nil,
			wantErr: true,
			setup: func() {
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			p := repository.NewPosts(pool, tp)

			got, err := p.StoreComment(context.Background(), tt.comment)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUpdateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name    string
		comment *domain.Comment
		want    *domain.Comment
		wantErr bool
		setup   func()
	}{
		{
			name:    "test case 1",
			comment: &domain.Comment{ID: 1, Content: "Test content 1"},
			want: &domain.Comment{
				ID:        1,
				PostID:    1,
				AuthorID:  1,
				Content:   "Test content 1",
				CreatedAt: customtime.CustomTime{Time: tp.Now()},
				UpdatedAt: customtime.CustomTime{Time: tp.Now()},
			},
			wantErr: false,
			setup: func() {
				row := pgxpoolmock.NewRow(uint(1), uint(1), uint(1), "Test content 1", tp.Now(), tp.Now())
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(row).AnyTimes()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			p := repository.NewPosts(pool, tp)

			got, err := p.UpdateComment(context.Background(), tt.comment)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDeleteComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name      string
		commentID uint
		wantErr   bool
		setup     func()
	}{
		{
			name:      "test case 1",
			commentID: 1,
			wantErr:   false,
			setup: func() {
				tag := pgconn.CommandTag("DELETE 1")
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(tag, nil)
			},
		},
		{
			name:      "test case 2",
			commentID: 1,
			wantErr:   true,
			setup: func() {
				tag := pgconn.CommandTag("DELETE 0")
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(tag, nil)
			},
		},
		{
			name:      "test case 3",
			commentID: 1,
			wantErr:   true,
			setup: func() {
				tag := pgconn.CommandTag("DELETE 2")
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(tag, nil)
			},
		},
		{
			name:      "test case 4",
			commentID: 1,
			wantErr:   true,
			setup: func() {
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			p := repository.NewPosts(pool, tp)

			err := p.DeleteComment(context.Background(), tt.commentID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStoreCommentLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name        string
		commentLike *domain.CommentLike
		want        *domain.CommentLike
		wantErr     bool
		setup       func()
	}{
		{
			name: "test case 1",
			commentLike: &domain.CommentLike{
				CommentID: 1,
				UserID:    1,
				CreatedAt: customtime.CustomTime{Time: tp.Now()},
			},
			want: &domain.CommentLike{
				ID:        1,
				CommentID: 1,
				UserID:    1,
				CreatedAt: customtime.CustomTime{Time: tp.Now()},
			},
			wantErr: false,
			setup: func() {
				row := pgxpoolmock.NewRow(uint(1), uint(1), uint(1), tp.Now())
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrRow{})
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(row)
			},
		},
		{
			name: "test case 2",
			commentLike: &domain.CommentLike{
				CommentID: 1,
				UserID:    1,
				CreatedAt: customtime.CustomTime{Time: tp.Now()},
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				row := pgxpoolmock.NewRow(uint(1), uint(1), uint(1), tp.Now())
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
			},
		},
		{
			name: "test case 3",
			commentLike: &domain.CommentLike{
				CommentID: 1,
				UserID:    1,
				CreatedAt: customtime.CustomTime{Time: tp.Now()},
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrRow{})
				pool.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			p := repository.NewPosts(pool, tp)

			got, err := p.StoreCommentLike(context.Background(), tt.commentLike)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDeleteCommentLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name        string
		commentLike *domain.CommentLike
		wantErr     bool
		setup       func()
	}{
		{
			name:        "test case 1",
			commentLike: &domain.CommentLike{CommentID: 1, UserID: 1},
			wantErr:     false,
			setup: func() {
				tag := pgconn.CommandTag("DELETE 1")
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tag, nil)
			},
		},
		{
			name:        "test case 2",
			commentLike: &domain.CommentLike{CommentID: 1, UserID: 1},
			wantErr:     true,
			setup: func() {
				tag := pgconn.CommandTag("DELETE 0")
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tag, nil)
			},
		},
		{
			name:        "test case 3",
			commentLike: &domain.CommentLike{CommentID: 1, UserID: 1},
			wantErr:     true,
			setup: func() {
				tag := pgconn.CommandTag("DELETE 2")
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tag, nil)
			},
		},
		{
			name:        "test case 4",
			commentLike: &domain.CommentLike{CommentID: 1, UserID: 1},
			wantErr:     true,
			setup: func() {
				pool.EXPECT().Exec(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			p := repository.NewPosts(pool, tp)

			err := p.DeleteCommentLike(context.Background(), tt.commentLike)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
