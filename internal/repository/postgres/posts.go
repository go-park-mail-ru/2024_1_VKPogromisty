package repository

import (
	"context"
	"mime/multipart"
	"socio/domain"
	"socio/errors"
	"socio/pkg/contextlogger"
	"socio/pkg/static"
	customtime "socio/pkg/time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

const (
	PostsByPage      = 20
	getPostByIDQuery = `
	SELECT p.id,
		p.author_id,
		p.content,
		p.created_at,
		p.updated_at,
		array_agg(pa.file_name) AS attachments
	FROM public.post AS p
		LEFT JOIN public.post_attachment AS pa ON p.id = pa.post_id
	WHERE p.id = $1
	GROUP BY p.id,
		p.author_id,
		p.content,
		p.created_at,
		p.updated_at;
	`
	getLastUserPostIDQuery = `
	SELECT COALESCE(MAX(id), 0) AS last_post_id
	FROM public.post
	WHERE author_id = $1;
	`
	getUserPostsQuery = `
	SELECT p.id,
		p.author_id,
		p.content,
		p.created_at,
		p.updated_at,
		array_agg(pa.file_name) AS attachments
	FROM public.post AS p
		LEFT JOIN public.post_attachment AS pa ON p.id = pa.post_id
	WHERE p.author_id = $1
		AND p.id < $2
	GROUP BY p.id,
		p.author_id,
		p.content,
		p.created_at,
		p.updated_at
	ORDER BY p.created_at DESC
	LIMIT $3;
	`
	getLastUserFriendsPostIDQuery = `
	SELECT COALESCE(MAX(p.id), 0) AS last_post_id
	FROM public.post AS p
		INNER JOIN public.subscription AS s ON p.author_id = s.subscribed_to_id
	WHERE s.subscriber_id = $1;
	`
	getUserFriendsPostsQuery = `
	SELECT p.id,
		p.author_id,
		p.content,
		p.created_at,
		p.updated_at,
		array_agg(pa.file_name) AS attachments,
		u.id AS user_id,
		u.first_name,
		u.last_name,
		u.email,
		u.avatar,
		u.date_of_birth,
		u.created_at AS user_created_at,
		u.updated_at AS user_updated_at
	FROM public.post AS p
		LEFT JOIN public.post_attachment AS pa ON p.id = pa.post_id
		INNER JOIN public.user AS u ON p.author_id = u.id
		INNER JOIN public.subscription AS s ON u.id = s.subscribed_to_id
	WHERE s.subscriber_id = $1
		AND p.id < $2
	GROUP BY p.id,
		p.author_id,
		p.content,
		p.created_at,
		p.updated_at,
		u.id,
		u.first_name,
		u.last_name,
		u.email,
		u.avatar,
		u.date_of_birth,
		u.created_at,
		u.updated_at
	ORDER BY p.created_at DESC
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

func textArrayIntoStringSlice(arr pgtype.TextArray) (res []string) {
	for _, v := range arr.Elements {
		if v.Status == pgtype.Present {
			res = append(res, v.String)
		}
	}

	return
}

func (p *Posts) GetPostByID(ctx context.Context, postID uint) (post *domain.Post, err error) {
	post = new(domain.Post)

	var attachments pgtype.TextArray

	contextlogger.LogSQL(ctx, getPostByIDQuery, postID)

	err = p.db.QueryRow(context.Background(), getPostByIDQuery, postID).Scan(
		&post.ID,
		&post.AuthorID,
		&post.Content,
		&post.CreatedAt.Time,
		&post.UpdatedAt.Time,
		&attachments,
	)
	if err != nil {
		return
	}

	post.Attachments = textArrayIntoStringSlice(attachments)

	return
}

func (p *Posts) GetUserPosts(ctx context.Context, userID uint, lastPostID uint) (posts []*domain.Post, err error) {
	if lastPostID == 0 {
		contextlogger.LogSQL(ctx, getLastUserPostIDQuery, userID)

		err = p.db.QueryRow(context.Background(), getLastUserPostIDQuery, userID).Scan(&lastPostID)
		if err != nil {
			return
		}
	}

	contextlogger.LogSQL(ctx, getUserPostsQuery, userID, lastPostID, PostsByPage)

	rows, err := p.db.Query(context.Background(), getUserPostsQuery, userID, lastPostID, PostsByPage)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		post := new(domain.Post)
		var attachments pgtype.TextArray
		err = rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Content,
			&post.CreatedAt.Time,
			&post.UpdatedAt.Time,
			&attachments,
		)
		if err != nil {
			return
		}

		post.Attachments = textArrayIntoStringSlice(attachments)

		posts = append(posts, post)
	}

	return
}

func (p *Posts) GetUserFriendsPosts(ctx context.Context, userID uint, lastPostID uint) (posts []domain.PostWithAuthor, err error) {
	if lastPostID == 0 {
		contextlogger.LogSQL(ctx, getLastUserFriendsPostIDQuery, userID)

		err = p.db.QueryRow(context.Background(), getLastUserFriendsPostIDQuery, userID).Scan(&lastPostID)
		if err != nil {
			return
		}
	}

	contextlogger.LogSQL(ctx, getUserFriendsPostsQuery, userID, lastPostID, PostsByPage)

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
		var attachments pgtype.TextArray
		err = rows.Scan(
			&post.Post.ID,
			&post.Post.AuthorID,
			&post.Post.Content,
			&post.Post.CreatedAt.Time,
			&post.Post.UpdatedAt.Time,
			&attachments,
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

		post.Post.Attachments = textArrayIntoStringSlice(attachments)

		posts = append(posts, post)
	}

	return
}

func (p *Posts) StorePost(ctx context.Context, post *domain.Post, attachments []*multipart.FileHeader) (newPost *domain.Post, err error) {
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

	contextlogger.LogSQL(ctx, storePostQuery, post.AuthorID, post.Content)

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
		fileName, err = static.SaveImage(attachment)
		if err != nil {
			return
		}

		contextlogger.LogSQL(ctx, storeAttachmentQuery, newPost.ID, fileName)

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

func (p *Posts) UpdatePost(ctx context.Context, post *domain.Post) (updatedPost *domain.Post, err error) {
	updatedPost = new(domain.Post)

	contextlogger.LogSQL(ctx, updatePostQuery, post.Content, post.ID)

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

func (p *Posts) DeletePost(ctx context.Context, postID uint) (err error) {
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

	var attachments pgtype.TextArray

	contextlogger.LogSQL(ctx, selectAttachmentsQuery, postID)

	err = tx.QueryRow(context.Background(), selectAttachmentsQuery, postID).Scan(&attachments)
	if err != nil && err != pgx.ErrNoRows {
		return
	}

	for _, v := range attachments.Elements {
		if v.Status == pgtype.Present {
			err = static.RemoveImage(v.String)
			if err != nil {
				return
			}
		}
	}

	contextlogger.LogSQL(ctx, deletePostQuery, postID)

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
