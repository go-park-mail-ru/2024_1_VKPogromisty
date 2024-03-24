package profile

import (
	"socio/domain"
	"socio/pkg/validators"
)

func (p *Service) ValidateUserInput(userInput UpdateUserInput, oldUser *domain.User) (err error) {
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
		if err = validators.CheckDuplicatedEmail(userInput.Email, p.UserStorage); err != nil {
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
