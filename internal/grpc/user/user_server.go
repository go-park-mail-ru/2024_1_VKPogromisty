package user

import (
	"context"
	uspb "socio/internal/grpc/user/proto"
	"socio/usecase/user"
)

type UserManager struct {
	uspb.UnimplementedUserServer

	UserService *user.Service
}

func NewUserManager(userStorage user.UserStorage) *UserManager {
	return &UserManager{
		UserService: user.NewUserService(userStorage),
	}
}

func (u *UserManager) GetByID(ctx context.Context, in *uspb.GetByIDRequest) (res *uspb.GetByIDResponse, err error) {
	userID := in.GetUserId()

	user, err := u.UserService.GetUserByID(ctx, uint(userID))
	if err != nil {
		return
	}

	res = &uspb.GetByIDResponse{
		User: uspb.ToUserResponse(user),
	}

	return

}

func (u *UserManager) GetByIDWithSubsInfo(ctx context.Context, in *uspb.GetByIDWithSubsInfoRequest) (res *uspb.GetByIDWithSubsInfoResponse, err error) {
	userID := in.GetUserId()
	authorizedUserID := in.GetAuthorizedUserId()

	userWithInfo, err := u.UserService.GetUserByIDWithSubsInfo(ctx, uint(userID), uint(authorizedUserID))
	if err != nil {
		return
	}

	res = &uspb.GetByIDWithSubsInfoResponse{
		User:         uspb.ToUserResponse(userWithInfo.User),
		IsSubscriber: userWithInfo.IsSubscriber,
		IsSubscribed: userWithInfo.IsSubscribedTo,
	}

	return
}

func (u *UserManager) Create(ctx context.Context, in *uspb.CreateRequest) (res *uspb.CreateResponse, err error) {
	userInput := uspb.ToCreateUserInput(in)

	user, err := u.UserService.CreateUser(ctx, *userInput)
	if err != nil {
		return
	}

	res = &uspb.CreateResponse{
		User: uspb.ToUserResponse(user),
	}

	return
}

func (u *UserManager) Update(ctx context.Context, in *uspb.UpdateRequest) (res *uspb.UpdateResponse, err error) {
	userInput := uspb.ToUpdateUserInput(in)

	user, err := u.UserService.UpdateUser(ctx, *userInput)
	if err != nil {
		return
	}

	res = &uspb.UpdateResponse{
		User: uspb.ToUserResponse(user),
	}

	return
}

func (u *UserManager) Delete(ctx context.Context, in *uspb.DeleteRequest) (res *uspb.DeleteResponse, err error) {
	userID := in.GetUserId()

	err = u.UserService.DeleteUser(ctx, uint(userID))
	if err != nil {
		return
	}

	res = &uspb.DeleteResponse{}

	return
}
