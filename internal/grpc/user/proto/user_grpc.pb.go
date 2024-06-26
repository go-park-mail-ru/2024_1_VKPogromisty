// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package user

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserClient interface {
	GetByID(ctx context.Context, in *GetByIDRequest, opts ...grpc.CallOption) (*GetByIDResponse, error)
	GetByEmail(ctx context.Context, in *GetByEmailRequest, opts ...grpc.CallOption) (*GetByEmailResponse, error)
	GetByIDWithSubsInfo(ctx context.Context, in *GetByIDWithSubsInfoRequest, opts ...grpc.CallOption) (*GetByIDWithSubsInfoResponse, error)
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	Upload(ctx context.Context, opts ...grpc.CallOption) (User_UploadClient, error)
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (*SubscribeResponse, error)
	Unsubscribe(ctx context.Context, in *UnsubscribeRequest, opts ...grpc.CallOption) (*UnsubscribeResponse, error)
	GetSubscriptions(ctx context.Context, in *GetSubscriptionsRequest, opts ...grpc.CallOption) (*GetSubscriptionsResponse, error)
	GetSubscribers(ctx context.Context, in *GetSubscribersRequest, opts ...grpc.CallOption) (*GetSubscribersResponse, error)
	GetFriends(ctx context.Context, in *GetFriendsRequest, opts ...grpc.CallOption) (*GetFriendsResponse, error)
	SearchByName(ctx context.Context, in *SearchByNameRequest, opts ...grpc.CallOption) (*SearchByNameResponse, error)
	GetSubscriptionIDs(ctx context.Context, in *GetSubscriptionIDsRequest, opts ...grpc.CallOption) (*GetSubscriptionIDsResponse, error)
	CreatePublicGroupAdmin(ctx context.Context, in *CreatePublicGroupAdminRequest, opts ...grpc.CallOption) (*CreatePublicGroupAdminResponse, error)
	DeletePublicGroupAdmin(ctx context.Context, in *DeletePublicGroupAdminRequest, opts ...grpc.CallOption) (*DeletePublicGroupAdminResponse, error)
	GetAdminsByPublicGroupID(ctx context.Context, in *GetAdminsByPublicGroupIDRequest, opts ...grpc.CallOption) (*GetAdminsByPublicGroupIDResponse, error)
	CheckIfUserIsAdmin(ctx context.Context, in *CheckIfUserIsAdminRequest, opts ...grpc.CallOption) (*CheckIfUserIsAdminResponse, error)
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

func (c *userClient) GetByID(ctx context.Context, in *GetByIDRequest, opts ...grpc.CallOption) (*GetByIDResponse, error) {
	out := new(GetByIDResponse)
	err := c.cc.Invoke(ctx, "/user.User/GetByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetByEmail(ctx context.Context, in *GetByEmailRequest, opts ...grpc.CallOption) (*GetByEmailResponse, error) {
	out := new(GetByEmailResponse)
	err := c.cc.Invoke(ctx, "/user.User/GetByEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetByIDWithSubsInfo(ctx context.Context, in *GetByIDWithSubsInfoRequest, opts ...grpc.CallOption) (*GetByIDWithSubsInfoResponse, error) {
	out := new(GetByIDWithSubsInfoResponse)
	err := c.cc.Invoke(ctx, "/user.User/GetByIDWithSubsInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/user.User/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, "/user.User/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, "/user.User/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) Upload(ctx context.Context, opts ...grpc.CallOption) (User_UploadClient, error) {
	stream, err := c.cc.NewStream(ctx, &User_ServiceDesc.Streams[0], "/user.User/Upload", opts...)
	if err != nil {
		return nil, err
	}
	x := &userUploadClient{stream}
	return x, nil
}

type User_UploadClient interface {
	Send(*UploadRequest) error
	CloseAndRecv() (*UploadResponse, error)
	grpc.ClientStream
}

type userUploadClient struct {
	grpc.ClientStream
}

func (x *userUploadClient) Send(m *UploadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *userUploadClient) CloseAndRecv() (*UploadResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *userClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (*SubscribeResponse, error) {
	out := new(SubscribeResponse)
	err := c.cc.Invoke(ctx, "/user.User/Subscribe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) Unsubscribe(ctx context.Context, in *UnsubscribeRequest, opts ...grpc.CallOption) (*UnsubscribeResponse, error) {
	out := new(UnsubscribeResponse)
	err := c.cc.Invoke(ctx, "/user.User/Unsubscribe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetSubscriptions(ctx context.Context, in *GetSubscriptionsRequest, opts ...grpc.CallOption) (*GetSubscriptionsResponse, error) {
	out := new(GetSubscriptionsResponse)
	err := c.cc.Invoke(ctx, "/user.User/GetSubscriptions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetSubscribers(ctx context.Context, in *GetSubscribersRequest, opts ...grpc.CallOption) (*GetSubscribersResponse, error) {
	out := new(GetSubscribersResponse)
	err := c.cc.Invoke(ctx, "/user.User/GetSubscribers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetFriends(ctx context.Context, in *GetFriendsRequest, opts ...grpc.CallOption) (*GetFriendsResponse, error) {
	out := new(GetFriendsResponse)
	err := c.cc.Invoke(ctx, "/user.User/GetFriends", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) SearchByName(ctx context.Context, in *SearchByNameRequest, opts ...grpc.CallOption) (*SearchByNameResponse, error) {
	out := new(SearchByNameResponse)
	err := c.cc.Invoke(ctx, "/user.User/SearchByName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetSubscriptionIDs(ctx context.Context, in *GetSubscriptionIDsRequest, opts ...grpc.CallOption) (*GetSubscriptionIDsResponse, error) {
	out := new(GetSubscriptionIDsResponse)
	err := c.cc.Invoke(ctx, "/user.User/GetSubscriptionIDs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) CreatePublicGroupAdmin(ctx context.Context, in *CreatePublicGroupAdminRequest, opts ...grpc.CallOption) (*CreatePublicGroupAdminResponse, error) {
	out := new(CreatePublicGroupAdminResponse)
	err := c.cc.Invoke(ctx, "/user.User/CreatePublicGroupAdmin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) DeletePublicGroupAdmin(ctx context.Context, in *DeletePublicGroupAdminRequest, opts ...grpc.CallOption) (*DeletePublicGroupAdminResponse, error) {
	out := new(DeletePublicGroupAdminResponse)
	err := c.cc.Invoke(ctx, "/user.User/DeletePublicGroupAdmin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetAdminsByPublicGroupID(ctx context.Context, in *GetAdminsByPublicGroupIDRequest, opts ...grpc.CallOption) (*GetAdminsByPublicGroupIDResponse, error) {
	out := new(GetAdminsByPublicGroupIDResponse)
	err := c.cc.Invoke(ctx, "/user.User/GetAdminsByPublicGroupID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) CheckIfUserIsAdmin(ctx context.Context, in *CheckIfUserIsAdminRequest, opts ...grpc.CallOption) (*CheckIfUserIsAdminResponse, error) {
	out := new(CheckIfUserIsAdminResponse)
	err := c.cc.Invoke(ctx, "/user.User/CheckIfUserIsAdmin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServer is the server API for User service.
// All implementations must embed UnimplementedUserServer
// for forward compatibility
type UserServer interface {
	GetByID(context.Context, *GetByIDRequest) (*GetByIDResponse, error)
	GetByEmail(context.Context, *GetByEmailRequest) (*GetByEmailResponse, error)
	GetByIDWithSubsInfo(context.Context, *GetByIDWithSubsInfoRequest) (*GetByIDWithSubsInfoResponse, error)
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Update(context.Context, *UpdateRequest) (*UpdateResponse, error)
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	Upload(User_UploadServer) error
	Subscribe(context.Context, *SubscribeRequest) (*SubscribeResponse, error)
	Unsubscribe(context.Context, *UnsubscribeRequest) (*UnsubscribeResponse, error)
	GetSubscriptions(context.Context, *GetSubscriptionsRequest) (*GetSubscriptionsResponse, error)
	GetSubscribers(context.Context, *GetSubscribersRequest) (*GetSubscribersResponse, error)
	GetFriends(context.Context, *GetFriendsRequest) (*GetFriendsResponse, error)
	SearchByName(context.Context, *SearchByNameRequest) (*SearchByNameResponse, error)
	GetSubscriptionIDs(context.Context, *GetSubscriptionIDsRequest) (*GetSubscriptionIDsResponse, error)
	CreatePublicGroupAdmin(context.Context, *CreatePublicGroupAdminRequest) (*CreatePublicGroupAdminResponse, error)
	DeletePublicGroupAdmin(context.Context, *DeletePublicGroupAdminRequest) (*DeletePublicGroupAdminResponse, error)
	GetAdminsByPublicGroupID(context.Context, *GetAdminsByPublicGroupIDRequest) (*GetAdminsByPublicGroupIDResponse, error)
	CheckIfUserIsAdmin(context.Context, *CheckIfUserIsAdminRequest) (*CheckIfUserIsAdminResponse, error)
	mustEmbedUnimplementedUserServer()
}

// UnimplementedUserServer must be embedded to have forward compatible implementations.
type UnimplementedUserServer struct {
}

func (UnimplementedUserServer) GetByID(context.Context, *GetByIDRequest) (*GetByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByID not implemented")
}
func (UnimplementedUserServer) GetByEmail(context.Context, *GetByEmailRequest) (*GetByEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByEmail not implemented")
}
func (UnimplementedUserServer) GetByIDWithSubsInfo(context.Context, *GetByIDWithSubsInfoRequest) (*GetByIDWithSubsInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByIDWithSubsInfo not implemented")
}
func (UnimplementedUserServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedUserServer) Update(context.Context, *UpdateRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedUserServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedUserServer) Upload(User_UploadServer) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedUserServer) Subscribe(context.Context, *SubscribeRequest) (*SubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedUserServer) Unsubscribe(context.Context, *UnsubscribeRequest) (*UnsubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unsubscribe not implemented")
}
func (UnimplementedUserServer) GetSubscriptions(context.Context, *GetSubscriptionsRequest) (*GetSubscriptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSubscriptions not implemented")
}
func (UnimplementedUserServer) GetSubscribers(context.Context, *GetSubscribersRequest) (*GetSubscribersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSubscribers not implemented")
}
func (UnimplementedUserServer) GetFriends(context.Context, *GetFriendsRequest) (*GetFriendsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFriends not implemented")
}
func (UnimplementedUserServer) SearchByName(context.Context, *SearchByNameRequest) (*SearchByNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchByName not implemented")
}
func (UnimplementedUserServer) GetSubscriptionIDs(context.Context, *GetSubscriptionIDsRequest) (*GetSubscriptionIDsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSubscriptionIDs not implemented")
}
func (UnimplementedUserServer) CreatePublicGroupAdmin(context.Context, *CreatePublicGroupAdminRequest) (*CreatePublicGroupAdminResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePublicGroupAdmin not implemented")
}
func (UnimplementedUserServer) DeletePublicGroupAdmin(context.Context, *DeletePublicGroupAdminRequest) (*DeletePublicGroupAdminResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePublicGroupAdmin not implemented")
}
func (UnimplementedUserServer) GetAdminsByPublicGroupID(context.Context, *GetAdminsByPublicGroupIDRequest) (*GetAdminsByPublicGroupIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAdminsByPublicGroupID not implemented")
}
func (UnimplementedUserServer) CheckIfUserIsAdmin(context.Context, *CheckIfUserIsAdminRequest) (*CheckIfUserIsAdminResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckIfUserIsAdmin not implemented")
}
func (UnimplementedUserServer) mustEmbedUnimplementedUserServer() {}

// UnsafeUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServer will
// result in compilation errors.
type UnsafeUserServer interface {
	mustEmbedUnimplementedUserServer()
}

func RegisterUserServer(s grpc.ServiceRegistrar, srv UserServer) {
	s.RegisterService(&User_ServiceDesc, srv)
}

func _User_GetByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/GetByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetByID(ctx, req.(*GetByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetByEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetByEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/GetByEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetByEmail(ctx, req.(*GetByEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetByIDWithSubsInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIDWithSubsInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetByIDWithSubsInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/GetByIDWithSubsInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetByIDWithSubsInfo(ctx, req.(*GetByIDWithSubsInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UserServer).Upload(&userUploadServer{stream})
}

type User_UploadServer interface {
	SendAndClose(*UploadResponse) error
	Recv() (*UploadRequest, error)
	grpc.ServerStream
}

type userUploadServer struct {
	grpc.ServerStream
}

func (x *userUploadServer) SendAndClose(m *UploadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *userUploadServer) Recv() (*UploadRequest, error) {
	m := new(UploadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _User_Subscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscribeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Subscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/Subscribe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Subscribe(ctx, req.(*SubscribeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_Unsubscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnsubscribeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Unsubscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/Unsubscribe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Unsubscribe(ctx, req.(*UnsubscribeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetSubscriptions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSubscriptionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetSubscriptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/GetSubscriptions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetSubscriptions(ctx, req.(*GetSubscriptionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetSubscribers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSubscribersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetSubscribers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/GetSubscribers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetSubscribers(ctx, req.(*GetSubscribersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetFriends_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFriendsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetFriends(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/GetFriends",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetFriends(ctx, req.(*GetFriendsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_SearchByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchByNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).SearchByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/SearchByName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).SearchByName(ctx, req.(*SearchByNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetSubscriptionIDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSubscriptionIDsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetSubscriptionIDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/GetSubscriptionIDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetSubscriptionIDs(ctx, req.(*GetSubscriptionIDsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_CreatePublicGroupAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePublicGroupAdminRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).CreatePublicGroupAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/CreatePublicGroupAdmin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).CreatePublicGroupAdmin(ctx, req.(*CreatePublicGroupAdminRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_DeletePublicGroupAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePublicGroupAdminRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).DeletePublicGroupAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/DeletePublicGroupAdmin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).DeletePublicGroupAdmin(ctx, req.(*DeletePublicGroupAdminRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetAdminsByPublicGroupID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAdminsByPublicGroupIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetAdminsByPublicGroupID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/GetAdminsByPublicGroupID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetAdminsByPublicGroupID(ctx, req.(*GetAdminsByPublicGroupIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_CheckIfUserIsAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckIfUserIsAdminRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).CheckIfUserIsAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/CheckIfUserIsAdmin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).CheckIfUserIsAdmin(ctx, req.(*CheckIfUserIsAdminRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// User_ServiceDesc is the grpc.ServiceDesc for User service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var User_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.User",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetByID",
			Handler:    _User_GetByID_Handler,
		},
		{
			MethodName: "GetByEmail",
			Handler:    _User_GetByEmail_Handler,
		},
		{
			MethodName: "GetByIDWithSubsInfo",
			Handler:    _User_GetByIDWithSubsInfo_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _User_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _User_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _User_Delete_Handler,
		},
		{
			MethodName: "Subscribe",
			Handler:    _User_Subscribe_Handler,
		},
		{
			MethodName: "Unsubscribe",
			Handler:    _User_Unsubscribe_Handler,
		},
		{
			MethodName: "GetSubscriptions",
			Handler:    _User_GetSubscriptions_Handler,
		},
		{
			MethodName: "GetSubscribers",
			Handler:    _User_GetSubscribers_Handler,
		},
		{
			MethodName: "GetFriends",
			Handler:    _User_GetFriends_Handler,
		},
		{
			MethodName: "SearchByName",
			Handler:    _User_SearchByName_Handler,
		},
		{
			MethodName: "GetSubscriptionIDs",
			Handler:    _User_GetSubscriptionIDs_Handler,
		},
		{
			MethodName: "CreatePublicGroupAdmin",
			Handler:    _User_CreatePublicGroupAdmin_Handler,
		},
		{
			MethodName: "DeletePublicGroupAdmin",
			Handler:    _User_DeletePublicGroupAdmin_Handler,
		},
		{
			MethodName: "GetAdminsByPublicGroupID",
			Handler:    _User_GetAdminsByPublicGroupID_Handler,
		},
		{
			MethodName: "CheckIfUserIsAdmin",
			Handler:    _User_CheckIfUserIsAdmin_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Upload",
			Handler:       _User_Upload_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "user.proto",
}
