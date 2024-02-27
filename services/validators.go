package services

import (
	"net/mail"
	"socio/utils"
	"time"
)

func ValidateUserInput(userInput RegistrationInput, service *AuthService) (err error) {
	if len(userInput.FirstName) == 0 || len(userInput.LastName) == 0 || len(userInput.Email) == 0 || len(userInput.Password) == 0 || len(userInput.RepeatPassword) == 0 {
		err = utils.ErrMissingFields
		return
	}

	_, err = mail.ParseAddress(userInput.Email)
	if err != nil {
		err = utils.ErrInvalidData
		return
	}

	if userInput.Password != userInput.RepeatPassword {
		err = utils.ErrInvalidData
		return
	}

	if _, ok := service.users.Load(userInput.Email); ok {
		err = utils.ErrInvalidData
		return
	}

	dateOfBirth, err := time.Parse(utils.DateFormat, userInput.DateOfBirth)
	if err != nil || dateOfBirth.After(time.Now()) {
		err = utils.ErrInvalidData
		return
	}

	return
}
