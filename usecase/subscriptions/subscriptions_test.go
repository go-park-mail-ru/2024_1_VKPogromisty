package subscriptions_test

import (
	"context"
	"reflect"
	"socio/domain"
	"socio/errors"
	mock_subscriptions "socio/mocks/usecase/subscriptions"
	customtime "socio/pkg/time"
	"socio/usecase/subscriptions"
	"testing"

	"github.com/golang/mock/gomock"
)

type fields struct {
	SubscriptionsStorage *mock_subscriptions.MockSubscriptionsStorage
	UserStorage          *mock_subscriptions.MockUserStorage
}

var timeProv = customtime.MockTimeProvider{}

func TestService_Subscribe(t *testing.T) {
	type args struct {
		ctx context.Context
		sub *domain.Subscription
	}
	tests := []struct {
		name             string
		args             args
		wantSubscription *domain.Subscription
		wantErr          bool
		prepareMock      func(*fields)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				sub: &domain.Subscription{
					SubscriberID:   1,
					SubscribedToID: 2,
				},
			},
			wantSubscription: &domain.Subscription{
				ID:             1,
				SubscriberID:   1,
				SubscribedToID: 2,
				CreatedAt: customtime.CustomTime{
					Time: timeProv.Now(),
				},
				UpdatedAt: customtime.CustomTime{
					Time: timeProv.Now(),
				},
			},
			wantErr: false,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&domain.User{}, nil)
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(2)).Return(&domain.User{}, nil)
				f.SubscriptionsStorage.EXPECT().Store(gomock.Any(), gomock.Any()).Return(&domain.Subscription{
					ID:             1,
					SubscriberID:   1,
					SubscribedToID: 2,
					CreatedAt: customtime.CustomTime{
						Time: timeProv.Now(),
					},
					UpdatedAt: customtime.CustomTime{
						Time: timeProv.Now(),
					},
				}, nil)
			},
		},
		{
			name: "error equal ids",
			args: args{
				ctx: context.Background(),
				sub: &domain.Subscription{
					SubscriberID:   1,
					SubscribedToID: 1,
				},
			},
			wantSubscription: nil,
			wantErr:          true,
			prepareMock: func(f *fields) {
			},
		},
		{
			name: "error first user",
			args: args{
				ctx: context.Background(),
				sub: &domain.Subscription{
					SubscriberID:   1,
					SubscribedToID: 2,
				},
			},
			wantSubscription: nil,
			wantErr:          true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name: "error second user",
			args: args{
				ctx: context.Background(),
				sub: &domain.Subscription{
					SubscriberID:   1,
					SubscribedToID: 2,
				},
			},
			wantSubscription: nil,
			wantErr:          true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&domain.User{}, nil)
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(2)).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name: "error store",
			args: args{
				ctx: context.Background(),
				sub: &domain.Subscription{
					SubscriberID:   1,
					SubscribedToID: 2,
				},
			},
			wantSubscription: nil,
			wantErr:          true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&domain.User{}, nil)
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(2)).Return(&domain.User{}, nil)
				f.SubscriptionsStorage.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				SubscriptionsStorage: mock_subscriptions.NewMockSubscriptionsStorage(ctrl),
				UserStorage:          mock_subscriptions.NewMockUserStorage(ctrl),
			}

			if tt.prepareMock != nil {
				tt.prepareMock(&f)
			}

			s := subscriptions.NewService(f.SubscriptionsStorage, f.UserStorage)

			gotSubscription, err := s.Subscribe(tt.args.ctx, tt.args.sub)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Subscribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSubscription, tt.wantSubscription) {
				t.Errorf("Service.Subscribe() = %v, want %v", gotSubscription, tt.wantSubscription)
			}
		})
	}
}

