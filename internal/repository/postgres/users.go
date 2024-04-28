package repository

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/pkg/contextlogger"
	"socio/pkg/hash"

	customtime "socio/pkg/time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
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
	updateUserQuery = `
	UPDATE public.user
	SET 
		first_name = $2,
		last_name = $3,
		email = $4,
		hashed_password = $5,
		salt = $6,
		avatar = $7,
		date_of_birth = $8
	WHERE id = $1
	RETURNING id,
		first_name,
		last_name,
		email,
		avatar,
		date_of_birth,
		created_at,
		updated_at;
	`
	deleteUserPostsQuery = `
	DELETE FROM public.post
	WHERE author_id=$1;
	`
	deleteUserMessagesQuery = `
	DELETE FROM public.personal_message
	WHERE sender_id = $1
		OR receiver_id = $1;
	`
	deleteUserQuery = `
	DELETE FROM public.user
	WHERE id = $1;
	`
	getUsersByNameQuery = `
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
	WHERE to_tsvector('english', first_name || ' ' || last_name) @@ plainto_tsquery('english', $1)
	OR to_tsvector('russian', first_name || ' ' || last_name) @@ plainto_tsquery('russian', $1);
	`
)

type Users struct {
	db DBPool
	TP customtime.TimeProvider
}

func NewUsers(db DBPool, tp customtime.TimeProvider) *Users {
	return &Users{
		db: db,
		TP: tp,
	}
}

func (s *Users) GetUserByID(ctx context.Context, userID uint) (user *domain.User, err error) {
	contextlogger.LogSQL(ctx, getUserByIDQuery, userID)

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

func (s *Users) GetUserByIDWithSubsInfo(ctx context.Context, userID, authorizedUserID uint) (user *domain.User, isSubscribedTo bool, isSubscriber bool, err error) {
	contextlogger.LogSQL(ctx, getUserByIDWithSubsInfoQuery, userID, authorizedUserID)

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
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
			return
		}

		return
	}

	return
}

func (s *Users) GetUserByEmail(ctx context.Context, email string) (user *domain.User, err error) {
	contextlogger.LogSQL(ctx, getUserByEmailQuery, email)

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

func (s *Users) StoreUser(ctx context.Context, user *domain.User) (err error) {
	salt := uuid.NewString()
	user.Password = hash.HashPassword(user.Password, []byte(salt))
	user.Salt = salt

	contextlogger.LogSQL(ctx, storeUserQuery,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Avatar,
		user.DateOfBirth.Time,
	)

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

func (s *Users) UpdateUser(ctx context.Context, user *domain.User, prevPassword string) (updatedUser *domain.User, err error) {
	updatedUser = &domain.User{}

	if user.Password != prevPassword {
		salt := uuid.NewString()
		user.Password = hash.HashPassword(user.Password, []byte(salt))
		user.Salt = salt
	}

	contextlogger.LogSQL(ctx, updateUserQuery,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Avatar,
		user.DateOfBirth.Time,
	)

	err = s.db.QueryRow(context.Background(), updateUserQuery,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Salt,
		user.Avatar,
		user.DateOfBirth.Time,
	).Scan(
		&updatedUser.ID,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.Email,
		&updatedUser.Avatar,
		&updatedUser.DateOfBirth.Time,
		&updatedUser.CreatedAt.Time,
		&updatedUser.UpdatedAt.Time,
	)
	if err != nil {
		return
	}

	return
}

func (s *Users) DeleteUser(ctx context.Context, userID uint) (err error) {
	tx, err := s.db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			return
		}
		if err = tx.Rollback(context.Background()); err != nil && err != pgx.ErrTxClosed {
			return
		}

		err = nil
	}()

	contextlogger.LogSQL(ctx, deleteUserMessagesQuery, userID)

	_, err = tx.Exec(context.Background(), deleteUserMessagesQuery, userID)
	if err != nil {
		return
	}

	contextlogger.LogSQL(ctx, deleteUserPostsQuery, userID)

	_, err = tx.Exec(context.Background(), deleteUserPostsQuery, userID)
	if err != nil {
		return
	}

	contextlogger.LogSQL(ctx, deleteUserQuery, userID)

	result, err := tx.Exec(context.Background(), deleteUserQuery, userID)
	if err != nil {
		return
	}

	if result.RowsAffected() != 1 {
		return errors.ErrRowsAffected
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return
	}

	return
}

func (s *Users) SearchByName(ctx context.Context, query string) (users []*domain.User, err error) {
	contextlogger.LogSQL(ctx, getUsersByNameQuery, query)

	rows, err := s.db.Query(context.Background(), getUsersByNameQuery, query)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		user := new(domain.User)

		err = rows.Scan(
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
			return
		}

		users = append(users, user)
	}

	return
}
