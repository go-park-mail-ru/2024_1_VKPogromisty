package auth

import (
	"context"
	"mime/multipart"
	"net/http"
	"socio/domain"
	"socio/errors"
	"socio/pkg/hash"
	"socio/pkg/sanitizer"
	"socio/pkg/static"
	customtime "socio/pkg/time"
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
	StoreUser(ctx context.Context, user *domain.User) (err error)
	GetUserByEmail(ctx context.Context, email string) (user *domain.User, err error)
	RefreshSaltAndRehashPassword(ctx context.Context, user *domain.User, password string) (err error)
}

type SessionStorage interface {
	CreateSession(ctx context.Context, userID uint) (sessionID string, err error)
	DeleteSession(ctx context.Context, sessionID string) (err error)
	GetUserIDBySession(ctx context.Context, sessionID string) (userID uint, err error)
}

type Service struct {
	UserStorage    UserStorage
	SessionStorage SessionStorage
	Sanitizer      *sanitizer.Sanitizer
}

type LoginResponse struct {
	User domain.User `json:"user"`
}

type IsAuthorizedResponse struct {
	IsAuthorized bool `json:"isAuthorized"`
}

func NewService(userStorage UserStorage, sessionStorage SessionStorage, sanitizer *sanitizer.Sanitizer) (a *Service) {
	return &Service{
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
		Sanitizer:      sanitizer,
	}
}

func (a *Service) newSession(ctx context.Context, userID uint) (session *http.Cookie, err error) {
	sessionID, err := a.SessionStorage.CreateSession(ctx, userID)

	if err != nil {
		return
	}

	session = &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		MaxAge:   10 * 60 * 60,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	}
	return
}

func (a *Service) RegistrateUser(ctx context.Context, userInput RegistrationInput) (user *domain.User, session *http.Cookie, err error) {
	if err = a.ValidateUserInput(ctx, userInput); err != nil {
		return
	}

	dateOfBirth, _ := time.Parse(customtime.DateFormat, userInput.DateOfBirth)

	fileName, err := static.SaveImage(userInput.Avatar)
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

	a.Sanitizer.SanitizeUser(user)

	err = a.UserStorage.StoreUser(ctx, user)
	if err != nil {
		return
	}

	session, err = a.newSession(ctx, user.ID)
	if err != nil {
		return
	}

	return
}

func (a *Service) Login(ctx context.Context, loginInput LoginInput) (user *domain.User, session *http.Cookie, err error) {
	user, err = a.UserStorage.GetUserByEmail(ctx, loginInput.Email)
	if err != nil {
		err = errors.ErrInvalidLoginData
		return
	}

	a.Sanitizer.SanitizeUser(user)

	if !hash.MatchPasswords(user.Password, loginInput.Password, []byte(user.Salt)) {
		err = errors.ErrInvalidLoginData
		return
	}

	a.UserStorage.RefreshSaltAndRehashPassword(ctx, user, loginInput.Password)

	session, err = a.newSession(ctx, user.ID)
	if err != nil {
		return
	}

	return
}

func (a *Service) Logout(ctx context.Context, session *http.Cookie) (err error) {
	if err = a.SessionStorage.DeleteSession(ctx, session.Value); err != nil {
		err = errors.ErrUnauthorized
		return
	}

	session.MaxAge = 0
	session.Value = ""
	session.Path = "/"
	session.HttpOnly = true
	session.Secure = true
	session.SameSite = http.SameSiteNoneMode

	return
}

func (a *Service) IsAuthorized(ctx context.Context, session *http.Cookie) (userID uint, err error) {
	userID, err = a.SessionStorage.GetUserIDBySession(ctx, session.Value)
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
