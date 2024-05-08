package user

import (
	"context"
	"socio/domain"
)

func (p *Service) CreatePublicGroupAdmin(ctx context.Context, publicGroupAdmin *domain.PublicGroupAdmin) (newPublicGroupAdmin *domain.PublicGroupAdmin, err error) {
	_, err = p.UserStorage.GetUserByID(ctx, publicGroupAdmin.UserID)
	if err != nil {
		return
	}

	newPublicGroupAdmin, err = p.UserStorage.StorePublicGroupAdmin(ctx, publicGroupAdmin)
	if err != nil {
		return
	}

	return
}

func (p *Service) DeletePublicGroupAdmin(ctx context.Context, publicGroupAdmin *domain.PublicGroupAdmin) (err error) {
	err = p.UserStorage.DeletePublicGroupAdmin(ctx, publicGroupAdmin)
	if err != nil {
		return
	}

	return
}

func (p *Service) GetAdminsByPublicGroupID(ctx context.Context, publicGroupID uint) (admins []*domain.User, err error) {
	admins, err = p.UserStorage.GetAdminsByPublicGroupID(ctx, publicGroupID)
	if err != nil {
		return
	}

	return
}

func (p *Service) CheckIfUserIsAdmin(ctx context.Context, publicGroupID, userID uint) (isAdmin bool, err error) {
	isAdmin, err = p.UserStorage.CheckIfUserIsAdmin(ctx, publicGroupID, userID)
	if err != nil {
		return
	}

	return
}
