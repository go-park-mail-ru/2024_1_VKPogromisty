package auth

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"socio/domain"
	"socio/errors"
	"socio/pkg/hash"
	customtime "socio/pkg/time"
	"socio/utils"
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
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserStorage interface {
	StoreUser(user *domain.User) (err error)
	GetUserByEmail(email string) (user *domain.User, err error)
	RefreshSaltAndRehashPassword(user *domain.User, password string) (err error)
}

type SessionStorage interface {
	CreateSession(userID uint) (sessionID string, err error)
	DeleteSession(sessionID string) (err error)
	GetUserIDBySession(sessionID string) (userID uint, err error)
}

type Service struct {
	UserStorage    UserStorage
	SessionStorage SessionStorage
}

type LoginResponse struct {
	SessionID string      `json:"sessionId"`
	User      domain.User `json:"user"`
}

type IsAuthorizedResponse struct {
	IsAuthorized bool `json:"isAuthorized"`
}

func NewService(userStorage UserStorage, sessionStorage SessionStorage) (a *Service) {
	return &Service{
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
	}
}

func (a *Service) newSession(userID uint) (session *http.Cookie, err error) {
	sessionID, err := a.SessionStorage.CreateSession(userID)

	if err != nil {
		return
	}

	session = &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		MaxAge:   10 * 60 * 60,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	}
	return
}

func (a *Service) RegistrateUser(userInput RegistrationInput) (user *domain.User, session *http.Cookie, err error) {
	if err = a.ValidateUserInput(userInput); err != nil {
		return
	}

	dateOfBirth, err := time.Parse(customtime.DateFormat, userInput.DateOfBirth)
	if err != nil {
		err = errors.ErrInvalidDate
		return
	}

	fileName, err := utils.SaveImage(userInput.Avatar)
	if err != nil {
		return
	}

	user = &domain.User{
		FirstName: userInput.FirstName,
		LastName:  userInput.LastName,
		Password:  userInput.Password,
		Email:     userInput.Email,
		Avatar:    fileName,
		DateOfBirth: customtime.CustomTime{
			Time: dateOfBirth,
		},
	}

	a.UserStorage.StoreUser(user)

	session, err = a.newSession(user.ID)
	if err != nil {
		return
	}

	return
}

func (a *Service) Login(loginInput LoginInput) (user *domain.User, session *http.Cookie, err error) {
	user, err = a.UserStorage.GetUserByEmail(loginInput.Email)
	if err != nil {
		fmt.Println(err)
		err = errors.ErrInvalidLoginData
		return
	}

	if !hash.MatchPasswords(user.Password, loginInput.Password, []byte(user.Salt)) {
		err = errors.ErrInvalidLoginData
		return
	}

	a.UserStorage.RefreshSaltAndRehashPassword(user, loginInput.Password)

	session, err = a.newSession(user.ID)
	if err != nil {
		return
	}

	return
}

func (a *Service) Logout(session *http.Cookie) (err error) {
	if err = a.SessionStorage.DeleteSession(session.Value); err != nil {
		err = errors.ErrUnauthorized
		return
	}

	session.MaxAge = 0
	session.Value = ""
	session.Path = "/"

	return
}

func (a *Service) IsAuthorized(session *http.Cookie) (userID uint, err error) {
	userID, err = a.SessionStorage.GetUserIDBySession(session.Value)
	if err != nil {
		err = errors.ErrUnauthorized
		return
	}

	if errCookie := session.Valid(); errCookie != nil {
		err = errors.ErrUnauthorized
		return
	}

	return
}
