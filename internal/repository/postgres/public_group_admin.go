package repository

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/pkg/contextlogger"
)

const (
	storePublicGroupAdminQuery = `
	INSERT INTO public.public_group_admin (public_group_id, user_id)
	VALUES ($1, $2)
	RETURNING id,
		public_group_id,
		user_id,
		created_at,
		updated_at;
	`
	deletePublicGroupAdminQuery = `
	DELETE FROM public.public_group_admin
	WHERE public_group_id = $1
		AND user_id = $2;
	`
	getAdminsByPublicGroupIDQuery = `
	SELECT u.id,
		u.first_name,
		u.last_name,
		u.email,
		u.avatar,
		u.hashed_password,
		u.salt,
		u.date_of_birth,
		u.created_at,
		u.updated_at
	FROM public.user AS u
		JOIN public.public_group_admin AS pga ON pga.user_id = u.id
	WHERE pga.public_group_id = $1;
	`
	checkIfUserIsAdminQuery = `
	SELECT EXISTS (
        SELECT 1
        FROM public.public_group_admin
        WHERE public_group_id = $1
            AND user_id = $2
    );
	`
)

func (s *Users) StorePublicGroupAdmin(ctx context.Context, publicGroupAdmin *domain.PublicGroupAdmin) (newPublicGroupAdmin *domain.PublicGroupAdmin, err error) {
	newPublicGroupAdmin = new(domain.PublicGroupAdmin)

	contextlogger.LogSQL(ctx, GetPostsByGroupSubIDsAndUserSubIDsQuery, publicGroupAdmin.PublicGroupID, publicGroupAdmin.UserID)

	err = s.db.QueryRow(
		context.Background(),
		storePublicGroupAdminQuery,
		publicGroupAdmin.PublicGroupID,
		publicGroupAdmin.UserID,
	).Scan(
		&newPublicGroupAdmin.ID,
		&newPublicGroupAdmin.PublicGroupID,
		&newPublicGroupAdmin.UserID,
		&newPublicGroupAdmin.CreatedAt.Time,
		&newPublicGroupAdmin.UpdatedAt.Time,
	)
	if err != nil {
		return
	}

	return
}

func (s *Users) DeletePublicGroupAdmin(ctx context.Context, publicGroupAdmin *domain.PublicGroupAdmin) (err error) {
	contextlogger.LogSQL(ctx, deletePublicGroupAdminQuery, publicGroupAdmin.PublicGroupID, publicGroupAdmin.UserID)

	result, err := s.db.Exec(context.Background(), deletePublicGroupAdminQuery, publicGroupAdmin.PublicGroupID, publicGroupAdmin.UserID)
	if err != nil {
		return
	}

	if result.RowsAffected() == 0 {
		err = errors.ErrNotFound
		return
	}

	return
}

func (s *Users) GetAdminsByPublicGroupID(ctx context.Context, publicGroupID uint) (admins []*domain.User, err error) {
	admins = make([]*domain.User, 0)

	contextlogger.LogSQL(ctx, getAdminsByPublicGroupIDQuery, publicGroupID)

	rows, err := s.db.Query(context.Background(), getAdminsByPublicGroupIDQuery, publicGroupID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		admin := new(domain.User)

		err = rows.Scan(
			&admin.ID,
			&admin.FirstName,
			&admin.LastName,
			&admin.Email,
			&admin.Avatar,
			&admin.Password,
			&admin.Salt,
			&admin.DateOfBirth.Time,
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

func (s *Users) CheckIfUserIsAdmin(ctx context.Context, publicGroupID, userID uint) (isAdmin bool, err error) {
	contextlogger.LogSQL(ctx, checkIfUserIsAdminQuery, publicGroupID, userID)

	err = s.db.QueryRow(context.Background(), checkIfUserIsAdminQuery, publicGroupID, userID).Scan(&isAdmin)
	if err != nil {
		return
	}

	return
}
