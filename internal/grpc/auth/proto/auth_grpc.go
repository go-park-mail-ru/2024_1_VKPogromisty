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
		Id:        uint64(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Avatar:    user.Avatar,
		DateOfBirth: &timestamppb.Timestamp{
			Seconds: user.DateOfBirth.Unix(),
			Nanos:   0,
		},
		CreatedAt: &timestamppb.Timestamp{
			Seconds: user.CreatedAt.Unix(),
			Nanos:   int32(user.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: user.UpdatedAt.Unix(),
			Nanos:   int32(user.UpdatedAt.Nanosecond()),
		},
		HashedPassword: user.Password,
		Salt:           user.Salt,
	}
}
