package repository_test

import (
	"context"
	"mime/multipart"
	"socio/domain"
	repository "socio/internal/repository/postgres"
	customtime "socio/pkg/time"
	"testing"

	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
)

var (
	postColumns           = []string{"id", "author_id", "content", "created_at", "updated_at", "attacments"}
	postWithAuthorColumns = []string{"id", "author_id", "content", "created_at", "updated_at", "attacments", "user_id", "first_name", "last_name", "email", "avatar", "date_of_birth", "created_at", "updated_at"}
)

func TestGetPostByID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	timeProv := customtime.MockTimeProvider{}

	arr := pgtype.TextArray{
		Elements: []pgtype.Text{{String: "file1"}, {String: "file2"}},
	}

	row := pgxpoolmock.NewRow(
		uint(1),
		uint(2),
		"content",
		timeProv.Now(),
		timeProv.Now(),
		arr,
	)

	pool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(row)

	repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

	post, err := repo.GetPostByID(context.Background(), 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if post.ID != 1 {
		t.Errorf("unexpected post id: %d", post.ID)
	}
}

func TestGetUserPosts(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	timeProv := customtime.MockTimeProvider{}

	arr := pgtype.TextArray{
		Elements: []pgtype.Text{{String: "file1"}, {String: "file2"}},
	}

	row := pgxpoolmock.NewRows(postColumns).AddRow(
		uint(1),
		uint(2),
		"content",
		timeProv.Now(),
		timeProv.Now(),
		arr,
	).ToPgxRows()

	pool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(row, nil)

	repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

	posts, err := repo.GetUserPosts(context.Background(), 1, 2, 20)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("unexpected posts length: %d", len(posts))
	}
}

func TestGetUserFriendsPosts(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	timeProv := customtime.MockTimeProvider{}

	arr := pgtype.TextArray{
		Elements: []pgtype.Text{{String: "file1"}, {String: "file2"}},
	}

	row := pgxpoolmock.NewRows(postWithAuthorColumns).AddRow(
		uint(1),
		uint(2),
		"content",
		timeProv.Now(),
		timeProv.Now(),
		arr,
		uint(2),
		"firstname",
		"lastname",
		"uemail@email",
		"avatar",
		timeProv.Now(),
		timeProv.Now(),
		timeProv.Now(),
	).ToPgxRows()

	pool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(row, nil)

	repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

	posts, err := repo.GetUserFriendsPosts(context.Background(), 1, 2, 20)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("unexpected posts length: %d", len(posts))
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
	}

	newPost, err := repo.StorePost(context.Background(), post, []*multipart.FileHeader{
		nil,
	})
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	tag := pgconn.CommandTag("DELETE 1")

	attachmentRows := pgxpoolmock.NewRow(
		pgtype.TextArray{
			Elements: []pgtype.Text{
				{String: "default_avatar.png", Status: pgtype.Present},
			},
		},
	)

	pool.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(pool, nil)
	pool.EXPECT().QueryRow(gomock.Any(), repository.SelectAttachmentsQuery, gomock.Any()).Return(attachmentRows)
	pool.EXPECT().Exec(gomock.Any(), repository.DeletePostQuery, gomock.Any()).Return(tag, nil)
	pool.EXPECT().Rollback(gomock.Any()).Return(nil)
	pool.EXPECT().Commit(gomock.Any()).Return(nil)

	repo := repository.NewPosts(pool, customtime.MockTimeProvider{})

	err := repo.DeletePost(context.Background(), 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
