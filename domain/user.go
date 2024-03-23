package domain

import customtime "socio/pkg/time"

type User struct {
	ID          uint                  `json:"userId"`
	FirstName   string                `json:"firstName"`
	LastName    string                `json:"lastName"`
	Password    string                `json:"-"`
	Salt        string                `json:"-"`
	Email       string                `json:"email"`
	Avatar      string                `json:"avatar" example:"default_avatar.png"`
	DateOfBirth customtime.CustomTime `json:"dateOfBirth,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	CreatedAt   customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt   customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}
