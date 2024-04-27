package user

import (
	"socio/domain"
	customtime "socio/pkg/time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToAdminResponse(admin *domain.Admin) (res *AdminResponse) {
	if admin == nil {
		return nil
	}

	return &AdminResponse{
		Id:        uint64(admin.ID),
		UserId:    uint64(admin.UserID),
		CreatedAt: timestamppb.New(admin.CreatedAt.Time),
		UpdatedAt: timestamppb.New(admin.UpdatedAt.Time),
	}
}

func ToAdminsResponse(admins []*domain.Admin) (res []*AdminResponse) {
	if admins == nil {
		return nil
	}

	res = make([]*AdminResponse, 0, len(admins))
	for _, admin := range admins {
		res = append(res, ToAdminResponse(admin))
	}

	return
}

func ToAdmin(res *AdminResponse) (admin *domain.Admin) {
	if res == nil {
		return nil
	}

	return &domain.Admin{
		ID:     uint(res.Id),
		UserID: uint(res.UserId),
		CreatedAt: customtime.CustomTime{
			Time: res.CreatedAt.AsTime(),
		},
		UpdatedAt: customtime.CustomTime{
			Time: res.UpdatedAt.AsTime(),
		},
	}
}

func ToAdmins(res []*AdminResponse) (admins []*domain.Admin) {
	if res == nil {
		return nil
	}

	admins = make([]*domain.Admin, 0, len(res))
	for _, admin := range res {
		admins = append(admins, ToAdmin(admin))
	}

	return
}
