package services

import (
	"net/mail"
	"socio/utils"
)

func ValidateUserInput(userInput RegistrationInput) (err error) {
	if len(userInput.FirstName) == 0 || len(userInput.LastName) == 0 || len(userInput.Email) == 0 || len(userInput.Password) == 0 {
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

	return
}
