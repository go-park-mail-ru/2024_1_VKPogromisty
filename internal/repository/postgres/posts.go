package repository

import (
	"database/sql"
	"socio/domain"
	"socio/errors"
	customtime "socio/pkg/time"
)

type Posts struct {
	db *sql.DB
	TP customtime.TimeProvider
}

func NewPosts(db *sql.DB, tp customtime.TimeProvider) *Posts {
	return &Posts{
		db: db,
		TP: tp,
	}
}

func (s *Posts) getAttachments(postID uint) (attachments []string, err error) {
	rows, err := s.db.Query("SELECT filename FROM post_attachments WHERE post_id = $1", postID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var attachment string
		err = rows.Scan(&attachment)
		if err != nil {
			return nil, errors.ErrInternal
		}
		attachments = append(attachments, attachment)
	}

	return
}

func (s *Posts) GetAll() (posts []*domain.Post, err error) {
	rows, err := s.db.Query("SELECT id, author_id, text, creation_date FROM posts")
	if err != nil {
		return nil, errors.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		post := domain.Post{}
		err = rows.Scan(&post.ID, &post.AuthorID, &post.Text, &post.CreationDate.Time)
		if err != nil {
			return nil, errors.ErrInternal
		}
		attachments, err := s.getAttachments(post.ID)
		if err != nil {
			return nil, err
		}
		post.Attachments = attachments

		posts = append(posts, &post)
	}

	return
}
