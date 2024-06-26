// Code generated by MockGen. DO NOT EDIT.
// Source: usecase/chat/client.go

// Package mock_chat is a generated GoMock package.
package mock_chat

import (
	context "context"
	reflect "reflect"
	domain "socio/domain"
	chat "socio/usecase/chat"

	gomock "github.com/golang/mock/gomock"
)

// MockPersonalMessagesRepository is a mock of PersonalMessagesRepository interface.
type MockPersonalMessagesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPersonalMessagesRepositoryMockRecorder
}

// MockPersonalMessagesRepositoryMockRecorder is the mock recorder for MockPersonalMessagesRepository.
type MockPersonalMessagesRepositoryMockRecorder struct {
	mock *MockPersonalMessagesRepository
}

// NewMockPersonalMessagesRepository creates a new mock instance.
func NewMockPersonalMessagesRepository(ctrl *gomock.Controller) *MockPersonalMessagesRepository {
	mock := &MockPersonalMessagesRepository{ctrl: ctrl}
	mock.recorder = &MockPersonalMessagesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPersonalMessagesRepository) EXPECT() *MockPersonalMessagesRepositoryMockRecorder {
	return m.recorder
}

// DeleteMessage mocks base method.
func (m *MockPersonalMessagesRepository) DeleteMessage(ctx context.Context, messageID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMessage", ctx, messageID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMessage indicates an expected call of DeleteMessage.
func (mr *MockPersonalMessagesRepositoryMockRecorder) DeleteMessage(ctx, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessage", reflect.TypeOf((*MockPersonalMessagesRepository)(nil).DeleteMessage), ctx, messageID)
}

// DeleteSticker mocks base method.
func (m *MockPersonalMessagesRepository) DeleteSticker(ctx context.Context, stickerID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSticker", ctx, stickerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSticker indicates an expected call of DeleteSticker.
func (mr *MockPersonalMessagesRepositoryMockRecorder) DeleteSticker(ctx, stickerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSticker", reflect.TypeOf((*MockPersonalMessagesRepository)(nil).DeleteSticker), ctx, stickerID)
}

// GetAllStickers mocks base method.
func (m *MockPersonalMessagesRepository) GetAllStickers(ctx context.Context) ([]*domain.Sticker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllStickers", ctx)
	ret0, _ := ret[0].([]*domain.Sticker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllStickers indicates an expected call of GetAllStickers.
func (mr *MockPersonalMessagesRepositoryMockRecorder) GetAllStickers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllStickers", reflect.TypeOf((*MockPersonalMessagesRepository)(nil).GetAllStickers), ctx)
}

// GetDialogsByUserID mocks base method.
func (m *MockPersonalMessagesRepository) GetDialogsByUserID(ctx context.Context, userID uint) ([]*domain.Dialog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDialogsByUserID", ctx, userID)
	ret0, _ := ret[0].([]*domain.Dialog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDialogsByUserID indicates an expected call of GetDialogsByUserID.
func (mr *MockPersonalMessagesRepositoryMockRecorder) GetDialogsByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDialogsByUserID", reflect.TypeOf((*MockPersonalMessagesRepository)(nil).GetDialogsByUserID), ctx, userID)
}

// GetLastMessageID mocks base method.
func (m *MockPersonalMessagesRepository) GetLastMessageID(ctx context.Context, senderID, receiverID uint) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastMessageID", ctx, senderID, receiverID)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastMessageID indicates an expected call of GetLastMessageID.
func (mr *MockPersonalMessagesRepositoryMockRecorder) GetLastMessageID(ctx, senderID, receiverID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastMessageID", reflect.TypeOf((*MockPersonalMessagesRepository)(nil).GetLastMessageID), ctx, senderID, receiverID)
}

// GetMessageByID mocks base method.
func (m *MockPersonalMessagesRepository) GetMessageByID(ctx context.Context, msgID uint) (*domain.PersonalMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageByID", ctx, msgID)
	ret0, _ := ret[0].(*domain.PersonalMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessageByID indicates an expected call of GetMessageByID.
func (mr *MockPersonalMessagesRepositoryMockRecorder) GetMessageByID(ctx, msgID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageByID", reflect.TypeOf((*MockPersonalMessagesRepository)(nil).GetMessageByID), ctx, msgID)
}

// GetMessagesByDialog mocks base method.
func (m *MockPersonalMessagesRepository) GetMessagesByDialog(ctx context.Context, senderID, receiverID, lastMessageID, messagesAmount uint) ([]*domain.PersonalMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessagesByDialog", ctx, senderID, receiverID, lastMessageID, messagesAmount)
	ret0, _ := ret[0].([]*domain.PersonalMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessagesByDialog indicates an expected call of GetMessagesByDialog.
func (mr *MockPersonalMessagesRepositoryMockRecorder) GetMessagesByDialog(ctx, senderID, receiverID, lastMessageID, messagesAmount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessagesByDialog", reflect.TypeOf((*MockPersonalMessagesRepository)(nil).GetMessagesByDialog), ctx, senderID, receiverID, lastMessageID, messagesAmount)
}

// GetStickerByID mocks base method.
func (m *MockPersonalMessagesRepository) GetStickerByID(ctx context.Context, stickerID uint) (*domain.Sticker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStickerByID", ctx, stickerID)
	ret0, _ := ret[0].(*domain.Sticker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStickerByID indicates an expected call of GetStickerByID.
func (mr *MockPersonalMessagesRepositoryMockRecorder) GetStickerByID(ctx, stickerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStickerByID", reflect.TypeOf((*MockPersonalMessagesRepository)(nil).GetStickerByID), ctx, stickerID)
}

// GetStickersByAuthorID mocks base method.
func (m *MockPersonalMessagesRepository) GetStickersByAuthorID(ctx context.Context, authorID uint) ([]*domain.Sticker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStickersByAuthorID", ctx, authorID)
	ret0, _ := ret[0].([]*domain.Sticker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStickersByAuthorID indicates an expected call of GetStickersByAuthorID.
func (mr *MockPersonalMessagesRepositoryMockRecorder) GetStickersByAuthorID(ctx, authorID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStickersByAuthorID", reflect.TypeOf((*MockPersonalMessagesRepository)(nil).GetStickersByAuthorID), ctx, authorID)
}

// StoreMessage mocks base method.
func (m *MockPersonalMessagesRepository) StoreMessage(ctx context.Context, message *domain.PersonalMessage) (*domain.PersonalMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreMessage", ctx, message)
	ret0, _ := ret[0].(*domain.PersonalMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreMessage indicates an expected call of StoreMessage.
func (mr *MockPersonalMessagesRepositoryMockRecorder) StoreMessage(ctx, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreMessage", reflect.TypeOf((*MockPersonalMessagesRepository)(nil).StoreMessage), ctx, message)
}

// StoreSticker mocks base method.
func (m *MockPersonalMessagesRepository) StoreSticker(ctx context.Context, sticker *domain.Sticker) (*domain.Sticker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreSticker", ctx, sticker)
	ret0, _ := ret[0].(*domain.Sticker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreSticker indicates an expected call of StoreSticker.
func (mr *MockPersonalMessagesRepositoryMockRecorder) StoreSticker(ctx, sticker interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreSticker", reflect.TypeOf((*MockPersonalMessagesRepository)(nil).StoreSticker), ctx, sticker)
}

// StoreStickerMessage mocks base method.
func (m *MockPersonalMessagesRepository) StoreStickerMessage(ctx context.Context, senderID, receiverID, stickerID uint) (*domain.PersonalMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreStickerMessage", ctx, senderID, receiverID, stickerID)
	ret0, _ := ret[0].(*domain.PersonalMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreStickerMessage indicates an expected call of StoreStickerMessage.
func (mr *MockPersonalMessagesRepositoryMockRecorder) StoreStickerMessage(ctx, senderID, receiverID, stickerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreStickerMessage", reflect.TypeOf((*MockPersonalMessagesRepository)(nil).StoreStickerMessage), ctx, senderID, receiverID, stickerID)
}

// UpdateMessage mocks base method.
func (m *MockPersonalMessagesRepository) UpdateMessage(ctx context.Context, msg *domain.PersonalMessage, attachmentsToDelete []string) (*domain.PersonalMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMessage", ctx, msg, attachmentsToDelete)
	ret0, _ := ret[0].(*domain.PersonalMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMessage indicates an expected call of UpdateMessage.
func (mr *MockPersonalMessagesRepositoryMockRecorder) UpdateMessage(ctx, msg, attachmentsToDelete interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMessage", reflect.TypeOf((*MockPersonalMessagesRepository)(nil).UpdateMessage), ctx, msg, attachmentsToDelete)
}

// MockPubSubRepository is a mock of PubSubRepository interface.
type MockPubSubRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPubSubRepositoryMockRecorder
}

// MockPubSubRepositoryMockRecorder is the mock recorder for MockPubSubRepository.
type MockPubSubRepositoryMockRecorder struct {
	mock *MockPubSubRepository
}

// NewMockPubSubRepository creates a new mock instance.
func NewMockPubSubRepository(ctrl *gomock.Controller) *MockPubSubRepository {
	mock := &MockPubSubRepository{ctrl: ctrl}
	mock.recorder = &MockPubSubRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPubSubRepository) EXPECT() *MockPubSubRepositoryMockRecorder {
	return m.recorder
}

// ReadActions mocks base method.
func (m *MockPubSubRepository) ReadActions(ctx context.Context, userID uint, ch chan *chat.Action) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadActions", ctx, userID, ch)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadActions indicates an expected call of ReadActions.
func (mr *MockPubSubRepositoryMockRecorder) ReadActions(ctx, userID, ch interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadActions", reflect.TypeOf((*MockPubSubRepository)(nil).ReadActions), ctx, userID, ch)
}

// WriteAction mocks base method.
func (m *MockPubSubRepository) WriteAction(ctx context.Context, action *chat.Action) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteAction", ctx, action)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteAction indicates an expected call of WriteAction.
func (mr *MockPubSubRepositoryMockRecorder) WriteAction(ctx, action interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteAction", reflect.TypeOf((*MockPubSubRepository)(nil).WriteAction), ctx, action)
}
