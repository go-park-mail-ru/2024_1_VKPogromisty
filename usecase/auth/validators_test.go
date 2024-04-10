package auth_test

import (
	"socio/errors"
	"socio/usecase/auth"

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
