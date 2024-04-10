package subscriptions

import (
	"context"
	"socio/domain"
	"socio/errors"
)

type SubscriptionsStorage interface {
	Store(ctx context.Context, sub *domain.Subscription) (subscription *domain.Subscription, err error)
	Delete(ctx context.Context, subscriberID uint, subscibedToID uint) (err error)
	GetBySubscriberAndSubscribedToID(ctx context.Context, subscriberID uint, subscribedToID uint) (subscription *domain.Subscription, err error)
	GetSubscriptions(ctx context.Context, userID uint) (subscriptions []*domain.User, err error)
	GetSubscribers(ctx context.Context, userID uint) (subscribers []*domain.User, err error)
	GetFriends(ctx context.Context, userID uint) (friends []*domain.User, err error)
}

type UserStorage interface {
	GetUserByID(ctx context.Context, userID uint) (user *domain.User, err error)
}

type Service struct {
	SubscriptionsStorage SubscriptionsStorage
	UserStorage          UserStorage
}

type GetSubscriptionsResponse struct {
	Subscriptions []*domain.User `json:"subscriptions"`
}

type GetSubscribersResponse struct {
	Subscribers []*domain.User `json:"subscribers"`
}

type GetFriendsResponse struct {
	Friends []*domain.User `json:"friends"`
}

func NewService(subStorage SubscriptionsStorage, userStorage UserStorage) (service *Service) {
	service = &Service{
		SubscriptionsStorage: subStorage,
		UserStorage:          userStorage,
	}
	return
}

func (s *Service) Subscribe(ctx context.Context, sub *domain.Subscription) (subscription *domain.Subscription, err error) {
	if sub.SubscriberID == sub.SubscribedToID {
		err = errors.ErrInvalidBody
		return
	}

	_, err = s.UserStorage.GetUserByID(ctx, sub.SubscriberID)
	if err != nil {
		err = errors.ErrInvalidBody
		return
	}

	_, err = s.UserStorage.GetUserByID(ctx, sub.SubscribedToID)
	if err != nil {
		err = errors.ErrInvalidBody
		return
	}

	subscription, err = s.SubscriptionsStorage.Store(ctx, sub)
	if err != nil {
		return
	}

	return
}

func (s *Service) Unsubscribe(ctx context.Context, sub *domain.Subscription) (err error) {
	_, err = s.UserStorage.GetUserByID(ctx, sub.SubscriberID)
	if err != nil {
		err = errors.ErrInvalidBody
		return
	}

	_, err = s.UserStorage.GetUserByID(ctx, sub.SubscribedToID)
	if err != nil {
		err = errors.ErrInvalidBody
		return
	}

	err = s.SubscriptionsStorage.Delete(ctx, sub.SubscriberID, sub.SubscribedToID)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetSubscriptions(ctx context.Context, userID uint) (subscriptions []*domain.User, err error) {
	_, err = s.UserStorage.GetUserByID(ctx, userID)
	if err != nil {
		err = errors.ErrNotFound
		return
	}

	subscriptions, err = s.SubscriptionsStorage.GetSubscriptions(ctx, userID)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetSubscribers(ctx context.Context, userID uint) (subscribers []*domain.User, err error) {
	_, err = s.UserStorage.GetUserByID(ctx, userID)
	if err != nil {
		err = errors.ErrNotFound
		return
	}

	subscribers, err = s.SubscriptionsStorage.GetSubscribers(ctx, userID)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetFriends(ctx context.Context, userID uint) (friends []*domain.User, err error) {
	_, err = s.UserStorage.GetUserByID(ctx, userID)
	if err != nil {
		err = errors.ErrNotFound
		return
	}

	friends, err = s.SubscriptionsStorage.GetFriends(ctx, userID)
	if err != nil {
		return
	}

	return
}
