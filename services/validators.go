package services

import (
	"net/mail"
	"socio/errors"
	"socio/utils"
	"time"
)

func ValidateUserInput(userInput RegistrationInput, service *AuthService) (err error) {
	if len(userInput.FirstName) == 0 || len(userInput.LastName) == 0 || len(userInput.Email) == 0 || len(userInput.Password) == 0 || len(userInput.RepeatPassword) == 0 {
		err = errors.ErrMissingFields
		return
	}

	_, err = mail.ParseAddress(userInput.Email)
	if err != nil {
		err = errors.ErrInvalidEmail
		return
	}

	if userInput.Password != userInput.RepeatPassword {
		err = errors.ErrNotMatchingPasswords
		return
	}

	if len(userInput.Password) < 6 {
		err = errors.ErrPasswordMinLength
		return
	}

	if _, ok := service.users.Load(userInput.Email); ok {
		err = errors.ErrEmailsDuplicate
		return
	}

	dateOfBirth, err := time.Parse(utils.DateFormat, userInput.DateOfBirth)
	if err != nil || dateOfBirth.After(time.Now()) {
		err = errors.ErrInvalidDate
		return
	}

	return
}
