package repository

import (
	"database/sql"
	"socio/domain"
	"socio/errors"
	customtime "socio/pkg/time"
)

type Subscriptions struct {
	db *sql.DB
	TP customtime.TimeProvider
}

type SubscriptionWithIsFriend struct {
	ID             uint                  `json:"subscriptionId"`
	SubscriberID   domain.User           `json:"subscriber"`
	SubscribedToID domain.User           `json:"subscribed_to"`
	CreationDate   customtime.CustomTime `json:"creationDate,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}

func NewSubscriptions(db *sql.DB, tp customtime.TimeProvider) *Subscriptions {
	return &Subscriptions{
		db: db,
		TP: tp,
	}
}

func (s *Subscriptions) serializeIntoUsers(rows *sql.Rows) (users []*domain.User, err error) {
	for rows.Next() {
		user := new(domain.User)

		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Avatar, &user.DateOfBirth.Time, &user.RegistrationDate.Time)
		if err != nil {
			err = errors.ErrInternal
			return
		}

		users = append(users, user)
	}

	return
}

func (s *Subscriptions) GetSubscriptions(userID uint) (subscriptions []*domain.User, err error) {
	rows, err := s.db.Query(`
		SELECT DISTINCT id, first_name, last_name, email, avatar, date_of_birth, registration_date FROM 
		(
			SELECT sub1.subscribed_to AS subscription_id
			FROM subscriptions sub1
			LEFT JOIN subscriptions sub2
			ON sub1.subscribed_to = sub2.subscriber
			AND sub1.subscriber = sub2.subscribed_to
			WHERE sub1.subscriber = $1 AND sub2.id IS null
		) AS sub_ids
		JOIN users ON users.id=sub_ids.subscription_id;
	`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		err = errors.ErrInternal
		return
	}

	subscriptions, err = s.serializeIntoUsers(rows)
	if err != nil {
		return
	}

	return
}

func (s *Subscriptions) GetFriends(userID uint) (friends []*domain.User, err error) {
	rows, err := s.db.Query(`
		SELECT DISTINCT id, first_name, last_name, email, avatar, date_of_birth, registration_date FROM 
		(
			SELECT sub1.subscribed_to AS friend_id
			FROM subscriptions sub1
			INNER JOIN subscriptions sub2
			ON sub1.subscribed_to = sub2.subscriber
			AND sub1.subscriber = sub2.subscribed_to
			WHERE sub1.subscriber = $1
		) AS friend_ids
		JOIN users ON users.id=friend_ids.friend_id;
	`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		err = errors.ErrInternal
		return
	}

	friends, err = s.serializeIntoUsers(rows)
	if err != nil {
		return
	}

	return
}

func (s *Subscriptions) GetSubscribers(userID uint) (subscribers []*domain.User, err error) {
	rows, err := s.db.Query(`
		SELECT DISTINCT id, first_name, last_name, email, avatar, date_of_birth, registration_date FROM 
		(
			SELECT sub2.subscriber AS subscriber_id
			FROM subscriptions sub1
			RIGHT JOIN subscriptions sub2
			ON sub1.subscribed_to = sub2.subscriber
			AND sub1.subscriber = sub2.subscribed_to
			WHERE sub2.subscribed_to = $1
			AND sub1.id IS NULL
		) AS subscriber_ids
		JOIN users ON users.id=subscriber_ids.subscriber_id;
	`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		err = errors.ErrInternal
		return
	}

	subscribers, err = s.serializeIntoUsers(rows)
	if err != nil {
		return
	}

	return
}

func (s *Subscriptions) Store(sub *domain.Subscription) (subscription *domain.Subscription, err error) {
	row := s.db.QueryRow("INSERT INTO subscriptions (subscriber, subscribed_to, creation_date) VALUES ($1, $2, $3) RETURNING id, subscriber, subscribed_to, creation_date", sub.SubscriberID, sub.SubscribedToID, s.TP.Now())

	subscription = new(domain.Subscription)

	err = row.Scan(&subscription.ID, &subscription.SubscriberID, &subscription.SubscribedToID, &subscription.CreationDate.Time)
	if err != nil {
		err = errors.ErrInternal
		return
	}

	return
}

func (s *Subscriptions) Delete(subscriptionID uint) (err error) {
	_, err = s.db.Exec("DELETE FROM subscriptions WHERE id = $1", subscriptionID)

	if err != nil {
		err = errors.ErrInternal
		return
	}

	return
}

func (s *Subscriptions) GetByData(subscriberID uint, subscribedToID uint) (subscription *domain.Subscription, err error) {
	subscription = new(domain.Subscription)

	err = s.db.QueryRow("SELECT id, subscriber, subscribed_to, creation_date FROM subscriptions WHERE subscriber = $1 AND subscribed_to = $2", subscriberID, subscribedToID).Scan(&subscription.ID, &subscription.SubscriberID, &subscription.SubscribedToID, &subscription.CreationDate.Time)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.ErrNotFound
			return
		}

		err = errors.ErrInternal
		return
	}

	return
}
