package repository

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/pkg/contextlogger"
	customtime "socio/pkg/time"
	"socio/pkg/utils"
	"socio/usecase/posts"

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
	GetLastUserPostIDQuery = `
	SELECT COALESCE(MAX(p.id), 0) AS last_post_id
	FROM public.post AS p
	LEFT JOIN public.public_group_post AS pgp ON pgp.post_id = p.id
	WHERE p.author_id = $1 
		AND pgp.post_id IS NULL;
	`
	GetUserPostsQuery = `
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
		LEFT JOIN public.public_group_post AS pgp ON pgp.post_id = p.id
    WHERE p.author_id = $1
        AND p.id < $2
		AND pgp.post_id IS NULL
    GROUP BY p.id,
        p.author_id,
        p.content,
        p.created_at,
        p.updated_at
    ORDER BY p.created_at DESC
    LIMIT $3;
	`
	GetLastUserFriendsPostIDQuery = `
	SELECT COALESCE(MAX(p.id), 0) AS last_post_id
	FROM public.post AS p
		INNER JOIN public.subscription AS s ON p.author_id = s.subscribed_to_id
	WHERE s.subscriber_id = $1;
	`
	GetUserFriendsPostsQuery = `
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
        INNER JOIN public.subscription AS s ON p.author_id = s.subscribed_to_id
		LEFT JOIN public.public_group_post AS pgp ON pgp.post_id = p.id
    WHERE s.subscriber_id = $1
        AND p.id < $2
		AND pgp.post_id IS NULL
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
	CreatePostLikeQuery = `
	INSERT INTO public.post_like (post_id, user_id)
	VALUES ($1, $2)
	RETURNING id,
		post_id,
		user_id,
		created_at;
	`
	DeletePostLikeQuery = `
	DELETE FROM public.post_like
	WHERE post_id = $1
		AND user_id = $2;
	`
	GetLikedPosts = `
	SELECT
		pl.id as like_id, 
		pl.post_id,
		pl.user_id,
		pl.created_at,
		p.author_id,
		p.content,
		p.created_at as post_created_at,
		p.updated_at as post_updated_at,
		array_agg(DISTINCT pa.file_name) AS attachments,
		array_agg(DISTINCT pl1.user_id) AS liked_by_users
	FROM public.post_like AS pl
	JOIN public.post AS p ON pl.post_id = p.id
	LEFT JOIN public.post_like AS pl1 ON p.id = pl1.post_id
	LEFT JOIN public.post_attachment AS pa ON p.id = pa.post_id
	WHERE p.author_id = $1
		AND pl.id < $2
	GROUP BY pl.id, 
        p.id,
		pl.post_id,
		pl.user_id,
		pl.created_at,
		p.author_id,
		p.content,
		p.created_at,
		p.updated_at
	ORDER BY pl.created_at DESC
	LIMIT $3;
	`
	GetLastPostLikeIDQuery = `
	SELECT COALESCE(MAX(pl.id), 0) AS last_like_id
	FROM public.post_like AS pl
	JOIN public.post AS p ON post_id = p.id
	WHERE p.author_id = $1;
	`
	StoreGroupPostQuery = `
	INSERT INTO public.public_group_post (post_id, public_group_id)
	VALUES ($1, $2)
	RETURNING id,
		post_id,
		public_group_id,
		created_at,
		updated_at;
	`
	DeleteGroupPostQuery = `
	DELETE FROM public.public_group_post
	WHERE post_id = $1;
	`
	GetLastPostOfGroupIDQuery = `
	SELECT COALESCE(MAX(p.id), 0) AS last_post_id
	FROM public.post AS p
		INNER JOIN public.public_group_post AS pgp ON p.id = pgp.post_id
		WHERE pgp.public_group_id = $1;
	`
	GetPostsOfGroupQuery = `
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
		LEFT JOIN public.public_group_post AS pgp ON p.id = pgp.post_id
		WHERE pgp.public_group_id = $1
			AND p.id < $2
		GROUP BY p.id,
			p.author_id,
			p.content,
			p.created_at,
			p.updated_at
		ORDER BY p.created_at DESC
		LIMIT $3;
	`
	GetGroupPostsBySubscriptionIDsQuery = `
	SELECT p.id,
		p.author_id,
		p.content,
		p.created_at,
		p.updated_at,
		array_agg(DISTINCT pa.file_name) AS attachments,
		array_agg(DISTINCT pl.user_id) AS liked_by_users,
		pgp.public_group_id
		FROM public.post AS p
		LEFT JOIN public.post_attachment AS pa ON p.id = pa.post_id
		LEFT JOIN public.post_like AS pl ON p.id = pl.post_id
		LEFT JOIN public.public_group_post AS pgp ON p.id = pgp.post_id
		WHERE pgp.public_group_id IN ($1)
			AND p.id < $2
			GROUP BY p.id,
			p.author_id,
			p.content,
			p.created_at,
			p.updated_at,
			pgp.public_group_id
			ORDER BY p.created_at DESC
			LIMIT $3;`
	GetLastGroupPostBySubscriptionIDsQuery = `
	SELECT COALESCE(MAX(p.id), 0) AS last_post_id
	FROM public.post AS p
		LEFT JOIN public.public_group_post AS pgp ON p.id = pgp.post_id
		WHERE pgp.public_group_id IN ($1);
	`
	GetPostsByGroupSubIDsAndUserSubIDsQuery = `
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
		LEFT JOIN public.public_group_post AS pgp ON p.id = pgp.post_id
		WHERE pgp.public_group_id IN ($1) OR p.author_id IN ($2)
			AND p.id < $3
			GROUP BY p.id,
			p.author_id,
			p.content,
			p.created_at,
			p.updated_at
			ORDER BY p.created_at DESC
			LIMIT $4;`
	GetLastPostByGroupSubIDsAndUserSubIDsQuery = `
	SELECT COALESCE(MAX(p.id), 0) AS last_post_id
	FROM public.post AS p
		LEFT JOIN public.public_group_post AS pgp ON p.id = pgp.post_id
		WHERE pgp.public_group_id IN ($1) OR p.author_id IN ($2);
	`
	GetLastPostIDQuery = `
	SELECT COALESCE(MAX(id), 0) AS last_post_id
	FROM public.post;
	`
	GetNewPostsQuery = `
	SELECT p.id,
		p.author_id,
		p.content,
		p.created_at,
		p.updated_at,
		array_agg(DISTINCT pa.file_name) AS attachments,
		array_agg(DISTINCT pl.user_id) AS liked_by_users,
		COALESCE(pgp.public_group_id, 0) AS group_id
		FROM public.post AS p
		LEFT JOIN public.post_attachment AS pa ON p.id = pa.post_id
		LEFT JOIN public.post_like AS pl ON p.id = pl.post_id
		LEFT JOIN public.public_group_post AS pgp ON p.id = pgp.post_id
		WHERE p.id < $1
		GROUP BY p.id,
			p.author_id,
			p.content,
			p.created_at,
			p.updated_at,
			pgp.public_group_id
			ORDER BY p.created_at DESC
			LIMIT $2;
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
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
		}

		return
	}

	post.Attachments = textArrayIntoStringSlice(attachments)
	post.LikedByIDs = int8ArrayIntoUintSlice(likedByUsers)

	return
}

