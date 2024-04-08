package posts_test

import (
	"context"
	"mime/multipart"
	"reflect"
	"socio/domain"
	"socio/errors"
	mock_posts "socio/mocks/usecase/posts"
	"socio/pkg/sanitizer"
	customtime "socio/pkg/time"
	"socio/usecase/posts"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/microcosm-cc/bluemonday"
)

type fields struct {
	PostsStorage *mock_posts.MockPostsStorage
	UserStorage  *mock_posts.MockUserStorage
	Sanitizer    *sanitizer.Sanitizer
}

var timeProv = customtime.MockTimeProvider{}

func TestService_GetUserPosts(t *testing.T) {
	type args struct {
		ctx        context.Context
		userID     uint
		lastPostID uint
	}

	tests := []struct {
		name        string
		args        args
		wantPosts   []*domain.Post
		wantAuthor  *domain.User
		wantErr     bool
		prepareMock func(*fields)
	}{
		{
			name: "success",
			args: args{
				ctx:        context.Background(),
				userID:     1,
				lastPostID: 0,
			},
			wantPosts: []*domain.Post{
				{
					ID:       1,
					AuthorID: 1,
					Content:  "content",
					CreatedAt: customtime.CustomTime{
						Time: timeProv.Now(),
					},
					UpdatedAt: customtime.CustomTime{
						Time: timeProv.Now(),
					},
				},
			},
			wantAuthor: &domain.User{
				ID:        1,
				FirstName: "firstName",
				LastName:  "lastName",
				Email:     "email",
				Avatar:    "avatar",
				CreatedAt: customtime.CustomTime{
					Time: timeProv.Now(),
				},
				UpdatedAt: customtime.CustomTime{
					Time: timeProv.Now(),
				},
			},
			wantErr: false,
			prepareMock: func(f *fields) {
				f.PostsStorage.EXPECT().GetUserPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*domain.Post{
					{
						ID:       1,
						AuthorID: 1,
						Content:  "content",
						CreatedAt: customtime.CustomTime{
							Time: timeProv.Now(),
						},
						UpdatedAt: customtime.CustomTime{
							Time: timeProv.Now(),
						},
					},
				}, nil)
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&domain.User{
					ID:        1,
					FirstName: "firstName",
					LastName:  "lastName",
					Email:     "email",
					Avatar:    "avatar",
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
			name: "no user",
			args: args{
				ctx:        context.Background(),
				userID:     1,
				lastPostID: 0,
			},
			wantPosts:  nil,
			wantAuthor: nil,
			wantErr:    true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name: "no posts",
			args: args{
				ctx:        context.Background(),
				userID:     1,
				lastPostID: 0,
			},
			wantPosts: nil,
			wantAuthor: &domain.User{
				ID:        1,
				FirstName: "firstName",
				LastName:  "lastName",
				Email:     "email",
				Avatar:    "avatar",
				CreatedAt: customtime.CustomTime{
					Time: timeProv.Now(),
				},
				UpdatedAt: customtime.CustomTime{
					Time: timeProv.Now(),
				},
			},
			wantErr: true,
			prepareMock: func(f *fields) {
				f.PostsStorage.EXPECT().GetUserPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&domain.User{
					ID:        1,
					FirstName: "firstName",
					LastName:  "lastName",
					Email:     "email",
					Avatar:    "avatar",
					CreatedAt: customtime.CustomTime{
						Time: timeProv.Now(),
					},
					UpdatedAt: customtime.CustomTime{
						Time: timeProv.Now(),
					},
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				PostsStorage: mock_posts.NewMockPostsStorage(ctrl),
				UserStorage:  mock_posts.NewMockUserStorage(ctrl),
				Sanitizer:    sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
			}

			if tt.prepareMock != nil {
				tt.prepareMock(&f)
			}

			s := posts.NewPostsService(f.PostsStorage, f.UserStorage, f.Sanitizer)

			gotPosts, gotAuthor, err := s.GetUserPosts(tt.args.ctx, tt.args.userID, tt.args.lastPostID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetUserPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPosts, tt.wantPosts) {
				t.Errorf("Service.GetUserPosts() gotPosts = %v, want %v", gotPosts, tt.wantPosts)
			}
			if !reflect.DeepEqual(gotAuthor, tt.wantAuthor) {
				t.Errorf("Service.GetUserPosts() gotAuthor = %v, want %v", gotAuthor, tt.wantAuthor)
			}
		})
	}
}

func TestService_GetUserFriendsPosts(t *testing.T) {
	type args struct {
		ctx        context.Context
		userID     uint
		lastPostID uint
	}
	tests := []struct {
		name        string
		args        args
		wantPosts   []domain.PostWithAuthor
		wantErr     bool
		prepareMock func(*fields)
	}{
		{
			name: "success",
			args: args{
				ctx:        context.Background(),
				userID:     1,
				lastPostID: 0,
			},
			wantPosts: []domain.PostWithAuthor{
				{
					Post: &domain.Post{
						ID:       1,
						AuthorID: 1,
						Content:  "content",
						CreatedAt: customtime.CustomTime{
							Time: timeProv.Now(),
						},
						UpdatedAt: customtime.CustomTime{
							Time: timeProv.Now(),
						},
					},
					Author: &domain.User{
						ID:        1,
						FirstName: "firstName",
						LastName:  "lastName",
						Email:     "email",
						Avatar:    "avatar",
						CreatedAt: customtime.CustomTime{
							Time: timeProv.Now(),
						},
						UpdatedAt: customtime.CustomTime{
							Time: timeProv.Now(),
						},
					},
				},
			},
			wantErr: false,
			prepareMock: func(f *fields) {
				f.PostsStorage.EXPECT().GetUserFriendsPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.PostWithAuthor{
					{
						Post: &domain.Post{
							ID:       1,
							AuthorID: 1,
							Content:  "content",
							CreatedAt: customtime.CustomTime{
								Time: timeProv.Now(),
							},
							UpdatedAt: customtime.CustomTime{
								Time: timeProv.Now(),
							},
						},
						Author: &domain.User{
							ID:        1,
							FirstName: "firstName",
							LastName:  "lastName",
							Email:     "email",
							Avatar:    "avatar",
							CreatedAt: customtime.CustomTime{
								Time: timeProv.Now(),
							},
							UpdatedAt: customtime.CustomTime{
								Time: timeProv.Now(),
							},
						},
					},
				}, nil)
			},
		},
		{
			name: "no posts",
			args: args{
				ctx:        context.Background(),
				userID:     1,
				lastPostID: 0,
			},
			wantPosts: nil,
			wantErr:   true,
			prepareMock: func(f *fields) {
				f.PostsStorage.EXPECT().GetUserFriendsPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				PostsStorage: mock_posts.NewMockPostsStorage(ctrl),
				UserStorage:  mock_posts.NewMockUserStorage(ctrl),
				Sanitizer:    sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
			}

			if tt.prepareMock != nil {
				tt.prepareMock(&f)
			}

			s := posts.NewPostsService(f.PostsStorage, f.UserStorage, f.Sanitizer)

			gotPosts, err := s.GetUserFriendsPosts(tt.args.ctx, tt.args.userID, tt.args.lastPostID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetUserFriendsPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPosts, tt.wantPosts) {
				t.Errorf("Service.GetUserFriendsPosts() = %v, want %v", gotPosts, tt.wantPosts)
			}
		})
	}
}

