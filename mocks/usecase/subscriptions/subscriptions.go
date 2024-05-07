// Code generated by MockGen. DO NOT EDIT.
// Source: usecase/subscriptions/subscriptions.go

// Package mock_subscriptions is a generated GoMock package.
package mock_subscriptions

import (
	context "context"
	reflect "reflect"
	domain "socio/domain"

	gomock "github.com/golang/mock/gomock"
)

// MockSubscriptionsStorage is a mock of SubscriptionsStorage interface.
type MockSubscriptionsStorage struct {
	ctrl     *gomock.Controller
	recorder *MockSubscriptionsStorageMockRecorder
}

// MockSubscriptionsStorageMockRecorder is the mock recorder for MockSubscriptionsStorage.
type MockSubscriptionsStorageMockRecorder struct {
	mock *MockSubscriptionsStorage
}

// NewMockSubscriptionsStorage creates a new mock instance.
func NewMockSubscriptionsStorage(ctrl *gomock.Controller) *MockSubscriptionsStorage {
	mock := &MockSubscriptionsStorage{ctrl: ctrl}
	mock.recorder = &MockSubscriptionsStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubscriptionsStorage) EXPECT() *MockSubscriptionsStorageMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockSubscriptionsStorage) Delete(ctx context.Context, subscriberID, subscibedToID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, subscriberID, subscibedToID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSubscriptionsStorageMockRecorder) Delete(ctx, subscriberID, subscibedToID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSubscriptionsStorage)(nil).Delete), ctx, subscriberID, subscibedToID)
}

// GetBySubscriberAndSubscribedToID mocks base method.
func (m *MockSubscriptionsStorage) GetBySubscriberAndSubscribedToID(ctx context.Context, subscriberID, subscribedToID uint) (*domain.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBySubscriberAndSubscribedToID", ctx, subscriberID, subscribedToID)
	ret0, _ := ret[0].(*domain.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBySubscriberAndSubscribedToID indicates an expected call of GetBySubscriberAndSubscribedToID.
func (mr *MockSubscriptionsStorageMockRecorder) GetBySubscriberAndSubscribedToID(ctx, subscriberID, subscribedToID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBySubscriberAndSubscribedToID", reflect.TypeOf((*MockSubscriptionsStorage)(nil).GetBySubscriberAndSubscribedToID), ctx, subscriberID, subscribedToID)
}

// GetFriends mocks base method.
func (m *MockSubscriptionsStorage) GetFriends(ctx context.Context, userID uint) ([]*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFriends", ctx, userID)
	ret0, _ := ret[0].([]*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFriends indicates an expected call of GetFriends.
func (mr *MockSubscriptionsStorageMockRecorder) GetFriends(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFriends", reflect.TypeOf((*MockSubscriptionsStorage)(nil).GetFriends), ctx, userID)
}

// GetSubscribers mocks base method.
func (m *MockSubscriptionsStorage) GetSubscribers(ctx context.Context, userID uint) ([]*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscribers", ctx, userID)
	ret0, _ := ret[0].([]*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscribers indicates an expected call of GetSubscribers.
func (mr *MockSubscriptionsStorageMockRecorder) GetSubscribers(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscribers", reflect.TypeOf((*MockSubscriptionsStorage)(nil).GetSubscribers), ctx, userID)
}

// GetSubscriptions mocks base method.
func (m *MockSubscriptionsStorage) GetSubscriptions(ctx context.Context, userID uint) ([]*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscriptions", ctx, userID)
	ret0, _ := ret[0].([]*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscriptions indicates an expected call of GetSubscriptions.
func (mr *MockSubscriptionsStorageMockRecorder) GetSubscriptions(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptions", reflect.TypeOf((*MockSubscriptionsStorage)(nil).GetSubscriptions), ctx, userID)
}

// Store mocks base method.
func (m *MockSubscriptionsStorage) Store(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", ctx, sub)
	ret0, _ := ret[0].(*domain.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Store indicates an expected call of Store.
func (mr *MockSubscriptionsStorageMockRecorder) Store(ctx, sub interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockSubscriptionsStorage)(nil).Store), ctx, sub)
}

// MockUserStorage is a mock of UserStorage interface.
type MockUserStorage struct {
	ctrl     *gomock.Controller
	recorder *MockUserStorageMockRecorder
}

// MockUserStorageMockRecorder is the mock recorder for MockUserStorage.
type MockUserStorageMockRecorder struct {
	mock *MockUserStorage
}

// NewMockUserStorage creates a new mock instance.
func NewMockUserStorage(ctrl *gomock.Controller) *MockUserStorage {
	mock := &MockUserStorage{ctrl: ctrl}
	mock.recorder = &MockUserStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserStorage) EXPECT() *MockUserStorageMockRecorder {
	return m.recorder
}

// GetUserByID mocks base method.
func (m *MockUserStorage) GetUserByID(ctx context.Context, userID uint) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, userID)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserStorageMockRecorder) GetUserByID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserStorage)(nil).GetUserByID), ctx, userID)
}