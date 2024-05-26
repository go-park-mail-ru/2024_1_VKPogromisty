package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"socio/domain"
	"socio/errors"
	postpb "socio/internal/grpc/post/proto"
	pgpb "socio/internal/grpc/public_group/proto"
	uspb "socio/internal/grpc/user/proto"
	mock_posts "socio/mocks/grpc/post_grpc"
	mock_public_group "socio/mocks/grpc/public_group_grpc"
	mock_user "socio/mocks/grpc/user_grpc"
	"socio/pkg/requestcontext"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var (
	validCtx = context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1))
)

type fields struct {
	PostsClient       *mock_posts.MockPostClient
	UserClient        *mock_user.MockUserClient
	PublicGroupClient *mock_public_group.MockPublicGroupClient
}

func TestHandleGetPostByID(t *testing.T) {
	tests := []struct {
		name           string
		request        *http.Request
		ctx            context.Context
		expectedStatus int
		prepare        func(f *fields)
	}{
		{
			name:           "TestHandleGetPostByID",
			request:        httptest.NewRequest("GET", "/posts/1", nil),
			expectedStatus: http.StatusOK,
			ctx:            validCtx,
			prepare: func(f *fields) {
				f.PostsClient.EXPECT().GetPostByID(gomock.Any(), gomock.Any()).Return(
					&postpb.GetPostByIDResponse{
						Post: &postpb.PostResponse{
							Id: 1,
						},
					}, nil,
				)
				f.UserClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&uspb.GetByIDResponse{
					User: &uspb.UserResponse{
						Id: 1,
					},
				}, nil)
				f.PostsClient.EXPECT().GetGroupPostByPostID(gomock.Any(), gomock.Any()).Return(
					&postpb.GetGroupPostByPostIDResponse{
						GroupPost: &postpb.GroupPostResponse{
							PostId: 1,
						},
					}, nil,
				)
				f.PublicGroupClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					&pgpb.GetByIDResponse{
						PublicGroup: &pgpb.PublicGroupWithInfoResponse{
							PublicGroup: &pgpb.PublicGroupResponse{
								Id: 1,
							},
						},
					}, nil,
				)
			},
		},
		{
			name:           "TestHandleGetPostByID",
			request:        httptest.NewRequest("GET", "/posts/", nil),
			ctx:            validCtx,
			expectedStatus: http.StatusNotFound,
			prepare: func(f *fields) {
			},
		},
		{
			name:           "TestHandleGetPostByID",
			request:        httptest.NewRequest("GET", "/posts/asd", nil),
			ctx:            validCtx,
			expectedStatus: http.StatusBadRequest,
			prepare: func(f *fields) {
			},
		},
		{
			name:           "TestHandleGetPostByID",
			request:        httptest.NewRequest("GET", "/posts/1", nil),
			expectedStatus: http.StatusNotFound,
			ctx:            validCtx,
			prepare: func(f *fields) {
				f.PostsClient.EXPECT().GetPostByID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound.GRPCStatus().Err())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := &fields{
				PostsClient:       mock_posts.NewMockPostClient(ctrl),
				UserClient:        mock_user.NewMockUserClient(ctrl),
				PublicGroupClient: mock_public_group.NewMockPublicGroupClient(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			h := NewPostsHandler(f.PostsClient, f.UserClient, f.PublicGroupClient)

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/posts/{postID}", h.HandleGetPostByID)

			tt.request = tt.request.WithContext(tt.ctx)

			router.ServeHTTP(rr, tt.request)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestHandleGetUserPosts(t *testing.T) {
	tests := []struct {
		name           string
		request        *http.Request
		expectedStatus int
		prepare        func(f *fields)
	}{
		{
			name:           "success",
			request:        httptest.NewRequest("GET", "/posts?userId=1&lastPostId=0&postsAmount=20", nil),
			expectedStatus: http.StatusOK,
			prepare: func(f *fields) {
				f.PostsClient.EXPECT().GetUserPosts(gomock.Any(), gomock.Any()).Return(nil, nil)
				f.UserClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&uspb.GetByIDResponse{
					User: &uspb.UserResponse{
						Id: 1,
					},
				}, nil)
			},
		},
		{
			name:           "err",
			request:        httptest.NewRequest("GET", "/posts?userId=1&lastPostId=0&postsAmount=20", nil),
			expectedStatus: http.StatusInternalServerError,
			prepare: func(f *fields) {
				f.PostsClient.EXPECT().GetUserPosts(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal.GRPCStatus().Err())
			},
		},
		{
			name:           "invalid user id",
			request:        httptest.NewRequest("GET", "/posts?userId=tyaazh&lastPostId=0&postsAmount=20", nil),
			expectedStatus: http.StatusBadRequest,
			prepare: func(f *fields) {
			},
		},
		{
			name:           "invalid last post id",
			request:        httptest.NewRequest("GET", "/posts?userId=1&lastPostId=opa&postsAmount=20", nil),
			expectedStatus: http.StatusBadRequest,
			prepare: func(f *fields) {
			},
		},
		{
			name:           "TestHandleGetUserPosts",
			request:        httptest.NewRequest("GET", "/posts?userId=1&lastPostId=0&postsAmount=opa", nil),
			expectedStatus: http.StatusBadRequest,
			prepare: func(f *fields) {
			},
		},
		{
			name:           "internal error",
			request:        httptest.NewRequest("GET", "/posts?userId=1&lastPostId=0&postsAmount=20", nil),
			expectedStatus: http.StatusInternalServerError,
			prepare: func(f *fields) {
				f.PostsClient.EXPECT().GetUserPosts(gomock.Any(), gomock.Any()).Return(nil, nil)
				f.UserClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := &fields{
				PostsClient:       mock_posts.NewMockPostClient(ctrl),
				UserClient:        mock_user.NewMockUserClient(ctrl),
				PublicGroupClient: mock_public_group.NewMockPublicGroupClient(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			h := NewPostsHandler(f.PostsClient, f.UserClient, f.PublicGroupClient)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.HandleGetUserPosts)

			handler.ServeHTTP(rr, tt.request)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestHandleGetUserFriendsPosts(t *testing.T) {
	validCtx := context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1))

	tests := []struct {
		name           string
		request        *http.Request
		expectedStatus int
		ctx            context.Context
		prepare        func(f *fields)
	}{
		{
			name:           "success",
			request:        httptest.NewRequest("GET", "/friends/posts?lastPostId=0&postsAmount=20", nil),
			expectedStatus: http.StatusOK,
			ctx:            validCtx,
			prepare: func(f *fields) {
				f.PostsClient.EXPECT().GetUserFriendsPosts(gomock.Any(), gomock.Any()).Return(&postpb.GetUserFriendsPostsResponse{
					Posts: []*postpb.PostResponse{
						{
							Id:       1,
							AuthorId: 1,
						},
					},
				}, nil)
				f.UserClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&uspb.GetByIDResponse{
					User: &uspb.UserResponse{
						Id: 1,
					},
				}, nil)
			},
		},
		{
			name:           "invalid last post id",
			request:        httptest.NewRequest("GET", "/friends/posts?lastPostId=tyazh&postsAmount=20", nil),
			expectedStatus: http.StatusBadRequest,
			ctx:            validCtx,
			prepare: func(f *fields) {
			},
		},
		{
			name:           "invalid posts amount",
			request:        httptest.NewRequest("GET", "/friends/posts?lastPostId=0&postsAmount=oppa", nil),
			expectedStatus: http.StatusBadRequest,
			ctx:            validCtx,
			prepare: func(f *fields) {
			},
		},
		{
			name:           "invalid context",
			request:        httptest.NewRequest("GET", "/friends/posts?lastPostId=0&postsAmount=20", nil),
			expectedStatus: http.StatusBadRequest,
			ctx:            context.Background(),
			prepare: func(f *fields) {
			},
		},
		{
			name:           "internal error",
			request:        httptest.NewRequest("GET", "/friends/posts?lastPostId=0&postsAmount=20", nil),
			expectedStatus: http.StatusInternalServerError,
			ctx:            validCtx,
			prepare: func(f *fields) {
				f.PostsClient.EXPECT().GetUserFriendsPosts(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal.GRPCStatus().Err())
			},
		},
		{
			name:           "internal error",
			request:        httptest.NewRequest("GET", "/friends/posts?lastPostId=0&postsAmount=20", nil),
			expectedStatus: http.StatusInternalServerError,
			ctx:            validCtx,
			prepare: func(f *fields) {
				f.PostsClient.EXPECT().GetUserFriendsPosts(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
		{
			name:           "internal error user",
			request:        httptest.NewRequest("GET", "/friends/posts?lastPostId=0&postsAmount=20", nil),
			expectedStatus: http.StatusInternalServerError,
			ctx:            validCtx,
			prepare: func(f *fields) {
				f.PostsClient.EXPECT().GetUserFriendsPosts(gomock.Any(), gomock.Any()).Return(&postpb.GetUserFriendsPostsResponse{
					Posts: []*postpb.PostResponse{
						{
							Id:       1,
							AuthorId: 1,
						},
					},
				}, nil)
				f.UserClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := &fields{
				PostsClient:       mock_posts.NewMockPostClient(ctrl),
				UserClient:        mock_user.NewMockUserClient(ctrl),
				PublicGroupClient: mock_public_group.NewMockPublicGroupClient(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			h := NewPostsHandler(f.PostsClient, f.UserClient, f.PublicGroupClient)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.HandleGetUserFriendsPosts)

			handler.ServeHTTP(rr, tt.request.WithContext(tt.ctx))

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestHandleDeletePost(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		request        *http.Request
		expectedStatus int
		prepare        func(f *fields)
	}{
		{
			name:           "TestHandleDeletePost",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			request:        httptest.NewRequest("DELETE", "/posts", bytes.NewBufferString(`{"postId": 1}`)),
			expectedStatus: http.StatusNoContent,
			prepare: func(f *fields) {
				f.PostsClient.EXPECT().DeletePost(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name:           "TestHandleDeletePost",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			request:        httptest.NewRequest("DELETE", "/posts", bytes.NewBufferString(`{"postID": 1}`)),
			expectedStatus: http.StatusInternalServerError,
			prepare: func(f *fields) {
				f.PostsClient.EXPECT().DeletePost(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
		{
			name:           "invalid context",
			ctx:            context.Background(),
			request:        httptest.NewRequest("DELETE", "/posts", bytes.NewBufferString(`{"postID": 1}`)),
			expectedStatus: http.StatusBadRequest,
			prepare: func(f *fields) {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := &fields{
				PostsClient:       mock_posts.NewMockPostClient(ctrl),
				UserClient:        mock_user.NewMockUserClient(ctrl),
				PublicGroupClient: mock_public_group.NewMockPublicGroupClient(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			h := NewPostsHandler(f.PostsClient, f.UserClient, f.PublicGroupClient)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.HandleDeletePost)

			handler.ServeHTTP(rr, tt.request.WithContext(tt.ctx))

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestHandleCreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		userID         int64
		content        string
		mockError      error
		expectedStatus int
		mock           func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient)
	}{
		{
			name:           "Successful post creation",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			content:        "Test content",
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {
				postsClient.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(&postpb.CreatePostResponse{
					Post: &postpb.PostResponse{
						Id: 1,
					},
				}, nil)
				userClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&uspb.GetByIDResponse{
					User: &uspb.UserResponse{
						Id: 1,
					},
				}, nil)
			},
		},
		{
			name:           "no user id in context",
			ctx:            context.Background(),
			content:        "Test content",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {
			},
		},
		{
			name:           "err internal",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			content:        "Test content",
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {
				postsClient.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
		{
			name:           "err internal user",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			content:        "Test content",
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {
				postsClient.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(&postpb.CreatePostResponse{
					Post: &postpb.PostResponse{
						Id: 1,
					},
				}, nil)
				userClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal.GRPCStatus().Err())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			_ = writer.WriteField("content", tt.content)
			_ = writer.Close()

			r := httptest.NewRequest("POST", "/", body)
			r.Header.Set("Content-Type", writer.FormDataContentType())
			r = r.WithContext(tt.ctx)

			rr := httptest.NewRecorder()

			mockPostsClient := mock_posts.NewMockPostClient(ctrl)
			mockUserClient := mock_user.NewMockUserClient(ctrl)
			publicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)

			tt.mock(mockPostsClient, mockUserClient, publicGroupClient)

			h := NewPostsHandler(mockPostsClient, mockUserClient, publicGroupClient)

			h.HandleCreatePost(rr, r)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleGetLikedPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		userID         int64
		lastLikeID     string
		postsAmount    string
		mockError      error
		expectedStatus int
		mock           func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient)
	}{
		{
			name:           "Successful get liked posts",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			lastLikeID:     "0",
			postsAmount:    "10",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {
				postsClient.EXPECT().GetLikedPosts(gomock.Any(), gomock.Any()).Return(&postpb.GetLikedPostsResponse{
					LikedPosts: []*postpb.LikedPostResponse{
						{
							Post: &postpb.PostResponse{},
							Like: &postpb.PostLikeResponse{},
						}},
				}, nil)
				userClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&uspb.GetByIDResponse{
					User: &uspb.UserResponse{
						Id: 1,
					},
				}, nil)
			},
		},
		{
			name:           "no user id in context",
			ctx:            context.Background(),
			userID:         1,
			lastLikeID:     "0",
			postsAmount:    "10",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "internal error",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			lastLikeID:     "0",
			postsAmount:    "10",
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {
				postsClient.EXPECT().GetLikedPosts(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
		{
			name:           "internal error user",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			lastLikeID:     "0",
			postsAmount:    "10",
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {
				postsClient.EXPECT().GetLikedPosts(gomock.Any(), gomock.Any()).Return(&postpb.GetLikedPostsResponse{
					LikedPosts: []*postpb.LikedPostResponse{
						{
							Post: &postpb.PostResponse{},
							Like: &postpb.PostLikeResponse{},
						}},
				}, nil)
				userClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
		{
			name:           "Successful get liked posts",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			lastLikeID:     "",
			postsAmount:    "",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {
				postsClient.EXPECT().GetLikedPosts(gomock.Any(), gomock.Any()).Return(&postpb.GetLikedPostsResponse{
					LikedPosts: []*postpb.LikedPostResponse{
						{
							Post: &postpb.PostResponse{},
							Like: &postpb.PostLikeResponse{},
						}},
				}, nil)
				userClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&uspb.GetByIDResponse{
					User: &uspb.UserResponse{
						Id: 1,
					},
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/?lastLikeId="+tt.lastLikeID+"&postsAmount="+tt.postsAmount, nil)
			r = r.WithContext(tt.ctx)

			rr := httptest.NewRecorder()

			mockPostsClient := mock_posts.NewMockPostClient(ctrl)
			mockUserClient := mock_user.NewMockUserClient(ctrl)
			publicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)

			tt.mock(mockPostsClient, mockUserClient, publicGroupClient)

			h := NewPostsHandler(mockPostsClient, mockUserClient, publicGroupClient)

			h.HandleGetLikedPosts(rr, r)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleLikePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		userID         uint64
		postID         uint64
		mockError      error
		expectedStatus int
		mock           func(postsClient *mock_posts.MockPostClient)
	}{
		{
			name:           "Successful post like",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			postID:         1,
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			mock: func(postsClient *mock_posts.MockPostClient) {
				postsClient.EXPECT().LikePost(gomock.Any(), gomock.Any()).Return(&postpb.LikePostResponse{
					Like: &postpb.PostLikeResponse{},
				}, nil)
			},
		},
		{
			name:           "no user id in context",
			ctx:            context.Background(),
			userID:         1,
			postID:         1,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(postsClient *mock_posts.MockPostClient) {

			},
		},
		{
			name:           "err internal",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			postID:         1,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient) {
				postsClient.EXPECT().LikePost(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := LikePostInput{
				PostID: uint(tt.postID),
			}
			inputBytes, _ := json.Marshal(input)
			r := httptest.NewRequest("POST", "/", bytes.NewBuffer(inputBytes))
			r = r.WithContext(tt.ctx)

			rr := httptest.NewRecorder()

			mockPostsClient := mock_posts.NewMockPostClient(ctrl)
			tt.mock(mockPostsClient)

			h := NewPostsHandler(mockPostsClient, nil, nil)

			h.HandleLikePost(rr, r)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleUnlikePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		userID         int64
		postID         int64
		mockError      error
		expectedStatus int
		mock           func(postsClient *mock_posts.MockPostClient)
	}{
		{
			name:           "Successful post unlike",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			postID:         1,
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			mock: func(postsClient *mock_posts.MockPostClient) {
				postsClient.EXPECT().UnlikePost(gomock.Any(), gomock.Any()).Return(
					&postpb.UnlikePostResponse{}, nil,
				)
			},
		},
		{
			name:           "no user id in context",
			ctx:            context.Background(),
			userID:         1,
			postID:         1,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(postsClient *mock_posts.MockPostClient) {

			},
		},
		{
			name:           "err internal",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			postID:         1,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient) {
				postsClient.EXPECT().UnlikePost(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := UnlikePostInput{
				PostID: uint(tt.postID),
			}
			inputBytes, _ := json.Marshal(input)
			r := httptest.NewRequest("POST", "/", bytes.NewBuffer(inputBytes))
			r = r.WithContext(tt.ctx)

			rr := httptest.NewRecorder()

			mockPostsClient := mock_posts.NewMockPostClient(ctrl)
			tt.mock(mockPostsClient)

			h := NewPostsHandler(mockPostsClient, nil, nil)

			h.HandleUnlikePost(rr, r)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleGetGroupPostsBySubscriptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		userID         int64
		lastPostID     string
		postsAmount    string
		mockError      error
		expectedStatus int
		mock           func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient)
	}{
		{
			name:           "Successful get group posts by subscriptions",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().GetBySubscriberID(gomock.Any(), gomock.Any()).Return(
					&pgpb.GetBySubscriberIDResponse{
						PublicGroups: []*pgpb.PublicGroupResponse{
							{
								Id: 1,
							},
						},
					}, nil,
				)
				postsClient.EXPECT().GetGroupPostsBySubscriptionIDs(gomock.Any(), gomock.Any()).Return(
					&postpb.GetGroupPostsBySubscriptionIDsResponse{
						Posts: []*postpb.PostResponse{
							{
								Id: 1,
							},
						},
					}, nil,
				)
				userClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&uspb.GetByIDResponse{
					User: &uspb.UserResponse{},
				}, nil)
			},
		},
		{
			name:           "no user id in context",
			ctx:            context.Background(),
			userID:         1,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "err internal",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().GetBySubscriberID(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
		{
			name:           "err",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().GetBySubscriberID(gomock.Any(), gomock.Any()).Return(
					&pgpb.GetBySubscriberIDResponse{
						PublicGroups: []*pgpb.PublicGroupResponse{
							{
								Id: 1,
							},
						},
					}, nil,
				)
				postsClient.EXPECT().GetGroupPostsBySubscriptionIDs(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
		{
			name:           "Successful get group posts by subscriptions",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().GetBySubscriberID(gomock.Any(), gomock.Any()).Return(
					&pgpb.GetBySubscriberIDResponse{
						PublicGroups: []*pgpb.PublicGroupResponse{
							{
								Id: 1,
							},
						},
					}, nil,
				)
				postsClient.EXPECT().GetGroupPostsBySubscriptionIDs(gomock.Any(), gomock.Any()).Return(
					&postpb.GetGroupPostsBySubscriptionIDsResponse{
						Posts: []*postpb.PostResponse{
							{
								Id: 1,
							},
						},
					}, nil,
				)
				userClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
		{
			name:           "Successful get group posts by subscriptions",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			lastPostID:     "0",
			postsAmount:    "10",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().GetBySubscriberID(gomock.Any(), gomock.Any()).Return(
					&pgpb.GetBySubscriberIDResponse{
						PublicGroups: []*pgpb.PublicGroupResponse{
							{
								Id: 1,
							},
						},
					}, nil,
				)
				postsClient.EXPECT().GetGroupPostsBySubscriptionIDs(gomock.Any(), gomock.Any()).Return(
					&postpb.GetGroupPostsBySubscriptionIDsResponse{
						Posts: []*postpb.PostResponse{
							{
								Id: 1,
							},
						},
					}, nil,
				)
				userClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&uspb.GetByIDResponse{
					User: &uspb.UserResponse{},
				}, nil)
			},
		},
		{
			name:           "Successful get group posts by subscriptions",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			lastPostID:     "asd",
			postsAmount:    "10",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful get group posts by subscriptions",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			lastPostID:     "10",
			postsAmount:    "asd",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(postsClient *mock_posts.MockPostClient, userClient *mock_user.MockUserClient, publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/?lastPostId="+tt.lastPostID+"&postsAmount="+tt.postsAmount, nil)
			r = r.WithContext(tt.ctx)

			rr := httptest.NewRecorder()

			mockPostsClient := mock_posts.NewMockPostClient(ctrl)
			userClient := mock_user.NewMockUserClient(ctrl)
			publicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)
			tt.mock(mockPostsClient, userClient, publicGroupClient)

			h := NewPostsHandler(mockPostsClient, userClient, publicGroupClient)

			h.HandleGetGroupPostsBySubscriptions(rr, r)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleGetPostsByGroupSubIDsAndUserSubIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		userID         int64
		lastPostID     string
		postsAmount    string
		mockError      error
		expectedStatus int
		mock           func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient)
	}{
		{
			name:           "Successful get posts by group subscription IDs and user subscription IDs",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			mock: func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {
				publicGroupClient.EXPECT().GetBySubscriberID(gomock.Any(), gomock.Any()).Return(
					&pgpb.GetBySubscriberIDResponse{
						PublicGroups: []*pgpb.PublicGroupResponse{
							{
								Id: 1,
							},
						},
					},
					nil,
				)
				userClient.EXPECT().GetSubscriptionIDs(gomock.Any(), gomock.Any()).Return(
					&uspb.GetSubscriptionIDsResponse{
						SubscriptionIds: []uint64{1},
					}, nil,
				)
				postsClient.EXPECT().GetPostsByGroupSubIDsAndUserSubIDs(gomock.Any(), gomock.Any()).Return(
					&postpb.GetPostsByGroupSubIDsAndUserSubIDsResponse{
						Posts: []*postpb.PostResponse{
							{
								Id:      1,
								GroupId: 1,
							},
						},
					}, nil,
				)
				userClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					&uspb.GetByIDResponse{
						User: &uspb.UserResponse{
							Id: 1,
						},
					}, nil,
				)
			},
		},
		{
			name:           "Successful get posts by group subscription IDs and user subscription IDs",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			lastPostID:     "asd",
			postsAmount:    "10",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {

			},
		},
		{
			name:           "Successful get posts by group subscription IDs and user subscription IDs",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			lastPostID:     "0",
			postsAmount:    "asd",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {

			},
		},
		{
			name:           "Successful get posts by group subscription IDs and user subscription IDs",
			ctx:            context.Background(),
			userID:         1,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {

			},
		},
		{
			name:           "Successful get posts by group subscription IDs and user subscription IDs",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {
				publicGroupClient.EXPECT().GetBySubscriberID(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
		{
			name:           "Successful get posts by group subscription IDs and user subscription IDs",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {
				publicGroupClient.EXPECT().GetBySubscriberID(gomock.Any(), gomock.Any()).Return(
					&pgpb.GetBySubscriberIDResponse{
						PublicGroups: []*pgpb.PublicGroupResponse{
							{
								Id: 1,
							},
						},
					},
					nil,
				)
				userClient.EXPECT().GetSubscriptionIDs(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
		{
			name:           "Successful get posts by group subscription IDs and user subscription IDs",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {
				publicGroupClient.EXPECT().GetBySubscriberID(gomock.Any(), gomock.Any()).Return(
					&pgpb.GetBySubscriberIDResponse{
						PublicGroups: []*pgpb.PublicGroupResponse{
							{
								Id: 1,
							},
						},
					},
					nil,
				)
				userClient.EXPECT().GetSubscriptionIDs(gomock.Any(), gomock.Any()).Return(
					&uspb.GetSubscriptionIDsResponse{
						SubscriptionIds: []uint64{1},
					}, nil,
				)
				postsClient.EXPECT().GetPostsByGroupSubIDsAndUserSubIDs(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
		{
			name:           "Successful get posts by group subscription IDs and user subscription IDs",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         1,
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {
				publicGroupClient.EXPECT().GetBySubscriberID(gomock.Any(), gomock.Any()).Return(
					&pgpb.GetBySubscriberIDResponse{
						PublicGroups: []*pgpb.PublicGroupResponse{
							{
								Id: 1,
							},
						},
					},
					nil,
				)
				userClient.EXPECT().GetSubscriptionIDs(gomock.Any(), gomock.Any()).Return(
					&uspb.GetSubscriptionIDsResponse{
						SubscriptionIds: []uint64{1},
					}, nil,
				)
				postsClient.EXPECT().GetPostsByGroupSubIDsAndUserSubIDs(gomock.Any(), gomock.Any()).Return(
					&postpb.GetPostsByGroupSubIDsAndUserSubIDsResponse{
						Posts: []*postpb.PostResponse{
							{
								Id:      1,
								GroupId: 1,
							},
						},
					}, nil,
				)
				userClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/?lastPostId="+tt.lastPostID+"&postsAmount="+tt.postsAmount, nil)
			r = r.WithContext(tt.ctx)

			rr := httptest.NewRecorder()

			mockPostsClient := mock_posts.NewMockPostClient(ctrl)
			mockPublicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)
			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockPostsClient, mockPublicGroupClient, mockUserClient)

			h := NewPostsHandler(mockPostsClient, mockUserClient, mockPublicGroupClient)

			h.HandleGetPostsByGroupSubIDsAndUserSubIDs(rr, r)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleGetNewPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		lastPostID     string
		postsAmount    string
		mockError      error
		expectedStatus int
		mock           func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient)
	}{
		{
			name:           "Successful get new posts",
			ctx:            context.Background(),
			lastPostID:     "",
			postsAmount:    "",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			mock: func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {
				postsClient.EXPECT().GetNewPosts(gomock.Any(), gomock.Any()).Return(
					&postpb.GetNewPostsResponse{
						Posts: []*postpb.PostResponse{
							{
								Id:      1,
								GroupId: 1,
							},
						},
					}, nil,
				)
				userClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					&uspb.GetByIDResponse{
						User: &uspb.UserResponse{
							Id: 1,
						},
					}, nil,
				)
				publicGroupClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					&pgpb.GetByIDResponse{
						PublicGroup: &pgpb.PublicGroupWithInfoResponse{
							PublicGroup: &pgpb.PublicGroupResponse{
								Id: 1,
							},
							IsSubscribed: true,
						},
					}, nil,
				)
			},
		},
		{
			name:           "Successful get new posts",
			ctx:            context.Background(),
			lastPostID:     "10",
			postsAmount:    "10",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			mock: func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {
				postsClient.EXPECT().GetNewPosts(gomock.Any(), gomock.Any()).Return(
					&postpb.GetNewPostsResponse{
						Posts: []*postpb.PostResponse{
							{
								Id:      1,
								GroupId: 1,
							},
						},
					}, nil,
				)
				userClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					&uspb.GetByIDResponse{
						User: &uspb.UserResponse{
							Id: 1,
						},
					}, nil,
				)
				publicGroupClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					&pgpb.GetByIDResponse{
						PublicGroup: &pgpb.PublicGroupWithInfoResponse{
							PublicGroup: &pgpb.PublicGroupResponse{
								Id: 1,
							},
							IsSubscribed: true,
						},
					}, nil,
				)
			},
		},
		{
			name:           "Successful get new posts",
			ctx:            context.Background(),
			lastPostID:     "asd",
			postsAmount:    "10",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {

			},
		},
		{
			name:           "Successful get new posts",
			ctx:            context.Background(),
			lastPostID:     "10",
			postsAmount:    "asd",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			mock: func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {

			},
		},
		{
			name:           "Successful get new posts",
			ctx:            context.Background(),
			lastPostID:     "",
			postsAmount:    "",
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {
				postsClient.EXPECT().GetNewPosts(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)

			},
		},
		{
			name:           "Successful get new posts",
			ctx:            context.Background(),
			lastPostID:     "",
			postsAmount:    "",
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {
				postsClient.EXPECT().GetNewPosts(gomock.Any(), gomock.Any()).Return(
					&postpb.GetNewPostsResponse{
						Posts: []*postpb.PostResponse{
							{
								Id:      1,
								GroupId: 1,
							},
						},
					}, nil,
				)
				userClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)

			},
		},
		{
			name:           "Successful get new posts",
			ctx:            context.Background(),
			lastPostID:     "",
			postsAmount:    "",
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			mock: func(postsClient *mock_posts.MockPostClient, publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {
				postsClient.EXPECT().GetNewPosts(gomock.Any(), gomock.Any()).Return(
					&postpb.GetNewPostsResponse{
						Posts: []*postpb.PostResponse{
							{
								Id:      1,
								GroupId: 1,
							},
						},
					}, nil,
				)
				userClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					&uspb.GetByIDResponse{
						User: &uspb.UserResponse{
							Id: 1,
						},
					}, nil,
				)
				publicGroupClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/?lastPostId="+tt.lastPostID+"&postsAmount="+tt.postsAmount, nil)
			r = r.WithContext(tt.ctx)

			rr := httptest.NewRecorder()

			mockPostsClient := mock_posts.NewMockPostClient(ctrl)
			mockPublicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)
			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockPostsClient, mockPublicGroupClient, mockUserClient)

			h := NewPostsHandler(mockPostsClient, mockUserClient, mockPublicGroupClient)

			h.HandleGetNewPosts(rr, r)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleGetCommentsByPostID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostsClient := mock_posts.NewMockPostClient(ctrl)
	mockUserClient := mock_user.NewMockUserClient(ctrl)

	tests := []struct {
		name         string
		postID       uint
		ctx          context.Context
		wantComments []*domain.CommentWithAuthor
		wantErr      error
		setup        func()
		muxVars      map[string]string
	}{
		{
			name:   "test case 1 - successful retrieval",
			postID: 1,
			wantComments: []*domain.CommentWithAuthor{
				{
					Comment: &domain.Comment{
						ID: 1,
					},
					Author: &domain.User{
						ID: 1,
					},
				},
			},
			wantErr: nil,
			ctx:     context.Background(),
			setup: func() {
				mockPostsClient.EXPECT().GetCommentsByPostID(gomock.Any(), gomock.Any()).Return(&postpb.GetCommentsByPostIDResponse{
					Comments: []*postpb.CommentResponse{
						{
							Id: 1,
						},
					},
				}, nil)
				mockUserClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&uspb.GetByIDResponse{
					User: &uspb.UserResponse{
						Id: 1,
					},
				}, nil)
			},
			muxVars: map[string]string{
				"postID": fmt.Sprintf("%d", 1),
			},
		},
		{
			name:         "test case 2",
			postID:       1,
			wantComments: nil,
			wantErr:      errors.ErrInvalidData,
			ctx:          context.Background(),
			setup: func() {

			},
			muxVars: map[string]string{},
		},
		{
			name:         "test case 3",
			postID:       1,
			wantComments: nil,
			wantErr:      errors.ErrInvalidData,
			ctx:          context.Background(),
			setup: func() {

			},
			muxVars: map[string]string{
				"postID": "asd",
			},
		},
		{
			name:         "test case 4",
			postID:       1,
			wantComments: nil,
			wantErr:      errors.ErrInvalidData,
			ctx:          context.Background(),
			setup: func() {
				mockPostsClient.EXPECT().GetCommentsByPostID(gomock.Any(), gomock.Any()).Return(&postpb.GetCommentsByPostIDResponse{
					Comments: []*postpb.CommentResponse{
						{
							Id: 1,
						},
					},
				}, nil)
				mockUserClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal.GRPCStatus().Err())
			},
			muxVars: map[string]string{
				"postID": fmt.Sprintf("%d", 1),
			},
		},
		{
			name:         "test case 2",
			postID:       1,
			wantComments: nil,
			wantErr:      errors.ErrInvalidData,
			ctx:          context.Background(),
			setup: func() {
				mockPostsClient.EXPECT().GetCommentsByPostID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal.GRPCStatus().Err())
			},
			muxVars: map[string]string{
				"postID": fmt.Sprintf("%d", 1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			h := NewPostsHandler(mockPostsClient, mockUserClient, nil)

			req, err := http.NewRequest("GET", "/{postID}/comments/", nil)
			if err != nil {
				t.Fatal(err)
			}

			req = req.WithContext(tt.ctx)

			req = mux.SetURLVars(req, tt.muxVars)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.HandleGetCommentsByPostID)

			handler.ServeHTTP(rr, req)

			if tt.wantErr != nil {
				assert.NotEqual(t, http.StatusOK, rr.Code)
			} else {
				assert.Equal(t, http.StatusOK, rr.Code)
			}
		})
	}
}
