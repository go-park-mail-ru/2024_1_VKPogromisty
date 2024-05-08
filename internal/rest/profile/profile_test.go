package rest

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"socio/errors"
	uspb "socio/internal/grpc/user/proto"
	mock_user "socio/mocks/grpc/user_grpc"
	"socio/pkg/requestcontext"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		userID         string
		mockError      error
		expectedStatus int
		mock           func(userClient *mock_user.MockUserClient)
	}{
		{
			name:           "Successful get profile",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         "1",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().GetByIDWithSubsInfo(gomock.Any(), gomock.Any()).Return(&uspb.GetByIDWithSubsInfoResponse{
					User: &uspb.UserResponse{
						Id: 1,
					},
				}, nil)
			},
		},
		{
			name:           "Successful get profile",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         "",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().GetByIDWithSubsInfo(gomock.Any(), gomock.Any()).Return(&uspb.GetByIDWithSubsInfoResponse{
					User: &uspb.UserResponse{
						Id: 1,
					},
				}, nil)
			},
		},
		{
			name:           "Successful get profile",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         "asd",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {

			},
		},
		{
			name:           "Successful get profile",
			ctx:            context.Background(),
			userID:         "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {

			},
		},
		{
			name:           "Successful get profile",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         "",
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().GetByIDWithSubsInfo(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal.GRPCStatus().Err())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("GET", "/"+tt.userID, nil)
			r = r.WithContext(tt.ctx)
			r = mux.SetURLVars(r, map[string]string{
				"userID": tt.userID,
			})

			// Set up the response recorder
			rr := httptest.NewRecorder()

			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockUserClient)

			// Set up the handler
			h := NewProfileHandler(mockUserClient)

			// Call the handler
			h.HandleGetProfile(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleUpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		body           string
		mockError      error
		expectedStatus int
		mock           func(userClient *mock_user.MockUserClient)
	}{
		{
			name:           "Successful update profile",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			mockError:      nil,
			expectedStatus: http.StatusOK,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().Update(gomock.Any(), gomock.Any()).Return(
					&uspb.UpdateResponse{
						User: &uspb.UserResponse{},
					}, nil,
				)
			},
		},
		{
			name:           "invalid user id",
			ctx:            context.Background(),
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {

			},
		},
		{
			name:           "Successful update profile",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().Update(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInvalidData.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			_ = writer.WriteField("firstName", "John")
			_ = writer.WriteField("lastName", "Doe")
			_ = writer.WriteField("email", "john.doe@example.com")
			_ = writer.WriteField("password", "123456")
			_ = writer.WriteField("repeatPassword", "123456")
			_ = writer.WriteField("dateOfBirth", "2000-01-01")
			_ = writer.Close()

			r := httptest.NewRequest("POST", "/", body)
			r.Header.Set("Content-Type", writer.FormDataContentType())
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockUserClient)

			// Set up the handler
			h := NewProfileHandler(mockUserClient)

			// Call the handler
			h.HandleUpdateProfile(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleDeleteProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		mockError      error
		expectedStatus int
		mock           func(userClient *mock_user.MockUserClient)
	}{
		{
			name:           "Successful delete profile",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name:           "invalid user id",
			ctx:            context.Background(),
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {

			},
		},
		{
			name:           "internal",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("DELETE", "/", nil)
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockUserClient)

			// Set up the handler
			h := NewProfileHandler(mockUserClient)

			// Call the handler
			h.HandleDeleteProfile(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleSearchByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		query          string
		mockError      error
		expectedStatus int
		mock           func(userClient *mock_user.MockUserClient)
	}{
		{
			name:           "Successful search by name",
			query:          "John",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().SearchByName(gomock.Any(), gomock.Any()).Return(&uspb.SearchByNameResponse{
					Users: []*uspb.UserResponse{},
				}, nil)
			},
		},
		{
			name:           "Empty query",
			query:          "",
			mockError:      errors.ErrInvalidData,
			expectedStatus: http.StatusBadRequest,
			mock:           func(userClient *mock_user.MockUserClient) {},
		},
		{
			name:           "err",
			query:          "John",
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().SearchByName(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal,
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("GET", "/?query="+tt.query, nil)
			r = mux.SetURLVars(r, map[string]string{
				"query": tt.query,
			})

			// Set up the response recorder
			rr := httptest.NewRecorder()

			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockUserClient)

			// Set up the handler
			h := NewProfileHandler(mockUserClient)

			// Call the handler
			h.HandleSearchByName(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
