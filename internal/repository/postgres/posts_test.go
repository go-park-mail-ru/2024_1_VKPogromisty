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
	"github.com/jackc/pgx/v4"
)

var (
	postColumns = []string{"id", "author_id", "content", "created_at", "updated_at", "attacments", "liked_by"}
)

type ErrRow struct{}

func (r ErrRow) Scan(...interface{}) error {
	return pgx.ErrNoRows
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
			name:   "Test OK",
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
			name:        "Test OK",
			userID:      1,
			lastPostID:  0,
			postsAmount: 5,
			mock: func(pool *pgxpoolmock.MockPgxIface, userID uint, lastPostID uint, postsAmount uint) {
				// Mock for QueryRow
				pool.EXPECT().QueryRow(gomock.Any(), repository.GetLastUserPostIDQuery, userID).Return(pgxpoolmock.NewRow(lastPostID))
				pool.EXPECT().Query(gomock.Any(), repository.GetUserPostsQuery, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, pgx.ErrNoRows)
			},
			expected: []*domain.Post{},
			err:      nil,
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

			if err != tt.err {
				t.Errorf("unexpected error: %v", err)
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
		err         error
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
			err:      nil,
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

			if err != tt.err {
				t.Errorf("unexpected error: %v", err)
			}

			// Add more assertions to compare the returned posts with the expected posts
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

func TestStorePost(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	timeProv := customtime.MockTimeProvider{}

	postRow := pgxpoolmock.NewRow(
		uint(1),
		uint(1),
		"content",
		timeProv.Now(),
		timeProv.Now(),
	)

	attachmentRow := pgxpoolmock.NewRow(
		"default_avatar.png",
	)

	pool.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(pool, nil)
	pool.EXPECT().QueryRow(gomock.Any(), repository.StorePostQuery, gomock.Any(), gomock.Any()).Return(postRow)
	pool.EXPECT().QueryRow(gomock.Any(), repository.StoreAttachmentQuery, gomock.Any(), gomock.Any()).Return(attachmentRow)
	pool.EXPECT().Rollback(gomock.Any()).Return(nil)
	pool.EXPECT().Commit(gomock.Any()).Return(nil)

	repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

	post := &domain.Post{
		AuthorID: 1,
		Content:  "content",
		Attachments: []string{
			"default_avatar.png",
		},
	}

	newPost, err := repo.StorePost(context.Background(), post)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if newPost.ID != 1 {
		t.Errorf("unexpected post id: %d", newPost.ID)
	}

	if newPost.AuthorID != 1 {
		t.Errorf("unexpected post author id: %d", newPost.AuthorID)
	}

	if newPost.Content != "content" {
		t.Errorf("unexpected post content: %s", newPost.Content)
	}

	if newPost.CreatedAt.Time != timeProv.Now() {
		t.Errorf("unexpected post created at: %v", newPost.CreatedAt)
	}

	if newPost.UpdatedAt.Time != timeProv.Now() {
		t.Errorf("unexpected post updated at: %v", newPost.UpdatedAt)
	}

	if len(newPost.Attachments) != 1 {
		t.Errorf("unexpected post attachments length: %d", len(newPost.Attachments))
	}

	if newPost.Attachments[0] != "default_avatar.png" {
		t.Errorf("unexpected post attachment: %s", newPost.Attachments[0])
	}
}

func TestUpdatePost(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	timeProv := customtime.MockTimeProvider{}

	row := pgxpoolmock.NewRow(
		uint(1),
		uint(1),
		"content",
		timeProv.Now(),
		timeProv.Now(),
	)

	pool.EXPECT().QueryRow(gomock.Any(), repository.UpdatePostQuery, gomock.Any(), gomock.Any()).Return(row)

	repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

	post := &domain.Post{
		ID:      1,
		Content: "content",
	}

	updatedPost, err := repo.UpdatePost(context.Background(), post)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if updatedPost.ID != 1 {
		t.Errorf("unexpected post id: %d", updatedPost.ID)
	}

	if updatedPost.AuthorID != 1 {
		t.Errorf("unexpected post author id: %d", updatedPost.AuthorID)
	}

	if updatedPost.Content != "content" {
		t.Errorf("unexpected post content: %s", updatedPost.Content)
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
				tag := pgconn.CommandTag("DELETE 0")
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

func TestErrRow_Scan(t *testing.T) {
	type args struct {
		in0 []interface{}
	}
	tests := []struct {
		name    string
		r       ErrRow
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Scan(tt.args.in0...); (err != nil) != tt.wantErr {
				t.Errorf("ErrRow.Scan() error = %v, wantErr %v", err, tt.wantErr)
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
			name:   "Test not found",
			postID: 2,
			mock: func(pool *pgxpoolmock.MockPgxIface, postID uint) {
				pool.EXPECT().Exec(gomock.Any(), repository.DeleteGroupPostQuery, postID).Return(pgconn.CommandTag("DELETE 0"), nil)
			},
			err: nil,
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
