package publicgroup

import (
	"socio/domain"
	publicgroup "socio/usecase/public_group"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToPublicGroupResponse(group *domain.PublicGroup) *PublicGroupResponse {
	return &PublicGroupResponse{
		Id:               uint64(group.ID),
		Name:             group.Name,
		Description:      group.Description,
		Avatar:           group.Avatar,
		SubscribersCount: uint64(group.SubscribersCount),
		CreatedAt:        timestamppb.New(group.CreatedAt.Time),
		UpdatedAt:        timestamppb.New(group.UpdatedAt.Time),
	}
}

func ToPublicGroupsResponse(groups []*domain.PublicGroup) (res []*PublicGroupResponse) {
	for _, group := range groups {
		res = append(res, ToPublicGroupResponse(group))
	}

	return
}

func ToPublicGroupWithInfoResponse(group *publicgroup.PublicGroupWithInfo) *PublicGroupWithInfoResponse {
	return &PublicGroupWithInfoResponse{
		PublicGroup:  ToPublicGroupResponse(group.PublicGroup),
		IsSubscribed: group.IsSubscribed,
	}
}

func ToPublicGroupsWithInfoResponse(groups []*publicgroup.PublicGroupWithInfo) (res []*PublicGroupWithInfoResponse) {
	for _, group := range groups {
		res = append(res, ToPublicGroupWithInfoResponse(group))
	}

	return
}

func ToSubscriptionResponse(subscription *domain.PublicGroupSubscription) *SubscriptionResponse {
	return &SubscriptionResponse{
		Id:            uint64(subscription.ID),
		SubscriberId:  uint64(subscription.SubscriberID),
		PublicGroupId: uint64(subscription.PublicGroupID),
		CreatedAt:     timestamppb.New(subscription.CreatedAt.Time),
		UpdatedAt:     timestamppb.New(subscription.UpdatedAt.Time),
	}
}
