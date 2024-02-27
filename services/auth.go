package services

import (
	"mime/multipart"
	"os"
	"socio/utils"
	"sync"
	"time"
)

type RegistrationInput struct {
	FirstName      string
	LastName       string
	Password       string
	RepeatPassword string
	Email          string
	Avatar         *multipart.FileHeader
	DateOfBirth    string
}

type LoginInput struct {
	Email    string
	Password string
}

type User struct {
	ID               uint             `json:"userId"`
	FirstName        string           `json:"firstName"`
	LastName         string           `json:"lastName"`
	Password         string           `json:"-"`
	Email            string           `json:"email"`
	RegistrationDate utils.CustomTime `json:"registrationDate,omitempty"`
	//Avatar stores the URL of the image that serves as user avatar
	Avatar      string           `json:"avatar"`
	DateOfBirth utils.CustomTime `json:"dateOfBirth,omitempty"`
}

type AuthService struct {
	users      sync.Map
	sessions   sync.Map
	nextUserId uint
}

func NewAuthService() (authService *AuthService) {
	authService = &AuthService{
		users:      sync.Map{},
		sessions:   sync.Map{},
		nextUserId: 2,
	}

	user1 := &User{
		ID:        0,
		FirstName: "Petr",
		LastName:  "Mitin",
		Password:  "admin",
		Email:     "petr09mitin@mail.ru",
		RegistrationDate: utils.CustomTime{
			Time: time.Now(),
		},
		Avatar: "",
		DateOfBirth: utils.CustomTime{
			Time: time.Now(),
		},
	}
	authService.users.Store(user1.Email, user1)

	user2 := &User{
		ID:        1,
		FirstName: "Alexey",
		LastName:  "Gorbunov",
		Password:  "admin2",
		Email:     "lexagorbunov14@gmail.com",
		RegistrationDate: utils.CustomTime{
			Time: time.Now(),
		},
		Avatar: "",
		DateOfBirth: utils.CustomTime{
			Time: time.Now(),
		},
	}
	authService.users.Store(user2.Email, user2)

	return
}

func (a *AuthService) newSession(userID uint) (sessionID string) {
	sessionID = utils.RandStringRunes(32)
	a.sessions.Store(sessionID, userID)
	return
}

func (a *AuthService) RegistrateUser(userInput RegistrationInput) (user *User, err error) {
	if err = ValidateUserInput(userInput, a); err != nil {
		return
	}

	dateOfBirth, err := time.Parse(utils.DateFormat, userInput.DateOfBirth)
	if err != nil {
		return
	}

	fileName, err := utils.SaveImage(userInput.Avatar)
	if err != nil {
		return
	}
	avatarURL := os.Getenv("PROTOCOL") + os.Getenv("HOST") + os.Getenv("PORT") + "/static/" + fileName

	user = &User{
		ID:        a.nextUserId,
		FirstName: userInput.FirstName,
		LastName:  userInput.LastName,
		Password:  userInput.Password,
		Email:     userInput.Email,
		RegistrationDate: utils.CustomTime{
			Time: time.Now(),
		},
		Avatar: avatarURL,
		DateOfBirth: utils.CustomTime{
			Time: dateOfBirth,
		},
	}
	a.nextUserId++

	a.users.Store(user.Email, user)

	return
}

func (a *AuthService) Login(loginInput LoginInput) (sessionID string, err error) {
	userData, ok := a.users.Load(loginInput.Email)
	if !ok {
		err = utils.ErrInvalidLoginData
		return
	}

	user, ok := userData.(*User)
	if !ok || user.Password != loginInput.Password {
		err = utils.ErrInvalidLoginData
		return
	}

	return a.newSession(user.ID), nil
}

func (a *AuthService) Logout(sessionID string) (err error) {
	_, ok := a.sessions.LoadAndDelete(sessionID)
	if !ok {
		err = utils.ErrInvalidData
		return
	}

	return
}
