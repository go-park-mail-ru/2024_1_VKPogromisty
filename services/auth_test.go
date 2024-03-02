package services_test

import (
	"net/http"
	"socio/errors"
	"socio/services"
	"socio/utils"
	"testing"
	"time"
)

func TestRegistrateUser(t *testing.T) {
	authService := services.NewAuthService(utils.MockTimeProvider{})

	tests := []struct {
		name    string
		input   services.RegistrationInput
		wantErr error
	}{
		{
			name: "Valid registration data",
			input: services.RegistrationInput{
				FirstName:      "John",
				LastName:       "Doe",
				Email:          "john@example.com",
				Password:       "password",
				RepeatPassword: "password",
				DateOfBirth:    "1990-01-01",
				Avatar:         nil, // Assuming Avatar is of type []byte
			},
			wantErr: nil,
		},
		{
			name: "Invalid email",
			input: services.RegistrationInput{
				FirstName:      "John",
				LastName:       "Doe",
				Email:          "invalid",
				Password:       "password",
				RepeatPassword: "password",
				DateOfBirth:    "1990-01-01",
				Avatar:         nil,
			},
			wantErr: errors.ErrInvalidEmail,
		},
		{
			name: "Invalid date of birth",
			input: services.RegistrationInput{
				FirstName:      "John",
				LastName:       "Doe",
				Email:          "john1@example.com",
				Password:       "password",
				RepeatPassword: "password",
				DateOfBirth:    "invalid",
				Avatar:         nil,
			},
			wantErr: errors.ErrInvalidDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := authService.RegistrateUser(tt.input)
			if err != tt.wantErr {
				t.Errorf("RegistrateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	authService := services.NewAuthService(utils.MockTimeProvider{})

	tests := []struct {
		name    string
		input   services.LoginInput
		wantErr error
	}{
		{
			name: "Valid login data",
			input: services.LoginInput{
				Email:    "test@example.com",
				Password: "password",
			},
			wantErr: nil,
		},
		{
			name: "Invalid email",
			input: services.LoginInput{
				Email:    "invalid@example.com",
				Password: "password",
			},
			wantErr: errors.ErrInvalidLoginData,
		},
		{
			name: "Invalid password",
			input: services.LoginInput{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			wantErr: errors.ErrInvalidLoginData,
		},
	}

	// Add a valid user to the AuthService
	authService.Users.Store("test@example.com", &services.User{
		Email:    "test@example.com",
		Password: utils.HashPassword("password", []byte("salt")),
		Salt:     "salt",
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := authService.Login(tt.input)
			if err != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	authService := services.NewAuthService(utils.MockTimeProvider{})

	tests := []struct {
		name    string
		session *http.Cookie
		wantErr error
	}{
		{
			name:    "Valid session",
			session: &http.Cookie{Name: "session", Value: "validSessionValue"},
			wantErr: nil,
		},
		{
			name:    "Invalid session",
			session: &http.Cookie{Name: "session", Value: "invalidSessionValue"},
			wantErr: errors.ErrUnauthorized,
		},
	}

	authService.Sessions.Store("validSessionValue", true)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.Logout(tt.session)
			if err != tt.wantErr {
				t.Errorf("Logout() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil {
				if _, ok := authService.Sessions.Load(tt.session.Value); ok {
					t.Errorf("Logout() did not delete session")
				}
			}

			if tt.wantErr == errors.ErrUnauthorized {
				if tt.session.Expires.After(time.Now()) {
					t.Errorf("Logout() did not set cookie expiry to a past date")
				}
			}
		})
	}
}

func TestIsAuthorized(t *testing.T) {
	authService := services.NewAuthService(utils.MockTimeProvider{})

	tests := []struct {
		name    string
		session *http.Cookie
		wantErr error
	}{
		{
			name:    "Valid session",
			session: &http.Cookie{Name: "session", Value: "validSessionValue"},
			wantErr: nil,
		},
		{
			name:    "Invalid session",
			session: &http.Cookie{Name: "session", Value: "invalidSessionValue"},
			wantErr: errors.ErrUnauthorized,
		},
	}

	authService.Sessions.Store("validSessionValue", true)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.IsAuthorized(tt.session)
			if err != tt.wantErr {
				t.Errorf("IsAuthorized() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
