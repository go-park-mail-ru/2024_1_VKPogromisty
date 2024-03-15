package domain

import "socio/utils"

type User struct {
	ID               uint             `json:"userId"`
	FirstName        string           `json:"firstName"`
	LastName         string           `json:"lastName"`
	Password         string           `json:"-"`
	Salt             string           `json:"-"`
	Email            string           `json:"email"`
	RegistrationDate utils.CustomTime `json:"registrationDate,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	Avatar           string           `json:"avatar" example:"default_avatar.png"`
	DateOfBirth      utils.CustomTime `json:"dateOfBirth,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}
