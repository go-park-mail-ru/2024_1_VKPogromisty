package services_test

import (
	"socio/errors"
	"socio/services"
	"socio/utils"
	"testing"
)

func TestCheckEmptyFields(t *testing.T) {
	tests := []struct {
		name  string
		input services.RegistrationInput
		want  error
	}{
		{"All fields filled", services.RegistrationInput{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Password: "password", RepeatPassword: "password"}, nil},
		{"Missing fields", services.RegistrationInput{FirstName: "John", LastName: "", Email: "john.doe@example.com", Password: "password", RepeatPassword: "password"}, errors.ErrMissingFields},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.CheckEmptyFields(tt.input); got != tt.want {
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
			if got := services.ValidateEmail(tt.email); got != tt.want {
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
			if got := services.ValidatePassword(tt.password, tt.repeatPassword); got != tt.want {
				t.Errorf("ValidatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckDuplicatedUser(t *testing.T) {
	service := services.NewAuthService(utils.MockTimeProvider{})
	service.Users.Store("petr09mitin@mail.ru", &services.User{})
	tests := []struct {
		name  string
		input services.RegistrationInput
		want  error
	}{
		{"Duplicated user", services.RegistrationInput{Email: "petr09mitin@mail.ru", Password: "password", RepeatPassword: "password"}, errors.ErrEmailsDuplicate},
		{"Unique user", services.RegistrationInput{Email: "petr01mitin@mail.ru", Password: "password", RepeatPassword: "password"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.CheckDuplicatedUser(tt.input, service); got != tt.want {
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
			if got := services.ValidateDateOfBirth(tt.date); got != tt.want {
				t.Errorf("ValidateDateOfBirth() = %v, want %v", got, tt.want)
			}
		})
	}
}
