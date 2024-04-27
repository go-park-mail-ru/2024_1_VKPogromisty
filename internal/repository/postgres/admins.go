package repository

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/usecase/user"

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
	SELECT a.id, a.user_id, a.created_at, a.updated_at, u.first_name, u.last_name, u.email, u.avatar, u.date_of_birth, u.created_at, u.updated_at
	FROM public.admin a
	JOIN public.user u ON a.user_id = u.id;
	`
)

func (c *Users) GetAdmins() (admins []user.AdminWithUser, err error) {
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
		adminUser := new(domain.User)
		err = rows.Scan(
			&admin.ID,
			&admin.UserID,
			&admin.CreatedAt.Time,
			&admin.UpdatedAt.Time,
			&adminUser.FirstName,
			&adminUser.LastName,
			&adminUser.Email,
			&adminUser.Avatar,
			&adminUser.DateOfBirth.Time,
			&adminUser.CreatedAt.Time,
			&adminUser.UpdatedAt.Time,
		)
		if err != nil {
			return
		}

		adminUser.ID = admin.UserID

		admins = append(admins, user.AdminWithUser{
			Admin: admin,
			User:  adminUser,
		})
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
