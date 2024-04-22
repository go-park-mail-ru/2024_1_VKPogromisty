package user

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/pkg/validators"
)

func (p *Service) ValidateUpdateUserInput(ctx context.Context, userInput UpdateUserInput, oldUser *domain.User) (err error) {
	if len(userInput.Email) > 0 {
		if err = validators.ValidateEmail(userInput.Email); err != nil {
			return
		}
	}

	if len(userInput.Password) > 0 {
		if err = validators.ValidatePassword(userInput.Password, userInput.RepeatPassword); err != nil {
			return
		}
	}

	if len(userInput.Email) > 0 && userInput.Email != oldUser.Email {
		if err = validators.CheckDuplicatedEmail(ctx, userInput.Email, p.UserStorage); err != nil {
			return
		}
	}

	if len(userInput.DateOfBirth) > 0 {
		if err = validators.ValidateDateOfBirth(userInput.DateOfBirth); err != nil {
			return
		}
	}

	return
}

func CheckEmptyFields(userInput CreateUserInput) (err error) {
	if len(userInput.FirstName) == 0 || len(userInput.LastName) == 0 || len(userInput.Email) == 0 || len(userInput.Password) == 0 || len(userInput.RepeatPassword) == 0 {
		err = errors.ErrMissingFields
		return
	}
	return
}

func (p *Service) ValidateCreateUserInput(ctx context.Context, userInput CreateUserInput) (err error) {
	if err = CheckEmptyFields(userInput); err != nil {
		return
	}

	if err = validators.ValidateEmail(userInput.Email); err != nil {
		return
	}

	if err = validators.ValidatePassword(userInput.Password, userInput.RepeatPassword); err != nil {
		return
	}

	if err = validators.CheckDuplicatedEmail(ctx, userInput.Email, p.UserStorage); err != nil {
		return
	}

	if err = validators.ValidateDateOfBirth(userInput.DateOfBirth); err != nil {
		return
	}

	return
}
