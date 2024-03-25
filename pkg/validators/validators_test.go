package validators_test

import (
	"socio/errors"
	repository "socio/internal/repository/map"
	customtime "socio/pkg/time"
	"socio/pkg/validators"
	"sync"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  error
	}{
		{"Valid email", "petr09mitin@mail.ru", nil},
		{"Invalid email", "uvy@", errors.ErrInvalidEmail},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validators.ValidateEmail(tt.email); got != tt.want {
				t.Errorf("ValidateEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		repeatPassword string
		want           error
	}{
		{"Matching passwords", "password", "password", nil},
		{"Not matching passwords", "password", "passwrd", errors.ErrNotMatchingPasswords},
		{"Short password", "pass", "pass", errors.ErrPasswordMinLength},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validators.ValidatePassword(tt.password, tt.repeatPassword); got != tt.want {
				t.Errorf("ValidatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckDuplicatedUser(t *testing.T) {
	userStorage := repository.NewUsers(customtime.MockTimeProvider{}, &sync.Map{})

	tests := []struct {
		name  string
		email string
		want  error
	}{
		{"Duplicated user", "petr09mitin@mail.ru", errors.ErrEmailsDuplicate},
		{"Unique user", "petr01mitin@mail.ru", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validators.CheckDuplicatedEmail(tt.email, userStorage); got != tt.want {
				t.Errorf("CheckDuplicatedUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateDateOfBirth(t *testing.T) {
	tests := []struct {
		name string
		date string
		want error
	}{
		{"Valid date", "1990-01-01", nil},
		{"Invalid date", "1990-13-01", errors.ErrInvalidDate},
		{"Future date", "2029-01-01", errors.ErrInvalidDate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validators.ValidateDateOfBirth(tt.date); got != tt.want {
				t.Errorf("ValidateDateOfBirth() = %v, want %v", got, tt.want)
			}
		})
	}
}
