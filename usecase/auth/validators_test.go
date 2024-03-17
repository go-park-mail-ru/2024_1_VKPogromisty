package auth_test

import (
	"socio/errors"
	repository "socio/internal/repository/map"
	customtime "socio/pkg/time"
	"socio/usecase/auth"
	"sync"

	"testing"
)

func TestCheckEmptyFields(t *testing.T) {
	tests := []struct {
		name  string
		input auth.RegistrationInput
		want  error
	}{
		{"All fields filled", auth.RegistrationInput{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Password: "password", RepeatPassword: "password"}, nil},
		{"Missing fields", auth.RegistrationInput{FirstName: "John", LastName: "", Email: "john.doe@example.com", Password: "password", RepeatPassword: "password"}, errors.ErrMissingFields},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := auth.CheckEmptyFields(tt.input); got != tt.want {
				t.Errorf("CheckEmptyFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			if got := auth.ValidateEmail(tt.email); got != tt.want {
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
			if got := auth.ValidatePassword(tt.password, tt.repeatPassword); got != tt.want {
				t.Errorf("ValidatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckDuplicatedUser(t *testing.T) {
	userStorage := repository.NewUsers(customtime.MockTimeProvider{}, &sync.Map{})
	sessionStorage, _ := repository.NewSessions(&sync.Map{})
	service := auth.NewService(customtime.MockTimeProvider{}, userStorage, sessionStorage)

	tests := []struct {
		name  string
		input auth.RegistrationInput
		want  error
	}{
		{"Duplicated user", auth.RegistrationInput{Email: "petr09mitin@mail.ru", Password: "password", RepeatPassword: "password"}, errors.ErrEmailsDuplicate},
		{"Unique user", auth.RegistrationInput{Email: "petr01mitin@mail.ru", Password: "password", RepeatPassword: "password"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service.CheckDuplicatedEmail(tt.input.Email); got != tt.want {
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
			if got := auth.ValidateDateOfBirth(tt.date); got != tt.want {
				t.Errorf("ValidateDateOfBirth() = %v, want %v", got, tt.want)
			}
		})
	}
}
