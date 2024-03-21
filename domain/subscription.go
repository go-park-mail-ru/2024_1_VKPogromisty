package domain

import customtime "socio/pkg/time"

type Subscription struct {
	ID             uint                  `json:"subscriptionId"`
	SubscriberID   uint                  `json:"subscriber"`
	SubscribedToID uint                  `json:"subscribed_to"`
	CreationDate   customtime.CustomTime `json:"creationDate,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}
