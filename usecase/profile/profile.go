package profile

import (
	"mime/multipart"
	"socio/domain"
	customtime "socio/pkg/time"
	"socio/utils"
	"time"
)

type UserStorage interface {
	GetUserByID(userID uint) (user *domain.User, err error)
	GetUserByEmail(email string) (user *domain.User, err error)
	GetUserByIDWithSubsInfo(userID, authorizedUserID uint) (user *domain.User, isSubscribedTo bool, isSubscriber bool, err error)
	UpdateUser(user *domain.User) (updatedUser *domain.User, err error)
	DeleteUser(userID uint) (err error)
}

type SessionStorage interface {
	DeleteSession(sessionID string) (err error)
}

type UserWithSubsInfo struct {
	User           *domain.User
	IsSubscriber   bool `json:"isSubscriber"`
	IsSubscribedTo bool `json:"isSubscribedTo"`
}

type UpdateUserInput struct {
	ID             uint
	FirstName      string
	LastName       string
	Password       string
	RepeatPassword string
	Email          string
	Avatar         *multipart.FileHeader
	DateOfBirth    string
}

type Service struct {
	UserStorage    UserStorage
	SessionStorage SessionStorage
}

func NewProfileService(userStorage UserStorage, sessionStorage SessionStorage) (p *Service) {
	return &Service{
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
	}
}

func (p *Service) GetUserByIDWithSubsInfo(userID, authorizedUserID uint) (userWithInfo UserWithSubsInfo, err error) {
	userWithInfo.User, userWithInfo.IsSubscribedTo, userWithInfo.IsSubscriber, err = p.UserStorage.GetUserByIDWithSubsInfo(userID, authorizedUserID)
	if err != nil {
		return
	}

	return
}

func (p *Service) UpdateUser(input UpdateUserInput) (updatedUser *domain.User, err error) {
	oldUser, err := p.UserStorage.GetUserByID(input.ID)
	if err != nil {
		return
	}

	updatedUser = oldUser

	if err = p.ValidateUserInput(input, oldUser); err != nil {
		return
	}

	if len(input.FirstName) > 0 {
		updatedUser.FirstName = input.FirstName
	}

	if len(input.LastName) > 0 {
		updatedUser.LastName = input.LastName
	}

	if len(input.Email) > 0 {
		updatedUser.Email = input.Email
	}

	if len(input.Password) > 0 {
		updatedUser.Password = input.Password
	}

	if len(input.DateOfBirth) > 0 {
		date, _ := time.Parse(customtime.DateFormat, input.DateOfBirth)
		updatedUser.DateOfBirth = customtime.CustomTime{Time: date}
	}

	if input.Avatar != nil {
		err = utils.RemoveImage(updatedUser.Avatar)
		if err != nil {
			return nil, err
		}

		fileName, err := utils.SaveImage(input.Avatar)
		if err != nil {
			return nil, err
		}

		updatedUser.Avatar = fileName
	}

	updatedUser, err = p.UserStorage.UpdateUser(updatedUser)
	if err != nil {
		return
	}

	return
}

func (p *Service) DeleteUser(userID uint, sessionID string) (err error) {
	err = p.UserStorage.DeleteUser(userID)
	if err != nil {
		return
	}

	err = p.SessionStorage.DeleteSession(sessionID)
	if err != nil {
		return
	}

	return
}
