package auth

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/pkg/hash"
	"socio/pkg/sanitizer"

	"github.com/microcosm-cc/bluemonday"
)

//easyjson:json
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SessionStorage interface {
	CreateSession(ctx context.Context, userID uint) (sessionID string, err error)
	DeleteSession(ctx context.Context, sessionID string) (err error)
	GetUserIDBySession(ctx context.Context, sessionID string) (userID uint, err error)
}

type Service struct {
	SessionStorage SessionStorage
	Sanitizer      *sanitizer.Sanitizer
}

//easyjson:json
type LoginResponse struct {
	User domain.User `json:"user"`
}

//easyjson:json
type IsAuthorizedResponse struct {
	IsAuthorized bool `json:"isAuthorized"`
}

func NewService(sessionStorage SessionStorage) (a *Service) {
	return &Service{
		SessionStorage: sessionStorage,
		Sanitizer:      sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
	}
}

func (a *Service) Login(ctx context.Context, loginInput LoginInput, user *domain.User) (sessionID string, err error) {
	if !hash.MatchPasswords(user.Password, loginInput.Password, []byte(user.Salt)) {
		err = errors.ErrInvalidLoginData
		return
	}

	sessionID, err = a.SessionStorage.CreateSession(ctx, user.ID)
	if err != nil {
		return
	}

	return
}

func (a *Service) Logout(ctx context.Context, sessionID string) (err error) {
	if err = a.SessionStorage.DeleteSession(ctx, sessionID); err != nil {
		err = errors.ErrUnauthorized
		return
	}

	return
}

func (a *Service) IsAuthorized(ctx context.Context, sessionID string) (userID uint, err error) {
	userID, err = a.SessionStorage.GetUserIDBySession(ctx, sessionID)
	if err != nil {
		err = errors.ErrUnauthorized
		return
	}

	return
}