func TestService_CreatePost(t *testing.T) {
	type args struct {
		ctx   context.Context
		input posts.PostInput
	}

	tests := []struct {
		name               string
		args               args
		wantPostWithAuthor domain.PostWithAuthor
		wantErr            bool
		prepareMock        func(*fields)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				input: posts.PostInput{
					AuthorID:    1,
					Content:     "content",
					Attachments: []*multipart.FileHeader{},
				},
			},
			wantPostWithAuthor: domain.PostWithAuthor{
				Post: &domain.Post{
					ID:          1,
					AuthorID:    1,
					Content:     "content",
					Attachments: []string{},
				},
				Author: &domain.User{
					ID:        1,
					FirstName: "firstName",
					LastName:  "lastName",
					Email:     "email",
					Avatar:    "avatar",
					DateOfBirth: customtime.CustomTime{
						Time: timeProv.Now(),
					},
					CreatedAt: customtime.CustomTime{
						Time: timeProv.Now(),
					},
					UpdatedAt: customtime.CustomTime{
						Time: timeProv.Now(),
					},
				},
			},
			wantErr: false,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&domain.User{
					ID:        1,
					FirstName: "firstName",
					LastName:  "lastName",
					Email:     "email",
					Avatar:    "avatar",
					DateOfBirth: customtime.CustomTime{
						Time: timeProv.Now(),
					},
					CreatedAt: customtime.CustomTime{
						Time: timeProv.Now(),
					},
					UpdatedAt: customtime.CustomTime{
						Time: timeProv.Now(),
					},
				}, nil)
				f.PostsStorage.EXPECT().StorePost(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Post{
					ID:          1,
					AuthorID:    1,
					Content:     "content",
					Attachments: []string{},
				}, nil)
			},
		},
		{
			name: "empty input",
			args: args{
				ctx: context.Background(),
				input: posts.PostInput{
					AuthorID:    1,
					Content:     "",
					Attachments: []*multipart.FileHeader{},
				},
			},
			wantPostWithAuthor: domain.PostWithAuthor{
				Post:   nil,
				Author: nil,
			},
			wantErr: true,
			prepareMock: func(f *fields) {
			},
		},
		{
			name: "no user",
			args: args{
				ctx: context.Background(),
				input: posts.PostInput{
					AuthorID:    1,
					Content:     "content",
					Attachments: []*multipart.FileHeader{},
				},
			},
			wantPostWithAuthor: domain.PostWithAuthor{
				Post:   nil,
				Author: nil,
			},
			wantErr: true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name: "error storing post",
			args: args{
				ctx: context.Background(),
				input: posts.PostInput{
					AuthorID:    1,
					Content:     "content",
					Attachments: []*multipart.FileHeader{},
				},
			},
			wantPostWithAuthor: domain.PostWithAuthor{
				Post:   nil,
				Author: nil,
			},
			wantErr: true,
			prepareMock: func(f *fields) {
				f.UserStorage.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&domain.User{
					ID:        1,
					FirstName: "firstName",
					LastName:  "lastName",
					Email:     "email",
					Avatar:    "avatar",
					DateOfBirth: customtime.CustomTime{
						Time: timeProv.Now(),
					},
					CreatedAt: customtime.CustomTime{
						Time: timeProv.Now(),
					},
					UpdatedAt: customtime.CustomTime{
						Time: timeProv.Now(),
					},
				}, nil)
				f.PostsStorage.EXPECT().StorePost(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				PostsStorage: mock_posts.NewMockPostsStorage(ctrl),
				UserStorage:  mock_posts.NewMockUserStorage(ctrl),
				Sanitizer:    sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
			}

			if tt.prepareMock != nil {
				tt.prepareMock(&f)
			}

			s := posts.NewPostsService(f.PostsStorage, f.UserStorage, f.Sanitizer)

			gotPostWithAuthor, err := s.CreatePost(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPostWithAuthor, tt.wantPostWithAuthor) {
				t.Errorf("Service.CreatePost() = %v, want %v", gotPostWithAuthor, tt.wantPostWithAuthor)
			}
		})
	}
}

