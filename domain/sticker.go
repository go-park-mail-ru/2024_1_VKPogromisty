package domain

import customtime "socio/pkg/time"

type Sticker struct {
	ID        uint                  `json:"id"`
	Name      string                `json:"name"`
	AuthorID  uint                  `json:"authorId"`
	FileName  string                `json:"fileName"`
	CreatedAt customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}
