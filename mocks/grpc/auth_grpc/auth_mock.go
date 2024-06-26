// Code generated by MockGen. DO NOT EDIT.
// Source: internal/grpc/auth/proto/auth_grpc.pb.go

// Package auth_grpc is a generated GoMock package.
package auth_grpc

import (
	context "context"
	reflect "reflect"
	auth "socio/internal/grpc/auth/proto"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockAuthClient is a mock of AuthClient interface.
type MockAuthClient struct {
	ctrl     *gomock.Controller
	recorder *MockAuthClientMockRecorder
}

// MockAuthClientMockRecorder is the mock recorder for MockAuthClient.
type MockAuthClientMockRecorder struct {
	mock *MockAuthClient
}

// NewMockAuthClient creates a new mock instance.
func NewMockAuthClient(ctrl *gomock.Controller) *MockAuthClient {
	mock := &MockAuthClient{ctrl: ctrl}
	mock.recorder = &MockAuthClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthClient) EXPECT() *MockAuthClientMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAuthClient) Login(ctx context.Context, in *auth.LoginRequest, opts ...grpc.CallOption) (*auth.LoginResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Login", varargs...)
	ret0, _ := ret[0].(*auth.LoginResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthClientMockRecorder) Login(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthClient)(nil).Login), varargs...)
}

// Logout mocks base method.
func (m *MockAuthClient) Logout(ctx context.Context, in *auth.LogoutRequest, opts ...grpc.CallOption) (*auth.LogoutResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Logout", varargs...)
	ret0, _ := ret[0].(*auth.LogoutResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Logout indicates an expected call of Logout.
func (mr *MockAuthClientMockRecorder) Logout(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockAuthClient)(nil).Logout), varargs...)
}

// ValidateSession mocks base method.
func (m *MockAuthClient) ValidateSession(ctx context.Context, in *auth.ValidateSessionRequest, opts ...grpc.CallOption) (*auth.ValidateSessionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ValidateSession", varargs...)
	ret0, _ := ret[0].(*auth.ValidateSessionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateSession indicates an expected call of ValidateSession.
func (mr *MockAuthClientMockRecorder) ValidateSession(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateSession", reflect.TypeOf((*MockAuthClient)(nil).ValidateSession), varargs...)
}

// MockAuthServer is a mock of AuthServer interface.
type MockAuthServer struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServerMockRecorder
}

// MockAuthServerMockRecorder is the mock recorder for MockAuthServer.
type MockAuthServerMockRecorder struct {
	mock *MockAuthServer
}

// NewMockAuthServer creates a new mock instance.
func NewMockAuthServer(ctrl *gomock.Controller) *MockAuthServer {
	mock := &MockAuthServer{ctrl: ctrl}
	mock.recorder = &MockAuthServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthServer) EXPECT() *MockAuthServerMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAuthServer) Login(arg0 context.Context, arg1 *auth.LoginRequest) (*auth.LoginResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(*auth.LoginResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthServerMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthServer)(nil).Login), arg0, arg1)
}

// Logout mocks base method.
func (m *MockAuthServer) Logout(arg0 context.Context, arg1 *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", arg0, arg1)
	ret0, _ := ret[0].(*auth.LogoutResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Logout indicates an expected call of Logout.
func (mr *MockAuthServerMockRecorder) Logout(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockAuthServer)(nil).Logout), arg0, arg1)
}

// ValidateSession mocks base method.
func (m *MockAuthServer) ValidateSession(arg0 context.Context, arg1 *auth.ValidateSessionRequest) (*auth.ValidateSessionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateSession", arg0, arg1)
	ret0, _ := ret[0].(*auth.ValidateSessionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateSession indicates an expected call of ValidateSession.
func (mr *MockAuthServerMockRecorder) ValidateSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateSession", reflect.TypeOf((*MockAuthServer)(nil).ValidateSession), arg0, arg1)
}

// mustEmbedUnimplementedAuthServer mocks base method.
func (m *MockAuthServer) mustEmbedUnimplementedAuthServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedAuthServer")
}

// mustEmbedUnimplementedAuthServer indicates an expected call of mustEmbedUnimplementedAuthServer.
func (mr *MockAuthServerMockRecorder) mustEmbedUnimplementedAuthServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedAuthServer", reflect.TypeOf((*MockAuthServer)(nil).mustEmbedUnimplementedAuthServer))
}

// MockUnsafeAuthServer is a mock of UnsafeAuthServer interface.
type MockUnsafeAuthServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeAuthServerMockRecorder
}

// MockUnsafeAuthServerMockRecorder is the mock recorder for MockUnsafeAuthServer.
type MockUnsafeAuthServerMockRecorder struct {
	mock *MockUnsafeAuthServer
}

// NewMockUnsafeAuthServer creates a new mock instance.
func NewMockUnsafeAuthServer(ctrl *gomock.Controller) *MockUnsafeAuthServer {
	mock := &MockUnsafeAuthServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeAuthServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeAuthServer) EXPECT() *MockUnsafeAuthServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedAuthServer mocks base method.
func (m *MockUnsafeAuthServer) mustEmbedUnimplementedAuthServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedAuthServer")
}

// mustEmbedUnimplementedAuthServer indicates an expected call of mustEmbedUnimplementedAuthServer.
func (mr *MockUnsafeAuthServerMockRecorder) mustEmbedUnimplementedAuthServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedAuthServer", reflect.TypeOf((*MockUnsafeAuthServer)(nil).mustEmbedUnimplementedAuthServer))
}
