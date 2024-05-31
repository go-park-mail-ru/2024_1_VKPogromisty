package repository_test

import (
	"context"
	"socio/domain"
	"socio/errors"
	repository "socio/internal/repository/postgres"
	customtime "socio/pkg/time"
	"socio/usecase/posts"
	"testing"
	"time"

	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

var (
	postColumns = []string{"id", "author_id", "content", "created_at", "updated_at", "attacments", "liked_by"}
)

type ErrRow struct{}

func (r ErrRow) Scan(...interface{}) error {
	return pgx.ErrNoRows
}

type ErrInternalRow struct{}

func (r ErrInternalRow) Scan(...interface{}) error {
	return errors.ErrInternal
}

func TestGetPostByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		postID   uint
		mock     func(pool *pgxpoolmock.MockPgxIface, postID uint)
		expected *domain.Post
		err      error
	}{
		{
			name:   "Test OK",
			postID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, postID uint) {
				timeProv := customtime.MockTimeProvider{}

				arr := pgtype.TextArray{
					Elements: []pgtype.Text{{String: "file1"}, {String: "file2"}},
				}

				likedBy := pgtype.Int8Array{
					Elements: []pgtype.Int8{{Int: 1}, {Int: 2}},
				}

				row := pgxpoolmock.NewRow(
					uint(1),
					uint(2),
					"content",
					timeProv.Now(),
					timeProv.Now(),
					arr,
					likedBy,
				)

				pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
			},
			expected: &domain.Post{ID: 1},
			err:      nil,
		},
		{
			name:   "Test err",
			postID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, postID uint) {
				pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
			expected: &domain.Post{},
			err:      errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.postID)

			post, err := repo.GetPostByID(context.Background(), tt.postID)

			if err != tt.err {
				t.Errorf("unexpected error: %v", err)
			}

			if post.ID != tt.expected.ID {
				t.Errorf("unexpected post id: %d", post.ID)
			}
		})
	}
}

