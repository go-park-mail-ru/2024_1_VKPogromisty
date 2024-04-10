package validators

import (
	"context"
	"fmt"
	"net/mail"
	"socio/domain"
	"socio/errors"
	customtime "socio/pkg/time"
	"time"
)

type UserStorage interface {
	GetUserByEmail(ctx context.Context, email string) (user *domain.User, err error)
}

func ValidateEmail(email string) (err error) {
	_, err = mail.ParseAddress(email)
	if err != nil {
		err = errors.ErrInvalidEmail
		return
	}
	return
}

func ValidatePassword(password string, repeatPassword string) (err error) {
	if password != repeatPassword {
		err = errors.ErrNotMatchingPasswords
		return
	}

	if len(password) < 6 {
		err = errors.ErrPasswordMinLength
		return
	}

	return
}

func ValidateDateOfBirth(date string) (err error) {
	dateOfBirth, err := time.Parse(customtime.DateFormat, date)
	if err != nil {
		err = errors.ErrInvalidDate
		return
	}

	leftTimeBound, _ := time.Parse(customtime.DateFormat, "1900-01-01")
	if dateOfBirth.Before(leftTimeBound) || dateOfBirth.After(time.Now()) {
		err = errors.ErrInvalidDate
		return
	}

	return
}

func CheckDuplicatedEmail(ctx context.Context, email string, userStorage UserStorage) (err error) {
	if _, err = userStorage.GetUserByEmail(ctx, email); err != errors.ErrNotFound {
		fmt.Println("CheckDuplicatedEmail: ", err)
		err = errors.ErrEmailsDuplicate
		return
	}

	return nil
}
