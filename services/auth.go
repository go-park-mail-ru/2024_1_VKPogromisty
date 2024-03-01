package services

import (
	"mime/multipart"
	"net/http"
	"socio/errors"
	"socio/utils"
	"sync"
	"time"

	"github.com/google/uuid"
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
	Salt             string           `json:"-"`
	Email            string           `json:"email"`
	RegistrationDate utils.CustomTime `json:"registrationDate,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	Avatar           string           `json:"avatar" example:"default_avatar.png"`
	DateOfBirth      utils.CustomTime `json:"dateOfBirth,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}

type AuthService struct {
	users      sync.Map
	sessions   sync.Map
	nextUserId uint
}

type LoginResponse struct {
	SessionID string `json:"sessionId"`
}

func NewAuthService() (authService *AuthService) {
	authService = &AuthService{
		users:      sync.Map{},
		sessions:   sync.Map{},
		nextUserId: 2,
	}

	salt1 := uuid.NewString()
	user1 := &User{
		ID:        0,
		FirstName: "Petr",
		LastName:  "Mitin",
		Password:  utils.HashPassword("admin", []byte(salt1)),
		Salt:      salt1,
		Email:     "petr09mitin@mail.ru",
		RegistrationDate: utils.CustomTime{
			Time: time.Now(),
		},
		Avatar: "default_avatar.png",
		DateOfBirth: utils.CustomTime{
			Time: time.Now(),
		},
	}
	authService.users.Store(user1.Email, user1)

	salt2 := uuid.NewString()
	user2 := &User{
		ID:        1,
		FirstName: "Alexey",
		LastName:  "Gorbunov",
		Password:  utils.HashPassword("admin2", []byte(salt2)),
		Salt:      salt2,
		Email:     "lexagorbunov14@gmail.com",
		RegistrationDate: utils.CustomTime{
			Time: time.Now(),
		},
		Avatar: "leha.jpg",
		DateOfBirth: utils.CustomTime{
			Time: time.Now(),
		},
	}
	authService.users.Store(user2.Email, user2)

	return
}

func (a *AuthService) newSession(userID uint) (session *http.Cookie) {
	sessionID := uuid.NewString()
	a.sessions.Store(sessionID, userID)

	session = &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		MaxAge:   10 * 60 * 60,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	return
}

func (a *AuthService) RegistrateUser(userInput RegistrationInput) (user *User, session *http.Cookie, err error) {
	if err = ValidateUserInput(userInput, a); err != nil {
		return
	}

	dateOfBirth, err := time.Parse(utils.DateFormat, userInput.DateOfBirth)
	if err != nil {
		err = errors.ErrInvalidDate
		return
	}

	fileName, err := utils.SaveImage(userInput.Avatar)
	if err != nil {
		return
	}

	salt := uuid.NewString()
	user = &User{
		ID:        a.nextUserId,
		FirstName: userInput.FirstName,
		LastName:  userInput.LastName,
		Password:  utils.HashPassword(userInput.Password, []byte(salt)),
		Salt:      salt,
		Email:     userInput.Email,
		RegistrationDate: utils.CustomTime{
			Time: time.Now(),
		},
		Avatar: fileName,
		DateOfBirth: utils.CustomTime{
			Time: dateOfBirth,
		},
	}
	a.nextUserId++

	a.users.Store(user.Email, user)

	session = a.newSession(user.ID)

	return
}

func (a *AuthService) Login(loginInput LoginInput) (session *http.Cookie, err error) {
	userData, ok := a.users.Load(loginInput.Email)
	if !ok {
		err = errors.ErrInvalidLoginData
		return
	}

	user, ok := userData.(*User)
	if !ok || !utils.MatchPasswords(user.Password, loginInput.Password, []byte(user.Salt)) {
		err = errors.ErrInvalidLoginData
		return
	}

	user.Salt = uuid.NewString()
	user.Password = utils.HashPassword(loginInput.Password, []byte(user.Salt))

	return a.newSession(user.ID), nil
}

func (a *AuthService) Logout(session *http.Cookie) (err error) {
	_, ok := a.sessions.LoadAndDelete(session.Value)
	if !ok {
		err = errors.ErrInvalidData
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	return
}

func (a *AuthService) IsAuthorized(session *http.Cookie) (err error) {
	_, ok := a.sessions.Load(session.Value)
	if !ok {
		err = errors.ErrUnauthorized
		return
	}

	if errCookie := session.Valid(); errCookie != nil {
		err = errors.ErrUnauthorized
		return
	}
	return
}
