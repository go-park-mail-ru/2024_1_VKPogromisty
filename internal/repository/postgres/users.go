package repository

import (
	"context"
	"fmt"
	"socio/domain"
	"socio/errors"
	"socio/pkg/hash"

	customtime "socio/pkg/time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

const (
	getUserByIDQuery = `
	SELECT id,
		first_name,
		last_name,
		email,
		hashed_password,
		salt,
		avatar,
		date_of_birth,
		created_at,
		updated_at
	FROM public.user
	WHERE id = $1;
	`
	getUserByIDWithSubsInfoQuery = `
	SELECT id,
		first_name,
		last_name,
		email,
		hashed_password,
		salt,
		avatar,
		date_of_birth,
		created_at,
		updated_at,
		CASE
			WHEN EXISTS (
				SELECT 1
				FROM subscription
				WHERE subscriber_id = $2
					AND subscribed_to_id = $1
			) THEN TRUE
			ELSE FALSE
		END AS is_subscribed_to,
		CASE
			WHEN EXISTS (
				SELECT 1
				FROM subscription
				WHERE subscriber_id = $1
					AND subscribed_to_id = $2
			) THEN TRUE
			ELSE FALSE
		END AS is_subscriber
	FROM public.user
	WHERE id = $1;
	`
	getUserByEmailQuery = `
		SELECT id,
		first_name,
		last_name,
		email,
		hashed_password,
		salt,
		avatar,
		date_of_birth,
		created_at,
		updated_at
	FROM public.user
	WHERE email = $1;
	`
	storeUserQuery = `
	INSERT INTO public.user (
			first_name,
			last_name,
			email,
			hashed_password,
			salt,
			avatar,
			date_of_birth
		)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id,
		first_name,
		last_name,
		email,
		avatar,
		date_of_birth,
		created_at,
		updated_at;
	`
	refreshSaltAndRehashPasswordQuery = `
	UPDATE public.user
	SET hashed_password = $1,
		salt = $2
	WHERE id = $3;
	`
)

type Users struct {
	db *pgxpool.Pool
	TP customtime.TimeProvider
}

func NewUsers(db *pgxpool.Pool, tp customtime.TimeProvider) *Users {
	return &Users{
		db: db,
		TP: tp,
	}
}

func (s *Users) GetUserByID(userID uint) (user *domain.User, err error) {
	user = &domain.User{}

	err = s.db.QueryRow(context.Background(), getUserByIDQuery, userID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Salt,
		&user.Avatar,
		&user.DateOfBirth.Time,
		&user.CreatedAt.Time,
		&user.UpdatedAt.Time,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
			return
		}

		return
	}

	return
}

func (s *Users) GetUserByIDWithSubsInfo(userID, authorizedUserID uint) (user *domain.User, isSubscribedTo bool, isSubscriber bool, err error) {
	user = &domain.User{}

	err = s.db.QueryRow(context.Background(), getUserByIDWithSubsInfoQuery, userID, authorizedUserID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Salt,
		&user.Avatar,
		&user.DateOfBirth.Time,
		&user.CreatedAt.Time,
		&user.UpdatedAt.Time,
		&isSubscribedTo,
		&isSubscriber,
	)
	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
			return
		}

		return
	}

	return
}

func (s *Users) GetUserByEmail(email string) (user *domain.User, err error) {
	user = &domain.User{}

	err = s.db.QueryRow(context.Background(), getUserByEmailQuery, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Salt,
		&user.Avatar,
		&user.DateOfBirth.Time,
		&user.CreatedAt.Time,
		&user.UpdatedAt.Time,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
			return
		}

		return
	}

	return
}

func (s *Users) StoreUser(user *domain.User) (err error) {
	salt := uuid.NewString()
	user.Password = hash.HashPassword(user.Password, []byte(salt))
	user.Salt = salt

	err = s.db.QueryRow(context.Background(), storeUserQuery,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Salt,
		user.Avatar,
		user.DateOfBirth.Time,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Avatar,
		&user.DateOfBirth.Time,
		&user.CreatedAt.Time,
		&user.UpdatedAt.Time,
	)
	if err != nil {
		return
	}

	return
}

func (s *Users) RefreshSaltAndRehashPassword(user *domain.User, password string) (err error) {
	salt := uuid.NewString()
	user.Password = hash.HashPassword(password, []byte(salt))
	user.Salt = salt

	result, err := s.db.Exec(context.Background(), refreshSaltAndRehashPasswordQuery,
		user.Password,
		user.Salt,
		user.ID,
	)
	if err != nil {
		return
	}

	if result.RowsAffected() != 1 {
		return errors.ErrRowsAffected
	}

	return
}
