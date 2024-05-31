package domain

import customtime "socio/pkg/time"

//easyjson:json
type PersonalMessage struct {
	ID          uint                  `json:"id"`
	SenderID    uint                  `json:"senderId"`
	ReceiverID  uint                  `json:"receiverId"`
	Content     string                `json:"content"`
	Sticker     *Sticker              `json:"sticker,omitempty"`
	CreatedAt   customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt   customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	Attachments []string              `json:"attachments"`
}

//easyjson:json
type MessageAttachment struct {
	ID           uint                  `json:"id"`
	MessageID    uint                  `json:"messageId"`
	FileName     string                `json:"fileName"`
	OriginalName string                `json:"originalName"`
	CreatedAt    customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}

//easyjson:json
type UnsentMessageAttachment struct {
	SenderID   uint   `json:"senderId"`
	ReceiverID uint   `json:"receiverId"`
	FileName   string `json:"fileName"`
}
