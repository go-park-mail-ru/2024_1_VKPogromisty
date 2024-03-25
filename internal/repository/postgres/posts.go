package repository

import (
	"context"
	"mime/multipart"
	"socio/domain"
	"socio/errors"
	customtime "socio/pkg/time"
	"socio/utils"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

const (
	storePostQuery = `
	INSERT INTO public.post (author_id, content)
	VALUES ($1, $2)
	RETURNING id,
		author_id,
		content,
		created_at,
		updated_at;
	`
	storeAttachmentQuery = `
	INSERT INTO public.post_attachment (post_id, file_name)
	VALUES ($1, $2)
	RETURNING file_name;
	`
	selectAttachmentsQuery = `
	SELECT array_agg(file_name) AS attachments
	FROM public.post_attachment
	WHERE post_id = $1;
	`
	deletePostQuery = `
	DELETE FROM public.post
	WHERE id = $1;
	`
)

type PostWithAuthor struct {
	Post   *domain.Post
	Author *domain.User
}

type Posts struct {
	db *pgxpool.Pool
	TP customtime.TimeProvider
}

func NewPosts(db *pgxpool.Pool, tp customtime.TimeProvider) *Posts {
	return &Posts{
		db: db,
		TP: tp,
	}
}

func (p *Posts) GetUserPosts(userID uint) (posts []domain.Post, err error) {
	return
}

func (p *Posts) GetUserFriendsPosts(userID uint) (posts []domain.Post, err error) {
	return
}

func (p *Posts) StorePost(post *domain.Post, attachments []*multipart.FileHeader) (newPost *domain.Post, err error) {
	newPost = new(domain.Post)

	tx, err := p.db.BeginTx(context.Background(), pgx.TxOptions{})

	if err != nil {
		return
	}

	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(), storePostQuery, post.AuthorID, post.Content).Scan(
		&newPost.ID,
		&newPost.AuthorID,
		&newPost.Content,
		&newPost.CreatedAt.Time,
		&newPost.UpdatedAt.Time,
	)
	if err != nil {
		return
	}

	for _, attachment := range attachments {
		var fileName string
		fileName, err = utils.SaveImage(attachment)
		if err != nil {
			return
		}

		err = tx.QueryRow(context.Background(), storeAttachmentQuery, newPost.ID, fileName).Scan(&fileName)
		if err != nil {
			return
		}

		newPost.Attachments = append(newPost.Attachments, fileName)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return
	}

	return
}

func (p *Posts) UpdatePost(post *domain.Post) (updatedPost *domain.Post, err error) {
	return
}

func (p *Posts) DeletePost(postID uint) (err error) {
	tx, err := p.db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return
	}
	defer tx.Rollback(context.Background())

	var attachments []string
	err = tx.QueryRow(context.Background(), selectAttachmentsQuery, postID).Scan(&attachments)
	if err != nil && err != pgx.ErrNoRows {
		return
	}

	for _, attachment := range attachments {
		err = utils.RemoveImage(attachment)
		if err != nil {
			return
		}
	}

	result, err := tx.Exec(context.Background(), deletePostQuery, postID)
	if err != nil {
		return
	}

	if result.RowsAffected() != 1 {
		return errors.ErrRowsAffected
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return
	}

	return
}
