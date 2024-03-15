package auth

import (
	"mime/multipart"
	"net/http"
	"socio/domain"
	"socio/errors"
	repository "socio/internal/repository/map"
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

type Service struct {
	UserStorage    *repository.Users
	SessionStorage *repository.Sessions
}

type LoginResponse struct {
	SessionID string      `json:"sessionId"`
	User      domain.User `json:"user"`
}

type IsAuthorizedResponse struct {
	IsAuthorized bool `json:"isAuthorized"`
}

func NewService(tp customtime.TimeProvider, userStorage *repository.Users, sessionStorage *repository.Sessions) (a *Service) {
	return &Service{
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
	}
}

func (a *Service) newSession(userID uint) (session *http.Cookie) {
	sessionID := a.SessionStorage.CreateSession(userID)

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

	session = a.newSession(user.ID)

	return
}

func (a *Service) Login(loginInput LoginInput) (user *domain.User, session *http.Cookie, err error) {
	user, err = a.UserStorage.GetUserByEmail(loginInput.Email)
	if err != nil {
		err = errors.ErrInvalidLoginData
		return
	}

	if !hash.MatchPasswords(user.Password, loginInput.Password, []byte(user.Salt)) {
		err = errors.ErrInvalidLoginData
		return
	}

	a.UserStorage.RefreshSaltAndRehashPassword(user)

	return user, a.newSession(user.ID), nil
}

func (a *Service) Logout(session *http.Cookie) (err error) {
	if err = a.SessionStorage.DeleteSession(session.Value); err != nil {
		err = errors.ErrUnauthorized
		return
	}

	session.Expires = time.Time{}
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
