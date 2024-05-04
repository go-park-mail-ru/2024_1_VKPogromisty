package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"socio/errors"
	"socio/pkg/requestcontext"
	"testing"

	uspb "socio/internal/grpc/user/proto"
	mock_user "socio/mocks/grpc/user_grpc"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandleSubscription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		body           *SubscriptionInput
		mockError      error
		expectedStatus int
		mock           func(userClient *mock_user.MockUserClient)
	}{
		{
			name:           "Successful subscription",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			body:           &SubscriptionInput{SubscribedToID: 2},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().Subscribe(gomock.Any(), gomock.Any()).Return(&uspb.SubscribeResponse{
					Subscription: &uspb.SubscriptionResponse{
						SubscriberId:   1,
						SubscribedToId: 2,
					},
				}, nil)
			},
		},
		{
			name:           "Successful subscription",
			ctx:            context.Background(),
			body:           &SubscriptionInput{SubscribedToID: 2},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {

			},
		},
		{
			name:           "Successful subscription",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			body:           &SubscriptionInput{SubscribedToID: 2},
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().Subscribe(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			body, _ := json.Marshal(tt.body)
			r := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
			r.Header.Set("Content-Type", "application/json")
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockUserClient)

			// Set up the handler
			h := NewSubscriptionsHandler(mockUserClient)

			// Call the handler
			h.HandleSubscription(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleUnsubscription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		body           *SubscriptionInput
		mockError      error
		expectedStatus int
		mock           func(userClient *mock_user.MockUserClient)
	}{
		{
			name:           "Successful unsubscription",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			body:           &SubscriptionInput{SubscribedToID: 2},
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().Unsubscribe(gomock.Any(), gomock.Any()).Return(&uspb.UnsubscribeResponse{}, nil)
			},
		},
		{
			name:           "Successful unsubscription",
			ctx:            context.Background(),
			body:           &SubscriptionInput{SubscribedToID: 2},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {
			},
		},
		{
			name:           "Successful unsubscription",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			body:           &SubscriptionInput{SubscribedToID: 2},
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().Unsubscribe(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			body, _ := json.Marshal(tt.body)
			r := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
			r.Header.Set("Content-Type", "application/json")
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockUserClient)

			// Set up the handler
			h := NewSubscriptionsHandler(mockUserClient)

			// Call the handler
			h.HandleUnsubscription(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleGetSubscriptions(t *testing.T) {
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
			name:           "Successful get subscriptions",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			mockError:      nil,
			expectedStatus: http.StatusOK,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().GetSubscriptions(gomock.Any(), gomock.Any()).Return(&uspb.GetSubscriptionsResponse{}, nil)
			},
		},
		{
			name:           "Successful get subscriptions",
			ctx:            context.Background(),
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {
			},
		},
		{
			name:           "Successful get subscriptions",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().GetSubscriptions(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("GET", "/", nil)
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockUserClient)

			// Set up the handler
			h := NewSubscriptionsHandler(mockUserClient)

			// Call the handler
			h.HandleGetSubscriptions(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleGetSubscribers(t *testing.T) {
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
			name:           "Successful get subscribers",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			mockError:      nil,
			expectedStatus: http.StatusOK,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().GetSubscribers(gomock.Any(), gomock.Any()).Return(&uspb.GetSubscribersResponse{}, nil)
			},
		},
		{
			name:           "Successful get subscribers",
			ctx:            context.Background(),
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {
			},
		},
		{
			name:           "Successful get subscribers",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().GetSubscribers(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("GET", "/", nil)
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockUserClient)

			// Set up the handler
			h := NewSubscriptionsHandler(mockUserClient)

			// Call the handler
			h.HandleGetSubscribers(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleGetFriends(t *testing.T) {
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
			name:           "Successful get friends",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			mockError:      nil,
			expectedStatus: http.StatusOK,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().GetFriends(gomock.Any(), gomock.Any()).Return(&uspb.GetFriendsResponse{}, nil)
			},
		},
		{
			name:           "Successful get friends",
			ctx:            context.Background(),
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {
			},
		},
		{
			name:           "Successful get friends",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().GetFriends(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("GET", "/", nil)
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockUserClient)

			// Set up the handler
			h := NewSubscriptionsHandler(mockUserClient)

			// Call the handler
			h.HandleGetFriends(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
