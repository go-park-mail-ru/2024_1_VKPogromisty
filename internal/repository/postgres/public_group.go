package repository

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/pkg/contextlogger"
	customtime "socio/pkg/time"
	publicgroup "socio/usecase/public_group"

	"github.com/jackc/pgx/v4"
)

const (
	getPublicGroupByIDWithInfoQuery = `
	SELECT pg.id, pg.name, pg.description, pg.avatar, pg.created_at, pg.updated_at, COUNT(pgs.id) AS subscribers_count, EXISTS (
		SELECT 1
		FROM public.public_group_subscription pgs
		WHERE pgs.public_group_id = pg.id AND pgs.subscriber_id = $2
	) AS is_subscribed
	FROM public.public_group pg
	LEFT JOIN public.public_group_subscription pgs ON pg.id = pgs.public_group_id
	WHERE pg.id = $1
	GROUP BY pg.id, pg.name, pg.description, pg.avatar, pg.created_at, pg.updated_at;
	`
	searchPublicGroupsByNameWithInfoQuery = `
	SELECT pg.id, pg.name, pg.description, pg.avatar, pg.created_at, pg.updated_at, COUNT(pgs.id) AS subscribers_count, EXISTS (
		SELECT 1
		FROM public.public_group_subscription pgs
		WHERE pgs.public_group_id = pg.id AND pgs.subscriber_id = $2
	) AS is_subscribed
	FROM public.public_group pg
	LEFT JOIN public.public_group_subscription pgs ON pg.id = pgs.public_group_id
	WHERE pg.name ILIKE '%' || $1 || '%'
	GROUP BY pg.id, pg.name, pg.description, pg.avatar, pg.created_at, pg.updated_at;
	`
	storePublicGroupQuery = `
	INSERT INTO public.public_group (name, description, avatar)
	VALUES ($1, $2, $3)
	RETURNING id, name, description, avatar, created_at, updated_at;
	`
	updatePublicGroupQuery = `
	UPDATE public.public_group
	SET name = $1, description = $2, avatar = $3
	WHERE id = $4
	RETURNING id, name, description, avatar, created_at, updated_at;
	`
	deletePublicGroupQuery = `
	DELETE FROM public.public_group
	WHERE id = $1;
	`
	getSubscriptionByPublicGroupIDAndSubscriberIDQuery = `
	SELECT id, public_group_id, subscriber_id, created_at, updated_at
	FROM public.public_group_subscription
	WHERE public_group_id = $1 AND subscriber_id = $2;
	`
	getPublicGroupsBySubscriberIDQuery = `
	SELECT pg.id,
		pg.name,
		pg.description,
		pg.avatar,
		pg.created_at,
		pg.updated_at,
		COUNT(DISTINCT pgs2.id) AS subscribers_count
	FROM public.public_group pg
		JOIN public.public_group_subscription pgs ON pg.id = pgs.public_group_id
		AND pgs.subscriber_id = $1
		LEFT JOIN public.public_group_subscription pgs2 ON pg.id = pgs2.public_group_id
	GROUP BY pg.id,
		pg.name,
		pg.description,
		pg.avatar,
		pg.created_at,
		pg.updated_at;
	`
	storePublicGroupSubscriptionQuery = `
	INSERT INTO public.public_group_subscription (public_group_id, subscriber_id)
	VALUES ($1, $2)
	RETURNING id, public_group_id, subscriber_id, created_at, updated_at;
	`
	deletePublicSubscriptionQuery = `
	DELETE FROM public.public_group_subscription
	WHERE public_group_id = $1 AND subscriber_id = $2;
	`
	getPublicGroupSubscriptionIDsQuery = `
	SELECT public_group_id
	FROM public.public_group_subscription
	WHERE subscriber_id = $1;
	`
)

type PublicGroup struct {
	db DBPool
	TP customtime.TimeProvider
}

func NewPublicGroup(db DBPool, tp customtime.TimeProvider) *PublicGroup {
	return &PublicGroup{
		db: db,
		TP: tp,
	}
}

