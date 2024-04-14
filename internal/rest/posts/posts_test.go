package rest

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"socio/domain"
	"socio/errors"
	mock_posts "socio/mocks/usecase/posts"
	"socio/pkg/requestcontext"
	"socio/pkg/sanitizer"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
)

type fields struct {
	PostsStorage *mock_posts.MockPostsStorage
	UserStorage  *mock_posts.MockUserStorage
	Sanitizer    *sanitizer.Sanitizer
}

func TestHandleGetPostByID(t *testing.T) {
	tests := []struct {
		name           string
		request        *http.Request
		expectedStatus int
		prepare        func(f *fields)
	}{
		{
			name:           "TestHandleGetPostByID",
			request:        httptest.NewRequest("GET", "/posts/1", nil),
			expectedStatus: http.StatusOK,
			prepare: func(f *fields) {
				f.PostsStorage.EXPECT().GetPostByID(gomock.Any(), gomock.Any()).Return(&domain.Post{}, nil)
			},
		},
		{
			name:           "TestHandleGetPostByID",
			request:        httptest.NewRequest("GET", "/posts/", nil),
			expectedStatus: http.StatusNotFound,
			prepare: func(f *fields) {
			},
		},
		{
			name:           "TestHandleGetPostByID",
			request:        httptest.NewRequest("GET", "/posts/asd", nil),
			expectedStatus: http.StatusBadRequest,
			prepare: func(f *fields) {
			},
		},
		{
			name:           "TestHandleGetPostByID",
			request:        httptest.NewRequest("GET", "/posts/1", nil),
			expectedStatus: http.StatusNotFound,
			prepare: func(f *fields) {
				f.PostsStorage.EXPECT().GetPostByID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := &fields{
				PostsStorage: mock_posts.NewMockPostsStorage(ctrl),
				UserStorage:  mock_posts.NewMockUserStorage(ctrl),
				Sanitizer:    sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			h := NewPostsHandler(f.PostsStorage, f.UserStorage, f.Sanitizer)

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/posts/{postID}", h.HandleGetPostByID)

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
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&domain.User{}, nil)
				f.PostsStorage.EXPECT().GetUserPosts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name:           "no user",
			request:        httptest.NewRequest("GET", "/posts?userId=1&lastPostId=0&postsAmount=20", nil),
			expectedStatus: http.StatusNotFound,
			prepare: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := &fields{
				PostsStorage: mock_posts.NewMockPostsStorage(ctrl),
				UserStorage:  mock_posts.NewMockUserStorage(ctrl),
				Sanitizer:    sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			h := NewPostsHandler(f.PostsStorage, f.UserStorage, f.Sanitizer)

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
				f.PostsStorage.EXPECT().GetUserFriendsPosts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
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
				f.PostsStorage.EXPECT().GetUserFriendsPosts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := &fields{
				PostsStorage: mock_posts.NewMockPostsStorage(ctrl),
				UserStorage:  mock_posts.NewMockUserStorage(ctrl),
				Sanitizer:    sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			h := NewPostsHandler(f.PostsStorage, f.UserStorage, f.Sanitizer)

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
		request        *http.Request
		expectedStatus int
		prepare        func(f *fields)
	}{
		{
			name:           "TestHandleDeletePost",
			request:        httptest.NewRequest("DELETE", "/posts", bytes.NewBufferString(`{"postID": 1}`)),
			expectedStatus: http.StatusNoContent,
			prepare: func(f *fields) {
				f.PostsStorage.EXPECT().DeletePost(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:           "TestHandleDeletePost",
			request:        httptest.NewRequest("DELETE", "/posts", bytes.NewBufferString(`{"postID": 1}`)),
			expectedStatus: http.StatusInternalServerError,
			prepare: func(f *fields) {
				f.PostsStorage.EXPECT().DeletePost(gomock.Any(), gomock.Any()).Return(errors.ErrInternal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := &fields{
				PostsStorage: mock_posts.NewMockPostsStorage(ctrl),
				UserStorage:  mock_posts.NewMockUserStorage(ctrl),
				Sanitizer:    sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			h := NewPostsHandler(f.PostsStorage, f.UserStorage, f.Sanitizer)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.HandleDeletePost)

			handler.ServeHTTP(rr, tt.request)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
