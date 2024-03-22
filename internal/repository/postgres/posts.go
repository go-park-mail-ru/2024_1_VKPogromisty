package repository

import (
	"database/sql"
	"socio/domain"
	customtime "socio/pkg/time"

	_ "github.com/lib/pq"
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
	rows, err := s.db.Query("SELECT filename FROM post_attachments WHERE post_id = $1;", postID)
	if err != nil {
		if err == sql.ErrNoRows {
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

	rerr := rows.Close()
	if rerr != nil {
		return
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (s *Posts) GetAll() (posts []*domain.Post, err error) {
	rows, err := s.db.Query("SELECT id, author_id, text, creation_date FROM posts;")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return
	}
	defer rows.Close()

	for rows.Next() {
		post := domain.Post{}
		err = rows.Scan(&post.ID, &post.AuthorID, &post.Text, &post.CreationDate.Time)
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

	rerr := rows.Close()
	if rerr != nil {
		return
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}