func TestService_Unsubscribe(t *testing.T) {
	type args struct {
		ctx context.Context
		sub *domain.Subscription
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		prepareMock func(*fields)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				sub: &domain.Subscription{
					SubscriberID:   1,
					SubscribedToID: 2,
				},
			},
			wantErr: false,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&domain.User{}, nil)
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(2)).Return(&domain.User{}, nil)
				f.SubscriptionsStorage.EXPECT().Delete(gomock.Any(), uint(1), uint(2)).Return(nil)
			},
		},
		{
			name: "error first user",
			args: args{
				ctx: context.Background(),
				sub: &domain.Subscription{
					SubscriberID:   1,
					SubscribedToID: 2,
				},
			},
			wantErr: true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name: "error second user",
			args: args{
				ctx: context.Background(),
				sub: &domain.Subscription{
					SubscriberID:   1,
					SubscribedToID: 2,
				},
			},
			wantErr: true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&domain.User{}, nil)
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(2)).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name: "error delete",
			args: args{
				ctx: context.Background(),
				sub: &domain.Subscription{
					SubscriberID:   1,
					SubscribedToID: 2,
				},
			},
			wantErr: true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&domain.User{}, nil)
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(2)).Return(&domain.User{}, nil)
				f.SubscriptionsStorage.EXPECT().Delete(gomock.Any(), uint(1), uint(2)).Return(errors.ErrInternal)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				SubscriptionsStorage: mock_subscriptions.NewMockSubscriptionsStorage(ctrl),
				UserStorage:          mock_subscriptions.NewMockUserStorage(ctrl),
			}

			if tt.prepareMock != nil {
				tt.prepareMock(&f)
			}

			s := subscriptions.NewService(f.SubscriptionsStorage, f.UserStorage)

			if err := s.Unsubscribe(tt.args.ctx, tt.args.sub); (err != nil) != tt.wantErr {
				t.Errorf("Service.Unsubscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetSubscriptions(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID uint
	}

	tests := []struct {
		name              string
		args              args
		wantSubscriptions []*domain.User
		wantErr           bool
		prepareMock       func(*fields)
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			wantSubscriptions: []*domain.User{
				{
					ID: 2,
				},
			},
			wantErr: false,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&domain.User{}, nil)
				f.SubscriptionsStorage.EXPECT().GetSubscriptions(gomock.Any(), uint(1)).Return([]*domain.User{
					{
						ID: 2,
					},
				}, nil)
			},
		},
		{
			name: "error user",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			wantSubscriptions: nil,
			wantErr:           true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name: "error subscriptions",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			wantSubscriptions: nil,
			wantErr:           true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&domain.User{}, nil)
				f.SubscriptionsStorage.EXPECT().GetSubscriptions(gomock.Any(), uint(1)).Return(nil, errors.ErrInternal)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				SubscriptionsStorage: mock_subscriptions.NewMockSubscriptionsStorage(ctrl),
				UserStorage:          mock_subscriptions.NewMockUserStorage(ctrl),
			}

			if tt.prepareMock != nil {
				tt.prepareMock(&f)
			}

			s := subscriptions.NewService(f.SubscriptionsStorage, f.UserStorage)

			gotSubscriptions, err := s.GetSubscriptions(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSubscriptions, tt.wantSubscriptions) {
				t.Errorf("Service.GetSubscriptions() = %v, want %v", gotSubscriptions, tt.wantSubscriptions)
			}
		})
	}
}

func TestService_GetSubscribers(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID uint
	}
	tests := []struct {
		name            string
		args            args
		wantSubscribers []*domain.User
		wantErr         bool
		prepareMock     func(*fields)
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			wantSubscribers: []*domain.User{
				{
					ID: 2,
				},
			},
			wantErr: false,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&domain.User{}, nil)
				f.SubscriptionsStorage.EXPECT().GetSubscribers(gomock.Any(), uint(1)).Return([]*domain.User{
					{
						ID: 2,
					},
				}, nil)
			},
		},
		{
			name: "error user",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			wantSubscribers: nil,
			wantErr:         true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name: "error subscribers",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			wantSubscribers: nil,
			wantErr:         true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&domain.User{}, nil)
				f.SubscriptionsStorage.EXPECT().GetSubscribers(gomock.Any(), uint(1)).Return(nil, errors.ErrInternal)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				SubscriptionsStorage: mock_subscriptions.NewMockSubscriptionsStorage(ctrl),
				UserStorage:          mock_subscriptions.NewMockUserStorage(ctrl),
			}

			if tt.prepareMock != nil {
				tt.prepareMock(&f)
			}

			s := subscriptions.NewService(f.SubscriptionsStorage, f.UserStorage)

			gotSubscribers, err := s.GetSubscribers(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetSubscribers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSubscribers, tt.wantSubscribers) {
				t.Errorf("Service.GetSubscribers() = %v, want %v", gotSubscribers, tt.wantSubscribers)
			}
		})
	}
}

func TestService_GetFriends(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID uint
	}

	tests := []struct {
		name        string
		args        args
		wantFriends []*domain.User
		wantErr     bool
		prepareMock func(*fields)
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			wantFriends: []*domain.User{
				{
					ID: 2,
				},
			},
			wantErr: false,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&domain.User{}, nil)
				f.SubscriptionsStorage.EXPECT().GetFriends(gomock.Any(), uint(1)).Return([]*domain.User{
					{
						ID: 2,
					},
				}, nil)
			},
		},
		{
			name: "error user",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			wantFriends: nil,
			wantErr:     true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name: "error friends",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			wantFriends: nil,
			wantErr:     true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&domain.User{}, nil)
				f.SubscriptionsStorage.EXPECT().GetFriends(gomock.Any(), uint(1)).Return(nil, errors.ErrInternal)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				SubscriptionsStorage: mock_subscriptions.NewMockSubscriptionsStorage(ctrl),
				UserStorage:          mock_subscriptions.NewMockUserStorage(ctrl),
			}

			if tt.prepareMock != nil {
				tt.prepareMock(&f)
			}

			s := subscriptions.NewService(f.SubscriptionsStorage, f.UserStorage)

			gotFriends, err := s.GetFriends(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetFriends() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFriends, tt.wantFriends) {
				t.Errorf("Service.GetFriends() = %v, want %v", gotFriends, tt.wantFriends)
			}
		})
	}
}
