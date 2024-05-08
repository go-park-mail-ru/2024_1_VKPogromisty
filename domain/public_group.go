package domain

import customtime "socio/pkg/time"

type PublicGroup struct {
	ID               uint                  `json:"id"`
	Name             string                `json:"name"`
	Description      string                `json:"description"`
	Avatar           string                `json:"avatar"`
	SubscribersCount uint                  `json:"subscribersCount"`
	CreatedAt        customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt        customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}

type PublicGroupSubscription struct {
	ID            uint                  `json:"id"`
	PublicGroupID uint                  `json:"publicGroupId"`
	SubscriberID  uint                  `json:"subscriberId"`
	CreatedAt     customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt     customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}

type PublicGroupAdmin struct {
	ID            uint                  `json:"id"`
	PublicGroupID uint                  `json:"publicGroupId"`
	UserID        uint                  `json:"adminId"`
	CreatedAt     customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt     customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}
