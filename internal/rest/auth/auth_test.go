package rest_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"socio/errors"
	authpb "socio/internal/grpc/auth/proto"
	userpb "socio/internal/grpc/user/proto"
	rest "socio/internal/rest/auth"
	auth_mocks "socio/mocks/grpc/auth_grpc"
	user_mocks "socio/mocks/grpc/user_grpc"
	auth "socio/usecase/auth"
	"socio/usecase/user"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandleLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := auth_mocks.NewMockAuthClient(ctrl)

	tests := []struct {
		name           string
		input          interface{}
		mockResponse   interface{}
		mockError      error
		expectedStatus int
	}{
		{
			name: "Successful login",
			input: &auth.LoginInput{
				Email:    "test@example.com",
				Password: "password",
			},
			mockResponse: &authpb.LoginResponse{
				SessionId: "some_session_id",
				User:      &authpb.UserResponse{ /* fill user details */ },
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "no body",
			input:          &auth.LoginInput{},
			mockResponse:   &authpb.LoginResponse{},
			mockError:      errors.ErrInvalidLoginData.GRPCStatus().Err(),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request body
			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			// Set up the mock expectations
			mockAuthClient.EXPECT().Login(gomock.Any(), gomock.Any()).Return(tt.mockResponse, tt.mockError)

			// Create the handler with the mock AuthClient
			handler := rest.NewAuthHandler(mockAuthClient, nil, nil)

			// Call the handler function
			handler.HandleLogin(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := auth_mocks.NewMockAuthClient(ctrl)

	tests := []struct {
		name           string
		cookieName     string
		sessionID      string
		mockError      error
		expectedStatus int
		mock           func(authClient *auth_mocks.MockAuthClient, sessionID string, err error)
	}{
		{
			name:           "Successful logout",
			cookieName:     "session_id",
			sessionID:      "some_session_id",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			mock: func(authClient *auth_mocks.MockAuthClient, sessionID string, err error) {
				authClient.EXPECT().Logout(gomock.Any(), &authpb.LogoutRequest{SessionId: sessionID}).Return(
					&authpb.LogoutResponse{}, err,
				)
			},
		},
		{
			name:           "no session",
			cookieName:     "opa",
			sessionID:      "some_session_id",
			mockError:      nil,
			expectedStatus: http.StatusUnauthorized,
			mock:           func(authClient *auth_mocks.MockAuthClient, sessionID string, err error) {},
		},
		{
			name:           "err",
			cookieName:     "session_id",
			sessionID:      "some_session_id",
			mockError:      errors.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
			mock: func(authClient *auth_mocks.MockAuthClient, sessionID string, err error) {
				authClient.EXPECT().Logout(gomock.Any(), &authpb.LogoutRequest{SessionId: sessionID}).Return(
					&authpb.LogoutResponse{}, err,
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			req, _ := http.NewRequest("POST", "/auth/logout", nil)
			req.AddCookie(&http.Cookie{Name: tt.cookieName, Value: tt.sessionID})
			rr := httptest.NewRecorder()

			tt.mock(mockAuthClient, tt.sessionID, tt.mockError)

			// Create the handler with the mock AuthClient
			// Create the handler with the mock AuthClient
			handler := rest.NewAuthHandler(mockAuthClient, nil, nil)

			// Call the handler function
			handler.HandleLogout(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleRegistration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := auth_mocks.NewMockAuthClient(ctrl)
	mockUserClient := user_mocks.NewMockUserClient(ctrl)

	tests := []struct {
		name           string
		input          *user.CreateUserInput
		mockCreateUser *userpb.CreateResponse
		mockLogin      *authpb.LoginResponse
		mockError      error
		expectedStatus int
		mock           func(userClient *user_mocks.MockUserClient, authClient *auth_mocks.MockAuthClient, input *user.CreateUserInput, err error)
	}{
		{
			name: "Successful registration",
			input: &user.CreateUserInput{
				FirstName:      "Test",
				LastName:       "User",
				Email:          "test@example.com",
				Password:       "password",
				RepeatPassword: "password",
				DateOfBirth:    "2000-01-01",
			},
			mockCreateUser: &userpb.CreateResponse{
				User: &userpb.UserResponse{
					Avatar: "default_avatar.png",
				},
			},
			mockLogin: &authpb.LoginResponse{
				SessionId: "some_session_id",
				User:      &authpb.UserResponse{ 
					
				 },
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			mock: func(userClient *user_mocks.MockUserClient, authClient *auth_mocks.MockAuthClient, input *user.CreateUserInput, err error) {
				userClient.EXPECT().Create(gomock.Any(), userpb.ToCreateRequest(input)).Return(&userpb.CreateResponse{
					User: &userpb.UserResponse{Avatar: "default_avatar.png"},
				}, err)
				authClient.EXPECT().Login(gomock.Any(), &authpb.LoginRequest{
					Email:    input.Email,
					Password: input.Password,
				}).Return(&authpb.LoginResponse{}, err)
			},
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request body
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			_ = writer.WriteField("firstName", tt.input.FirstName)
			_ = writer.WriteField("lastName", tt.input.LastName)
			_ = writer.WriteField("email", tt.input.Email)
			_ = writer.WriteField("password", tt.input.Password)
			_ = writer.WriteField("repeatPassword", tt.input.RepeatPassword)
			_ = writer.WriteField("dateOfBirth", tt.input.DateOfBirth)

			// Close the multipart writer
			_ = writer.Close()

			req, _ := http.NewRequest("POST", "/auth/register", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			rr := httptest.NewRecorder()

			// Set up the mock expectations
			tt.mock(mockUserClient, mockAuthClient, tt.input, tt.mockError)
			// Create the handler with the mock AuthClient and UserClient
			handler := rest.NewAuthHandler(mockAuthClient, mockUserClient, nil)

			// Call the handler function
			handler.HandleRegistration(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
