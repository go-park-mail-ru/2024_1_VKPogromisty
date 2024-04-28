package user

import (
	"socio/domain"
	customtime "socio/pkg/time"
	"socio/usecase/user"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserResponse(user *domain.User) (res *UserResponse) {
	res = &UserResponse{
		Id:             uint64(user.ID),
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		HashedPassword: user.Password,
		Salt:           user.Salt,
		Avatar:         user.Avatar,
		DateOfBirth:    timestamppb.New(user.DateOfBirth.Time),
		CreatedAt:      timestamppb.New(user.CreatedAt.Time),
		UpdatedAt:      timestamppb.New(user.UpdatedAt.Time),
	}

	return
}

func ToUsersResponse(users []*domain.User) (res []*UserResponse) {
	for _, user := range users {
		res = append(res, ToUserResponse(user))
	}

	return
}

func ToUser(user *UserResponse) *domain.User {
	return &domain.User{
		ID:        uint(user.Id),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.HashedPassword,
		Salt:      user.Salt,
		Avatar:    user.Avatar,
		DateOfBirth: customtime.CustomTime{
			Time: user.DateOfBirth.AsTime(),
		},
		CreatedAt: customtime.CustomTime{
			Time: user.CreatedAt.AsTime(),
		},
		UpdatedAt: customtime.CustomTime{
			Time: user.UpdatedAt.AsTime(),
		},
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
