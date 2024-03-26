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
	PostsByPage      = 20
	getPostByIDQuery = `
	SELECT id,
		author_id,
		content,
		created_at,
		updated_at,
		attachments
	FROM public.post
		LEFT JOIN (
			SELECT post_id,
				array_agg(file_name) AS attachments
			FROM public.post_attachment
			GROUP BY post_id
		) AS post_attachments ON public.post.id = post_attachments.post_id
	WHERE id = $1;
	`
	getUserPostsQuery = `
	SELECT id,
		author_id,
		content,
		created_at,
		updated_at,
		attachments
	FROM public.post
		LEFT JOIN (
			SELECT post_id,
				array_agg(file_name) AS attachments
			FROM public.post_attachment
			GROUP BY post_id
		) AS post_attachments ON public.post.id = post_attachments.post_id
	WHERE author_id = $1
		AND id > $2
	ORDER BY created_at DESC
	LIMIT $3;
	`
	getUserFriendsPostsQuery = `
	SELECT public.post.id,
		author_id,
		content,
		public.post.created_at,
		public.post.updated_at,
		attachments,
		public.user.id AS user_id,
		public.user.first_name,
		public.user.last_name,
		public.user.email,
		public.user.avatar,
		public.user.date_of_birth,
		public.user.created_at AS user_created_at,
		public.user.updated_at AS user_updated_at
	FROM public.post
		LEFT JOIN (
			SELECT post_id,
				array_agg(file_name) AS attachments
			FROM public.post_attachment
			GROUP BY post_id
		) AS post_attachments ON public.post.id = post_attachments.post_id
		LEFT JOIN public.user ON public.post.author_id = public.user.id
	WHERE author_id IN (
			SELECT subscribed_to_id
			FROM public.subscription
			WHERE subscriber_id = $1
		)
		AND public.post.id > $2
	ORDER BY public.post.created_at DESC
	LIMIT $3;
	`
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
	updatePostQuery = `
	UPDATE public.post
	SET content = $1
	WHERE id = $2
	RETURNING id,
		author_id,
		content,
		created_at,
		updated_at;
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

func (p *Posts) GetPostByID(postID uint) (post *domain.Post, err error) {
	post = new(domain.Post)

	err = p.db.QueryRow(context.Background(), getPostByIDQuery, postID).Scan(
		&post.ID,
		&post.AuthorID,
		&post.Content,
		&post.CreatedAt.Time,
		&post.UpdatedAt.Time,
		&post.Attachments,
	)
	if err != nil {
		return
	}

	return
}

func (p *Posts) GetUserPosts(userID uint, lastPostID uint) (posts []*domain.Post, err error) {
	rows, err := p.db.Query(context.Background(), getUserPostsQuery, userID, lastPostID, PostsByPage)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		post := new(domain.Post)
		err = rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Content,
			&post.CreatedAt.Time,
			&post.UpdatedAt.Time,
			&post.Attachments,
		)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}

	return
}

func (p *Posts) GetUserFriendsPosts(userID uint, lastPostID uint) (posts []domain.PostWithAuthor, err error) {
	rows, err := p.db.Query(context.Background(), getUserFriendsPostsQuery, userID, lastPostID, PostsByPage)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		post := domain.PostWithAuthor{
			Post:   new(domain.Post),
			Author: new(domain.User),
		}
		err = rows.Scan(
			&post.Post.ID,
			&post.Post.AuthorID,
			&post.Post.Content,
			&post.Post.CreatedAt.Time,
			&post.Post.UpdatedAt.Time,
			&post.Post.Attachments,
			&post.Author.ID,
			&post.Author.FirstName,
			&post.Author.LastName,
			&post.Author.Email,
			&post.Author.Avatar,
			&post.Author.DateOfBirth.Time,
			&post.Author.CreatedAt.Time,
			&post.Author.UpdatedAt.Time,
		)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}

	return
}

func (p *Posts) StorePost(post *domain.Post, attachments []*multipart.FileHeader) (newPost *domain.Post, err error) {
	newPost = new(domain.Post)

	tx, err := p.db.BeginTx(context.Background(), pgx.TxOptions{})

	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			return
		}
		if err = tx.Rollback(context.Background()); err != nil && err != pgx.ErrTxClosed {
			return
		}

		err = nil
	}()

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
	updatedPost = new(domain.Post)

	err = p.db.QueryRow(context.Background(), updatePostQuery, post.Content, post.ID).Scan(
		&updatedPost.ID,
		&updatedPost.AuthorID,
		&updatedPost.Content,
		&updatedPost.CreatedAt.Time,
		&updatedPost.UpdatedAt.Time,
	)
	if err != nil {
		return
	}

	return
}

func (p *Posts) DeletePost(postID uint) (err error) {
	tx, err := p.db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			return
		}
		if err = tx.Rollback(context.Background()); err != nil && err != pgx.ErrTxClosed {
			return
		}

		err = nil
	}()

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
