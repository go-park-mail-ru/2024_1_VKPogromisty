package auth

import (
	"context"
	"socio/errors"
	authpb "socio/internal/grpc/auth/proto"
	uspb "socio/internal/grpc/user/proto"
	"socio/usecase/auth"
)

type AuthManager struct {
	authpb.UnimplementedAuthServer

	AuthService *auth.Service
	UserClient  uspb.UserClient
}

func NewAuthManager(userClient uspb.UserClient, sessionStorage auth.SessionStorage) *AuthManager {
	return &AuthManager{
		AuthService: auth.NewService(sessionStorage),
		UserClient:  userClient,
	}
}

func (a *AuthManager) Login(ctx context.Context, in *authpb.LoginRequest) (res *authpb.LoginResponse, err error) {
	loginInput := auth.LoginInput{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}

	userRes, err := a.UserClient.GetByEmail(ctx, &uspb.GetByEmailRequest{Email: loginInput.Email})
	if err != nil {
		return
	}

	user := uspb.ToUser(userRes.User)

	sessionID, err := a.AuthService.Login(ctx, loginInput, user)
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &authpb.LoginResponse{
		User:      authpb.ToUserResponse(user),
		SessionId: sessionID,
	}

	return
}

func (a *AuthManager) Logout(ctx context.Context, in *authpb.LogoutRequest) (res *authpb.LogoutResponse, err error) {
	sessionID := in.GetSessionId()

	err = a.AuthService.Logout(ctx, sessionID)
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &authpb.LogoutResponse{}

	return
}

func (a *AuthManager) ValidateSession(ctx context.Context, in *authpb.ValidateSessionRequest) (res *authpb.ValidateSessionResponse, err error) {
	sessionID := in.GetSessionId()

	userID, err := a.AuthService.IsAuthorized(ctx, sessionID)
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &authpb.ValidateSessionResponse{
		UserId: uint64(userID),
	}

	return
}
