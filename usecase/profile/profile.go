package profile

import "socio/domain"

type UserStorage interface {
	GetUserByIDWithSubsInfo(userID, authorizedUserID uint) (user *domain.User, isSubscribedTo bool, isSubscriber bool, err error)
}

type UserWithSubsInfo struct {
	User           *domain.User
	IsSubscriber   bool `json:"is_subscriber"`
	IsSubscribedTo bool `json:"is_subscribed_to"`
}

type Service struct {
	UserStorage UserStorage
}

func NewProfileService(userStorage UserStorage) (p *Service) {
	return &Service{
		UserStorage: userStorage,
	}
}

func (p *Service) GetUserByIDWithSubsInfo(userID, authorizedUserID uint) (userWithInfo UserWithSubsInfo, err error) {
	userWithInfo.User, userWithInfo.IsSubscribedTo, userWithInfo.IsSubscriber, err = p.UserStorage.GetUserByIDWithSubsInfo(userID, authorizedUserID)
	if err != nil {
		return
	}

	return
}