func (p *PublicGroup) GetPublicGroupByID(ctx context.Context, groupID uint, userID uint) (publicGroupWithInfo *publicgroup.PublicGroupWithInfo, err error) {
	contextlogger.LogSQL(ctx, getPublicGroupByIDWithInfoQuery, groupID)

	publicGroupWithInfo = &publicgroup.PublicGroupWithInfo{
		PublicGroup:  new(domain.PublicGroup),
		IsSubscribed: false,
	}

	err = p.db.QueryRow(
		context.Background(),
		getPublicGroupByIDWithInfoQuery,
		groupID,
		userID,
	).Scan(
		&publicGroupWithInfo.PublicGroup.ID,
		&publicGroupWithInfo.PublicGroup.Name,
		&publicGroupWithInfo.PublicGroup.Description,
		&publicGroupWithInfo.PublicGroup.Avatar,
		&publicGroupWithInfo.PublicGroup.CreatedAt.Time,
		&publicGroupWithInfo.PublicGroup.UpdatedAt.Time,
		&publicGroupWithInfo.PublicGroup.SubscribersCount,
		&publicGroupWithInfo.IsSubscribed,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
		}

		return
	}

	return
}

func (p *PublicGroup) SearchPublicGroupsByNameWithInfo(ctx context.Context, query string, userID uint) (publicGroups []*publicgroup.PublicGroupWithInfo, err error) {
	contextlogger.LogSQL(ctx, searchPublicGroupsByNameWithInfoQuery, query, userID)

	rows, err := p.db.Query(
		context.Background(),
		searchPublicGroupsByNameWithInfoQuery,
		query,
		userID,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		publicGroupWithInfo := &publicgroup.PublicGroupWithInfo{
			PublicGroup: new(domain.PublicGroup),
		}

		err = rows.Scan(
			&publicGroupWithInfo.PublicGroup.ID,
			&publicGroupWithInfo.PublicGroup.Name,
			&publicGroupWithInfo.PublicGroup.Description,
			&publicGroupWithInfo.PublicGroup.Avatar,
			&publicGroupWithInfo.PublicGroup.CreatedAt.Time,
			&publicGroupWithInfo.PublicGroup.UpdatedAt.Time,
			&publicGroupWithInfo.PublicGroup.SubscribersCount,
			&publicGroupWithInfo.IsSubscribed,
		)
		if err != nil {
			return
		}

		publicGroups = append(publicGroups, publicGroupWithInfo)
	}

	return
}

func (p *PublicGroup) StorePublicGroup(ctx context.Context, publicGroup *domain.PublicGroup) (newGroup *domain.PublicGroup, err error) {
	contextlogger.LogSQL(ctx, storePublicGroupQuery, publicGroup.Name, publicGroup.Description, publicGroup.Avatar)

	newGroup = new(domain.PublicGroup)

	err = p.db.QueryRow(
		context.Background(),
		storePublicGroupQuery,
		publicGroup.Name,
		publicGroup.Description,
		publicGroup.Avatar,
	).Scan(
		&newGroup.ID,
		&newGroup.Name,
		&newGroup.Description,
		&newGroup.Avatar,
		&newGroup.CreatedAt.Time,
		&newGroup.UpdatedAt.Time,
	)
	if err != nil {
		return
	}

	return
}

func (p *PublicGroup) UpdatePublicGroup(ctx context.Context, publicGroup *domain.PublicGroup) (updatedGroup *domain.PublicGroup, err error) {
	contextlogger.LogSQL(ctx, updatePublicGroupQuery, publicGroup.Name, publicGroup.Description, publicGroup.Avatar, publicGroup.ID)

	updatedGroup = new(domain.PublicGroup)

	err = p.db.QueryRow(
		context.Background(),
		updatePublicGroupQuery,
		publicGroup.Name,
		publicGroup.Description,
		publicGroup.Avatar,
		publicGroup.ID,
	).Scan(
		&updatedGroup.ID,
		&updatedGroup.Name,
		&updatedGroup.Description,
		&updatedGroup.Avatar,
		&updatedGroup.CreatedAt.Time,
		&updatedGroup.UpdatedAt.Time,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
		}

		return
	}

	return
}

