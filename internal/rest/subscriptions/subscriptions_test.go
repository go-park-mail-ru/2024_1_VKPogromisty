package rest

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"socio/domain"
// 	"socio/errors"
// 	mock_subscriptions "socio/mocks/usecase/subscriptions"
// 	"socio/pkg/requestcontext"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// )

// type fields struct {
// 	UserStorage         *mock_subscriptions.MockUserStorage
// 	SubscriptionStorage *mock_subscriptions.MockSubscriptionsStorage
// }

// func TestHandleSubscription(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	tests := []struct {
// 		name           string
// 		input          *SubscriptionInput
// 		expectedStatus int
// 		ctx            context.Context
// 		setupMocks     func(*fields)
// 	}{
// 		{
// 			name: "valid input",
// 			input: &SubscriptionInput{
// 				SubscribedToID: 1,
// 			},
// 			expectedStatus: http.StatusCreated,
// 			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&domain.User{
// 					ID: 1,
// 				}, nil).AnyTimes()
// 				fields.SubscriptionStorage.EXPECT().Store(gomock.Any(), gomock.Any()).Return(&domain.Subscription{
// 					SubscriberID: 1,
// 				}, nil)
// 			},
// 		},
// 		{
// 			name: "invalid context",
// 			input: &SubscriptionInput{
// 				SubscribedToID: 1,
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			ctx:            context.Background(),
// 			setupMocks: func(fields *fields) {
// 			},
// 		},
// 		{
// 			name: "err get user",
// 			input: &SubscriptionInput{
// 				SubscribedToID: 1,
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound).AnyTimes()
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			fields := &fields{
// 				SubscriptionStorage: mock_subscriptions.NewMockSubscriptionsStorage(ctrl),
// 				UserStorage:         mock_subscriptions.NewMockUserStorage(ctrl),
// 			}

// 			tt.setupMocks(fields)

// 			handler := NewSubscriptionsHandler(fields.SubscriptionStorage, fields.UserStorage)

// 			b, err := json.Marshal(tt.input)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			req, err := http.NewRequest("POST", "/subscribe", bytes.NewBuffer(b))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Set("Content-Type", "application/json")

// 			req = req.WithContext(tt.ctx)

// 			rr := httptest.NewRecorder()
// 			handler.HandleSubscription(rr, req)

// 			assert.Equal(t, tt.expectedStatus, rr.Code)
// 		})
// 	}
// }

// func TestHandleUnsubscription(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	tests := []struct {
// 		name           string
// 		input          *SubscriptionInput
// 		expectedStatus int
// 		ctx            context.Context
// 		setupMocks     func(*fields)
// 	}{
// 		{
// 			name: "valid input",
// 			input: &SubscriptionInput{
// 				SubscribedToID: 1,
// 			},
// 			expectedStatus: http.StatusNoContent,
// 			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&domain.User{
// 					ID: 1,
// 				}, nil).AnyTimes()
// 				fields.SubscriptionStorage.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
// 			},
// 		},
// 		{
// 			name: "invalid context",
// 			input: &SubscriptionInput{
// 				SubscribedToID: 1,
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			ctx:            context.Background(),
// 			setupMocks: func(fields *fields) {
// 			},
// 		},
// 		{
// 			name: "err get user",
// 			input: &SubscriptionInput{
// 				SubscribedToID: 1,
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound).AnyTimes()
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			fields := &fields{
// 				SubscriptionStorage: mock_subscriptions.NewMockSubscriptionsStorage(ctrl),
// 				UserStorage:         mock_subscriptions.NewMockUserStorage(ctrl),
// 			}

// 			tt.setupMocks(fields)

// 			handler := NewSubscriptionsHandler(fields.SubscriptionStorage, fields.UserStorage)

// 			b, err := json.Marshal(tt.input)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			req, err := http.NewRequest("POST", "/unsubscribe", bytes.NewBuffer(b))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Set("Content-Type", "application/json")

// 			req = req.WithContext(tt.ctx)

// 			rr := httptest.NewRecorder()
// 			handler.HandleUnsubscription(rr, req)

// 			assert.Equal(t, tt.expectedStatus, rr.Code)
// 		})
// 	}
// }

// func TestHandleGetSubscriptions(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	tests := []struct {
// 		name           string
// 		expectedStatus int
// 		ctx            context.Context
// 		setupMocks     func(*fields)
// 	}{
// 		{
// 			name:           "valid context",
// 			expectedStatus: http.StatusOK,
// 			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&domain.User{}, nil)
// 				fields.SubscriptionStorage.EXPECT().GetSubscriptions(gomock.Any(), gomock.Any()).Return([]*domain.User{}, nil)
// 			},
// 		},
// 		{
// 			name:           "invalid context",
// 			expectedStatus: http.StatusBadRequest,
// 			ctx:            context.Background(),
// 			setupMocks: func(fields *fields) {
// 			},
// 		},
// 		{
// 			name:           "error getting subscriptions",
// 			expectedStatus: http.StatusNotFound,
// 			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			fields := &fields{
// 				SubscriptionStorage: mock_subscriptions.NewMockSubscriptionsStorage(ctrl),
// 				UserStorage:         mock_subscriptions.NewMockUserStorage(ctrl),
// 			}

// 			tt.setupMocks(fields)

// 			handler := NewSubscriptionsHandler(fields.SubscriptionStorage, fields.UserStorage)

// 			req, err := http.NewRequest("GET", "/get-subs", bytes.NewBuffer(nil))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Set("Content-Type", "application/json")

// 			req = req.WithContext(tt.ctx)

// 			rr := httptest.NewRecorder()
// 			handler.HandleGetSubscriptions(rr, req)

// 			assert.Equal(t, tt.expectedStatus, rr.Code)
// 		})
// 	}
// }

// func TestHandleGetSubscribers(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	tests := []struct {
// 		name           string
// 		expectedStatus int
// 		ctx            context.Context
// 		setupMocks     func(*fields)
// 	}{
// 		{
// 			name:           "valid context",
// 			expectedStatus: http.StatusOK,
// 			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&domain.User{}, nil)
// 				fields.SubscriptionStorage.EXPECT().GetSubscribers(gomock.Any(), uint(2)).Return([]*domain.User{}, nil)
// 			},
// 		},
// 		{
// 			name:           "invalid context",
// 			expectedStatus: http.StatusBadRequest,
// 			ctx:            context.Background(),
// 			setupMocks: func(fields *fields) {
// 			},
// 		},
// 		{
// 			name:           "error",
// 			expectedStatus: http.StatusNotFound,
// 			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			fields := &fields{
// 				SubscriptionStorage: mock_subscriptions.NewMockSubscriptionsStorage(ctrl),
// 				UserStorage:         mock_subscriptions.NewMockUserStorage(ctrl),
// 			}

// 			tt.setupMocks(fields)

// 			handler := NewSubscriptionsHandler(fields.SubscriptionStorage, fields.UserStorage)

// 			req, err := http.NewRequest("GET", "/subscribers", nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			req = req.WithContext(tt.ctx)

// 			rr := httptest.NewRecorder()
// 			handler.HandleGetSubscribers(rr, req)

// 			assert.Equal(t, tt.expectedStatus, rr.Code)
// 		})
// 	}
// }

// func TestHandleGetFriends(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	tests := []struct {
// 		name           string
// 		expectedStatus int
// 		ctx            context.Context
// 		setupMocks     func(*fields)
// 	}{
// 		{
// 			name:           "valid context",
// 			expectedStatus: http.StatusOK,
// 			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&domain.User{}, nil)
// 				fields.SubscriptionStorage.EXPECT().GetFriends(gomock.Any(), uint(2)).Return([]*domain.User{}, nil)
// 			},
// 		},
// 		{
// 			name:           "invalid context",
// 			expectedStatus: http.StatusBadRequest,
// 			ctx:            context.Background(),
// 			setupMocks: func(fields *fields) {
// 			},
// 		},
// 		{
// 			name:           "error",
// 			expectedStatus: http.StatusNotFound,
// 			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			fields := &fields{
// 				SubscriptionStorage: mock_subscriptions.NewMockSubscriptionsStorage(ctrl),
// 				UserStorage:         mock_subscriptions.NewMockUserStorage(ctrl),
// 			}

// 			tt.setupMocks(fields)

// 			handler := NewSubscriptionsHandler(fields.SubscriptionStorage, fields.UserStorage)

// 			req, err := http.NewRequest("GET", "/friends", nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			req = req.WithContext(tt.ctx)

// 			rr := httptest.NewRecorder()
// 			handler.HandleGetFriends(rr, req)

// 			assert.Equal(t, tt.expectedStatus, rr.Code)
// 		})
// 	}
// }
