package validators_test

import (
	"context"
	"socio/domain"
	"socio/errors"
	mock_profile "socio/mocks/usecase/profile"
	"socio/pkg/validators"
	"testing"

	"github.com/golang/mock/gomock"
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

func TestCheckDuplicatedEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStorage := mock_profile.NewMockUserStorage(ctrl)
	ctx := context.Background()
	email := "test@example.com"

	// Test case: email not found
	mockUserStorage.EXPECT().GetUserByEmail(ctx, email).Return(&domain.User{}, errors.ErrNotFound)
	err := validators.CheckDuplicatedEmail(ctx, email, mockUserStorage)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	// Test case: email found
	mockUserStorage.EXPECT().GetUserByEmail(ctx, email).Return(&domain.User{}, nil)
	err = validators.CheckDuplicatedEmail(ctx, email, mockUserStorage)
	if err == nil || err.Error() != errors.ErrEmailsDuplicate.Error() {
		t.Errorf("Expected 'email duplicate', got %v", err)
	}
}