func TestService_UpdatePost(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID uint
		input  posts.PostUpdateInput
	}

	tests := []struct {
		name        string
		args        args
		wantPost    *domain.Post
		wantErr     bool
		prepareMock func(*fields)
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				input: posts.PostUpdateInput{
					PostID:  1,
					Content: "newContent",
				},
			},
			wantPost: &domain.Post{
				ID:       1,
				AuthorID: 1,
				Content:  "newContent",
			},
			wantErr: false,
			prepareMock: func(f *fields) {
				f.PostsStorage.EXPECT().GetPostByID(gomock.Any(), gomock.Any()).Return(&domain.Post{
					ID:       1,
					AuthorID: 1,
					Content:  "content",
				}, nil)
				f.PostsStorage.EXPECT().UpdatePost(gomock.Any(), gomock.Any()).Return(&domain.Post{
					ID:       1,
					AuthorID: 1,
					Content:  "newContent",
				}, nil)
			},
		},
		{
			name: "error no post",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				input: posts.PostUpdateInput{
					PostID:  1,
					Content: "newContent",
				},
			},
			wantPost: nil,
			wantErr:  true,
			prepareMock: func(f *fields) {
				f.PostsStorage.EXPECT().GetPostByID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name: "error forbidden",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				input: posts.PostUpdateInput{
					PostID:  1,
					Content: "newContent",
				},
			},
			wantPost: nil,
			wantErr:  true,
			prepareMock: func(f *fields) {
				f.PostsStorage.EXPECT().GetPostByID(gomock.Any(), gomock.Any()).Return(&domain.Post{
					ID:       1,
					AuthorID: 2,
					Content:  "content",
				}, nil)
			},
		},
		{
			name: "error empty content",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				input: posts.PostUpdateInput{
					PostID:  1,
					Content: "",
				},
			},
			wantPost: nil,
			wantErr:  true,
			prepareMock: func(f *fields) {
				f.PostsStorage.EXPECT().GetPostByID(gomock.Any(), gomock.Any()).Return(&domain.Post{
					ID:       1,
					AuthorID: 1,
					Content:  "content",
				}, nil)
			},
		},
		{
			name: "error updating post",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				input: posts.PostUpdateInput{
					PostID:  1,
					Content: "newContent",
				},
			},
			wantPost: nil,
			wantErr:  true,
			prepareMock: func(f *fields) {
				f.PostsStorage.EXPECT().GetPostByID(gomock.Any(), gomock.Any()).Return(&domain.Post{
					ID:       1,
					AuthorID: 1,
					Content:  "content",
				}, nil)
				f.PostsStorage.EXPECT().UpdatePost(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				PostsStorage: mock_posts.NewMockPostsStorage(ctrl),
				UserStorage:  mock_posts.NewMockUserStorage(ctrl),
				Sanitizer:    sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
			}

			if tt.prepareMock != nil {
				tt.prepareMock(&f)
			}

			s := posts.NewPostsService(f.PostsStorage, f.UserStorage, f.Sanitizer)

			gotPost, err := s.UpdatePost(tt.args.ctx, tt.args.userID, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPost, tt.wantPost) {
				t.Errorf("Service.UpdatePost() = %v, want %v", gotPost, tt.wantPost)
			}
		})
	}
}

func TestService_DeletePost(t *testing.T) {
	type args struct {
		ctx    context.Context
		postID uint
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
				ctx:    context.Background(),
				postID: 1,
			},
			wantErr: false,
			prepareMock: func(f *fields) {
				f.PostsStorage.EXPECT().DeletePost(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "error deleting post",
			args: args{
				ctx:    context.Background(),
				postID: 1,
			},
			wantErr: true,
			prepareMock: func(f *fields) {
				f.PostsStorage.EXPECT().DeletePost(gomock.Any(), gomock.Any()).Return(errors.ErrInternal)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				PostsStorage: mock_posts.NewMockPostsStorage(ctrl),
				UserStorage:  mock_posts.NewMockUserStorage(ctrl),
				Sanitizer:    sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
			}

			if tt.prepareMock != nil {
				tt.prepareMock(&f)
			}

			s := posts.NewPostsService(f.PostsStorage, f.UserStorage, f.Sanitizer)

			if err := s.DeletePost(tt.args.ctx, tt.args.postID); (err != nil) != tt.wantErr {
				t.Errorf("Service.DeletePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