func (p *Posts) GetUserPosts(ctx context.Context, userID uint, lastPostID uint, postsAmount uint) (posts []*domain.Post, err error) {
	if lastPostID == 0 {
		contextlogger.LogSQL(ctx, GetLastUserPostIDQuery, userID)

		err = p.db.QueryRow(context.Background(), GetLastUserPostIDQuery, userID).Scan(&lastPostID)
		if err != nil {
			return
		}

		lastPostID += 1
	}

	contextlogger.LogSQL(ctx, GetUserPostsQuery, userID, lastPostID, postsAmount)

	rows, err := p.db.Query(context.Background(), GetUserPostsQuery, userID, lastPostID, postsAmount)
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
		contextlogger.LogSQL(ctx, GetLastUserFriendsPostIDQuery, userID)

		err = p.db.QueryRow(context.Background(), GetLastUserFriendsPostIDQuery, userID).Scan(&lastPostID)
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

func (p *Posts) StoreGroupPost(ctx context.Context, groupPost *domain.GroupPost) (newGroupPost *domain.GroupPost, err error) {
	newGroupPost = new(domain.GroupPost)

	contextlogger.LogSQL(ctx, StoreGroupPostQuery, groupPost.PostID, groupPost.GroupID)

	err = p.db.QueryRow(context.Background(), StoreGroupPostQuery, groupPost.PostID, groupPost.GroupID).Scan(
		&newGroupPost.ID,
		&newGroupPost.PostID,
		&newGroupPost.GroupID,
		&newGroupPost.CreatedAt.Time,
		&newGroupPost.UpdatedAt.Time,
	)
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

	if result.RowsAffected() > 1 {
		return errors.ErrRowsAffected
	}

	return
}

func (p *Posts) DeleteGroupPost(ctx context.Context, postID uint) (err error) {
	contextlogger.LogSQL(ctx, DeleteGroupPostQuery, postID)

	result, err := p.db.Exec(context.Background(), DeleteGroupPostQuery, postID)
	if err != nil {
		return
	}

	if result.RowsAffected() > 1 {
		return errors.ErrRowsAffected
	}

	return
}

func (p *Posts) GetLikedPosts(ctx context.Context, userID uint, lastLikeID uint, limit uint) (likedPosts []posts.LikeWithPost, err error) {
	contextlogger.LogSQL(ctx, GetUserFriendsPostsQuery, userID, lastLikeID, limit)
	if lastLikeID == 0 {
		contextlogger.LogSQL(ctx, GetLastPostLikeIDQuery, userID)

		err = p.db.QueryRow(context.Background(), GetLastPostLikeIDQuery, userID).Scan(&lastLikeID)
		if err != nil {
			return
		}

		lastLikeID += 1
	}

	contextlogger.LogSQL(ctx, GetLikedPosts, userID, lastLikeID, limit)

	rows, err := p.db.Query(context.Background(), GetLikedPosts, userID, lastLikeID, limit)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		post := new(domain.Post)
		like := new(domain.PostLike)

		var attachments pgtype.TextArray
		var likedByUsers pgtype.Int8Array

		err = rows.Scan(
			&like.ID,
			&like.PostID,
			&like.UserID,
			&like.CreatedAt.Time,
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

		post.ID = like.PostID
		post.Attachments = textArrayIntoStringSlice(attachments)
		post.LikedByIDs = int8ArrayIntoUintSlice(likedByUsers)

		likedPosts = append(likedPosts, posts.LikeWithPost{
			Like: like,
			Post: post,
		})
	}

	return
}

func (p *Posts) StorePostLike(ctx context.Context, likeData *domain.PostLike) (like *domain.PostLike, err error) {
	like = new(domain.PostLike)

	contextlogger.LogSQL(ctx, CreatePostLikeQuery, likeData.PostID, likeData.UserID)

	err = p.db.QueryRow(context.Background(), CreatePostLikeQuery, likeData.PostID, likeData.UserID).Scan(
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
	contextlogger.LogSQL(ctx, DeletePostLikeQuery, likeData.PostID, likeData.UserID)

	result, err := p.db.Exec(context.Background(), DeletePostLikeQuery, likeData.PostID, likeData.UserID)
	if err != nil {
		return
	}

	if result.RowsAffected() > 1 {
		return errors.ErrRowsAffected
	}

	return
}

func (p *Posts) GetPostsOfGroup(ctx context.Context, groupID, lastPostID, postsAmount uint) (posts []*domain.Post, err error) {
	if lastPostID == 0 {
		contextlogger.LogSQL(ctx, GetLastPostOfGroupIDQuery, groupID)

		err = p.db.QueryRow(context.Background(), GetLastPostOfGroupIDQuery, groupID).Scan(&lastPostID)
		if err != nil {
			return
		}

		lastPostID += 1
	}

	contextlogger.LogSQL(ctx, GetPostsOfGroupQuery, groupID, lastPostID, postsAmount)

	rows, err := p.db.Query(context.Background(), GetPostsOfGroupQuery, groupID, lastPostID, postsAmount)
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

func (p *Posts) GetGroupPostsBySubscriptionIDs(ctx context.Context, subIDs []uint, lastPostID, postsAmount uint) (posts []*domain.Post, err error) {
	if len(subIDs) == 0 {
		subIDs = append(subIDs, 0)
	}

	subIDsStr := utils.UintArrayIntoString(subIDs)

	if lastPostID == 0 {
		contextlogger.LogSQL(ctx, GetLastGroupPostBySubscriptionIDsQuery, subIDsStr)

		err = p.db.QueryRow(context.Background(), GetLastGroupPostBySubscriptionIDsQuery, subIDsStr).Scan(&lastPostID)
		if err != nil {
			return
		}

		lastPostID += 1
	}

	contextlogger.LogSQL(ctx, GetGroupPostsBySubscriptionIDsQuery, subIDsStr, lastPostID, postsAmount)

	rows, err := p.db.Query(context.Background(), GetGroupPostsBySubscriptionIDsQuery, subIDsStr, lastPostID, postsAmount)
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
			&post.GroupID,
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

func (p *Posts) GetPostsByGroupSubIDsAndUserSubIDs(ctx context.Context, groupSubIDs, userSubIDs []uint, lastPostID, postsAmount uint) (posts []*domain.Post, err error) {
	if len(groupSubIDs) == 0 {
		groupSubIDs = append(groupSubIDs, 0)
	}

	if len(userSubIDs) == 0 {
		userSubIDs = append(userSubIDs, 0)
	}

	groupSubIDsStr := utils.UintArrayIntoString(groupSubIDs)
	userSubIDsStr := utils.UintArrayIntoString(userSubIDs)

	if lastPostID == 0 {
		contextlogger.LogSQL(ctx, GetLastPostByGroupSubIDsAndUserSubIDsQuery, groupSubIDsStr, userSubIDsStr)

		err = p.db.QueryRow(context.Background(), GetLastPostByGroupSubIDsAndUserSubIDsQuery, groupSubIDsStr, userSubIDsStr).Scan(&lastPostID)
		if err != nil {
			return
		}

		lastPostID += 1
	}

	contextlogger.LogSQL(ctx, GetPostsByGroupSubIDsAndUserSubIDsQuery, groupSubIDsStr, userSubIDsStr, lastPostID, postsAmount)

	rows, err := p.db.Query(context.Background(), GetPostsByGroupSubIDsAndUserSubIDsQuery, groupSubIDsStr, userSubIDsStr, lastPostID, postsAmount)
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

func (p *Posts) GetNewPosts(ctx context.Context, lastPostID, postsAmount uint) (posts []*domain.Post, err error) {
	if lastPostID == 0 {
		contextlogger.LogSQL(ctx, GetLastPostIDQuery)

		err = p.db.QueryRow(context.Background(), GetLastPostIDQuery).Scan(&lastPostID)
		if err != nil {
			return
		}

		lastPostID += 1
	}

	contextlogger.LogSQL(ctx, GetNewPostsQuery, lastPostID, postsAmount)

	rows, err := p.db.Query(context.Background(), GetNewPostsQuery, lastPostID, postsAmount)
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
			&post.GroupID,
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
