package repository

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/pkg/contextlogger"
	customtime "socio/pkg/time"

	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

const (
	getSubscriptionsQuery = `
	SELECT DISTINCT id,
		first_name,
		last_name,
		email,
		avatar,
		date_of_birth,
		created_at,
		updated_at
	FROM (
			SELECT sub1.subscribed_to_id AS subscription_id
			FROM public.subscription sub1
				LEFT JOIN public.subscription sub2 ON sub1.subscribed_to_id = sub2.subscriber_id
				AND sub1.subscriber_id = sub2.subscribed_to_id
			WHERE sub1.subscriber_id = $1
				AND sub2.id IS null
		) AS sub_id
		JOIN public.user AS u ON u.id = sub_id.subscription_id;
	`
	getFriendsQuery = `
	SELECT DISTINCT id,
		first_name,
		last_name,
		email,
		avatar,
		date_of_birth,
		created_at,
		updated_at
	FROM (
			SELECT sub1.subscribed_to_id AS friend_id
			FROM public.subscription sub1
				INNER JOIN public.subscription sub2 ON sub1.subscribed_to_id = sub2.subscriber_id
				AND sub1.subscriber_id = sub2.subscribed_to_id
			WHERE sub1.subscriber_id = $1
		) AS friend_id
		JOIN public.user AS u ON u.id = friend_id.friend_id;
	`
	getSubscribersQuery = `
	SELECT DISTINCT id,
		first_name,
		last_name,
		email,
		avatar,
		date_of_birth,
		created_at,
		updated_at
	FROM (
			SELECT sub2.subscriber_id AS subscriber_id
			FROM public.subscription sub1
				RIGHT JOIN public.subscription sub2 ON sub1.subscribed_to_id = sub2.subscriber_id
				AND sub1.subscriber_id = sub2.subscribed_to_id
			WHERE sub2.subscribed_to_id = $1
				AND sub1.id IS NULL
		) AS subscriber_id
		JOIN public.user AS u ON u.id = subscriber_id.subscriber_id;
	`
	storeSubscriptionQuery = `
	INSERT INTO public.subscription (subscriber_id, subscribed_to_id)
	VALUES ($1, $2) 
	ON CONFLICT (subscriber_id, subscribed_to_id) DO NOTHING
	RETURNING id,
		subscriber_id,
		subscribed_to_id,
		created_at,
		updated_at;
	`
	deleteSubscriptionQuery = `
	DELETE FROM public.subscription
	WHERE subscriber_id = $1
		AND subscribed_to_id = $2;
	`
	getSubscriptionBySubscriberAndSubscribedToIDQuery = `
	SELECT id,
		subscriber_id,
		subscribed_to_id,
		created_at,
		updated_at
	FROM public.subscription
	WHERE subscriber_id = $1
		AND subscribed_to_id = $2;
	`
)

type Subscriptions struct {
	db DBPool
	TP customtime.TimeProvider
}

func NewSubscriptions(db DBPool, tp customtime.TimeProvider) *Subscriptions {
	return &Subscriptions{
		db: db,
		TP: tp,
	}
}

func (s *Subscriptions) serializeIntoUsers(rows pgx.Rows) (users []*domain.User, err error) {
	for rows.Next() {
		user := new(domain.User)

		err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Avatar,
			&user.DateOfBirth.Time,
			&user.CreatedAt.Time,
			&user.UpdatedAt.Time,
		)
		if err != nil {
			return
		}

		users = append(users, user)
	}

	return
}

func (s *Subscriptions) GetSubscriptions(ctx context.Context, userID uint) (subscriptions []*domain.User, err error) {
	contextlogger.LogSQL(ctx, getSubscriptionsQuery, userID)

	rows, err := s.db.Query(context.Background(), getSubscriptionsQuery, userID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return
	}

	defer rows.Close()

	subscriptions, err = s.serializeIntoUsers(rows)
	if err != nil {
		return
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (s *Subscriptions) GetFriends(ctx context.Context, userID uint) (friends []*domain.User, err error) {
	contextlogger.LogSQL(ctx, getFriendsQuery, userID)

	rows, err := s.db.Query(context.Background(), getFriendsQuery, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return
	}

	defer rows.Close()

	friends, err = s.serializeIntoUsers(rows)
	if err != nil {
		return
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (s *Subscriptions) GetSubscribers(ctx context.Context, userID uint) (subscribers []*domain.User, err error) {
	contextlogger.LogSQL(ctx, getSubscribersQuery, userID)

	rows, err := s.db.Query(context.Background(), getSubscribersQuery, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return
	}

	defer rows.Close()

	subscribers, err = s.serializeIntoUsers(rows)
	if err != nil {
		return
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (s *Subscriptions) Store(ctx context.Context, sub *domain.Subscription) (subscription *domain.Subscription, err error) {
	subscription = new(domain.Subscription)

	contextlogger.LogSQL(ctx, storeSubscriptionQuery, sub.SubscriberID, sub.SubscribedToID)

	err = s.db.QueryRow(context.Background(), storeSubscriptionQuery,
		sub.SubscriberID,
		sub.SubscribedToID,
	).Scan(
		&subscription.ID,
		&subscription.SubscriberID,
		&subscription.SubscribedToID,
		&subscription.CreatedAt.Time,
		&subscription.UpdatedAt.Time,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return
		}

		return
	}

	return
}

func (s *Subscriptions) Delete(ctx context.Context, subsciberID uint, subscribedToID uint) (err error) {
	contextlogger.LogSQL(ctx, deleteSubscriptionQuery, subsciberID, subscribedToID)

	result, err := s.db.Exec(context.Background(), deleteSubscriptionQuery, subsciberID, subscribedToID)
	if err != nil {
		return
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return errors.ErrInvalidBody
	}

	if rowsAffected != 1 {
		return errors.ErrRowsAffected
	}

	return
}

func (s *Subscriptions) GetBySubscriberAndSubscribedToID(ctx context.Context, subscriberID uint, subscribedToID uint) (subscription *domain.Subscription, err error) {
	subscription = new(domain.Subscription)

	contextlogger.LogSQL(ctx, getSubscriptionBySubscriberAndSubscribedToIDQuery, subscriberID, subscribedToID)

	err = s.db.QueryRow(context.Background(), getSubscriptionBySubscriberAndSubscribedToIDQuery,
		subscriberID,
		subscribedToID,
	).Scan(
		&subscription.ID,
		&subscription.SubscriberID,
		&subscription.SubscribedToID,
		&subscription.CreatedAt.Time,
		&subscription.UpdatedAt.Time,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
			return
		}

		return
	}

	return
}
