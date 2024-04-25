// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package post

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

// PostClient is the client API for Post service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostClient interface {
	GetPostByID(ctx context.Context, in *GetPostByIDRequest, opts ...grpc.CallOption) (*GetPostByIDResponse, error)
	GetUserPosts(ctx context.Context, in *GetUserPostsRequest, opts ...grpc.CallOption) (*GetUserPostsResponse, error)
	GetUserFriendsPosts(ctx context.Context, in *GetUserFriendsPostsRequest, opts ...grpc.CallOption) (*GetUserFriendsPostsResponse, error)
	CreatePost(ctx context.Context, in *CreatePostRequest, opts ...grpc.CallOption) (*CreatePostResponse, error)
	UpdatePost(ctx context.Context, in *UpdatePostRequest, opts ...grpc.CallOption) (*UpdatePostResponse, error)
	DeletePost(ctx context.Context, in *DeletePostRequest, opts ...grpc.CallOption) (*DeletePostResponse, error)
	GetLikedPosts(ctx context.Context, in *GetLikedPostsRequest, opts ...grpc.CallOption) (*GetLikedPostsResponse, error)
	LikePost(ctx context.Context, in *LikePostRequest, opts ...grpc.CallOption) (*LikePostResponse, error)
	UnlikePost(ctx context.Context, in *UnlikePostRequest, opts ...grpc.CallOption) (*UnlikePostResponse, error)
	Upload(ctx context.Context, opts ...grpc.CallOption) (Post_UploadClient, error)
}

type postClient struct {
	cc grpc.ClientConnInterface
}

func NewPostClient(cc grpc.ClientConnInterface) PostClient {
	return &postClient{cc}
}

