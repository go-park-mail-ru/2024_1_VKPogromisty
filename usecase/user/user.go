package user

import (
	"context"
	"socio/domain"
	"socio/pkg/sanitizer"
	customtime "socio/pkg/time"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

type UserStorage interface {
	GetUserByID(ctx context.Context, userID uint) (user *domain.User, err error)
	GetUserByEmail(ctx context.Context, email string) (user *domain.User, err error)
	GetUserByIDWithSubsInfo(ctx context.Context, userID, authorizedUserID uint) (user *domain.User, isSubscribedTo bool, isSubscriber bool, err error)
	StoreUser(ctx context.Context, user *domain.User) (err error)
	UpdateUser(ctx context.Context, user *domain.User, prevPassword string) (updatedUser *domain.User, err error)
	DeleteUser(ctx context.Context, userID uint) (err error)
	SearchByName(ctx context.Context, query string) (users []*domain.User, err error)
	GetSubscriptionIDs(ctx context.Context, userID uint) (subscribedToIDs []uint, err error)
	StorePublicGroupAdmin(ctx context.Context, publicGroupAdmin *domain.PublicGroupAdmin) (newPublicGroupAdmin *domain.PublicGroupAdmin, err error)
	DeletePublicGroupAdmin(ctx context.Context, publicGroupAdmin *domain.PublicGroupAdmin) (err error)
	GetAdminsByPublicGroupID(ctx context.Context, publicGroupID uint) (admins []*domain.User, err error)
	CheckIfUserIsAdmin(ctx context.Context, publicGroupID, userID uint) (isAdmin bool, err error)
}

type AvatarStorage interface {
	Store(fileName string, filePath string, contentType string) (err error)
	Delete(fileName string) (err error)
}

//easyjson:json
type UserWithSubsInfo struct {
	User           *domain.User
	IsSubscriber   bool `json:"isSubscriber"`
	IsSubscribedTo bool `json:"isSubscribedTo"`
}

type CreateUserInput struct {
	FirstName      string
	LastName       string
	Password       string
	RepeatPassword string
	Email          string
	Avatar         string
	DateOfBirth    string
}

type UpdateUserInput struct {
	ID             uint
	FirstName      string
	LastName       string
	Password       string
	RepeatPassword string
	Email          string
	Avatar         string
	DateOfBirth    string
}

type Service struct {
	UserStorage   UserStorage
	AvatarStorage AvatarStorage
	Sanitizer     *sanitizer.Sanitizer
}

func NewUserService(userStorage UserStorage, avatarStorage AvatarStorage) (p *Service) {
	return &Service{
		UserStorage:   userStorage,
		AvatarStorage: avatarStorage,
		Sanitizer:     sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
	}
}

func (p *Service) GetUserByID(ctx context.Context, userID uint) (user *domain.User, err error) {
	user, err = p.UserStorage.GetUserByID(ctx, userID)
	if err != nil {
		return
	}

	return
}

func (p *Service) GetUserByEmail(ctx context.Context, email string) (user *domain.User, err error) {
	user, err = p.UserStorage.GetUserByEmail(ctx, email)
	if err != nil {
		return
	}

	return
}

func (p *Service) UploadAvatar(fileName string, filePath string, contentType string) (err error) {
	err = p.AvatarStorage.Store(fileName, filePath, contentType)
	if err != nil {
		return
	}

	return
}

func (p *Service) CreateUser(ctx context.Context, userInput CreateUserInput) (user *domain.User, err error) {
	if err = p.ValidateCreateUserInput(ctx, userInput); err != nil {
		return
	}

	dateOfBirth, _ := time.Parse(customtime.DateFormat, userInput.DateOfBirth)

	user = &domain.User{
		FirstName: userInput.FirstName,
		LastName:  userInput.LastName,
		Password:  userInput.Password,
		Email:     userInput.Email,
		Avatar:    userInput.Avatar,
		DateOfBirth: customtime.CustomTime{
			Time: dateOfBirth,
		},
	}

	p.Sanitizer.SanitizeUser(user)

	err = p.UserStorage.StoreUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return
}

func (p *Service) GetUserByIDWithSubsInfo(ctx context.Context, userID uint, authorizedUserID uint) (userWithInfo UserWithSubsInfo, err error) {
	userWithInfo = UserWithSubsInfo{
		User:           &domain.User{},
		IsSubscriber:   false,
		IsSubscribedTo: false,
	}

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

	prevPassword := oldUser.Password

	updatedUser = oldUser

	if err = p.ValidateUpdateUserInput(ctx, input, oldUser); err != nil {
		return nil, err
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

	if len(input.Avatar) > 0 {
		err = p.AvatarStorage.Delete(updatedUser.Avatar)
		if err != nil {
			return nil, err
		}

		updatedUser.Avatar = input.Avatar
	}

	updatedUser, err = p.UserStorage.UpdateUser(ctx, updatedUser, prevPassword)
	if err != nil {
		return nil, err
	}

	p.Sanitizer.SanitizeUser(updatedUser)

	return
}

func (p *Service) DeleteUser(ctx context.Context, userID uint) (err error) {
	user, err := p.UserStorage.GetUserByID(ctx, userID)
	if err != nil {
		return
	}

	err = p.AvatarStorage.Delete(user.Avatar)
	if err != nil {
		return
	}

	err = p.UserStorage.DeleteUser(ctx, userID)
	if err != nil {
		return
	}

	return
}

func (p *Service) SearchByName(ctx context.Context, query string) (users []*domain.User, err error) {
	users, err = p.UserStorage.SearchByName(ctx, query)
	if err != nil {
		return
	}

	for _, user := range users {
		p.Sanitizer.SanitizeUser(user)
	}

	return
}

func (p *Service) GetSubscriptionIDs(ctx context.Context, userID uint) (subIDs []uint, err error) {
	_, err = p.UserStorage.GetUserByID(ctx, userID)
	if err != nil {
		return
	}

	subIDs, err = p.UserStorage.GetSubscriptionIDs(ctx, userID)
	if err != nil {
		return
	}

	return
}
