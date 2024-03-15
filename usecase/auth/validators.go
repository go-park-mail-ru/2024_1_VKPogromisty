package auth

import (
	"net/mail"
	"socio/errors"
	customtime "socio/pkg/time"
	"time"
)

func CheckEmptyFields(userInput RegistrationInput) (err error) {
	if len(userInput.FirstName) == 0 || len(userInput.LastName) == 0 || len(userInput.Email) == 0 || len(userInput.Password) == 0 || len(userInput.RepeatPassword) == 0 {
		err = errors.ErrMissingFields
		return
	}
	return
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

func (a *Service) CheckDuplicatedEmail(email string) (err error) {
	if _, err = a.UserStorage.GetUserByEmail(email); err != errors.ErrNotFound {
		err = errors.ErrEmailsDuplicate
		return
	}

	return nil
}

func (a *Service) ValidateUserInput(userInput RegistrationInput) (err error) {
	if err = CheckEmptyFields(userInput); err != nil {
		return
	}

	if err = ValidateEmail(userInput.Email); err != nil {
		return
	}

	if err = ValidatePassword(userInput.Password, userInput.RepeatPassword); err != nil {
		return
	}

	if err = a.CheckDuplicatedEmail(userInput.Email); err != nil {
		return
	}

	if err = ValidateDateOfBirth(userInput.DateOfBirth); err != nil {
		return
	}

	return
}
