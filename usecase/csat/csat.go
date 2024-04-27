package csat

import (
	"socio/domain"
)

type CSATStorage interface {
	GetAdmins() (admins []*domain.Admin, err error)
	GetAdminByUserID(userID uint) (admin *domain.Admin, err error)
	StoreAdmin(admin *domain.Admin) (newAdmin *domain.Admin, err error)
	DeleteAdmin(adminID uint) (err error)
}

type Service struct {
	storage CSATStorage
}

func NewService(storage CSATStorage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) GetAdminByUserID(userID uint) (admin *domain.Admin, err error) {
	admin, err = s.storage.GetAdminByUserID(userID)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetAdmins() (admins []*domain.Admin, err error) {
	admins, err = s.storage.GetAdmins()
	if err != nil {
		return
	}

	return
}

func (s *Service) CreateAdmin(creatorID uint, admin *domain.Admin) (newAdmin *domain.Admin, err error) {
	newAdmin, err = s.storage.StoreAdmin(admin)
	if err != nil {
		return
	}

	return
}

func (s *Service) DeleteAdmin(creatorID, adminID uint) (err error) {
	err = s.storage.DeleteAdmin(adminID)
	if err != nil {
		return
	}

	return
}
