package publicgroup

import (
	"socio/domain"
	customtime "socio/pkg/time"
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

func ToPublicGroup(group *PublicGroupResponse) *domain.PublicGroup {
	return &domain.PublicGroup{
		ID:               uint(group.Id),
		Name:             group.Name,
		Description:      group.Description,
		Avatar:           group.Avatar,
		SubscribersCount: uint(group.SubscribersCount),
		CreatedAt: customtime.CustomTime{
			Time: group.CreatedAt.AsTime(),
		},
		UpdatedAt: customtime.CustomTime{
			Time: group.UpdatedAt.AsTime(),
		},
	}
}

func ToPublicGroupsResponse(groups []*domain.PublicGroup) (res []*PublicGroupResponse) {
	for _, group := range groups {
		res = append(res, ToPublicGroupResponse(group))
	}

	return
}

func ToPublicGroups(groups []*PublicGroupResponse) (res []*domain.PublicGroup) {
	for _, group := range groups {
		res = append(res, ToPublicGroup(group))
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

func ToPublicGroupWithInfo(group *PublicGroupWithInfoResponse) *publicgroup.PublicGroupWithInfo {
	return &publicgroup.PublicGroupWithInfo{
		PublicGroup:  ToPublicGroup(group.PublicGroup),
		IsSubscribed: group.IsSubscribed,
	}
}

func ToPublicGroupsWithInfo(groups []*PublicGroupWithInfoResponse) (res []*publicgroup.PublicGroupWithInfo) {
	for _, group := range groups {
		res = append(res, ToPublicGroupWithInfo(group))
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

func ToSubscription(subscription *SubscriptionResponse) *domain.PublicGroupSubscription {
	return &domain.PublicGroupSubscription{
		ID:            uint(subscription.Id),
		SubscriberID:  uint(subscription.SubscriberId),
		PublicGroupID: uint(subscription.PublicGroupId),
		CreatedAt: customtime.CustomTime{
			Time: subscription.CreatedAt.AsTime(),
		},
		UpdatedAt: customtime.CustomTime{
			Time: subscription.UpdatedAt.AsTime(),
		},
	}
}
