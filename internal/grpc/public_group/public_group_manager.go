package publicgroup

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"socio/domain"
	"socio/errors"
	pgpb "socio/internal/grpc/public_group/proto"
	publicgroup "socio/usecase/public_group"

	"github.com/google/uuid"
)

const (
	staticFilePath = "."
)

type PublicGroupManager struct {
	pgpb.UnimplementedPublicGroupServer

	PublicGroupService *publicgroup.Service
}

func NewPublicGroupManager(publicGroupStorage publicgroup.PublicGroupStorage, avatarStorage publicgroup.AvatarStorage) *PublicGroupManager {
	return &PublicGroupManager{
		PublicGroupService: publicgroup.NewService(publicGroupStorage, avatarStorage),
	}
}

func (p *PublicGroupManager) GetByID(ctx context.Context, in *pgpb.GetByIDRequest) (res *pgpb.GetByIDResponse, err error) {
	groupID := in.GetId()
	userID := in.GetUserId()

	group, err := p.PublicGroupService.GetByID(ctx, uint(groupID), uint(userID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &pgpb.GetByIDResponse{
		PublicGroup: pgpb.ToPublicGroupWithInfoResponse(group),
	}

	return
}

func (p *PublicGroupManager) SearchByName(ctx context.Context, in *pgpb.SearchByNameRequest) (res *pgpb.SearchByNameResponse, err error) {
	query := in.GetQuery()
	userID := in.GetUserId()

	groups, err := p.PublicGroupService.SearchByName(ctx, query, uint(userID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &pgpb.SearchByNameResponse{
		PublicGroups: pgpb.ToPublicGroupsWithInfoResponse(groups),
	}

	return
}

func (p *PublicGroupManager) Create(ctx context.Context, in *pgpb.CreateRequest) (res *pgpb.CreateResponse, err error) {
	group := &domain.PublicGroup{
		Name:        in.GetName(),
		Description: in.GetDescription(),
		Avatar:      in.GetAvatar(),
	}

	newGroup, err := p.PublicGroupService.Create(ctx, group)
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &pgpb.CreateResponse{
		PublicGroup: pgpb.ToPublicGroupResponse(newGroup),
	}

	return
}

func (p *PublicGroupManager) Update(ctx context.Context, in *pgpb.UpdateRequest) (res *pgpb.UpdateResponse, err error) {
	group := &domain.PublicGroup{
		ID:          uint(in.GetId()),
		Name:        in.GetName(),
		Description: in.GetDescription(),
		Avatar:      in.GetAvatar(),
	}

	updatedGroup, err := p.PublicGroupService.Update(ctx, group)
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &pgpb.UpdateResponse{
		PublicGroup: pgpb.ToPublicGroupResponse(updatedGroup),
	}

	return
}

func (p *PublicGroupManager) Delete(ctx context.Context, in *pgpb.DeleteRequest) (res *pgpb.DeleteResponse, err error) {
	groupID := in.GetId()

	err = p.PublicGroupService.Delete(ctx, uint(groupID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &pgpb.DeleteResponse{}

	return
}

func (p *PublicGroupManager) GetSubscriptionByPublicGroupIDAndSubscriberID(ctx context.Context, in *pgpb.GetSubscriptionByPublicGroupIDAndSubscriberIDRequest) (res *pgpb.GetSubscriptionByPublicGroupIDAndSubscriberIDResponse, err error) {
	publicGroupID := in.GetPublicGroupId()
	subscriberID := in.GetSubscriberId()

	subscription, err := p.PublicGroupService.GetSubscriptionByPublicGroupIDAndSubscriberID(ctx, uint(publicGroupID), uint(subscriberID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &pgpb.GetSubscriptionByPublicGroupIDAndSubscriberIDResponse{
		Subscription: pgpb.ToSubscriptionResponse(subscription),
	}

	return
}

func (p *PublicGroupManager) GetBySubscriberID(ctx context.Context, in *pgpb.GetBySubscriberIDRequest) (res *pgpb.GetBySubscriberIDResponse, err error) {
	subscriberID := in.GetSubscriberId()

	groups, err := p.PublicGroupService.GetBySubscriberID(ctx, uint(subscriberID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &pgpb.GetBySubscriberIDResponse{
		PublicGroups: pgpb.ToPublicGroupsResponse(groups),
	}

	return
}

func (p *PublicGroupManager) Subscribe(ctx context.Context, in *pgpb.SubscribeRequest) (res *pgpb.SubscribeResponse, err error) {
	subscriberID := in.GetSubscriberId()
	groupID := in.GetPublicGroupId()

	sub, err := p.PublicGroupService.Subscribe(ctx, &domain.PublicGroupSubscription{
		SubscriberID:  uint(subscriberID),
		PublicGroupID: uint(groupID),
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &pgpb.SubscribeResponse{
		Subscription: pgpb.ToSubscriptionResponse(sub),
	}

	return
}

func (p *PublicGroupManager) Unsubscribe(ctx context.Context, in *pgpb.UnsubscribeRequest) (res *pgpb.UnsubscribeResponse, err error) {
	subscriberID := in.GetSubscriberId()
	groupID := in.GetPublicGroupId()

	err = p.PublicGroupService.Unsubscribe(ctx, &domain.PublicGroupSubscription{
		SubscriberID:  uint(subscriberID),
		PublicGroupID: uint(groupID),
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &pgpb.UnsubscribeResponse{}

	return
}

func (p *PublicGroupManager) Upload(stream pgpb.PublicGroup_UploadServer) (err error) {
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

	p.PublicGroupService.UploadAvatar(fileName, file.Name(), contentType)

	if err = os.Remove(file.Name()); err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	return stream.SendAndClose(&pgpb.UploadResponse{
		FileName: fileName,
		Size:     fileSize,
	})
}

func (p *PublicGroupManager) GetSubscriptionIDs(ctx context.Context, in *pgpb.GetSubscriptionIDsRequest) (res *pgpb.GetSubscriptionIDsResponse, err error) {
	groupID := in.GetUserId()

	subIDs, err := p.PublicGroupService.GetSubscriptionIDs(ctx, uint(groupID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &pgpb.GetSubscriptionIDsResponse{
		PublicGroupIds: pgpb.UintToUint64Slice(subIDs),
	}

	return
}
