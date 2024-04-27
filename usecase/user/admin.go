package user

import (
	"context"
	"fmt"
	"socio/domain"
	"socio/errors"
)

type AdminWithUser struct {
	Admin *domain.Admin `json:"admin"`
	User  *domain.User  `json:"user"`
}

func (s *Service) GetAdminByUserID(ctx context.Context, userID uint) (admin *domain.Admin, err error) {
	admin, err = s.UserStorage.GetAdminByUserID(userID)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetAdmins(ctx context.Context) (admins []AdminWithUser, err error) {
	admins, err = s.UserStorage.GetAdmins()
	if err != nil {
		return
	}

	return
}

func (s *Service) CreateAdmin(ctx context.Context, admin *domain.Admin) (newAdmin *domain.Admin, err error) {
	fmt.Println(admin)

	_, err = s.UserStorage.GetUserByID(ctx, admin.UserID)
	if err != nil {
		err = errors.ErrInvalidBody
		return
	}

	newAdmin, err = s.UserStorage.StoreAdmin(admin)
	if err != nil {
		return
	}

	return
}

func (s *Service) DeleteAdmin(ctx context.Context, adminID uint) (err error) {
	err = s.UserStorage.DeleteAdmin(adminID)
	if err != nil {
		return
	}

	return
}