package user

import (
	"socio/domain"
	customtime "socio/pkg/time"
	"socio/usecase/user"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

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
	}
}

func ToUser(user *UserResponse) *domain.User {
	return &domain.User{
		ID:        uint(user.Id),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		DateOfBirth: customtime.CustomTime{
			Time: user.DateOfBirth.AsTime(),
		},
	}
}

func ToCreateUserInput(req *CreateRequest) (userInput *user.CreateUserInput) {
	return &user.CreateUserInput{
		FirstName:      req.GetFirstName(),
		LastName:       req.GetLastName(),
		Email:          req.GetEmail(),
		Password:       req.GetPassword(),
		RepeatPassword: req.GetRepeatPassword(),
		DateOfBirth:    req.DateOfBirth.AsTime().Format(customtime.DateFormat),
	}
}

func ToUpdateRequest(input *user.UpdateUserInput) *UpdateRequest {
	dateOfBirth, _ := time.Parse(customtime.DateFormat, input.DateOfBirth)
	return &UpdateRequest{
		UserId:         uint64(input.ID),
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		Password:       input.Password,
		RepeatPassword: input.RepeatPassword,
		Avatar:         input.Avatar.Filename,
		DateOfBirth:    timestamppb.New(dateOfBirth),
	}
}

func ToUpdateUserInput(req *UpdateRequest) (userInput *user.UpdateUserInput) {
	return &user.UpdateUserInput{
		ID:             uint(req.GetUserId()),
		FirstName:      req.GetFirstName(),
		LastName:       req.GetLastName(),
		Email:          req.GetEmail(),
		Password:       req.GetPassword(),
		RepeatPassword: req.GetRepeatPassword(),
		DateOfBirth:    req.DateOfBirth.AsTime().Format(customtime.DateFormat),
	}
}