func (c *postClient) GetPostByID(ctx context.Context, in *GetPostByIDRequest, opts ...grpc.CallOption) (*GetPostByIDResponse, error) {
	out := new(GetPostByIDResponse)
	err := c.cc.Invoke(ctx, "/post.Post/GetPostByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postClient) GetUserPosts(ctx context.Context, in *GetUserPostsRequest, opts ...grpc.CallOption) (*GetUserPostsResponse, error) {
	out := new(GetUserPostsResponse)
	err := c.cc.Invoke(ctx, "/post.Post/GetUserPosts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postClient) GetUserFriendsPosts(ctx context.Context, in *GetUserFriendsPostsRequest, opts ...grpc.CallOption) (*GetUserFriendsPostsResponse, error) {
	out := new(GetUserFriendsPostsResponse)
	err := c.cc.Invoke(ctx, "/post.Post/GetUserFriendsPosts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postClient) CreatePost(ctx context.Context, in *CreatePostRequest, opts ...grpc.CallOption) (*CreatePostResponse, error) {
	out := new(CreatePostResponse)
	err := c.cc.Invoke(ctx, "/post.Post/CreatePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postClient) UpdatePost(ctx context.Context, in *UpdatePostRequest, opts ...grpc.CallOption) (*UpdatePostResponse, error) {
	out := new(UpdatePostResponse)
	err := c.cc.Invoke(ctx, "/post.Post/UpdatePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postClient) DeletePost(ctx context.Context, in *DeletePostRequest, opts ...grpc.CallOption) (*DeletePostResponse, error) {
	out := new(DeletePostResponse)
	err := c.cc.Invoke(ctx, "/post.Post/DeletePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postClient) GetLikedPosts(ctx context.Context, in *GetLikedPostsRequest, opts ...grpc.CallOption) (*GetLikedPostsResponse, error) {
	out := new(GetLikedPostsResponse)
	err := c.cc.Invoke(ctx, "/post.Post/GetLikedPosts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postClient) LikePost(ctx context.Context, in *LikePostRequest, opts ...grpc.CallOption) (*LikePostResponse, error) {
	out := new(LikePostResponse)
	err := c.cc.Invoke(ctx, "/post.Post/LikePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postClient) UnlikePost(ctx context.Context, in *UnlikePostRequest, opts ...grpc.CallOption) (*UnlikePostResponse, error) {
	out := new(UnlikePostResponse)
	err := c.cc.Invoke(ctx, "/post.Post/UnlikePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postClient) Upload(ctx context.Context, opts ...grpc.CallOption) (Post_UploadClient, error) {
	stream, err := c.cc.NewStream(ctx, &Post_ServiceDesc.Streams[0], "/post.Post/Upload", opts...)
	if err != nil {
		return nil, err
	}
	x := &postUploadClient{stream}
	return x, nil
}

type Post_UploadClient interface {
	Send(*UploadRequest) error
	CloseAndRecv() (*UploadResponse, error)
	grpc.ClientStream
}

type postUploadClient struct {
	grpc.ClientStream
}

func (x *postUploadClient) Send(m *UploadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *postUploadClient) CloseAndRecv() (*UploadResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PostServer is the server API for Post service.
// All implementations must embed UnimplementedPostServer
// for forward compatibility
type PostServer interface {
	GetPostByID(context.Context, *GetPostByIDRequest) (*GetPostByIDResponse, error)
	GetUserPosts(context.Context, *GetUserPostsRequest) (*GetUserPostsResponse, error)
	GetUserFriendsPosts(context.Context, *GetUserFriendsPostsRequest) (*GetUserFriendsPostsResponse, error)
	CreatePost(context.Context, *CreatePostRequest) (*CreatePostResponse, error)
	UpdatePost(context.Context, *UpdatePostRequest) (*UpdatePostResponse, error)
	DeletePost(context.Context, *DeletePostRequest) (*DeletePostResponse, error)
	GetLikedPosts(context.Context, *GetLikedPostsRequest) (*GetLikedPostsResponse, error)
	LikePost(context.Context, *LikePostRequest) (*LikePostResponse, error)
	UnlikePost(context.Context, *UnlikePostRequest) (*UnlikePostResponse, error)
	Upload(Post_UploadServer) error
	mustEmbedUnimplementedPostServer()
}

// UnimplementedPostServer must be embedded to have forward compatible implementations.
type UnimplementedPostServer struct {
}

func (UnimplementedPostServer) GetPostByID(context.Context, *GetPostByIDRequest) (*GetPostByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPostByID not implemented")
}
func (UnimplementedPostServer) GetUserPosts(context.Context, *GetUserPostsRequest) (*GetUserPostsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserPosts not implemented")
}
func (UnimplementedPostServer) GetUserFriendsPosts(context.Context, *GetUserFriendsPostsRequest) (*GetUserFriendsPostsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserFriendsPosts not implemented")
}
func (UnimplementedPostServer) CreatePost(context.Context, *CreatePostRequest) (*CreatePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePost not implemented")
}
func (UnimplementedPostServer) UpdatePost(context.Context, *UpdatePostRequest) (*UpdatePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePost not implemented")
}
func (UnimplementedPostServer) DeletePost(context.Context, *DeletePostRequest) (*DeletePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePost not implemented")
}
func (UnimplementedPostServer) GetLikedPosts(context.Context, *GetLikedPostsRequest) (*GetLikedPostsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLikedPosts not implemented")
}
func (UnimplementedPostServer) LikePost(context.Context, *LikePostRequest) (*LikePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LikePost not implemented")
}
func (UnimplementedPostServer) UnlikePost(context.Context, *UnlikePostRequest) (*UnlikePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnlikePost not implemented")
}
func (UnimplementedPostServer) Upload(Post_UploadServer) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedPostServer) mustEmbedUnimplementedPostServer() {}

// UnsafePostServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostServer will
// result in compilation errors.
type UnsafePostServer interface {
	mustEmbedUnimplementedPostServer()
}

func RegisterPostServer(s grpc.ServiceRegistrar, srv PostServer) {
	s.RegisterService(&Post_ServiceDesc, srv)
}

func _Post_GetPostByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPostByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServer).GetPostByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post.Post/GetPostByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServer).GetPostByID(ctx, req.(*GetPostByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Post_GetUserPosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserPostsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServer).GetUserPosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post.Post/GetUserPosts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServer).GetUserPosts(ctx, req.(*GetUserPostsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Post_GetUserFriendsPosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserFriendsPostsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServer).GetUserFriendsPosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post.Post/GetUserFriendsPosts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServer).GetUserFriendsPosts(ctx, req.(*GetUserFriendsPostsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Post_CreatePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServer).CreatePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post.Post/CreatePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServer).CreatePost(ctx, req.(*CreatePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Post_UpdatePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServer).UpdatePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post.Post/UpdatePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServer).UpdatePost(ctx, req.(*UpdatePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Post_DeletePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServer).DeletePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post.Post/DeletePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServer).DeletePost(ctx, req.(*DeletePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Post_GetLikedPosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLikedPostsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServer).GetLikedPosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post.Post/GetLikedPosts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServer).GetLikedPosts(ctx, req.(*GetLikedPostsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Post_LikePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LikePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServer).LikePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post.Post/LikePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServer).LikePost(ctx, req.(*LikePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Post_UnlikePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnlikePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServer).UnlikePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post.Post/UnlikePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServer).UnlikePost(ctx, req.(*UnlikePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Post_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PostServer).Upload(&postUploadServer{stream})
}

type Post_UploadServer interface {
	SendAndClose(*UploadResponse) error
	Recv() (*UploadRequest, error)
	grpc.ServerStream
}

type postUploadServer struct {
	grpc.ServerStream
}

func (x *postUploadServer) SendAndClose(m *UploadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *postUploadServer) Recv() (*UploadRequest, error) {
	m := new(UploadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Post_ServiceDesc is the grpc.ServiceDesc for Post service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Post_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "post.Post",
	HandlerType: (*PostServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPostByID",
			Handler:    _Post_GetPostByID_Handler,
		},
		{
			MethodName: "GetUserPosts",
			Handler:    _Post_GetUserPosts_Handler,
		},
		{
			MethodName: "GetUserFriendsPosts",
			Handler:    _Post_GetUserFriendsPosts_Handler,
		},
		{
			MethodName: "CreatePost",
			Handler:    _Post_CreatePost_Handler,
		},
		{
			MethodName: "UpdatePost",
			Handler:    _Post_UpdatePost_Handler,
		},
		{
			MethodName: "DeletePost",
			Handler:    _Post_DeletePost_Handler,
		},
		{
			MethodName: "GetLikedPosts",
			Handler:    _Post_GetLikedPosts_Handler,
		},
		{
			MethodName: "LikePost",
			Handler:    _Post_LikePost_Handler,
		},
		{
			MethodName: "UnlikePost",
			Handler:    _Post_UnlikePost_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Upload",
			Handler:       _Post_Upload_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "post.proto",
}