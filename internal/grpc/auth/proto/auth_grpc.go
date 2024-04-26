package auth

import (
	"socio/domain"

	customtime "socio/pkg/time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUser(user *UserResponse) *domain.User {
	return &domain.User{
		ID:        uint(user.Id),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Avatar:    user.Avatar,
		DateOfBirth: customtime.CustomTime{
			Time: user.DateOfBirth.AsTime(),
		},
		Salt:     user.Salt,
		Password: user.HashedPassword,
	}
}

func ToUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		Id:             uint64(user.ID),
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		Avatar:         user.Avatar,
		DateOfBirth:    timestamppb.New(user.DateOfBirth.Time),
		CreatedAt:      timestamppb.New(user.CreatedAt.Time),
		UpdatedAt:      timestamppb.New(user.UpdatedAt.Time),
		HashedPassword: user.Password,
		Salt:           user.Salt,
	}
}
