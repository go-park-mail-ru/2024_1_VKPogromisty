package auth_test

import (
	"net/http"
	"socio/domain"
	"socio/errors"
	repository "socio/internal/repository/map"
	customtime "socio/pkg/time"
	"socio/usecase/auth"
	"sync"
	"testing"
	"time"
)

func TestRegistrateUser(t *testing.T) {
	userStorage := repository.NewUsers(customtime.MockTimeProvider{}, &sync.Map{})
	sessionStorage, _ := repository.NewSessions(&sync.Map{})
	authService := auth.NewService(userStorage, sessionStorage)

	tests := []struct {
		name    string
		input   auth.RegistrationInput
		wantErr error
	}{
		{
			name: "Valid registration data",
			input: auth.RegistrationInput{
				FirstName:      "John",
				LastName:       "Doe",
				Email:          "john@example.com",
				Password:       "password",
				RepeatPassword: "password",
				DateOfBirth:    "1990-01-01",
				Avatar:         nil,
			},
			wantErr: nil,
		},
		{
			name: "Invalid email",
			input: auth.RegistrationInput{
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
			input: auth.RegistrationInput{
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
	userStorage := repository.NewUsers(customtime.MockTimeProvider{}, &sync.Map{})
	sessionStorage, _ := repository.NewSessions(&sync.Map{})
	authService := auth.NewService(userStorage, sessionStorage)

	tests := []struct {
		name    string
		input   auth.LoginInput
		wantErr error
	}{
		{
			name: "Valid login data",
			input: auth.LoginInput{
				Email:    "test@example.com",
				Password: "password",
			},
			wantErr: nil,
		},
		{
			name: "Invalid email",
			input: auth.LoginInput{
				Email:    "invalid@example.com",
				Password: "password",
			},
			wantErr: errors.ErrInvalidLoginData,
		},
		{
			name: "Invalid password",
			input: auth.LoginInput{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			wantErr: errors.ErrInvalidLoginData,
		},
	}

	authService.UserStorage.StoreUser(&domain.User{
		Email:    "test@example.com",
		Password: "password",
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := authService.Login(tt.input)
			if err != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	userStorage := repository.NewUsers(customtime.MockTimeProvider{}, &sync.Map{})
	sessionStorage, _ := repository.NewSessions(&sync.Map{})
	authService := auth.NewService(userStorage, sessionStorage)

	sessionID, _ := authService.SessionStorage.CreateSession(0)

	tests := []struct {
		name    string
		session *http.Cookie
		wantErr error
	}{
		{
			name:    "Valid session",
			session: &http.Cookie{Name: "session", Value: sessionID},
			wantErr: nil,
		},
		{
			name:    "Invalid session",
			session: &http.Cookie{Name: "session", Value: "invalidSessionValue"},
			wantErr: errors.ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.Logout(tt.session)
			if err != tt.wantErr {
				t.Errorf("Logout() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil {
				if _, err := authService.SessionStorage.GetUserIDBySession(tt.session.Value); err == nil {
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
	userStorage := repository.NewUsers(customtime.MockTimeProvider{}, &sync.Map{})
	sessionStorage, _ := repository.NewSessions(&sync.Map{})
	authService := auth.NewService(userStorage, sessionStorage)

	sessionID, _ := authService.SessionStorage.CreateSession(0)

	tests := []struct {
		name    string
		session *http.Cookie
		wantErr error
	}{
		{
			name:    "Valid session",
			session: &http.Cookie{Name: "session", Value: sessionID},
			wantErr: nil,
		},
		{
			name:    "Invalid session",
			session: &http.Cookie{Name: "session", Value: "invalidSessionValue"},
			wantErr: errors.ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := authService.IsAuthorized(tt.session)
			if err != tt.wantErr {
				t.Errorf("IsAuthorized() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
