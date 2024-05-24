package domain

import customtime "socio/pkg/time"

//easyjson:json
type Subscription struct {
	ID             uint                  `json:"subscriptionId"`
	SubscriberID   uint                  `json:"subscriber"`
	SubscribedToID uint                  `json:"subscribedTo"`
	CreatedAt      customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt      customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}
