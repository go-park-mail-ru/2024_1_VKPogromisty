package user

import (
	"socio/domain"
	customtime "socio/pkg/time"
)

func ToSubscriptionResponse(sub *domain.Subscription) *SubscriptionResponse {
	return &SubscriptionResponse{
		Id:             uint64(sub.ID),
		SubscriberId:   uint64(sub.SubscriberID),
		SubscribedToId: uint64(sub.SubscribedToID),
	}
}

func ToSubscriptionsResponse(subs []*domain.User) (res []*UserResponse) {
	res = make([]*UserResponse, 0)

	for _, sub := range subs {
		res = append(res, ToUserResponse(sub))
	}

	return
}

func ToSubscription(sub *SubscriptionResponse) *domain.Subscription {
	return &domain.Subscription{
		ID:             uint(sub.Id),
		SubscriberID:   uint(sub.SubscriberId),
		SubscribedToID: uint(sub.SubscribedToId),
		CreatedAt: customtime.CustomTime{
			Time: sub.CreatedAt.AsTime(),
		},
		UpdatedAt: customtime.CustomTime{
			Time: sub.UpdatedAt.AsTime(),
		},
	}
}

func ToSubscriptions(subs []*UserResponse) (res []*domain.User) {
	res = make([]*domain.User, 0)

	for _, sub := range subs {
		res = append(res, ToUser(sub))
	}

	return
}
