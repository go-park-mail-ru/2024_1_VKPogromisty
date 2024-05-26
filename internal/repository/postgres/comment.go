package repository

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/pkg/contextlogger"
	"socio/pkg/utils"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const (
	getCommentsByPostIDQuery = `
	SELECT c.id,
		c.post_id,
		c.author_id,
		c.content,
		c.created_at,
		c.updated_at,
		array_agg(cl.user_id) AS liked_by
	FROM public.comment AS c
	LEFT JOIN public.comment_like cl ON cl.comment_id = c.id
	WHERE post_id = $1
    GROUP BY c.id, 
        c.post_id, 
        c.author_id, 
        c.content, 
        c.created_at, 
        c.updated_at;
	`
	getCommentByIDQuery = `
	SELECT id,
		post_id,
		author_id,
		content,
		created_at,
		updated_at
	FROM public.comment
	WHERE id = $1;
	`
	storeCommentQuery = `
	INSERT INTO public.comment (post_id, author_id, content)
	VALUES ($1, $2, $3)
	RETURNING id,
		post_id,
		author_id,
		content,
		created_at,
		updated_at;
	`
	updateCommentQuery = `
	UPDATE public.comment
	SET content = $1
	WHERE id = $2
	RETURNING id,
		post_id,
		author_id,
		content,
		created_at,
		updated_at;
	`
	deleteCommentQuery = `
	DELETE FROM public.comment
	WHERE id = $1;
	`
	getCommentLikeByCommentIDAndUserIDQuery = `
	SELECT id,
		comment_id,
		user_id,
		created_at
	FROM public.comment_like
	WHERE comment_id = $1
	AND user_id = $2;
	`
	storeCommentLikeQuery = `
	INSERT INTO public.comment_like (comment_id, user_id)
	VALUES ($1, $2)
	RETURNING id,
		comment_id,
		user_id,
		created_at;
	`
	deleteCommentLikeQuery = `
	DELETE FROM public.comment_like
	WHERE comment_id = $1
	AND user_id = $2;
	`
)

func (p *Posts) GetCommentsByPostID(ctx context.Context, postID uint) (comments []*domain.Comment, err error) {
	contextlogger.LogSQL(ctx, getCommentsByPostIDQuery, postID)

	rows, err := p.db.Query(context.Background(), getCommentsByPostIDQuery, postID)
	if err != nil {
		return
	}
	defer rows.Close()

	var likedByIDS pgtype.Int8Array

	for rows.Next() {
		comment := new(domain.Comment)

		err = rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.AuthorID,
			&comment.Content,
			&comment.CreatedAt.Time,
			&comment.UpdatedAt.Time,
			&likedByIDS,
		)
		if err != nil {
			return
		}

		comment.LikedByIDs = utils.Int8ArrayIntoUintSlice(likedByIDS)

		comments = append(comments, comment)
	}

	return
}

func (p *Posts) GetCommentByID(ctx context.Context, id uint) (comment *domain.Comment, err error) {
	comment = new(domain.Comment)

	err = p.db.QueryRow(
		context.Background(),
		getCommentByIDQuery,
		id,
	).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.AuthorID,
		&comment.Content,
		&comment.CreatedAt.Time,
		&comment.UpdatedAt.Time,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
		}

		return
	}

	return
}

func (p *Posts) StoreComment(ctx context.Context, comment *domain.Comment) (newComment *domain.Comment, err error) {
	newComment = new(domain.Comment)

	err = p.db.QueryRow(
		context.Background(),
		storeCommentQuery,
		comment.PostID,
		comment.AuthorID,
		comment.Content,
	).Scan(
		&newComment.ID,
		&newComment.PostID,
		&newComment.AuthorID,
		&newComment.Content,
		&newComment.CreatedAt.Time,
		&newComment.UpdatedAt.Time,
	)
	if err != nil {
		return
	}

	return
}

func (p *Posts) UpdateComment(ctx context.Context, comment *domain.Comment) (updatedComment *domain.Comment, err error) {
	_, err = p.GetCommentByID(ctx, comment.ID)
	if err != nil {
		return
	}

	updatedComment = new(domain.Comment)

	err = p.db.QueryRow(
		context.Background(),
		updateCommentQuery,
		comment.Content,
		comment.ID,
	).Scan(
		&updatedComment.ID,
		&updatedComment.PostID,
		&updatedComment.AuthorID,
		&updatedComment.Content,
		&updatedComment.CreatedAt.Time,
		&updatedComment.UpdatedAt.Time,
	)
	if err != nil {
		return
	}

	return
}

func (p *Posts) DeleteComment(ctx context.Context, id uint) (err error) {
	result, err := p.db.Exec(
		context.Background(),
		deleteCommentQuery,
		id,
	)
	if err != nil {
		return
	}

	if result.RowsAffected() == 0 {
		err = errors.ErrNotFound
		return
	}

	if result.RowsAffected() > 1 {
		err = errors.ErrRowsAffected
		return
	}

	return
}

func (p *Posts) GetCommentLikeByCommentIDAndUserID(ctx context.Context, data *domain.CommentLike) (commentLike *domain.CommentLike, err error) {
	commentLike = new(domain.CommentLike)

	err = p.db.QueryRow(
		context.Background(),
		getCommentLikeByCommentIDAndUserIDQuery,
		data.CommentID,
		data.UserID,
	).Scan(
		&commentLike.ID,
		&commentLike.CommentID,
		&commentLike.UserID,
		&commentLike.CreatedAt.Time,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
		}

		return
	}

	return
}

func (p *Posts) StoreCommentLike(ctx context.Context, commentLike *domain.CommentLike) (newCommentLike *domain.CommentLike, err error) {
	_, err = p.GetCommentLikeByCommentIDAndUserID(ctx, commentLike)
	if err == nil {
		err = errors.ErrInvalidData
		return
	}

	newCommentLike = new(domain.CommentLike)

	err = p.db.QueryRow(
		context.Background(),
		storeCommentLikeQuery,
		commentLike.CommentID,
		commentLike.UserID,
	).Scan(
		&newCommentLike.ID,
		&newCommentLike.CommentID,
		&newCommentLike.UserID,
		&newCommentLike.CreatedAt.Time,
	)
	if err != nil {
		return
	}

	return
}

func (p *Posts) DeleteCommentLike(ctx context.Context, commentLike *domain.CommentLike) (err error) {
	result, err := p.db.Exec(
		context.Background(),
		deleteCommentLikeQuery,
		commentLike.CommentID,
		commentLike.UserID,
	)
	if err != nil {
		return
	}

	if result.RowsAffected() == 0 {
		err = errors.ErrNotFound
		return
	}

	if result.RowsAffected() > 1 {
		err = errors.ErrRowsAffected
		return
	}

	return
}
