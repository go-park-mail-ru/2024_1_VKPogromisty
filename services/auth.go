package services

import (
	"math/rand"
	"socio/utils"
	"sync"
	"time"
)

type RegistrationInput struct {
	FirstName      string           `json:"firstName,string"`
	LastName       string           `json:"lastName,string"`
	Password       string           `json:"-"`
	RepeatPassword string           `json:"-"`
	Email          string           `json:"email,string"`
	Avatar         string           `json:"avatar,string"`
	DateOfBirth    utils.CustomTime `json:"dateOfBirth,omitempty"`
}

type LoginInput struct {
	Email    string
	Password string
}

type User struct {
	ID               uint             `json:"userId,string"`
	FirstName        string           `json:"firstName,string"`
	LastName         string           `json:"lastName,string"`
	Password         string           `json:"-"`
	Email            string           `json:"email,string"`
	RegistrationDate utils.CustomTime `json:"registrationDate,omitempty"`
	//Avatar stores the URL of the image that serves as user avatar
	Avatar      string           `json:"avatar,string"`
	DateOfBirth utils.CustomTime `json:"dateOfBirth,omitempty"`
}

type AuthService struct {
	users      sync.Map
	sessions   sync.Map
	nextUserId uint
}

var (
	runes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
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
	sessionID = RandStringRunes(32)
	a.sessions.Store(sessionID, userID)
	return
}

func (a *AuthService) RegistrateUser(userInput RegistrationInput) (err error) {
	if err = ValidateUserInput(userInput); err != nil {
		return
	}

	newUser := &User{
		ID:        a.nextUserId,
		FirstName: userInput.FirstName,
		LastName:  userInput.LastName,
		Password:  userInput.Password,
		Email:     userInput.Email,
		RegistrationDate: utils.CustomTime{
			Time: time.Now(),
		},
		Avatar:      userInput.Avatar,
		DateOfBirth: userInput.DateOfBirth,
	}
	a.nextUserId++

	a.users.Store(newUser.Email, newUser)

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
