package domain

import (
	customtime "socio/pkg/time"
)

type Post struct {
	ID           uint                  `json:"postId"`
	AuthorID     uint                  `json:"authorId"`
	Text         string                `json:"text"`
	Attachments  []string              `json:"attachments"`
	CreationDate customtime.CustomTime `json:"creationDate,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}
