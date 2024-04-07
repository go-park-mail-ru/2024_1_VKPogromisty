package profile

import (
	"context"
	"mime/multipart"
	"socio/domain"
	"socio/pkg/sanitizer"
	"socio/pkg/static"
	customtime "socio/pkg/time"
	"time"
)

type UserStorage interface {
	GetUserByID(ctx context.Context, userID uint) (user *domain.User, err error)
	GetUserByEmail(ctx context.Context, email string) (user *domain.User, err error)
	GetUserByIDWithSubsInfo(ctx context.Context, userID, authorizedUserID uint) (user *domain.User, isSubscribedTo bool, isSubscriber bool, err error)
	UpdateUser(ctx context.Context, user *domain.User) (updatedUser *domain.User, err error)
	DeleteUser(ctx context.Context, userID uint) (err error)
}

type SessionStorage interface {
	DeleteSession(ctx context.Context, sessionID string) (err error)
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
	Sanitizer      *sanitizer.Sanitizer
}

func NewProfileService(userStorage UserStorage, sessionStorage SessionStorage, sanitizer *sanitizer.Sanitizer) (p *Service) {
	return &Service{
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
		Sanitizer:      sanitizer,
	}
}

func (p *Service) GetUserByIDWithSubsInfo(ctx context.Context, userID, authorizedUserID uint) (userWithInfo UserWithSubsInfo, err error) {
	userWithInfo.User, userWithInfo.IsSubscribedTo, userWithInfo.IsSubscriber, err = p.UserStorage.GetUserByIDWithSubsInfo(ctx, userID, authorizedUserID)
	if err != nil {
		return
	}

	p.Sanitizer.SanitizeUser(userWithInfo.User)

	return
}

func (p *Service) UpdateUser(ctx context.Context, input UpdateUserInput) (updatedUser *domain.User, err error) {
	oldUser, err := p.UserStorage.GetUserByID(ctx, input.ID)
	if err != nil {
		return
	}

	updatedUser = oldUser

	if err = p.ValidateUserInput(ctx, input, oldUser); err != nil {
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
		err = static.RemoveImage(updatedUser.Avatar)
		if err != nil {
			return nil, err
		}

		fileName, err := static.SaveImage(input.Avatar)
		if err != nil {
			return nil, err
		}

		updatedUser.Avatar = fileName
	}

	updatedUser, err = p.UserStorage.UpdateUser(ctx, updatedUser)
	if err != nil {
		return
	}

	p.Sanitizer.SanitizeUser(updatedUser)

	return
}

func (p *Service) DeleteUser(ctx context.Context, userID uint, sessionID string) (err error) {
	err = p.UserStorage.DeleteUser(ctx, userID)
	if err != nil {
		return
	}

	err = p.SessionStorage.DeleteSession(ctx, sessionID)
	if err != nil {
		return
	}

	return
}
