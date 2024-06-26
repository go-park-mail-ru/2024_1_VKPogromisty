package user

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"socio/domain"
	"socio/errors"
	uspb "socio/internal/grpc/user/proto"
	"socio/pkg/utils"
	"socio/usecase/subscriptions"
	"socio/usecase/user"

	"github.com/google/uuid"
)

const (
	staticFilePath = "."
)

type UserManager struct {
	uspb.UnimplementedUserServer

	UserService          *user.Service
	SubscriptionsService *subscriptions.Service
}

func NewUserManager(userStorage user.UserStorage, subscriptionsStorage subscriptions.SubscriptionsStorage, avatarStorage user.AvatarStorage) *UserManager {
	return &UserManager{
		UserService:          user.NewUserService(userStorage, avatarStorage),
		SubscriptionsService: subscriptions.NewService(subscriptionsStorage, userStorage),
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

func (u *UserManager) GetByEmail(ctx context.Context, in *uspb.GetByEmailRequest) (res *uspb.GetByEmailResponse, err error) {
	email := in.GetEmail()

	user, err := u.UserService.GetUserByEmail(ctx, email)
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &uspb.GetByEmailResponse{
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
	contentType := ""

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
		contentType = req.GetContentType()
	}

	err = u.UserService.UploadAvatar(fileName, file.Name(), contentType)
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

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

func (u *UserManager) Subscribe(ctx context.Context, in *uspb.SubscribeRequest) (res *uspb.SubscribeResponse, err error) {
	subscriberID := in.GetSubscriberId()
	subscribedToID := in.GetSubscribedToId()

	sub, err := u.SubscriptionsService.Subscribe(ctx, &domain.Subscription{
		SubscriberID:   uint(subscriberID),
		SubscribedToID: uint(subscribedToID),
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &uspb.SubscribeResponse{
		Subscription: uspb.ToSubscriptionResponse(sub),
	}

	return
}

func (u *UserManager) Unsubscribe(ctx context.Context, in *uspb.UnsubscribeRequest) (res *uspb.UnsubscribeResponse, err error) {
	subscriberID := in.GetSubscriberId()
	subscribedToID := in.GetSubscribedToId()

	err = u.SubscriptionsService.Unsubscribe(ctx, &domain.Subscription{
		SubscriberID:   uint(subscriberID),
		SubscribedToID: uint(subscribedToID),
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &uspb.UnsubscribeResponse{}

	return
}

func (u *UserManager) GetSubscriptions(ctx context.Context, in *uspb.GetSubscriptionsRequest) (res *uspb.GetSubscriptionsResponse, err error) {
	userID := in.GetUserId()

	subs, err := u.SubscriptionsService.GetSubscriptions(ctx, uint(userID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &uspb.GetSubscriptionsResponse{
		Subscriptions: uspb.ToSubscriptionsResponse(subs),
	}

	return
}

func (u *UserManager) GetSubscribers(ctx context.Context, in *uspb.GetSubscribersRequest) (res *uspb.GetSubscribersResponse, err error) {
	userID := in.GetUserId()

	subs, err := u.SubscriptionsService.GetSubscribers(ctx, uint(userID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &uspb.GetSubscribersResponse{
		Subscribers: uspb.ToSubscriptionsResponse(subs),
	}

	return
}

func (u *UserManager) GetFriends(ctx context.Context, in *uspb.GetFriendsRequest) (res *uspb.GetFriendsResponse, err error) {
	userID := in.GetUserId()

	friends, err := u.SubscriptionsService.GetFriends(ctx, uint(userID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &uspb.GetFriendsResponse{
		Friends: uspb.ToSubscriptionsResponse(friends),
	}

	return
}

func (u *UserManager) SearchByName(ctx context.Context, in *uspb.SearchByNameRequest) (res *uspb.SearchByNameResponse, err error) {
	query := in.GetQuery()

	users, err := u.UserService.SearchByName(ctx, query)
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &uspb.SearchByNameResponse{
		Users: uspb.ToUsersResponse(users),
	}

	return
}

func (u *UserManager) GetSubscriptionIDs(ctx context.Context, in *uspb.GetSubscriptionIDsRequest) (res *uspb.GetSubscriptionIDsResponse, err error) {
	userID := in.GetUserId()

	subIDs, err := u.UserService.GetSubscriptionIDs(ctx, uint(userID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &uspb.GetSubscriptionIDsResponse{
		SubscriptionIds: utils.UintToUint64Slice(subIDs),
	}

	return
}

func (u *UserManager) CreatePublicGroupAdmin(ctx context.Context, in *uspb.CreatePublicGroupAdminRequest) (res *uspb.CreatePublicGroupAdminResponse, err error) {
	userID := in.GetUserId()
	publicGroupID := in.GetPublicGroupId()

	_, err = u.UserService.CreatePublicGroupAdmin(ctx, &domain.PublicGroupAdmin{
		UserID:        uint(userID),
		PublicGroupID: uint(publicGroupID),
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &uspb.CreatePublicGroupAdminResponse{}

	return
}

func (u *UserManager) DeletePublicGroupAdmin(ctx context.Context, in *uspb.DeletePublicGroupAdminRequest) (res *uspb.DeletePublicGroupAdminResponse, err error) {
	userID := in.GetUserId()
	publicGroupID := in.GetPublicGroupId()

	err = u.UserService.DeletePublicGroupAdmin(ctx, &domain.PublicGroupAdmin{
		UserID:        uint(userID),
		PublicGroupID: uint(publicGroupID),
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &uspb.DeletePublicGroupAdminResponse{}

	return
}

func (u *UserManager) GetAdminsByPublicGroupID(ctx context.Context, in *uspb.GetAdminsByPublicGroupIDRequest) (res *uspb.GetAdminsByPublicGroupIDResponse, err error) {
	publicGroupID := in.GetPublicGroupId()

	admins, err := u.UserService.GetAdminsByPublicGroupID(ctx, uint(publicGroupID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &uspb.GetAdminsByPublicGroupIDResponse{
		Admins: uspb.ToUsersResponse(admins),
	}

	return
}

func (u *UserManager) CheckIfUserIsAdmin(ctx context.Context, in *uspb.CheckIfUserIsAdminRequest) (res *uspb.CheckIfUserIsAdminResponse, err error) {
	userID := in.GetUserId()
	publicGroupID := in.GetPublicGroupId()

	isAdmin, err := u.UserService.CheckIfUserIsAdmin(ctx, uint(publicGroupID), uint(userID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &uspb.CheckIfUserIsAdminResponse{
		IsAdmin: isAdmin,
	}

	return
}
