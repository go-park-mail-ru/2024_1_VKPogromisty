package repository

import (
	"context"
	"socio/domain"
	"socio/errors"

	"github.com/jackc/pgx/v4"
)

const (
	storeAdminQuery = `
	INSERT INTO public.admin (user_id)
	VALUES ($1)
	RETURNING id, user_id, created_at, updated_at;
	`
	deleteAdminQuery = `
	DELETE FROM public.admin
	WHERE id = $1;
	`
	getAdminByUserIDQuery = `
	SELECT id, user_id, created_at, updated_at
	FROM public.admin
	WHERE user_id = $1;
	`
	getAdminsQuery = `
	SELECT id, user_id, created_at, updated_at
	FROM public.admin;
	`
)

func (c *Users) GetAdmins() (admins []*domain.Admin, err error) {
	rows, err := c.db.Query(context.Background(), getAdminsQuery)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
			return
		}

		return
	}

	for rows.Next() {
		admin := new(domain.Admin)
		err = rows.Scan(
			&admin.ID,
			&admin.UserID,
			&admin.CreatedAt.Time,
			&admin.UpdatedAt.Time,
		)
		if err != nil {
			return
		}

		admins = append(admins, admin)
	}

	return
}

func (c *Users) StoreAdmin(admin *domain.Admin) (newAdmin *domain.Admin, err error) {
	newAdmin = new(domain.Admin)
	err = c.db.QueryRow(context.Background(), storeAdminQuery, admin.UserID).Scan(
		&newAdmin.ID,
		&newAdmin.UserID,
		&newAdmin.CreatedAt.Time,
		&newAdmin.UpdatedAt.Time,
	)

	if err != nil {
		return
	}

	return
}

func (c *Users) GetAdminByUserID(userID uint) (admin *domain.Admin, err error) {
	admin = new(domain.Admin)
	err = c.db.QueryRow(context.Background(), getAdminByUserIDQuery, userID).Scan(
		&admin.ID,
		&admin.UserID,
		&admin.CreatedAt.Time,
		&admin.UpdatedAt.Time,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ErrForbidden
			return
		}

		return
	}

	return

}

func (c *Users) DeleteAdmin(adminID uint) (err error) {
	_, err = c.db.Exec(context.Background(), deleteAdminQuery, adminID)
	if err != nil {
		return
	}

	return
}
