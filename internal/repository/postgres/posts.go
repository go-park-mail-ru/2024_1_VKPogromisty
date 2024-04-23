package repository

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/pkg/contextlogger"
	customtime "socio/pkg/time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

const (
	GetPostByIDQuery = `
	SELECT p.id,
        p.author_id,
        p.content,
        p.created_at,
        p.updated_at,
        array_agg(DISTINCT pa.file_name) AS attachments,
        array_agg(DISTINCT pl.user_id) AS liked_by_users
    FROM public.post AS p
        LEFT JOIN public.post_attachment AS pa ON p.id = pa.post_id
        LEFT JOIN public.post_like AS pl ON p.id = pl.post_id
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
        array_agg(DISTINCT pa.file_name) AS attachments,
        array_agg(DISTINCT pl.user_id) AS liked_by_users
    FROM public.post AS p
        LEFT JOIN public.post_attachment AS pa ON p.id = pa.post_id
        LEFT JOIN public.post_like AS pl ON p.id = pl.post_id
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
	GetUserFriendsPostsQuery = `
	SELECT p.id,
        p.content,
        p.created_at,
        p.updated_at,
        array_agg(DISTINCT pa.file_name) AS attachments,
        array_agg(DISTINCT pl.user_id) AS liked_by_users
    FROM public.post AS p
        LEFT JOIN public.post_attachment AS pa ON p.id = pa.post_id
        LEFT JOIN public.post_like AS pl ON p.id = pl.post_id
        INNER JOIN public.subscription AS s ON p.author_id = s.subscribed_to_id
    WHERE s.subscriber_id = $1
        AND p.id < $2
    GROUP BY p.id,
        p.content,
        p.created_at,
        p.updated_at
    ORDER BY p.created_at DESC
    LIMIT $3;
	`
	StorePostQuery = `
	INSERT INTO public.post (author_id, content)
	VALUES ($1, $2)
	RETURNING id,
		author_id,
		content,
		created_at,
		updated_at;
	`
	StoreAttachmentQuery = `
	INSERT INTO public.post_attachment (post_id, file_name)
	VALUES ($1, $2)
	RETURNING file_name;
	`
	UpdatePostQuery = `
	UPDATE public.post
	SET content = $1
	WHERE id = $2
	RETURNING id,
		author_id,
		content,
		created_at,
		updated_at;
	`
	SelectAttachmentsQuery = `
	SELECT array_agg(file_name) AS attachments
	FROM public.post_attachment
	WHERE post_id = $1;
	`
	DeletePostQuery = `
	DELETE FROM public.post
	WHERE id = $1;
	`
	createPostLikeQuery = `
	INSERT INTO public.post_like (post_id, user_id)
	VALUES ($1, $2)
	RETURNING id,
		post_id,
		user_id,
		created_at;
	`
	deletePostLikeQuery = `
	DELETE FROM public.post_like
	WHERE post_id = $1
		AND user_id = $2;
	`
)

type Posts struct {
	db DBPool
	TP customtime.TimeProvider
}

func NewPosts(db DBPool, tp customtime.TimeProvider) *Posts {
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

func int8ArrayIntoUintSlice(arr pgtype.Int8Array) (res []uint64) {
	for _, v := range arr.Elements {
		if v.Status == pgtype.Present {
			res = append(res, uint64(v.Int))
		}
	}

	return
}

func (p *Posts) GetPostByID(ctx context.Context, postID uint) (post *domain.Post, err error) {
	post = new(domain.Post)

	var attachments pgtype.TextArray
	var likedByUsers pgtype.Int8Array

	contextlogger.LogSQL(ctx, GetPostByIDQuery, postID)

	err = p.db.QueryRow(context.Background(), GetPostByIDQuery, postID).Scan(
		&post.ID,
		&post.AuthorID,
		&post.Content,
		&post.CreatedAt.Time,
		&post.UpdatedAt.Time,
		&attachments,
		&likedByUsers,
	)
	if err != nil {
		return
	}

	post.Attachments = textArrayIntoStringSlice(attachments)
	post.LikedByIDs = int8ArrayIntoUintSlice(likedByUsers)

	return
}

func (p *Posts) GetUserPosts(ctx context.Context, userID uint, lastPostID uint, postsAmount uint) (posts []*domain.Post, err error) {
	if lastPostID == 0 {
		contextlogger.LogSQL(ctx, getLastUserPostIDQuery, userID)

		err = p.db.QueryRow(context.Background(), getLastUserPostIDQuery, userID).Scan(&lastPostID)
		if err != nil {
			return
		}

		lastPostID += 1
	}

	contextlogger.LogSQL(ctx, getUserPostsQuery, userID, lastPostID, postsAmount)

	rows, err := p.db.Query(context.Background(), getUserPostsQuery, userID, lastPostID, postsAmount)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		post := new(domain.Post)

		var attachments pgtype.TextArray
		var likedByUsers pgtype.Int8Array

		err = rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Content,
			&post.CreatedAt.Time,
			&post.UpdatedAt.Time,
			&attachments,
			&likedByUsers,
		)
		if err != nil {
			return
		}

		post.Attachments = textArrayIntoStringSlice(attachments)
		post.LikedByIDs = int8ArrayIntoUintSlice(likedByUsers)

		posts = append(posts, post)
	}

	return
}

func (p *Posts) GetUserFriendsPosts(ctx context.Context, userID uint, lastPostID uint, postsAmount uint) (posts []*domain.Post, err error) {
	contextlogger.LogSQL(ctx, GetUserFriendsPostsQuery, userID, lastPostID, postsAmount)
	if lastPostID == 0 {
		contextlogger.LogSQL(ctx, getLastUserFriendsPostIDQuery, userID)

		err = p.db.QueryRow(context.Background(), getLastUserFriendsPostIDQuery, userID).Scan(&lastPostID)
		if err != nil {
			return
		}

		lastPostID += 1
	}

	contextlogger.LogSQL(ctx, GetUserFriendsPostsQuery, userID, lastPostID, postsAmount)

	rows, err := p.db.Query(context.Background(), GetUserFriendsPostsQuery, userID, lastPostID, postsAmount)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		post := new(domain.Post)

		var attachments pgtype.TextArray
		var likedByUsers pgtype.Int8Array

		err = rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Content,
			&post.CreatedAt.Time,
			&post.UpdatedAt.Time,
			&attachments,
			&likedByUsers,
		)
		if err != nil {
			return
		}

		post.Attachments = textArrayIntoStringSlice(attachments)
		post.LikedByIDs = int8ArrayIntoUintSlice(likedByUsers)

		posts = append(posts, post)
	}

	return
}

func (p *Posts) StorePost(ctx context.Context, post *domain.Post) (newPost *domain.Post, err error) {
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

	contextlogger.LogSQL(ctx, StorePostQuery, post.AuthorID, post.Content)

	err = tx.QueryRow(context.Background(), StorePostQuery, post.AuthorID, post.Content).Scan(
		&newPost.ID,
		&newPost.AuthorID,
		&newPost.Content,
		&newPost.CreatedAt.Time,
		&newPost.UpdatedAt.Time,
	)
	if err != nil {
		return
	}

	for _, attachment := range post.Attachments {
		contextlogger.LogSQL(ctx, StoreAttachmentQuery, newPost.ID, attachment)

		err = tx.QueryRow(context.Background(), StoreAttachmentQuery, newPost.ID, attachment).Scan(&attachment)
		if err != nil {
			return
		}

		newPost.Attachments = append(newPost.Attachments, attachment)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return
	}

	return
}

func (p *Posts) UpdatePost(ctx context.Context, post *domain.Post) (updatedPost *domain.Post, err error) {
	updatedPost = new(domain.Post)

	contextlogger.LogSQL(ctx, UpdatePostQuery, post.Content, post.ID)

	err = p.db.QueryRow(context.Background(), UpdatePostQuery, post.Content, post.ID).Scan(
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
	contextlogger.LogSQL(ctx, DeletePostQuery, postID)

	result, err := p.db.Exec(context.Background(), DeletePostQuery, postID)
	if err != nil {
		return
	}

	if result.RowsAffected() != 1 {
		return errors.ErrRowsAffected
	}

	return
}

func (p *Posts) StorePostLike(ctx context.Context, likeData *domain.PostLike) (like *domain.PostLike, err error) {
	like = new(domain.PostLike)

	contextlogger.LogSQL(ctx, createPostLikeQuery, likeData.PostID, likeData.UserID)

	err = p.db.QueryRow(context.Background(), createPostLikeQuery, likeData.PostID, likeData.UserID).Scan(
		&like.ID,
		&like.PostID,
		&like.UserID,
		&like.CreatedAt.Time,
	)
	if err != nil {
		return
	}

	return
}

func (p *Posts) DeletePostLike(ctx context.Context, likeData *domain.PostLike) (err error) {
	contextlogger.LogSQL(ctx, deletePostLikeQuery, likeData.PostID, likeData.UserID)

	result, err := p.db.Exec(context.Background(), deletePostLikeQuery, likeData.PostID, likeData.UserID)
	if err != nil {
		return
	}

	if result.RowsAffected() != 1 {
		return errors.ErrRowsAffected
	}

	return
}
