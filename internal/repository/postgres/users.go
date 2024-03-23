package repository

import (
	"database/sql"
	"socio/domain"
	"socio/errors"
	"socio/pkg/hash"

	customtime "socio/pkg/time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Users struct {
	db *sql.DB
	TP customtime.TimeProvider
}

func NewUsers(db *sql.DB, tp customtime.TimeProvider) *Users {
	return &Users{
		db: db,
		TP: tp,
	}
}

func (s *Users) GetUserByID(userID uint) (user *domain.User, err error) {
	user = &domain.User{}

	err = s.db.QueryRow(`
		SELECT 
			id, 
			first_name, 
			last_name, 
			email, 
			avatar, 
			date_of_birth, 
			registration_date 
		FROM users WHERE id = $1;`,
		userID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Avatar,
		&user.DateOfBirth.Time,
		&user.RegistrationDate.Time,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.ErrNotFound
			return
		}

		return
	}

	return
}

func (s *Users) GetUserByEmail(email string) (user *domain.User, err error) {
	user = &domain.User{}

	err = s.db.QueryRow(`
		SELECT 
			id, 
			first_name, 
			last_name, 
			email, 
			avatar, 
			date_of_birth, 
			registration_date 
		FROM users WHERE email = $1;`,
		email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Avatar,
		&user.DateOfBirth.Time,
		&user.RegistrationDate.Time,
	)
	if err != nil {
		if err == sql.ErrNoRows {
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

	err = s.db.QueryRow(`
		INSERT INTO users 
		(first_name, last_name, email, password, salt, avatar, date_of_birth) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, first_name, last_name, email, avatar, date_of_birth, registration_date;`,
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
		&user.RegistrationDate.Time,
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

	result, err := s.db.Exec(`
		UPDATE users 
		SET 
			password = $1, 
			salt = $2 
		WHERE id = $3;`,
		user.Password, user.Salt, user.ID)
	if err != nil {
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return
	}
	if rows != 1 {
		return errors.ErrRowsAffected
	}

	return
}
