// Code generated by MockGen. DO NOT EDIT.
// Source: usecase/public_group/public_group.go

// Package mock_publicgroup is a generated GoMock package.
package mock_publicgroup

import (
	context "context"
	reflect "reflect"
	domain "socio/domain"
	publicgroup "socio/usecase/public_group"

	gomock "github.com/golang/mock/gomock"
)

// MockPublicGroupStorage is a mock of PublicGroupStorage interface.
type MockPublicGroupStorage struct {
	ctrl     *gomock.Controller
	recorder *MockPublicGroupStorageMockRecorder
}

// MockPublicGroupStorageMockRecorder is the mock recorder for MockPublicGroupStorage.
type MockPublicGroupStorageMockRecorder struct {
	mock *MockPublicGroupStorage
}

// NewMockPublicGroupStorage creates a new mock instance.
func NewMockPublicGroupStorage(ctrl *gomock.Controller) *MockPublicGroupStorage {
	mock := &MockPublicGroupStorage{ctrl: ctrl}
	mock.recorder = &MockPublicGroupStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPublicGroupStorage) EXPECT() *MockPublicGroupStorageMockRecorder {
	return m.recorder
}

// DeletePublicGroup mocks base method.
func (m *MockPublicGroupStorage) DeletePublicGroup(ctx context.Context, publicGroupID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePublicGroup", ctx, publicGroupID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePublicGroup indicates an expected call of DeletePublicGroup.
func (mr *MockPublicGroupStorageMockRecorder) DeletePublicGroup(ctx, publicGroupID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePublicGroup", reflect.TypeOf((*MockPublicGroupStorage)(nil).DeletePublicGroup), ctx, publicGroupID)
}

// DeletePublicGroupSubscription mocks base method.
func (m *MockPublicGroupStorage) DeletePublicGroupSubscription(ctx context.Context, subscription *domain.PublicGroupSubscription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePublicGroupSubscription", ctx, subscription)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePublicGroupSubscription indicates an expected call of DeletePublicGroupSubscription.
func (mr *MockPublicGroupStorageMockRecorder) DeletePublicGroupSubscription(ctx, subscription interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePublicGroupSubscription", reflect.TypeOf((*MockPublicGroupStorage)(nil).DeletePublicGroupSubscription), ctx, subscription)
}

// GetPublicGroupByID mocks base method.
func (m *MockPublicGroupStorage) GetPublicGroupByID(ctx context.Context, groupID, userID uint) (*publicgroup.PublicGroupWithInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPublicGroupByID", ctx, groupID, userID)
	ret0, _ := ret[0].(*publicgroup.PublicGroupWithInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPublicGroupByID indicates an expected call of GetPublicGroupByID.
func (mr *MockPublicGroupStorageMockRecorder) GetPublicGroupByID(ctx, groupID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPublicGroupByID", reflect.TypeOf((*MockPublicGroupStorage)(nil).GetPublicGroupByID), ctx, groupID, userID)
}

// GetPublicGroupSubscriptionIDs mocks base method.
func (m *MockPublicGroupStorage) GetPublicGroupSubscriptionIDs(ctx context.Context, userID uint) ([]uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPublicGroupSubscriptionIDs", ctx, userID)
	ret0, _ := ret[0].([]uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPublicGroupSubscriptionIDs indicates an expected call of GetPublicGroupSubscriptionIDs.
func (mr *MockPublicGroupStorageMockRecorder) GetPublicGroupSubscriptionIDs(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPublicGroupSubscriptionIDs", reflect.TypeOf((*MockPublicGroupStorage)(nil).GetPublicGroupSubscriptionIDs), ctx, userID)
}

// GetPublicGroupsBySubscriberID mocks base method.
func (m *MockPublicGroupStorage) GetPublicGroupsBySubscriberID(ctx context.Context, subscriberID uint) ([]*domain.PublicGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPublicGroupsBySubscriberID", ctx, subscriberID)
	ret0, _ := ret[0].([]*domain.PublicGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPublicGroupsBySubscriberID indicates an expected call of GetPublicGroupsBySubscriberID.
func (mr *MockPublicGroupStorageMockRecorder) GetPublicGroupsBySubscriberID(ctx, subscriberID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPublicGroupsBySubscriberID", reflect.TypeOf((*MockPublicGroupStorage)(nil).GetPublicGroupsBySubscriberID), ctx, subscriberID)
}

// GetSubscriptionByPublicGroupIDAndSubscriberID mocks base method.
func (m *MockPublicGroupStorage) GetSubscriptionByPublicGroupIDAndSubscriberID(ctx context.Context, publicGroupID, subscriberID uint) (*domain.PublicGroupSubscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscriptionByPublicGroupIDAndSubscriberID", ctx, publicGroupID, subscriberID)
	ret0, _ := ret[0].(*domain.PublicGroupSubscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscriptionByPublicGroupIDAndSubscriberID indicates an expected call of GetSubscriptionByPublicGroupIDAndSubscriberID.
func (mr *MockPublicGroupStorageMockRecorder) GetSubscriptionByPublicGroupIDAndSubscriberID(ctx, publicGroupID, subscriberID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptionByPublicGroupIDAndSubscriberID", reflect.TypeOf((*MockPublicGroupStorage)(nil).GetSubscriptionByPublicGroupIDAndSubscriberID), ctx, publicGroupID, subscriberID)
}

// SearchPublicGroupsByNameWithInfo mocks base method.
func (m *MockPublicGroupStorage) SearchPublicGroupsByNameWithInfo(ctx context.Context, query string, userID uint) ([]*publicgroup.PublicGroupWithInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchPublicGroupsByNameWithInfo", ctx, query, userID)
	ret0, _ := ret[0].([]*publicgroup.PublicGroupWithInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchPublicGroupsByNameWithInfo indicates an expected call of SearchPublicGroupsByNameWithInfo.
func (mr *MockPublicGroupStorageMockRecorder) SearchPublicGroupsByNameWithInfo(ctx, query, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchPublicGroupsByNameWithInfo", reflect.TypeOf((*MockPublicGroupStorage)(nil).SearchPublicGroupsByNameWithInfo), ctx, query, userID)
}

// StorePublicGroup mocks base method.
func (m *MockPublicGroupStorage) StorePublicGroup(ctx context.Context, publicGroup *domain.PublicGroup) (*domain.PublicGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorePublicGroup", ctx, publicGroup)
	ret0, _ := ret[0].(*domain.PublicGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StorePublicGroup indicates an expected call of StorePublicGroup.
func (mr *MockPublicGroupStorageMockRecorder) StorePublicGroup(ctx, publicGroup interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorePublicGroup", reflect.TypeOf((*MockPublicGroupStorage)(nil).StorePublicGroup), ctx, publicGroup)
}

// StorePublicGroupSubscription mocks base method.
func (m *MockPublicGroupStorage) StorePublicGroupSubscription(ctx context.Context, publicGroupSubscription *domain.PublicGroupSubscription) (*domain.PublicGroupSubscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorePublicGroupSubscription", ctx, publicGroupSubscription)
	ret0, _ := ret[0].(*domain.PublicGroupSubscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StorePublicGroupSubscription indicates an expected call of StorePublicGroupSubscription.
func (mr *MockPublicGroupStorageMockRecorder) StorePublicGroupSubscription(ctx, publicGroupSubscription interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorePublicGroupSubscription", reflect.TypeOf((*MockPublicGroupStorage)(nil).StorePublicGroupSubscription), ctx, publicGroupSubscription)
}

// UpdatePublicGroup mocks base method.
func (m *MockPublicGroupStorage) UpdatePublicGroup(ctx context.Context, publicGroup *domain.PublicGroup) (*domain.PublicGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePublicGroup", ctx, publicGroup)
	ret0, _ := ret[0].(*domain.PublicGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePublicGroup indicates an expected call of UpdatePublicGroup.
func (mr *MockPublicGroupStorageMockRecorder) UpdatePublicGroup(ctx, publicGroup interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePublicGroup", reflect.TypeOf((*MockPublicGroupStorage)(nil).UpdatePublicGroup), ctx, publicGroup)
}

// MockAvatarStorage is a mock of AvatarStorage interface.
type MockAvatarStorage struct {
	ctrl     *gomock.Controller
	recorder *MockAvatarStorageMockRecorder
}

// MockAvatarStorageMockRecorder is the mock recorder for MockAvatarStorage.
type MockAvatarStorageMockRecorder struct {
	mock *MockAvatarStorage
}

// NewMockAvatarStorage creates a new mock instance.
func NewMockAvatarStorage(ctrl *gomock.Controller) *MockAvatarStorage {
	mock := &MockAvatarStorage{ctrl: ctrl}
	mock.recorder = &MockAvatarStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAvatarStorage) EXPECT() *MockAvatarStorageMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockAvatarStorage) Delete(fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAvatarStorageMockRecorder) Delete(fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAvatarStorage)(nil).Delete), fileName)
}

// Store mocks base method.
func (m *MockAvatarStorage) Store(fileName, filePath, contentType string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", fileName, filePath, contentType)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockAvatarStorageMockRecorder) Store(fileName, filePath, contentType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockAvatarStorage)(nil).Store), fileName, filePath, contentType)
}