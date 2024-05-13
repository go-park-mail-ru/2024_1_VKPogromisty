// Code generated by MockGen. DO NOT EDIT.
// Source: usecase/posts/posts.go

// Package mock_posts is a generated GoMock package.
package mock_posts

import (
	context "context"
	reflect "reflect"
	domain "socio/domain"
	posts "socio/usecase/posts"

	gomock "github.com/golang/mock/gomock"
)

// MockPostsStorage is a mock of PostsStorage interface.
type MockPostsStorage struct {
	ctrl     *gomock.Controller
	recorder *MockPostsStorageMockRecorder
}

// MockPostsStorageMockRecorder is the mock recorder for MockPostsStorage.
type MockPostsStorageMockRecorder struct {
	mock *MockPostsStorage
}

// NewMockPostsStorage creates a new mock instance.
func NewMockPostsStorage(ctrl *gomock.Controller) *MockPostsStorage {
	mock := &MockPostsStorage{ctrl: ctrl}
	mock.recorder = &MockPostsStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostsStorage) EXPECT() *MockPostsStorageMockRecorder {
	return m.recorder
}

// DeleteComment mocks base method.
func (m *MockPostsStorage) DeleteComment(ctx context.Context, id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteComment", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteComment indicates an expected call of DeleteComment.
func (mr *MockPostsStorageMockRecorder) DeleteComment(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteComment", reflect.TypeOf((*MockPostsStorage)(nil).DeleteComment), ctx, id)
}

// DeleteCommentLike mocks base method.
func (m *MockPostsStorage) DeleteCommentLike(ctx context.Context, commentLike *domain.CommentLike) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCommentLike", ctx, commentLike)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCommentLike indicates an expected call of DeleteCommentLike.
func (mr *MockPostsStorageMockRecorder) DeleteCommentLike(ctx, commentLike interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCommentLike", reflect.TypeOf((*MockPostsStorage)(nil).DeleteCommentLike), ctx, commentLike)
}

// DeleteGroupPost mocks base method.
func (m *MockPostsStorage) DeleteGroupPost(ctx context.Context, postID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGroupPost", ctx, postID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGroupPost indicates an expected call of DeleteGroupPost.
func (mr *MockPostsStorageMockRecorder) DeleteGroupPost(ctx, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroupPost", reflect.TypeOf((*MockPostsStorage)(nil).DeleteGroupPost), ctx, postID)
}

// DeletePost mocks base method.
func (m *MockPostsStorage) DeletePost(ctx context.Context, postID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePost", ctx, postID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePost indicates an expected call of DeletePost.
func (mr *MockPostsStorageMockRecorder) DeletePost(ctx, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePost", reflect.TypeOf((*MockPostsStorage)(nil).DeletePost), ctx, postID)
}

// DeletePostLike mocks base method.
func (m *MockPostsStorage) DeletePostLike(ctx context.Context, likeData *domain.PostLike) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePostLike", ctx, likeData)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePostLike indicates an expected call of DeletePostLike.
func (mr *MockPostsStorageMockRecorder) DeletePostLike(ctx, likeData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePostLike", reflect.TypeOf((*MockPostsStorage)(nil).DeletePostLike), ctx, likeData)
}

// GetCommentByID mocks base method.
func (m *MockPostsStorage) GetCommentByID(ctx context.Context, id uint) (*domain.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommentByID", ctx, id)
	ret0, _ := ret[0].(*domain.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentByID indicates an expected call of GetCommentByID.
func (mr *MockPostsStorageMockRecorder) GetCommentByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentByID", reflect.TypeOf((*MockPostsStorage)(nil).GetCommentByID), ctx, id)
}

// GetCommentLikeByCommentIDAndUserID mocks base method.
func (m *MockPostsStorage) GetCommentLikeByCommentIDAndUserID(ctx context.Context, data *domain.CommentLike) (*domain.CommentLike, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommentLikeByCommentIDAndUserID", ctx, data)
	ret0, _ := ret[0].(*domain.CommentLike)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentLikeByCommentIDAndUserID indicates an expected call of GetCommentLikeByCommentIDAndUserID.
func (mr *MockPostsStorageMockRecorder) GetCommentLikeByCommentIDAndUserID(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentLikeByCommentIDAndUserID", reflect.TypeOf((*MockPostsStorage)(nil).GetCommentLikeByCommentIDAndUserID), ctx, data)
}

// GetCommentsByPostID mocks base method.
func (m *MockPostsStorage) GetCommentsByPostID(ctx context.Context, postID uint) ([]*domain.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommentsByPostID", ctx, postID)
	ret0, _ := ret[0].([]*domain.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentsByPostID indicates an expected call of GetCommentsByPostID.
func (mr *MockPostsStorageMockRecorder) GetCommentsByPostID(ctx, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentsByPostID", reflect.TypeOf((*MockPostsStorage)(nil).GetCommentsByPostID), ctx, postID)
}

// GetGroupPostsBySubscriptionIDs mocks base method.
func (m *MockPostsStorage) GetGroupPostsBySubscriptionIDs(ctx context.Context, subIDs []uint, lastPostID, postsAmount uint) ([]*domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupPostsBySubscriptionIDs", ctx, subIDs, lastPostID, postsAmount)
	ret0, _ := ret[0].([]*domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupPostsBySubscriptionIDs indicates an expected call of GetGroupPostsBySubscriptionIDs.
func (mr *MockPostsStorageMockRecorder) GetGroupPostsBySubscriptionIDs(ctx, subIDs, lastPostID, postsAmount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupPostsBySubscriptionIDs", reflect.TypeOf((*MockPostsStorage)(nil).GetGroupPostsBySubscriptionIDs), ctx, subIDs, lastPostID, postsAmount)
}

// GetLikedPosts mocks base method.
func (m *MockPostsStorage) GetLikedPosts(ctx context.Context, userID, lastLikeID, limit uint) ([]posts.LikeWithPost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikedPosts", ctx, userID, lastLikeID, limit)
	ret0, _ := ret[0].([]posts.LikeWithPost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikedPosts indicates an expected call of GetLikedPosts.
func (mr *MockPostsStorageMockRecorder) GetLikedPosts(ctx, userID, lastLikeID, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikedPosts", reflect.TypeOf((*MockPostsStorage)(nil).GetLikedPosts), ctx, userID, lastLikeID, limit)
}

// GetNewPosts mocks base method.
func (m *MockPostsStorage) GetNewPosts(ctx context.Context, lastPostID, postsAmount uint) ([]*domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNewPosts", ctx, lastPostID, postsAmount)
	ret0, _ := ret[0].([]*domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNewPosts indicates an expected call of GetNewPosts.
func (mr *MockPostsStorageMockRecorder) GetNewPosts(ctx, lastPostID, postsAmount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNewPosts", reflect.TypeOf((*MockPostsStorage)(nil).GetNewPosts), ctx, lastPostID, postsAmount)
}

// GetPostByID mocks base method.
func (m *MockPostsStorage) GetPostByID(ctx context.Context, postID uint) (*domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostByID", ctx, postID)
	ret0, _ := ret[0].(*domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostByID indicates an expected call of GetPostByID.
func (mr *MockPostsStorageMockRecorder) GetPostByID(ctx, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostByID", reflect.TypeOf((*MockPostsStorage)(nil).GetPostByID), ctx, postID)
}

// GetPostsByGroupSubIDsAndUserSubIDs mocks base method.
func (m *MockPostsStorage) GetPostsByGroupSubIDsAndUserSubIDs(ctx context.Context, groupSubIDs, userSubIDs []uint, lastPostID, postsAmount uint) ([]*domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostsByGroupSubIDsAndUserSubIDs", ctx, groupSubIDs, userSubIDs, lastPostID, postsAmount)
	ret0, _ := ret[0].([]*domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostsByGroupSubIDsAndUserSubIDs indicates an expected call of GetPostsByGroupSubIDsAndUserSubIDs.
func (mr *MockPostsStorageMockRecorder) GetPostsByGroupSubIDsAndUserSubIDs(ctx, groupSubIDs, userSubIDs, lastPostID, postsAmount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostsByGroupSubIDsAndUserSubIDs", reflect.TypeOf((*MockPostsStorage)(nil).GetPostsByGroupSubIDsAndUserSubIDs), ctx, groupSubIDs, userSubIDs, lastPostID, postsAmount)
}

// GetPostsOfGroup mocks base method.
func (m *MockPostsStorage) GetPostsOfGroup(ctx context.Context, groupID, lastPostID, postsAmount uint) ([]*domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostsOfGroup", ctx, groupID, lastPostID, postsAmount)
	ret0, _ := ret[0].([]*domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostsOfGroup indicates an expected call of GetPostsOfGroup.
func (mr *MockPostsStorageMockRecorder) GetPostsOfGroup(ctx, groupID, lastPostID, postsAmount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostsOfGroup", reflect.TypeOf((*MockPostsStorage)(nil).GetPostsOfGroup), ctx, groupID, lastPostID, postsAmount)
}

// GetUserFriendsPosts mocks base method.
func (m *MockPostsStorage) GetUserFriendsPosts(ctx context.Context, userID, lastPostID, postsAmount uint) ([]*domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserFriendsPosts", ctx, userID, lastPostID, postsAmount)
	ret0, _ := ret[0].([]*domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserFriendsPosts indicates an expected call of GetUserFriendsPosts.
func (mr *MockPostsStorageMockRecorder) GetUserFriendsPosts(ctx, userID, lastPostID, postsAmount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserFriendsPosts", reflect.TypeOf((*MockPostsStorage)(nil).GetUserFriendsPosts), ctx, userID, lastPostID, postsAmount)
}

// GetUserPosts mocks base method.
func (m *MockPostsStorage) GetUserPosts(ctx context.Context, userID, lastPostID, postsAmount uint) ([]*domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserPosts", ctx, userID, lastPostID, postsAmount)
	ret0, _ := ret[0].([]*domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserPosts indicates an expected call of GetUserPosts.
func (mr *MockPostsStorageMockRecorder) GetUserPosts(ctx, userID, lastPostID, postsAmount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserPosts", reflect.TypeOf((*MockPostsStorage)(nil).GetUserPosts), ctx, userID, lastPostID, postsAmount)
}

// StoreComment mocks base method.
func (m *MockPostsStorage) StoreComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreComment", ctx, comment)
	ret0, _ := ret[0].(*domain.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreComment indicates an expected call of StoreComment.
func (mr *MockPostsStorageMockRecorder) StoreComment(ctx, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreComment", reflect.TypeOf((*MockPostsStorage)(nil).StoreComment), ctx, comment)
}

// StoreCommentLike mocks base method.
func (m *MockPostsStorage) StoreCommentLike(ctx context.Context, commentLike *domain.CommentLike) (*domain.CommentLike, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreCommentLike", ctx, commentLike)
	ret0, _ := ret[0].(*domain.CommentLike)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreCommentLike indicates an expected call of StoreCommentLike.
func (mr *MockPostsStorageMockRecorder) StoreCommentLike(ctx, commentLike interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreCommentLike", reflect.TypeOf((*MockPostsStorage)(nil).StoreCommentLike), ctx, commentLike)
}

// StoreGroupPost mocks base method.
func (m *MockPostsStorage) StoreGroupPost(ctx context.Context, groupPost *domain.GroupPost) (*domain.GroupPost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreGroupPost", ctx, groupPost)
	ret0, _ := ret[0].(*domain.GroupPost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreGroupPost indicates an expected call of StoreGroupPost.
func (mr *MockPostsStorageMockRecorder) StoreGroupPost(ctx, groupPost interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreGroupPost", reflect.TypeOf((*MockPostsStorage)(nil).StoreGroupPost), ctx, groupPost)
}

// StorePost mocks base method.
func (m *MockPostsStorage) StorePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorePost", ctx, post)
	ret0, _ := ret[0].(*domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StorePost indicates an expected call of StorePost.
func (mr *MockPostsStorageMockRecorder) StorePost(ctx, post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorePost", reflect.TypeOf((*MockPostsStorage)(nil).StorePost), ctx, post)
}

// StorePostLike mocks base method.
func (m *MockPostsStorage) StorePostLike(ctx context.Context, likeData *domain.PostLike) (*domain.PostLike, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorePostLike", ctx, likeData)
	ret0, _ := ret[0].(*domain.PostLike)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StorePostLike indicates an expected call of StorePostLike.
func (mr *MockPostsStorageMockRecorder) StorePostLike(ctx, likeData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorePostLike", reflect.TypeOf((*MockPostsStorage)(nil).StorePostLike), ctx, likeData)
}

// UpdateComment mocks base method.
func (m *MockPostsStorage) UpdateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateComment", ctx, comment)
	ret0, _ := ret[0].(*domain.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateComment indicates an expected call of UpdateComment.
func (mr *MockPostsStorageMockRecorder) UpdateComment(ctx, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateComment", reflect.TypeOf((*MockPostsStorage)(nil).UpdateComment), ctx, comment)
}

// UpdatePost mocks base method.
func (m *MockPostsStorage) UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePost", ctx, post)
	ret0, _ := ret[0].(*domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePost indicates an expected call of UpdatePost.
func (mr *MockPostsStorageMockRecorder) UpdatePost(ctx, post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePost", reflect.TypeOf((*MockPostsStorage)(nil).UpdatePost), ctx, post)
}

// MockAttachmentStorage is a mock of AttachmentStorage interface.
type MockAttachmentStorage struct {
	ctrl     *gomock.Controller
	recorder *MockAttachmentStorageMockRecorder
}

// MockAttachmentStorageMockRecorder is the mock recorder for MockAttachmentStorage.
type MockAttachmentStorageMockRecorder struct {
	mock *MockAttachmentStorage
}

// NewMockAttachmentStorage creates a new mock instance.
func NewMockAttachmentStorage(ctrl *gomock.Controller) *MockAttachmentStorage {
	mock := &MockAttachmentStorage{ctrl: ctrl}
	mock.recorder = &MockAttachmentStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAttachmentStorage) EXPECT() *MockAttachmentStorageMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockAttachmentStorage) Delete(fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAttachmentStorageMockRecorder) Delete(fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAttachmentStorage)(nil).Delete), fileName)
}

// Store mocks base method.
func (m *MockAttachmentStorage) Store(fileName, filePath, contentType string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", fileName, filePath, contentType)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockAttachmentStorageMockRecorder) Store(fileName, filePath, contentType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockAttachmentStorage)(nil).Store), fileName, filePath, contentType)
}