func TestGetUserPosts(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name        string
		userID      uint
		lastPostID  uint
		postsAmount uint
		mock        func(pool *pgxpoolmock.MockPgxIface, userID uint, lastPostID uint, postsAmount uint)
		expected    []*domain.Post
		err         error
	}{
		{
			name:        "Test OK",
			userID:      1,
			lastPostID:  0,
			postsAmount: 5,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint, lastPostID uint, postsAmount uint) {
				// Mock for QueryRow
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastUserPostIDQuery, userID).Return(pgxpoolmock.NewRow(lastPostID))

				tp := customtime.MockTimeProvider{}

				arr := pgtype.TextArray{
					Elements: []pgtype.Text{{String: "file1"}, {String: "file2"}},
				}

				likedBy := pgtype.Int8Array{
					Elements: []pgtype.Int8{{Int: 1}, {Int: 2}},
				}

				// Mock for Query
				rows := pgxpoolmock.NewRows(postColumns)
				for i := 0; i < int(postsAmount); i++ {
					rows.AddRow(uint(1), userID, "content", tp.Now(), tp.Now(), arr, likedBy)
				}
				pool.EXPECT().Query(gomock.Any(), repository.GetUserPostsQuery, gomock.Any(), gomock.Any(), gomock.Any()).Return(rows.ToPgxRows(), nil)
			},
			expected: []*domain.Post{
				{
					ID:       1,
					AuthorID: 1,
					Content:  "content",
					CreatedAt: customtime.CustomTime{
						Time: tp.Now(),
					},
					UpdatedAt: customtime.CustomTime{
						Time: tp.Now(),
					},
					Attachments: []string{
						"file1",
						"file2",
					},
				},
			},
			err: nil,
		},
		{
			name:        "Test ErrNoRows LastPostID",
			userID:      1,
			lastPostID:  0,
			postsAmount: 5,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint, lastPostID uint, postsAmount uint) {
				// Mock for QueryRow
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastUserPostIDQuery, gomock.Any()).Return(ErrRow{})
			},
			expected: []*domain.Post{},
			err:      pgx.ErrNoRows,
		},
		{
			name:        "Test err no rows",
			userID:      1,
			lastPostID:  0,
			postsAmount: 5,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint, lastPostID uint, postsAmount uint) {
				// Mock for QueryRow
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastUserPostIDQuery, userID).Return(pgxpoolmock.NewRow(lastPostID))
				pool.EXPECT().Query(gomock.Any(), repository.GetUserPostsQuery, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, pgx.ErrNoRows)
			},
			expected: []*domain.Post{},
			err:      pgx.ErrNoRows,
		},
		{
			name:        "Test err scan",
			userID:      1,
			lastPostID:  0,
			postsAmount: 5,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint, lastPostID uint, postsAmount uint) {
				// Mock for QueryRow
				pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(lastPostID))

				// Mock for Query
				rows := pgxpoolmock.NewRows([]string{"err"}).AddRow(ErrRow{}).ToPgxRows()
				pool.EXPECT().Query(gomock.Any(), repository.GetUserPostsQuery, gomock.Any(), gomock.Any(), gomock.Any()).Return(rows, nil)
			},
			expected: []*domain.Post{},
			err:      pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.userID, tt.lastPostID, tt.postsAmount)

			_, err := repo.GetUserPosts(context.Background(), tt.userID, tt.lastPostID, tt.postsAmount)

			if (tt.err != nil) != (err != nil) {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

func TestGetUserFriendsPosts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		userID      uint
		lastPostID  uint
		postsAmount uint
		mock        func(pool *pgxpoolmock.MockPgxIface, userID uint, lastPostID uint, postsAmount uint)
		expected    []*domain.Post
		err         bool
	}{
		{
			name:        "Test OK",
			userID:      1,
			lastPostID:  0,
			postsAmount: 20,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint, lastPostID uint, postsAmount uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastUserFriendsPostIDQuery, gomock.Any()).Return(pgxpoolmock.NewRow(lastPostID))

				timeProv := customtime.MockTimeProvider{}

				arr := pgtype.TextArray{
					Elements: []pgtype.Text{{String: "file1"}, {String: "file2"}},
				}

				likedBy := pgtype.Int8Array{
					Elements: []pgtype.Int8{{Int: 1}, {Int: 2}},
				}

				row := pgxpoolmock.NewRows(postColumns).AddRow(
					uint(1),
					uint(2),
					"content",
					timeProv.Now(),
					timeProv.Now(),
					arr,
					likedBy,
				).ToPgxRows()

				pool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(row, nil)
			},
			expected: []*domain.Post{},
			err:      false,
		},
		{
			name:        "test 2",
			userID:      1,
			lastPostID:  0,
			postsAmount: 20,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint, lastPostID uint, postsAmount uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastUserFriendsPostIDQuery, gomock.Any()).Return(pgxpoolmock.NewRow(lastPostID))
				pool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
			expected: []*domain.Post{},
			err:      true,
		},
		{
			name:        "test 3",
			userID:      1,
			lastPostID:  0,
			postsAmount: 20,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint, lastPostID uint, postsAmount uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastUserFriendsPostIDQuery, gomock.Any()).Return(ErrRow{})
			},
			expected: []*domain.Post{},
			err:      true,
		},
		{
			name:        "Test err scan",
			userID:      1,
			lastPostID:  0,
			postsAmount: 20,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint, lastPostID uint, postsAmount uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastUserFriendsPostIDQuery, gomock.Any()).Return(pgxpoolmock.NewRow(lastPostID))

				row := pgxpoolmock.NewRows([]string{"err"}).AddRow(ErrRow{}).ToPgxRows()

				pool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(row, nil)
			},
			expected: []*domain.Post{},
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.userID, tt.lastPostID, tt.postsAmount)

			_, err := repo.GetUserFriendsPosts(context.Background(), tt.userID, tt.lastPostID, tt.postsAmount)

			if (err != nil) != tt.err {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestStorePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name    string
		post    *domain.Post
		want    *domain.Post
		wantErr bool
		setup   func()
	}{
		{
			name: "test case 1",
			post: &domain.Post{
				AuthorID:    1,
				Content:     "Test content",
				Attachments: []string{"attachment1", "attachment2"},
			},
			want: &domain.Post{
				ID:          1,
				AuthorID:    1,
				Content:     "Test content",
				CreatedAt:   customtime.CustomTime{Time: tp.Now()},
				UpdatedAt:   customtime.CustomTime{Time: tp.Now()},
				Attachments: []string{"attachment1", "attachment2"},
			},
			wantErr: false,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), "Test content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment1"))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment2"))
				mockDB.EXPECT().Commit(context.Background()).Return(nil)
				mockDB.EXPECT().Rollback(context.Background()).Return(nil)
			},
		},
		{
			name: "test case 2",
			post: &domain.Post{
				AuthorID:    1,
				Content:     "Test content",
				Attachments: []string{"attachment1", "attachment2"},
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), "Test content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment1"))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment2"))
				mockDB.EXPECT().Commit(context.Background()).Return(nil)
				mockDB.EXPECT().Rollback(context.Background()).Return(errors.ErrInternal)
			},
		},
		{
			name: "test case 3",
			post: &domain.Post{
				AuthorID:    1,
				Content:     "Test content",
				Attachments: []string{"attachment1", "attachment2"},
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), "Test content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment1"))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment2"))
				mockDB.EXPECT().Commit(context.Background()).Return(errors.ErrInternal)
			},
		},
		{
			name: "test case 4",
			post: &domain.Post{
				AuthorID:    1,
				Content:     "Test content",
				Attachments: []string{"attachment1"},
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), "Test content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
		{
			name: "test case 5",
			post: &domain.Post{
				AuthorID:    1,
				Content:     "Test content",
				Attachments: []string{"attachment1"},
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
		{
			name: "test case 6",
			post: &domain.Post{
				AuthorID:    1,
				Content:     "Test content",
				Attachments: []string{"attachment1"},
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

			p := repository.NewPosts(mockDB, tp)

			got, err := p.StorePost(context.Background(), tt.post)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestStoreGroupPost(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name      string
		groupPost *domain.GroupPost
		mock      func(pool *pgxpoolmock.MockPgxIface, groupPost *domain.GroupPost)
		expected  *domain.GroupPost
		err       error
	}{
		{
			name: "Test OK",
			groupPost: &domain.GroupPost{
				PostID:  1,
				GroupID: 2,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, groupPost *domain.GroupPost) {
				newGroupPost := &domain.GroupPost{
					ID:      1,
					PostID:  groupPost.PostID,
					GroupID: groupPost.GroupID,
					CreatedAt: customtime.CustomTime{
						Time: tp.Now(),
					},
					UpdatedAt: customtime.CustomTime{
						Time: tp.Now(),
					},
				}

				pool.EXPECT().QueryRow(gomock.Any(), repository.StoreGroupPostQuery, gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(newGroupPost.ID, newGroupPost.PostID, newGroupPost.GroupID, newGroupPost.CreatedAt.Time, newGroupPost.UpdatedAt.Time))
			},
			expected: &domain.GroupPost{},
			err:      nil,
		},
		{
			name: "Test OK",
			groupPost: &domain.GroupPost{
				PostID:  1,
				GroupID: 2,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, groupPost *domain.GroupPost) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.StoreGroupPostQuery, gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
			expected: &domain.GroupPost{},
			err:      pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.groupPost)

			_, err := repo.StoreGroupPost(context.Background(), tt.groupPost)

			if err != tt.err {
				t.Errorf("unexpected error: %v", err)
			}

		})
	}
}

func TestGetGroupPostByPostID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name    string
		postID  uint
		want    *domain.GroupPost
		wantErr bool
		setup   func()
	}{
		{
			name:   "test case 1 - post found",
			postID: 1,
			want: &domain.GroupPost{
				ID:        1,
				PostID:    1,
				GroupID:   1,
				CreatedAt: customtime.CustomTime{Time: tp.Now()},
				UpdatedAt: customtime.CustomTime{Time: tp.Now()},
			},
			wantErr: false,
			setup: func() {
				mockDB.EXPECT().QueryRow(context.Background(), repository.GetGroupPostByPostIDQuery, gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), uint(1), tp.Now(), tp.Now()))
			},
		},
		{
			name:    "test case 2 - post not found",
			postID:  2,
			want:    nil,
			wantErr: true,
			setup: func() {
				mockDB.EXPECT().QueryRow(context.Background(), repository.GetGroupPostByPostIDQuery, gomock.Any()).Return(ErrRow{})
			},
			// Add more test cases here
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			p := repository.NewPosts(mockDB, customtime.MockTimeProvider{})

			got, err := p.GetGroupPostByPostID(context.Background(), tt.postID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUpdatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := pgxpoolmock.NewMockPgxIface(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name                string
		post                *domain.Post
		attachmentsToDelete []string
		want                *domain.Post
		wantErr             bool
		setup               func()
	}{
		{
			name: "test case 1",
			post: &domain.Post{
				ID:          1,
				AuthorID:    1,
				Content:     "Test content",
				Attachments: []string{"attachment1", "attachment2"},
			},
			attachmentsToDelete: []string{"attachment3"},
			want: &domain.Post{
				ID:          1,
				AuthorID:    1,
				Content:     "Test content",
				CreatedAt:   customtime.CustomTime{Time: tp.Now()},
				UpdatedAt:   customtime.CustomTime{Time: tp.Now()},
				Attachments: []string{"attachment1", "attachment2"},
			},
			wantErr: false,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), repository.UpdatePostQuery, gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), "Test content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), repository.StorePostAttachmentQuery, gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment1"))
				mockDB.EXPECT().QueryRow(context.Background(), repository.StorePostAttachmentQuery, gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment2"))
				mockDB.EXPECT().Exec(context.Background(), repository.DeletePostAttachmentQuery, gomock.Any()).Return(pgconn.CommandTag("DELETE 1"), nil)
				mockDB.EXPECT().Commit(context.Background()).Return(nil)
				mockDB.EXPECT().Rollback(context.Background()).Return(nil)
			},
		},
		{
			name: "test case 2",
			post: &domain.Post{
				ID:          1,
				AuthorID:    1,
				Content:     "Test content",
				Attachments: []string{"attachment1", "attachment2"},
			},
			attachmentsToDelete: []string{"attachment3"},
			want:                nil,
			wantErr:             true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), repository.UpdatePostQuery, gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), "Test content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), repository.StorePostAttachmentQuery, gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment1"))
				mockDB.EXPECT().QueryRow(context.Background(), repository.StorePostAttachmentQuery, gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment2"))
				mockDB.EXPECT().Exec(context.Background(), repository.DeletePostAttachmentQuery, gomock.Any()).Return(pgconn.CommandTag("DELETE 1"), nil)
				mockDB.EXPECT().Commit(context.Background()).Return(nil)
				mockDB.EXPECT().Rollback(context.Background()).Return(errors.ErrInternal)
			},
		},
		{
			name: "test case 3",
			post: &domain.Post{
				ID:          1,
				AuthorID:    1,
				Content:     "Test content",
				Attachments: []string{"attachment1", "attachment2"},
			},
			attachmentsToDelete: []string{"attachment3"},
			want:                nil,
			wantErr:             true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), repository.UpdatePostQuery, gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), "Test content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), repository.StorePostAttachmentQuery, gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment1"))
				mockDB.EXPECT().QueryRow(context.Background(), repository.StorePostAttachmentQuery, gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment2"))
				mockDB.EXPECT().Exec(context.Background(), repository.DeletePostAttachmentQuery, gomock.Any()).Return(pgconn.CommandTag("DELETE 1"), nil)
				mockDB.EXPECT().Commit(context.Background()).Return(errors.ErrInternal)
			},
		},
		{
			name: "test case 4",
			post: &domain.Post{
				ID:          1,
				AuthorID:    1,
				Content:     "Test content",
				Attachments: []string{"attachment1", "attachment2"},
			},
			attachmentsToDelete: []string{"attachment3"},
			want:                nil,
			wantErr:             true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), repository.UpdatePostQuery, gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), "Test content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), repository.StorePostAttachmentQuery, gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment1"))
				mockDB.EXPECT().QueryRow(context.Background(), repository.StorePostAttachmentQuery, gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow("attachment2"))
				mockDB.EXPECT().Exec(context.Background(), repository.DeletePostAttachmentQuery, gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
		{
			name: "test case 5",
			post: &domain.Post{
				ID:          1,
				AuthorID:    1,
				Content:     "Test content",
				Attachments: []string{"attachment1"},
			},
			attachmentsToDelete: []string{"attachment3"},
			want:                nil,
			wantErr:             true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), repository.UpdatePostQuery, gomock.Any(), gomock.Any()).Return(pgxpoolmock.NewRow(uint(1), uint(1), "Test content", tp.Now(), tp.Now()))
				mockDB.EXPECT().QueryRow(context.Background(), repository.StorePostAttachmentQuery, gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
		{
			name: "test case 6",
			post: &domain.Post{
				ID:          1,
				AuthorID:    1,
				Content:     "Test content",
				Attachments: []string{"attachment1"},
			},
			attachmentsToDelete: []string{"attachment3"},
			want:                nil,
			wantErr:             true,
			setup: func() {
				mockDB.EXPECT().BeginTx(context.Background(), gomock.Any()).Return(mockDB, nil)
				mockDB.EXPECT().QueryRow(context.Background(), repository.UpdatePostQuery, gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
		},
		{
			name: "test case 7",
			post: &domain.Post{
				ID:          1,
				AuthorID:    1,
				Content:     "Test content",
				Attachments: []string{"attachment1"},
			},
			attachmentsToDelete: []string{"attachment3"},
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

			p := repository.NewPosts(mockDB, customtime.MockTimeProvider{})

			got, err := p.UpdatePost(context.Background(), tt.post, tt.attachmentsToDelete)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDeletePost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		postID   uint
		mock     func(pool *pgxpoolmock.MockPgxIface, postID uint)
		expected error
	}{
		{
			name:   "Test OK",
			postID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, postID uint) {
				tag := pgconn.CommandTag("DELETE 1")
				pool.EXPECT().Exec(gomock.Any(), repository.DeletePostQuery, gomock.Any()).Return(tag, nil)
			},
			expected: nil,
		},
		{
			name:   "Test err",
			postID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, postID uint) {
				pool.EXPECT().Exec(gomock.Any(), repository.DeletePostQuery, gomock.Any()).Return(nil, errors.ErrInternal)
			},
			expected: errors.ErrInternal,
		},
		{
			name:   "Test err rows affected",
			postID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, postID uint) {
				tag := pgconn.CommandTag("DELETE 2")
				pool.EXPECT().Exec(gomock.Any(), repository.DeletePostQuery, gomock.Any()).Return(tag, nil)
			},
			expected: errors.ErrRowsAffected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.postID)

			err := repo.DeletePost(context.Background(), tt.postID)

			if err != tt.expected {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestDeleteGroupPost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		postID uint
		mock   func(pool *pgxpoolmock.MockPgxIface, postID uint)
		err    error
	}{
		{
			name:   "Test OK",
			postID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, postID uint) {
				pool.EXPECT().Exec(gomock.Any(), repository.DeleteGroupPostQuery, postID).Return(pgconn.CommandTag("DELETE 1"), nil)
			},
			err: nil,
		},
		{
			name:   "Test rows affected err",
			postID: 2,
			mock: func(pool *pgxpoolmock.MockPgxIface, postID uint) {
				pool.EXPECT().Exec(gomock.Any(), repository.DeleteGroupPostQuery, postID).Return(pgconn.CommandTag("DELETE 2"), nil)
			},
			err: errors.ErrRowsAffected,
		},
		{
			name:   "Test err",
			postID: 2,
			mock: func(pool *pgxpoolmock.MockPgxIface, postID uint) {
				pool.EXPECT().Exec(gomock.Any(), repository.DeleteGroupPostQuery, postID).Return(nil, errors.ErrInternal)
			},
			err: errors.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.postID)

			err := repo.DeleteGroupPost(context.Background(), tt.postID)

			if err != tt.err {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetLikedPosts(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name       string
		userID     uint
		lastLikeID uint
		limit      uint
		mock       func(pool *pgxpoolmock.MockPgxIface, userID uint, lastLikeID uint, limit uint)
		expected   []posts.LikeWithPost
		err        bool
	}{
		{
			name:       "Test OK",
			userID:     1,
			lastLikeID: 0,
			limit:      10,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint, lastLikeID uint, limit uint) {
				// Mock the GetLastPostLikeIDQuery
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostLikeIDQuery, gomock.Any()).Return(pgxpoolmock.NewRow(uint(1)))

				arr := pgtype.TextArray{
					Elements: []pgtype.Text{{String: "file1", Status: pgtype.Present}, {String: "file2", Status: pgtype.Present}},
					Status:   pgtype.Present,
				}

				likedBy := pgtype.Int8Array{
					Elements: []pgtype.Int8{{Int: 1, Status: pgtype.Present}, {Int: 2, Status: pgtype.Present}},
					Status:   pgtype.Present,
				}

				// Mock the GetLikedPosts query
				rows := pgxpoolmock.NewRows([]string{"like_id", "post_id", "user_id", "created_at", "author_id", "content", "created_at", "updated_at", "attachments", "liked_by_ids"}).AddRow(
					uint(1), uint(1), uint(1), tp.Now(), uint(1), "content", tp.Now(), tp.Now(), arr, likedBy).ToPgxRows()
				pool.EXPECT().Query(gomock.Any(), repository.GetLikedPosts, gomock.Any(), gomock.Any(), gomock.Any()).Return(rows, nil)
			},
			expected: []posts.LikeWithPost{
				{
					Like: &domain.PostLike{
						ID:     1,
						PostID: 1,
						UserID: 1,
						CreatedAt: customtime.CustomTime{
							Time: tp.Now(),
						},
					},
					Post: &domain.Post{
						ID:       1,
						AuthorID: 1,
						Content:  "content",
						CreatedAt: customtime.CustomTime{
							Time: tp.Now(),
						},
						UpdatedAt: customtime.CustomTime{
							Time: tp.Now(),
						},
						Attachments: []string{"file1", "file2"},
						LikedByIDs:  []uint64{1, 2},
					},
				},
			},
			err: false,
		},
		{
			name:       "Test 2",
			userID:     1,
			lastLikeID: 0,
			limit:      10,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint, lastLikeID uint, limit uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostLikeIDQuery, gomock.Any()).Return(pgxpoolmock.NewRow(uint(1)))
				pool.EXPECT().Query(gomock.Any(), repository.GetLikedPosts, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
			expected: nil,
			err:      true,
		},
		{
			name:       "Test 3",
			userID:     1,
			lastLikeID: 0,
			limit:      10,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint, lastLikeID uint, limit uint) {
				// Mock the GetLastPostLikeIDQuery
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostLikeIDQuery, gomock.Any()).Return(pgxpoolmock.NewRow(uint(1)))
				rows := pgxpoolmock.NewRows([]string{"err"}).AddRow(ErrRow{}).ToPgxRows()
				pool.EXPECT().Query(gomock.Any(), repository.GetLikedPosts, gomock.Any(), gomock.Any(), gomock.Any()).Return(rows, nil)
			},
			expected: nil,
			err:      true,
		},
		{
			name:       "Test 4",
			userID:     1,
			lastLikeID: 0,
			limit:      10,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint, lastLikeID uint, limit uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostLikeIDQuery, gomock.Any()).Return(ErrRow{})
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

			repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.userID, tt.lastLikeID, tt.limit)

			got, err := repo.GetLikedPosts(context.Background(), tt.userID, tt.lastLikeID, tt.limit)

			if tt.err != (err != nil) {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !tt.err {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestGetPostLikeByUserIDAndPostID(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name     string
		postID   uint
		userID   uint
		mock     func(pool *pgxpoolmock.MockPgxIface, postID uint)
		expected *domain.PostLike
		err      error
	}{
		{
			name:   "Test OK",
			postID: 1,
			userID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, postID uint) {
				pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(1), uint(1), uint(1), tp.Now()),
				)
			},
			expected: &domain.PostLike{
				ID:        1,
				PostID:    1,
				UserID:    1,
				CreatedAt: customtime.CustomTime{Time: tp.Now()},
			},
			err: nil,
		},
		{
			name:   "Test err",
			postID: 1,
			userID: 1,
			mock: func(pool *pgxpoolmock.MockPgxIface, postID uint) {
				pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrRow{})
			},
			expected: nil,
			err:      errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.postID)

			like, err := repo.GetPostLikeByUserIDAndPostID(context.Background(), tt.userID, tt.postID)

			if err != tt.err {
				t.Errorf("unexpected error: %v", err)
			}

			if err == nil {
				assert.Equal(t, tt.expected, like)
			}
		})
	}
}