func (p *PublicGroup) DeletePublicGroup(ctx context.Context, groupID uint) (err error) {
	contextlogger.LogSQL(ctx, deletePublicGroupQuery, groupID)

	result, err := p.db.Exec(context.Background(), deletePublicGroupQuery, groupID)
	if err != nil {
		return
	}

	if result.RowsAffected() > 1 {
		err = errors.ErrRowsAffected
		return
	}

	return
}

func (p *PublicGroup) GetSubscriptionByPublicGroupIDAndSubscriberID(ctx context.Context, publicGroupID, subscriberID uint) (subscription *domain.PublicGroupSubscription, err error) {
	contextlogger.LogSQL(ctx, getSubscriptionByPublicGroupIDAndSubscriberIDQuery, publicGroupID, subscriberID)

	subscription = new(domain.PublicGroupSubscription)

	err = p.db.QueryRow(
		context.Background(),
		getSubscriptionByPublicGroupIDAndSubscriberIDQuery,
		publicGroupID,
		subscriberID,
	).Scan(
		&subscription.ID,
		&subscription.PublicGroupID,
		&subscription.SubscriberID,
		&subscription.CreatedAt.Time,
		&subscription.UpdatedAt.Time,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
		}

		return
	}

	return
}

func (p *PublicGroup) GetPublicGroupsBySubscriberID(ctx context.Context, subscriberID uint) (groups []*domain.PublicGroup, err error) {
	contextlogger.LogSQL(ctx, getPublicGroupsBySubscriberIDQuery, subscriberID)

	rows, err := p.db.Query(
		context.Background(),
		getPublicGroupsBySubscriberIDQuery,
		subscriberID,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		group := new(domain.PublicGroup)

		err = rows.Scan(
			&group.ID,
			&group.Name,
			&group.Description,
			&group.Avatar,
			&group.CreatedAt.Time,
			&group.UpdatedAt.Time,
			&group.SubscribersCount,
		)
		if err != nil {
			return
		}

		groups = append(groups, group)
	}

	return
}

func (p *PublicGroup) StorePublicGroupSubscription(ctx context.Context, publicGroupSubscription *domain.PublicGroupSubscription) (newSubscription *domain.PublicGroupSubscription, err error) {
	contextlogger.LogSQL(ctx, storePublicGroupSubscriptionQuery, publicGroupSubscription.PublicGroupID, publicGroupSubscription.SubscriberID)

	newSubscription = new(domain.PublicGroupSubscription)

	err = p.db.QueryRow(
		context.Background(),
		storePublicGroupSubscriptionQuery,
		publicGroupSubscription.PublicGroupID,
		publicGroupSubscription.SubscriberID,
	).Scan(
		&newSubscription.ID,
		&newSubscription.PublicGroupID,
		&newSubscription.SubscriberID,
		&newSubscription.CreatedAt.Time,
		&newSubscription.UpdatedAt.Time,
	)
	if err != nil {
		return
	}

	return
}

func (p *PublicGroup) DeletePublicGroupSubscription(ctx context.Context, subscription *domain.PublicGroupSubscription) (err error) {
	contextlogger.LogSQL(ctx, deletePublicSubscriptionQuery, subscription.PublicGroupID, subscription.SubscriberID)

	result, err := p.db.Exec(context.Background(), deletePublicSubscriptionQuery, subscription.PublicGroupID, subscription.SubscriberID)
	if err != nil {
		return
	}

	if result.RowsAffected() > 1 {
		err = errors.ErrRowsAffected
		return
	}

	return
}

func (p *PublicGroup) GetPublicGroupSubscriptionIDs(ctx context.Context, userID uint) (subIDs []uint, err error) {
	contextlogger.LogSQL(ctx, getPublicGroupSubscriptionIDsQuery, userID)

	rows, err := p.db.Query(context.Background(), getPublicGroupSubscriptionIDsQuery, userID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var subID uint

		err = rows.Scan(&subID)
		if err != nil {
			return
		}

		subIDs = append(subIDs, subID)
	}

	return
}
