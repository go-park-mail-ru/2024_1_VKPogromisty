package subscriptions

import (
	"socio/domain"
	"socio/errors"
)

type SubscriptionsStorage interface {
	Store(sub *domain.Subscription) (subscription *domain.Subscription, err error)
	Delete(subscriberID uint, subscibedToID uint) (err error)
	GetBySubscriberAndSubscribedToID(subscriberID uint, subscribedToID uint) (subscription *domain.Subscription, err error)
	GetSubscriptions(userID uint) (subscriptions []*domain.User, err error)
	GetSubscribers(userID uint) (subscribers []*domain.User, err error)
	GetFriends(userID uint) (friends []*domain.User, err error)
}

type UserStorage interface {
	GetUserByID(userID uint) (user *domain.User, err error)
}

type Service struct {
	SubscriptionsStorage SubscriptionsStorage
	UserStorage          UserStorage
}

func NewService(subStorage SubscriptionsStorage, userStorage UserStorage) (service *Service) {
	service = &Service{
		SubscriptionsStorage: subStorage,
		UserStorage:          userStorage,
	}
	return
}

func (s *Service) Subscribe(sub *domain.Subscription) (subscription *domain.Subscription, err error) {
	if sub.SubscriberID == sub.SubscribedToID {
		err = errors.ErrInvalidBody
		return
	}

	_, err = s.UserStorage.GetUserByID(sub.SubscriberID)
	if err != nil {
		err = errors.ErrInvalidBody
		return
	}

	_, err = s.UserStorage.GetUserByID(sub.SubscribedToID)
	if err != nil {
		err = errors.ErrInvalidBody
		return
	}

	subscription, err = s.SubscriptionsStorage.Store(sub)
	if err != nil {
		return
	}

	return
}

func (s *Service) Unsubscribe(sub *domain.Subscription) (err error) {
	_, err = s.UserStorage.GetUserByID(sub.SubscriberID)
	if err != nil {
		err = errors.ErrInvalidBody
		return
	}

	_, err = s.UserStorage.GetUserByID(sub.SubscribedToID)
	if err != nil {
		err = errors.ErrInvalidBody
		return
	}

	err = s.SubscriptionsStorage.Delete(sub.SubscriberID, sub.SubscribedToID)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetSubscriptions(userID uint) (subscriptions []*domain.User, err error) {
	_, err = s.UserStorage.GetUserByID(userID)
	if err != nil {
		err = errors.ErrNotFound
		return
	}

	subscriptions, err = s.SubscriptionsStorage.GetSubscriptions(userID)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetSubscribers(userID uint) (subscribers []*domain.User, err error) {
	_, err = s.UserStorage.GetUserByID(userID)
	if err != nil {
		err = errors.ErrNotFound
		return
	}

	subscribers, err = s.SubscriptionsStorage.GetSubscribers(userID)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetFriends(userID uint) (friends []*domain.User, err error) {
	_, err = s.UserStorage.GetUserByID(userID)
	if err != nil {
		err = errors.ErrNotFound
		return
	}

	friends, err = s.SubscriptionsStorage.GetFriends(userID)
	if err != nil {
		return
	}

	return
}