func TestStorePostLike(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name     string
		likeData *domain.PostLike
		mock     func(pool *pgxpoolmock.MockPgxIface, likeData *domain.PostLike)
		expected *domain.PostLike
		err      error
	}{
		{
			name: "Test OK",
			likeData: &domain.PostLike{
				PostID: 1,
				UserID: 1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, likeData *domain.PostLike) {
				// Mock the CreatePostLikeQuery
				pool.EXPECT().QueryRow(gomock.Any(), repository.CreatePostLikeQuery, likeData.PostID, likeData.UserID).Return(pgxpoolmock.NewRow(uint(1), likeData.PostID, likeData.UserID, time.Now()))
			},
			expected: &domain.PostLike{
				ID:     1,
				PostID: 1,
				UserID: 1,
				CreatedAt: customtime.CustomTime{
					Time: tp.Now(),
				},
			},
			err: nil,
		},
		{
			name: "Test err",
			likeData: &domain.PostLike{
				PostID: 1,
				UserID: 1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, likeData *domain.PostLike) {
				// Mock the CreatePostLikeQuery
				pool.EXPECT().QueryRow(gomock.Any(), repository.CreatePostLikeQuery, likeData.PostID, likeData.UserID).Return(
					ErrRow{},
				)
			},
			expected: &domain.PostLike{},
			err:      pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pool := pgxpoolmock.NewMockPgxIface(ctrl)

			repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.likeData)

			_, err := repo.StorePostLike(context.Background(), tt.likeData)

			if err != tt.err {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestDeletePostLike(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		likeData *domain.PostLike
		mock     func(pool *pgxpoolmock.MockPgxIface, likeData *domain.PostLike)
		err      error
	}{
		{
			name: "Test OK",
			likeData: &domain.PostLike{
				PostID: 1,
				UserID: 1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, likeData *domain.PostLike) {
				// Mock the DeletePostLikeQuery
				pool.EXPECT().Exec(gomock.Any(), repository.DeletePostLikeQuery, likeData.PostID, likeData.UserID).Return(
					pgconn.CommandTag("DELETE 1"),
					nil,
				)
			},
			err: nil,
		},
		{
			name: "Test ErrInternal",
			likeData: &domain.PostLike{
				PostID: 1,
				UserID: 1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, likeData *domain.PostLike) {
				// Mock the DeletePostLikeQuery
				pool.EXPECT().Exec(gomock.Any(), repository.DeletePostLikeQuery, likeData.PostID, likeData.UserID).Return(
					nil,
					errors.ErrInternal,
				)
			},
			err: errors.ErrInternal,
		},
		{
			name: "Test ErrRowsAffected",
			likeData: &domain.PostLike{
				PostID: 1,
				UserID: 1,
			},
			mock: func(pool *pgxpoolmock.MockPgxIface, likeData *domain.PostLike) {
				// Mock the DeletePostLikeQuery
				pool.EXPECT().Exec(gomock.Any(), repository.DeletePostLikeQuery, likeData.PostID, likeData.UserID).Return(
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

			repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.likeData)

			err := repo.DeletePostLike(context.Background(), tt.likeData)

			if err != tt.err {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetPostsOfGroup(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name        string
		groupID     uint
		lastPostID  uint
		postsAmount uint
		mock        func(pool *pgxpoolmock.MockPgxIface, groupID uint, lastPostID uint, postsAmount uint)
		expected    []*domain.Post
		err         bool
	}{
		{
			name:        "Test OK",
			groupID:     1,
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupID uint, lastPostID uint, postsAmount uint) {
				// Mock the GetLastPostOfGroupIDQuery
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostOfGroupIDQuery, gomock.Any()).Return(pgxpoolmock.NewRow(uint(1)))

				arr := pgtype.TextArray{
					Elements: []pgtype.Text{{String: "file1", Status: pgtype.Present}, {String: "file2", Status: pgtype.Present}},
					Status:   pgtype.Present,
				}

				likedBy := pgtype.Int8Array{
					Elements: []pgtype.Int8{{Int: 1, Status: pgtype.Present}, {Int: 2, Status: pgtype.Present}},
					Status:   pgtype.Present,
				}

				// Mock the GetPostsOfGroupQuery
				rows := pgxpoolmock.NewRows([]string{"id", "author_id", "content", "created_at", "updated_at", "attachments", "liked_by_ids"})
				rows.AddRow(uint(1), uint(1), "content", tp.Now(), tp.Now(), arr, likedBy)
				pool.EXPECT().Query(gomock.Any(), repository.GetPostsOfGroupQuery, gomock.Any(), gomock.Any(), gomock.Any()).Return(
					rows.ToPgxRows(),
					nil,
				)
			},
			expected: []*domain.Post{
				{
					ID:       1,
					AuthorID: 1,
					Content:  "content",
					CreatedAt: customtime.CustomTime{
						Time: tp.Now(),
					},
					UpdatedAt: customtime.CustomTime{
						Time: tp.Now(),
					},
					Attachments: []string{"file1", "file2"},
					LikedByIDs:  []uint64{1, 2},
				},
			},
			err: false,
		},
		{
			name:        "Test ErrQueryRow",
			groupID:     1,
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupID uint, lastPostID uint, postsAmount uint) {
				// Mock the GetLastPostOfGroupIDQuery
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostOfGroupIDQuery, gomock.Any()).Return(
					ErrRow{},
				)
			},
			expected: nil,
			err:      true,
		},
		{
			name:        "Test ErrQuery",
			groupID:     1,
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupID uint, lastPostID uint, postsAmount uint) {
				// Mock the GetLastPostOfGroupIDQuery
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostOfGroupIDQuery, gomock.Any()).Return(pgxpoolmock.NewRow(uint(1)))
				pool.EXPECT().Query(gomock.Any(), repository.GetPostsOfGroupQuery, gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil,
					errors.ErrInternal,
				)
			},
			expected: nil,
			err:      true,
		},
		{
			name:        "Test OK",
			groupID:     1,
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupID uint, lastPostID uint, postsAmount uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostOfGroupIDQuery, gomock.Any()).Return(pgxpoolmock.NewRow(uint(1)))
				rows := pgxpoolmock.NewRows([]string{"err"})
				rows.AddRow(ErrRow{})
				pool.EXPECT().Query(gomock.Any(), repository.GetPostsOfGroupQuery, gomock.Any(), gomock.Any(), gomock.Any()).Return(
					rows.ToPgxRows(),
					nil,
				)
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

			repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.groupID, tt.lastPostID, tt.postsAmount)

			got, err := repo.GetPostsOfGroup(context.Background(), tt.groupID, tt.lastPostID, tt.postsAmount)

			if tt.err != (err != nil) {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !tt.err {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestGetGroupPostsBySubscriptionIDs(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name        string
		subIDs      []uint
		lastPostID  uint
		postsAmount uint
		mock        func(pool *pgxpoolmock.MockPgxIface, subIDs []uint, lastPostID uint, postsAmount uint)
		expected    []*domain.Post
		err         bool
	}{
		{
			name:        "Test OK",
			subIDs:      []uint{1},
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, subIDs []uint, lastPostID uint, postsAmount uint) {
				// Mock the GetLastGroupPostBySubscriptionIDsQuery
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastGroupPostBySubscriptionIDsQuery, gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(1)),
				)

				arr := pgtype.TextArray{
					Elements: []pgtype.Text{{String: "file1", Status: pgtype.Present}, {String: "file2", Status: pgtype.Present}},
					Status:   pgtype.Present,
				}

				likedBy := pgtype.Int8Array{
					Elements: []pgtype.Int8{{Int: 1, Status: pgtype.Present}, {Int: 2, Status: pgtype.Present}},
					Status:   pgtype.Present,
				}
				// Mock the GetGroupPostsBySubscriptionIDsQuery
				rows := pgxpoolmock.NewRows([]string{"id", "author_id", "content", "created_at", "updated_at", "attachments", "liked_by_ids", "group_id"})
				rows.AddRow(uint(1), uint(1), "content", tp.Now(), tp.Now(), arr, likedBy, uint(1))
				pool.EXPECT().Query(gomock.Any(), repository.GetGroupPostsBySubscriptionIDsQuery, gomock.Any(), gomock.Any(), gomock.Any()).Return(
					rows.ToPgxRows(),
					nil,
				)
			},
			expected: []*domain.Post{
				{
					ID:       1,
					AuthorID: 1,
					Content:  "content",
					CreatedAt: customtime.CustomTime{
						Time: tp.Now(),
					},
					UpdatedAt: customtime.CustomTime{
						Time: tp.Now(),
					},
					Attachments: []string{"file1", "file2"},
					LikedByIDs:  []uint64{1, 2},
					GroupID:     1,
				},
			},
			err: false,
		},
		{
			name:        "Test 2",
			subIDs:      []uint{},
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, subIDs []uint, lastPostID uint, postsAmount uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastGroupPostBySubscriptionIDsQuery, gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(1)),
				)
				rows := pgxpoolmock.NewRows([]string{"err"}).AddRow(ErrRow{}).ToPgxRows()
				pool.EXPECT().Query(gomock.Any(), repository.GetGroupPostsBySubscriptionIDsQuery, gomock.Any(), gomock.Any(), gomock.Any()).Return(
					rows,
					nil,
				)
			},
			expected: nil,
			err:      true,
		},
		{
			name:        "Test 3",
			subIDs:      []uint{},
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, subIDs []uint, lastPostID uint, postsAmount uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastGroupPostBySubscriptionIDsQuery, gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(1)),
				)
				pool.EXPECT().Query(gomock.Any(), repository.GetGroupPostsBySubscriptionIDsQuery, gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal,
				)
			},
			expected: nil,
			err:      true,
		},
		{
			name:        "Test 4",
			subIDs:      []uint{},
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, subIDs []uint, lastPostID uint, postsAmount uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastGroupPostBySubscriptionIDsQuery, gomock.Any()).Return(
					ErrRow{},
				)
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

			repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.subIDs, tt.lastPostID, tt.postsAmount)

			got, err := repo.GetGroupPostsBySubscriptionIDs(context.Background(), tt.subIDs, tt.lastPostID, tt.postsAmount)

			if tt.err != (err != nil) {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !tt.err {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestGetPostsByGroupSubIDsAndUserSubIDs(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name        string
		groupSubIDs []uint
		userSubIDs  []uint
		lastPostID  uint
		postsAmount uint
		mock        func(pool *pgxpoolmock.MockPgxIface, groupSubIDs, userSubIDs []uint, lastPostID, postsAmount uint)
		expected    []*domain.Post
		err         bool
	}{
		{
			name:        "Test OK",
			groupSubIDs: []uint{1},
			userSubIDs:  []uint{1},
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupSubIDs, userSubIDs []uint, lastPostID, postsAmount uint) {
				// Mock the GetLastPostByGroupSubIDsAndUserSubIDsQuery
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostByGroupSubIDsAndUserSubIDsQuery, gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(1)),
				)

				arr := pgtype.TextArray{
					Elements: []pgtype.Text{{String: "file1", Status: pgtype.Present}, {String: "file2", Status: pgtype.Present}},
					Status:   pgtype.Present,
				}

				likedBy := pgtype.Int8Array{
					Elements: []pgtype.Int8{{Int: 1, Status: pgtype.Present}, {Int: 2, Status: pgtype.Present}},
					Status:   pgtype.Present,
				}

				// Mock the GetPostsByGroupSubIDsAndUserSubIDsQuery
				rows := pgxpoolmock.NewRows([]string{"id", "author_id", "content", "created_at", "updated_at", "attachments", "liked_by_ids", "public_group_id"})
				rows.AddRow(uint(1), uint(1), "content", tp.Now(), tp.Now(), arr, likedBy, uint(1))
				pool.EXPECT().Query(gomock.Any(), repository.GetPostsByGroupSubIDsAndUserSubIDsQuery, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					rows.ToPgxRows(),
					nil,
				)
			},
			expected: []*domain.Post{
				{
					ID:       1,
					AuthorID: 1,
					GroupID:  1,
					Content:  "content",
					CreatedAt: customtime.CustomTime{
						Time: tp.Now(),
					},
					UpdatedAt: customtime.CustomTime{
						Time: tp.Now(),
					},
					Attachments: []string{"file1", "file2"},
					LikedByIDs:  []uint64{1, 2},
				},
			},
			err: false,
		},
		{
			name:        "Test 2",
			groupSubIDs: []uint{},
			userSubIDs:  []uint{},
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupSubIDs, userSubIDs []uint, lastPostID, postsAmount uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostByGroupSubIDsAndUserSubIDsQuery, gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(1)),
				)
				rows := pgxpoolmock.NewRows([]string{"err"}).AddRow(
					ErrRow{},
				).ToPgxRows()
				pool.EXPECT().Query(gomock.Any(), repository.GetPostsByGroupSubIDsAndUserSubIDsQuery, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					rows,
					nil,
				)
			},
			expected: nil,
			err:      true,
		},
		{
			name:        "Test 3",
			groupSubIDs: []uint{},
			userSubIDs:  []uint{},
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupSubIDs, userSubIDs []uint, lastPostID, postsAmount uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostByGroupSubIDsAndUserSubIDsQuery, gomock.Any(), gomock.Any()).Return(
					pgxpoolmock.NewRow(uint(1)),
				)
				pool.EXPECT().Query(gomock.Any(), repository.GetPostsByGroupSubIDsAndUserSubIDsQuery, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal,
				)
			},
			expected: nil,
			err:      true,
		},
		{
			name:        "Test 4",
			groupSubIDs: []uint{},
			userSubIDs:  []uint{},
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, groupSubIDs, userSubIDs []uint, lastPostID, postsAmount uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostByGroupSubIDsAndUserSubIDsQuery, gomock.Any(), gomock.Any()).Return(
					ErrRow{},
				)
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

			repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.groupSubIDs, tt.userSubIDs, tt.lastPostID, tt.postsAmount)

			got, err := repo.GetPostsByGroupSubIDsAndUserSubIDs(context.Background(), tt.groupSubIDs, tt.userSubIDs, tt.lastPostID, tt.postsAmount)

			if tt.err != (err != nil) {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !tt.err {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestGetNewPosts(t *testing.T) {
	t.Parallel()

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name        string
		lastPostID  uint
		postsAmount uint
		mock        func(pool *pgxpoolmock.MockPgxIface, lastPostID, postsAmount uint)
		expected    []*domain.Post
		err         bool
	}{
		{
			name:        "Test OK",
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, lastPostID, postsAmount uint) {
				// Mock the GetLastPostIDQuery
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostIDQuery).Return(pgxpoolmock.NewRow(uint(1)))

				arr := pgtype.TextArray{
					Elements: []pgtype.Text{{String: "file1", Status: pgtype.Present}, {String: "file2", Status: pgtype.Present}},
					Status:   pgtype.Present,
				}

				likedBy := pgtype.Int8Array{
					Elements: []pgtype.Int8{{Int: 1, Status: pgtype.Present}, {Int: 2, Status: pgtype.Present}},
					Status:   pgtype.Present,
				}

				// Mock the GetNewPostsQuery
				rows := pgxpoolmock.NewRows([]string{"id", "author_id", "content", "created_at", "updated_at", "attachments", "liked_by_ids", "group_id"})
				rows.AddRow(uint(1), uint(1), "content", tp.Now(), tp.Now(), arr, likedBy, uint(1))
				pool.EXPECT().Query(gomock.Any(), repository.GetNewPostsQuery, gomock.Any(), gomock.Any()).Return(rows.ToPgxRows(), nil)
			},
			expected: []*domain.Post{
				{
					ID:       1,
					AuthorID: 1,
					Content:  "content",
					CreatedAt: customtime.CustomTime{
						Time: tp.Now(),
					},
					UpdatedAt: customtime.CustomTime{
						Time: tp.Now(),
					},
					Attachments: []string{"file1", "file2"},
					LikedByIDs:  []uint64{1, 2},
					GroupID:     1,
				},
			},
			err: false,
		},
		{
			name:        "Test 2",
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, lastPostID, postsAmount uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostIDQuery).Return(pgxpoolmock.NewRow(uint(1)))
				rows := pgxpoolmock.NewRows([]string{"err"}).AddRow(ErrRow{}).ToPgxRows()
				pool.EXPECT().Query(gomock.Any(), repository.GetNewPostsQuery, gomock.Any(), gomock.Any()).Return(rows, nil)
			},
			expected: nil,
			err:      true,
		},
		{
			name:        "Test 3",
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, lastPostID, postsAmount uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostIDQuery).Return(pgxpoolmock.NewRow(uint(1)))
				pool.EXPECT().Query(gomock.Any(), repository.GetNewPostsQuery, gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
			expected: nil,
			err:      true,
		},
		{
			name:        "Test 2",
			lastPostID:  0,
			postsAmount: 10,
			mock: func(pool *pgxpoolmock.MockPgxIface, lastPostID, postsAmount uint) {
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastPostIDQuery).Return(ErrRow{})
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

			repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

			tt.mock(pool, tt.lastPostID, tt.postsAmount)

			got, err := repo.GetNewPosts(context.Background(), tt.lastPostID, tt.postsAmount)

			if tt.err != (err != nil) {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !tt.err {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
