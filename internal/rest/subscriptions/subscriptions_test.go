package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"socio/domain"
	mock_subscriptions "socio/mocks/usecase/subscriptions"
	"socio/pkg/requestcontext"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type fields struct {
	UserStorage         *mock_subscriptions.MockUserStorage
	SubscriptionStorage *mock_subscriptions.MockSubscriptionsStorage
}

func TestHandleSubscription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		input          *SubscriptionInput
		expectedStatus int
		ctx            context.Context
		setupMocks     func(*fields)
	}{
		{
			name: "valid input",
			input: &SubscriptionInput{
				SubscribedToID: 1,
			},
			expectedStatus: http.StatusOK,
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
			setupMocks: func(fields *fields) {
				fields.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&domain.User{
					ID: 1,
				}, nil).AnyTimes()
				fields.SubscriptionStorage.EXPECT().Store(gomock.Any(), gomock.Any()).Return(&domain.Subscription{
					SubscriberID: 1,
				}, nil)
			},
		},
		{
			name: "valid input",
			input: &SubscriptionInput{
				SubscribedToID: 1,
			},
			expectedStatus: http.StatusOK,
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
			setupMocks: func(fields *fields) {
				fields.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&domain.User{
					ID: 1,
				}, nil).AnyTimes()
				fields.SubscriptionStorage.EXPECT().Store(gomock.Any(), gomock.Any()).Return(&domain.Subscription{
					SubscriberID: 1,
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := &fields{
				SubscriptionStorage: mock_subscriptions.NewMockSubscriptionsStorage(ctrl),
				UserStorage:         mock_subscriptions.NewMockUserStorage(ctrl),
			}

			tt.setupMocks(fields)

			handler := NewSubscriptionsHandler(fields.SubscriptionStorage, fields.UserStorage)

			b, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/subscribe", bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			req = req.WithContext(tt.ctx)

			rr := httptest.NewRecorder()
			handler.HandleSubscription(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
