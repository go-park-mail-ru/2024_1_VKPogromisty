package domain

import customtime "socio/pkg/time"

type PersonalMessage struct {
	ID         uint                  `json:"id"`
	SenderID   uint                  `json:"senderId"`
	ReceiverID uint                  `json:"receiverId"`
	Content    string                `json:"content"`
	CreatedAt  customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt  customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}
