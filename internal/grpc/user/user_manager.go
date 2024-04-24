package user

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"socio/errors"
	uspb "socio/internal/grpc/user/proto"
	"socio/usecase/user"

	"github.com/google/uuid"
)

const (
	staticFilePath = "."
)

type UserManager struct {
	uspb.UnimplementedUserServer

	UserService *user.Service
}

func NewUserManager(userStorage user.UserStorage, avatarStorage user.AvatarStorage) *UserManager {
	return &UserManager{
		UserService: user.NewUserService(userStorage, avatarStorage),
	}
}

func (u *UserManager) GetByID(ctx context.Context, in *uspb.GetByIDRequest) (res *uspb.GetByIDResponse, err error) {
	userID := in.GetUserId()

	user, err := u.UserService.GetUserByID(ctx, uint(userID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
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
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
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
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &uspb.CreateResponse{
		User: uspb.ToUserResponse(user),
	}

	return
}

func (u *UserManager) Upload(stream uspb.User_UploadServer) (err error) {
	file, err := os.Create(filepath.Join(staticFilePath, uuid.NewString()))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	fileName := ""

	var fileSize uint64
	fileSize = 0
	defer func() {
		if err = file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	for {
		req, err := stream.Recv()
		if fileName == "" {
			fileName = req.GetFileName()
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			customErr := errors.NewCustomError(err)
			err = customErr.GRPCStatus().Err()
			return err
		}
		chunk := req.GetChunk()
		fileSize += uint64(len(chunk))
		if _, err = file.Write(chunk); err != nil {
			customErr := errors.NewCustomError(err)
			err = customErr.GRPCStatus().Err()
			return err
		}
	}

	u.UserService.UploadAvatar(fileName, file.Name())

	if err = os.Remove(file.Name()); err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	return stream.SendAndClose(&uspb.UploadResponse{
		FileName: fileName,
		Size:     fileSize,
	})
}

func (u *UserManager) Update(ctx context.Context, in *uspb.UpdateRequest) (res *uspb.UpdateResponse, err error) {
	userInput := uspb.ToUpdateUserInput(in)

	user, err := u.UserService.UpdateUser(ctx, *userInput)
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
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
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &uspb.DeleteResponse{}

	return
}
