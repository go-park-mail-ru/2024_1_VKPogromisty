package user

import (
	"socio/domain"
	customtime "socio/pkg/time"
	"socio/usecase/user"

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
		HashedPassword: user.Password,
		Salt:           user.Salt,
	}
}

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

func ToCreateUserInput(req *CreateRequest) (userInput *user.CreateUserInput) {
	return &user.CreateUserInput{
		FirstName:      req.GetFirstName(),
		LastName:       req.GetLastName(),
		Email:          req.GetEmail(),
		Avatar:         req.GetAvatar(),
		Password:       req.GetPassword(),
		RepeatPassword: req.GetRepeatPassword(),
		DateOfBirth:    req.GetDateOfBirth(),
	}
}

func ToCreateRequest(input *user.CreateUserInput) *CreateRequest {
	return &CreateRequest{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		Password:       input.Password,
		RepeatPassword: input.RepeatPassword,
		Avatar:         input.Avatar,
		DateOfBirth:    input.DateOfBirth,
	}
}

func ToUpdateRequest(input *user.UpdateUserInput) *UpdateRequest {
	return &UpdateRequest{
		UserId:         uint64(input.ID),
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		Password:       input.Password,
		RepeatPassword: input.RepeatPassword,
		Avatar:         input.Avatar,
		DateOfBirth:    input.DateOfBirth,
	}
}

func ToUpdateUserInput(req *UpdateRequest) (userInput *user.UpdateUserInput) {
	return &user.UpdateUserInput{
		ID:             uint(req.GetUserId()),
		FirstName:      req.GetFirstName(),
		LastName:       req.GetLastName(),
		Email:          req.GetEmail(),
		Avatar:         req.GetAvatar(),
		Password:       req.GetPassword(),
		RepeatPassword: req.GetRepeatPassword(),
		DateOfBirth:    req.GetDateOfBirth(),
	}
}

func ToUserWithInfo(res *GetByIDWithSubsInfoResponse) (userWithInfo *user.UserWithSubsInfo) {
	return &user.UserWithSubsInfo{
		User:           ToUser(res.GetUser()),
		IsSubscriber:   res.GetIsSubscriber(),
		IsSubscribedTo: res.GetIsSubscribed(),
	}
}
