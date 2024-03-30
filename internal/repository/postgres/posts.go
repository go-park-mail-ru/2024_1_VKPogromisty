package repository

import (
	"context"
	"socio/domain"
	customtime "socio/pkg/time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

const (
	getAllPostsQuery = `
	SELECT id,
		author_id,
		content,
		created_at,
		updated_at
	FROM public.post;
	`
	getAttachmentFilenameQuery = `
	SELECT file_name
	FROM public.post_attachments
	WHERE post_id = $1;
	`
)

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

func (s *Posts) getAttachments(postID uint) (attachments []string, err error) {
	rows, err := s.db.Query(context.Background(), getAttachmentFilenameQuery, postID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return
	}

	defer rows.Close()

	for rows.Next() {
		var attachment string
		err = rows.Scan(&attachment)
		if err != nil {
			return
		}
		attachments = append(attachments, attachment)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (s *Posts) GetAll() (posts []*domain.Post, err error) {
	rows, err := s.db.Query(context.Background(), getAllPostsQuery)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return
	}
	defer rows.Close()

	for rows.Next() {
		post := domain.Post{}
		err = rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.CreatedAt.Time, &post.UpdatedAt.Time)
		if err != nil {
			return
		}
		attachments, err := s.getAttachments(post.ID)
		if err != nil {
			return nil, err
		}
		post.Attachments = attachments

		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}
