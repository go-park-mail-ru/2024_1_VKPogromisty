package auth

import (
	"context"
	"socio/errors"
	"socio/pkg/validators"
)

func CheckEmptyFields(userInput RegistrationInput) (err error) {
	if len(userInput.FirstName) == 0 || len(userInput.LastName) == 0 || len(userInput.Email) == 0 || len(userInput.Password) == 0 || len(userInput.RepeatPassword) == 0 {
		err = errors.ErrMissingFields
		return
	}
	return
}

func (a *Service) ValidateUserInput(ctx context.Context, userInput RegistrationInput) (err error) {
	if err = CheckEmptyFields(userInput); err != nil {
		return
	}

	if err = validators.ValidateEmail(userInput.Email); err != nil {
		return
	}

	if err = validators.ValidatePassword(userInput.Password, userInput.RepeatPassword); err != nil {
		return
	}

	if err = validators.CheckDuplicatedEmail(ctx, userInput.Email, a.UserStorage); err != nil {
		return
	}

	if err = validators.ValidateDateOfBirth(userInput.DateOfBirth); err != nil {
		return
	}

	return
}
