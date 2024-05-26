package posts_test

import (
	"context"
	"reflect"
	"socio/domain"
	"socio/errors"
	mock_posts "socio/mocks/usecase/posts"
	"socio/usecase/posts"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetPostByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		postID   uint
		mock     func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, postID uint)
		wantPost *domain.Post
		wantErr  bool
	}{
		{
			name:   "Test OK",
			postID: 1,
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, postID uint) {
				post := &domain.Post{ID: postID, Content: "Test Content"}
				postsStorage.EXPECT().GetPostByID(gomock.Any(), postID).Return(post, nil)
			},
			wantPost: &domain.Post{ID: 1, Content: "Test Content"},
			wantErr:  false,
		},
		{
			name:   "Test Error",
			postID: 1,
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, postID uint) {
				postsStorage.EXPECT().GetPostByID(gomock.Any(), postID).Return(nil, errors.ErrInternal)
			},
			wantPost: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postsStorage := mock_posts.NewMockPostsStorage(ctrl)
			attachmentStorage := mock_posts.NewMockAttachmentStorage(ctrl)

			s := posts.NewPostsService(postsStorage, attachmentStorage)

			tt.mock(postsStorage, attachmentStorage, tt.postID)

			gotPost, err := s.GetPostByID(context.Background(), tt.postID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPostByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotPost, tt.wantPost) {
				t.Errorf("GetPostByID() gotPost = %v, want %v", gotPost, tt.wantPost)
			}
		})
	}
}

func TestGetUserPosts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		userID      uint
		lastPostID  uint
		postsAmount uint
		mock        func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, lastPostID uint, postsAmount uint)
		wantPosts   []*domain.Post
		wantErr     bool
	}{
		{
			name:        "Test OK",
			userID:      1,
			lastPostID:  0,
			postsAmount: 0,
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, lastPostID uint, postsAmount uint) {
				testPosts := []*domain.Post{{ID: 1, Content: "Test Content"}}
				postsStorage.EXPECT().GetUserPosts(gomock.Any(), userID, lastPostID, posts.DefaultPostsAmount).Return(testPosts, nil)
			},
			wantPosts: []*domain.Post{{ID: 1, Content: "Test Content"}},
			wantErr:   false,
		},
		{
			name:        "Test Error",
			userID:      1,
			lastPostID:  0,
			postsAmount: 0,
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, lastPostID uint, postsAmount uint) {
				postsStorage.EXPECT().GetUserPosts(gomock.Any(), userID, lastPostID, posts.DefaultPostsAmount).Return(nil, errors.ErrInternal)
			},
			wantPosts: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postsStorage := mock_posts.NewMockPostsStorage(ctrl)
			attachmentStorage := mock_posts.NewMockAttachmentStorage(ctrl)

			s := posts.NewPostsService(postsStorage, attachmentStorage)

			tt.mock(postsStorage, attachmentStorage, tt.userID, tt.lastPostID, tt.postsAmount)

			gotPosts, err := s.GetUserPosts(context.Background(), tt.userID, tt.lastPostID, tt.postsAmount)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserPosts() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotPosts, tt.wantPosts) {
				t.Errorf("GetUserPosts() gotPosts = %v, want %v", gotPosts, tt.wantPosts)
			}
		})
	}
}

func TestCreatePost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    posts.PostInput
		mock     func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, input posts.PostInput)
		wantPost *domain.Post
		wantErr  bool
	}{
		{
			name: "Test OK",
			input: posts.PostInput{
				AuthorID:    1,
				Content:     "Test Content",
				Attachments: []string{"attachment1", "attachment2"},
			},
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, input posts.PostInput) {
				post := &domain.Post{AuthorID: input.AuthorID, Content: input.Content, Attachments: input.Attachments}
				postsStorage.EXPECT().StorePost(gomock.Any(), post).Return(post, nil)
			},
			wantPost: &domain.Post{AuthorID: 1, Content: "Test Content", Attachments: []string{"attachment1", "attachment2"}},
			wantErr:  false,
		},
		{
			name: "Test Error",
			input: posts.PostInput{
				AuthorID:    1,
				Content:     "",
				Attachments: []string{},
			},
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, input posts.PostInput) {
			},
			wantPost: nil,
			wantErr:  true,
		},
		{
			name: "Test err internal",
			input: posts.PostInput{
				AuthorID:    1,
				Content:     "Test Content",
				Attachments: []string{"attachment1", "attachment2"},
			},
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, input posts.PostInput) {
				postsStorage.EXPECT().StorePost(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
			wantPost: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postsStorage := mock_posts.NewMockPostsStorage(ctrl)
			attachmentStorage := mock_posts.NewMockAttachmentStorage(ctrl)

			s := posts.NewPostsService(postsStorage, attachmentStorage)

			tt.mock(postsStorage, attachmentStorage, tt.input)

			gotPost, err := s.CreatePost(context.Background(), tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePost() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotPost, tt.wantPost) {
				t.Errorf("CreatePost() gotPost = %v, want %v", gotPost, tt.wantPost)
			}
		})
	}
}

func TestUpdatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostsStorage := mock_posts.NewMockPostsStorage(ctrl)
	mockAttachmentStorage := mock_posts.NewMockAttachmentStorage(ctrl)

	tests := []struct {
		name    string
		userID  uint
		input   posts.PostUpdateInput
		want    *domain.Post
		wantErr error
		setup   func()
	}{
		{
			name:   "test case 1 - successful update",
			userID: 1,
			input: posts.PostUpdateInput{
				PostID:              1,
				Content:             "Updated content",
				AttachmentsToAdd:    []string{"new_attachment.jpg"},
				AttachmentsToDelete: []string{"old_attachment.jpg"},
			},
			want: &domain.Post{
				ID:          1,
				AuthorID:    1,
				Content:     "Updated content",
				Attachments: []string{"new_attachment.jpg"},
			},
			wantErr: nil,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(&domain.Post{
					ID:          1,
					AuthorID:    1,
					Content:     "Old content",
					Attachments: []string{"old_attachment.jpg"},
				}, nil).Times(2)

				mockPostsStorage.EXPECT().UpdatePost(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Post{
					ID:          1,
					AuthorID:    1,
					Content:     "Updated content",
					Attachments: []string{"new_attachment.jpg"},
				}, nil)

				mockAttachmentStorage.EXPECT().Delete("old_attachment.jpg").Return(nil)
			},
		},
		{
			name:   "test case 2",
			userID: 1,
			input: posts.PostUpdateInput{
				PostID:              1,
				Content:             "Updated content",
				AttachmentsToAdd:    []string{"new_attachment.jpg"},
				AttachmentsToDelete: []string{"old_attachment.jpg"},
			},
			want:    &domain.Post{},
			wantErr: errors.ErrNotFound,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(&domain.Post{
					ID:          1,
					AuthorID:    1,
					Content:     "Old content",
					Attachments: []string{"old_attachment.jpg"},
				}, nil)

				mockPostsStorage.EXPECT().UpdatePost(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Post{
					ID:          1,
					AuthorID:    1,
					Content:     "Updated content",
					Attachments: []string{"new_attachment.jpg"},
				}, nil)

				mockAttachmentStorage.EXPECT().Delete("old_attachment.jpg").Return(nil)

				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name:   "test case 2",
			userID: 1,
			input: posts.PostUpdateInput{
				PostID:              1,
				Content:             "Updated content",
				AttachmentsToAdd:    []string{"new_attachment.jpg"},
				AttachmentsToDelete: []string{"old_attachment.jpg"},
			},
			want:    nil,
			wantErr: errors.ErrNotFound,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(&domain.Post{
					ID:          1,
					AuthorID:    1,
					Content:     "Old content",
					Attachments: []string{"old_attachment.jpg"},
				}, nil)

				mockPostsStorage.EXPECT().UpdatePost(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Post{
					ID:          1,
					AuthorID:    1,
					Content:     "Updated content",
					Attachments: []string{"new_attachment.jpg"},
				}, nil)

				mockAttachmentStorage.EXPECT().Delete("old_attachment.jpg").Return(errors.ErrNotFound)
			},
		},
		{
			name:   "test case 3",
			userID: 1,
			input: posts.PostUpdateInput{
				PostID:              1,
				Content:             "Updated content",
				AttachmentsToAdd:    []string{"new_attachment.jpg"},
				AttachmentsToDelete: []string{"old_attachment.jpg"},
			},
			want:    nil,
			wantErr: errors.ErrNotFound,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(&domain.Post{
					ID:          1,
					AuthorID:    1,
					Content:     "Old content",
					Attachments: []string{"old_attachment.jpg"},
				}, nil)

				mockPostsStorage.EXPECT().UpdatePost(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name:   "test case 4",
			userID: 1,
			input: posts.PostUpdateInput{
				PostID:              1,
				Content:             "Updated content",
				AttachmentsToAdd:    []string{"new_attachment.jpg"},
				AttachmentsToDelete: []string{"old_attachment.jpg"},
			},
			want:    nil,
			wantErr: errors.ErrNotFound,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name:   "test case 5",
			userID: 1,
			input: posts.PostUpdateInput{
				PostID:              1,
				Content:             "",
				AttachmentsToAdd:    []string{},
				AttachmentsToDelete: []string{},
			},
			want:    nil,
			wantErr: errors.ErrInvalidBody,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(&domain.Post{
					ID:          1,
					AuthorID:    1,
					Content:     "Old content",
					Attachments: []string{},
				}, nil)
			},
		},
		{
			name:   "test case 6",
			userID: 1,
			input: posts.PostUpdateInput{
				PostID:              1,
				Content:             "Updated content",
				AttachmentsToAdd:    []string{"new_attachment.jpg"},
				AttachmentsToDelete: []string{"old_attachment.jpg"},
			},
			want:    nil,
			wantErr: errors.ErrForbidden,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(&domain.Post{
					ID:          1,
					AuthorID:    2,
					Content:     "Old content",
					Attachments: []string{"old_attachment.jpg"},
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := posts.NewPostsService(mockPostsStorage, mockAttachmentStorage)

			got, err := s.UpdatePost(context.Background(), tt.userID, tt.input)
			assert.Equal(t, tt.wantErr, err)

			if tt.wantErr == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetUserFriendsPosts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		userID      uint
		lastPostID  uint
		postsAmount uint
		mock        func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, lastPostID uint, postsAmount uint)
		wantPosts   []*domain.Post
		wantErr     bool
	}{
		{
			name:        "Test OK",
			userID:      1,
			lastPostID:  0,
			postsAmount: 0,
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, lastPostID uint, postsAmount uint) {
				testPosts := []*domain.Post{{ID: 1, Content: "Test Content"}}
				postsStorage.EXPECT().GetUserFriendsPosts(gomock.Any(), userID, lastPostID, posts.DefaultPostsAmount).Return(testPosts, nil)
			},
			wantPosts: []*domain.Post{{ID: 1, Content: "Test Content"}},
			wantErr:   false,
		},
		{
			name:        "Test Error",
			userID:      1,
			lastPostID:  0,
			postsAmount: 0,
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, lastPostID uint, postsAmount uint) {
				postsStorage.EXPECT().GetUserFriendsPosts(gomock.Any(), userID, lastPostID, posts.DefaultPostsAmount).Return(nil, errors.ErrInternal)
			},
			wantPosts: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postsStorage := mock_posts.NewMockPostsStorage(ctrl)
			attachmentStorage := mock_posts.NewMockAttachmentStorage(ctrl)

			s := posts.NewPostsService(postsStorage, attachmentStorage)

			tt.mock(postsStorage, attachmentStorage, tt.userID, tt.lastPostID, tt.postsAmount)

			gotPosts, err := s.GetUserFriendsPosts(context.Background(), tt.userID, tt.lastPostID, tt.postsAmount)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserFriendsPosts() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotPosts, tt.wantPosts) {
				t.Errorf("GetUserFriendsPosts() gotPosts = %v, want %v", gotPosts, tt.wantPosts)
			}
		})
	}
}
func TestDeletePost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userID  uint
		postID  uint
		mock    func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, postID uint)
		wantErr bool
	}{
		{
			name:   "Test OK",
			userID: 1,
			postID: 1,
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, postID uint) {
				post := &domain.Post{ID: postID, AuthorID: userID, Attachments: []string{"attachment1", "attachment2"}}
				postsStorage.EXPECT().GetPostByID(gomock.Any(), postID).Return(post, nil)
				attachmentStorage.EXPECT().Delete(gomock.Any()).Times(len(post.Attachments))
				postsStorage.EXPECT().DeleteGroupPost(gomock.Any(), postID).Return(nil)
				postsStorage.EXPECT().DeletePost(gomock.Any(), postID).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Test Error",
			userID: 1,
			postID: 1,
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, postID uint) {
				post := &domain.Post{ID: postID, AuthorID: userID + 1}
				postsStorage.EXPECT().GetPostByID(gomock.Any(), postID).Return(post, nil)
			},
			wantErr: true,
		},
		{
			name:   "Test not found",
			userID: 1,
			postID: 1,
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, postID uint) {
				postsStorage.EXPECT().GetPostByID(gomock.Any(), postID).Return(nil, errors.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name:   "Test err deleting attachments",
			userID: 1,
			postID: 1,
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, postID uint) {
				post := &domain.Post{ID: postID, AuthorID: userID, Attachments: []string{"attachment1", "attachment2"}}
				postsStorage.EXPECT().GetPostByID(gomock.Any(), postID).Return(post, nil)
				attachmentStorage.EXPECT().Delete(gomock.Any()).Return(errors.ErrInternal)
			},
			wantErr: true,
		},
		{
			name:   "Test err deleting group post",
			userID: 1,
			postID: 1,
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, postID uint) {
				post := &domain.Post{ID: postID, AuthorID: userID, Attachments: []string{"attachment1", "attachment2"}}
				postsStorage.EXPECT().GetPostByID(gomock.Any(), postID).Return(post, nil)
				attachmentStorage.EXPECT().Delete(gomock.Any()).Times(len(post.Attachments))
				postsStorage.EXPECT().DeleteGroupPost(gomock.Any(), postID).Return(errors.ErrInternal)
			},
			wantErr: true,
		},
		{
			name:   "Test err internal",
			userID: 1,
			postID: 1,
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, postID uint) {
				post := &domain.Post{ID: postID, AuthorID: userID, Attachments: []string{"attachment1", "attachment2"}}
				postsStorage.EXPECT().GetPostByID(gomock.Any(), postID).Return(post, nil)
				attachmentStorage.EXPECT().Delete(gomock.Any()).Times(len(post.Attachments))
				postsStorage.EXPECT().DeleteGroupPost(gomock.Any(), postID).Return(nil)
				postsStorage.EXPECT().DeletePost(gomock.Any(), postID).Return(errors.ErrInternal)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postsStorage := mock_posts.NewMockPostsStorage(ctrl)
			attachmentStorage := mock_posts.NewMockAttachmentStorage(ctrl)

			s := posts.NewPostsService(postsStorage, attachmentStorage)

			tt.mock(postsStorage, attachmentStorage, tt.userID, tt.postID)

			err := s.DeletePost(context.Background(), tt.userID, tt.postID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeletePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetLikedPosts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		userID     uint
		lastLikeID uint
		limit      uint
		mock       func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, lastLikeID uint, limit uint)
		wantPosts  []posts.LikeWithPost
		wantErr    bool
	}{
		{
			name:       "Test OK",
			userID:     1,
			lastLikeID: 0,
			limit:      0,
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, lastLikeID uint, limit uint) {
				likedPosts := []posts.LikeWithPost{{Post: &domain.Post{ID: 1, Content: "Test Content"}}}
				postsStorage.EXPECT().GetLikedPosts(gomock.Any(), userID, lastLikeID, posts.DefaultLikedPostsAmount).Return(likedPosts, nil)
			},
			wantPosts: []posts.LikeWithPost{{Post: &domain.Post{ID: 1, Content: "Test Content"}}},
			wantErr:   false,
		},
		{
			name:       "Test Error",
			userID:     1,
			lastLikeID: 0,
			limit:      0,
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, userID uint, lastLikeID uint, limit uint) {
				postsStorage.EXPECT().GetLikedPosts(gomock.Any(), userID, lastLikeID, posts.DefaultLikedPostsAmount).Return(nil, errors.ErrInternal)
			},
			wantPosts: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postsStorage := mock_posts.NewMockPostsStorage(ctrl)
			attachmentStorage := mock_posts.NewMockAttachmentStorage(ctrl)

			s := posts.NewPostsService(postsStorage, attachmentStorage)

			tt.mock(postsStorage, attachmentStorage, tt.userID, tt.lastLikeID, tt.limit)

			gotPosts, err := s.GetLikedPosts(context.Background(), tt.userID, tt.lastLikeID, tt.limit)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLikedPosts() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotPosts, tt.wantPosts) {
				t.Errorf("GetLikedPosts() gotPosts = %v, want %v", gotPosts, tt.wantPosts)
			}
		})
	}
}

func TestLikePost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		likeData *domain.PostLike
		mock     func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, likeData *domain.PostLike)
		wantLike *domain.PostLike
		wantErr  bool
	}{
		{
			name:     "Test OK",
			likeData: &domain.PostLike{UserID: 1, PostID: 1},
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, likeData *domain.PostLike) {
				postsStorage.EXPECT().StorePostLike(gomock.Any(), likeData).Return(likeData, nil)
			},
			wantLike: &domain.PostLike{UserID: 1, PostID: 1},
			wantErr:  false,
		},
		{
			name:     "Test Error",
			likeData: &domain.PostLike{UserID: 1, PostID: 1},
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, likeData *domain.PostLike) {
				postsStorage.EXPECT().StorePostLike(gomock.Any(), likeData).Return(nil, errors.ErrInternal)
			},
			wantLike: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postsStorage := mock_posts.NewMockPostsStorage(ctrl)
			attachmentStorage := mock_posts.NewMockAttachmentStorage(ctrl)

			s := posts.NewPostsService(postsStorage, attachmentStorage)

			tt.mock(postsStorage, attachmentStorage, tt.likeData)

			gotLike, err := s.LikePost(context.Background(), tt.likeData)

			if (err != nil) != tt.wantErr {
				t.Errorf("LikePost() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotLike, tt.wantLike) {
				t.Errorf("LikePost() gotLike = %v, want %v", gotLike, tt.wantLike)
			}
		})
	}
}

func TestUnlikePost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		likeData *domain.PostLike
		mock     func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, likeData *domain.PostLike)
		wantErr  bool
	}{
		{
			name:     "Test OK",
			likeData: &domain.PostLike{UserID: 1, PostID: 1},
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, likeData *domain.PostLike) {
				postsStorage.EXPECT().DeletePostLike(gomock.Any(), likeData).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Test Error",
			likeData: &domain.PostLike{UserID: 1, PostID: 1},
			mock: func(postsStorage *mock_posts.MockPostsStorage, attachmentStorage *mock_posts.MockAttachmentStorage, likeData *domain.PostLike) {
				postsStorage.EXPECT().DeletePostLike(gomock.Any(), likeData).Return(errors.ErrInternal)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postsStorage := mock_posts.NewMockPostsStorage(ctrl)
			attachmentStorage := mock_posts.NewMockAttachmentStorage(ctrl)

			s := posts.NewPostsService(postsStorage, attachmentStorage)

			tt.mock(postsStorage, attachmentStorage, tt.likeData)

			err := s.UnlikePost(context.Background(), tt.likeData)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnlikePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUploadAttachment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		fileName    string
		filePath    string
		contentType string
		mock        func(attachmentStorage *mock_posts.MockAttachmentStorage, fileName string, filePath string, contentType string)
		wantErr     bool
	}{
		{
			name:        "Test OK",
			fileName:    "test.jpg",
			filePath:    "/path/to/test.jpg",
			contentType: "image/jpeg",
			mock: func(attachmentStorage *mock_posts.MockAttachmentStorage, fileName string, filePath string, contentType string) {
				attachmentStorage.EXPECT().Store(fileName, filePath, contentType).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "Test Error",
			fileName:    "test.jpg",
			filePath:    "/path/to/test.jpg",
			contentType: "image/jpeg",
			mock: func(attachmentStorage *mock_posts.MockAttachmentStorage, fileName string, filePath string, contentType string) {
				attachmentStorage.EXPECT().Store(fileName, filePath, contentType).Return(errors.ErrInternal)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			attachmentStorage := mock_posts.NewMockAttachmentStorage(ctrl)

			s := posts.NewPostsService(nil, attachmentStorage)

			tt.mock(attachmentStorage, tt.fileName, tt.filePath, tt.contentType)

			err := s.UploadAttachment(tt.fileName, tt.filePath, tt.contentType)

			if (err != nil) != tt.wantErr {
				t.Errorf("UploadAttachment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateGroupPost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		groupPost *domain.GroupPost
		mock      func(postsStorage *mock_posts.MockPostsStorage, groupPost *domain.GroupPost)
		wantPost  *domain.GroupPost
		wantErr   bool
	}{
		{
			name:      "Test OK",
			groupPost: &domain.GroupPost{GroupID: 1, PostID: 1},
			mock: func(postsStorage *mock_posts.MockPostsStorage, groupPost *domain.GroupPost) {
				postsStorage.EXPECT().StoreGroupPost(gomock.Any(), groupPost).Return(groupPost, nil)
			},
			wantPost: &domain.GroupPost{GroupID: 1, PostID: 1},
			wantErr:  false,
		},
		{
			name:      "Test Error",
			groupPost: &domain.GroupPost{GroupID: 1, PostID: 1},
			mock: func(postsStorage *mock_posts.MockPostsStorage, groupPost *domain.GroupPost) {
				postsStorage.EXPECT().StoreGroupPost(gomock.Any(), groupPost).Return(nil, errors.ErrInternal)
			},
			wantPost: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postsStorage := mock_posts.NewMockPostsStorage(ctrl)

			s := posts.NewPostsService(postsStorage, nil)

			tt.mock(postsStorage, tt.groupPost)

			gotPost, err := s.CreateGroupPost(context.Background(), tt.groupPost)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateGroupPost() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotPost, tt.wantPost) {
				t.Errorf("CreateGroupPost() gotPost = %v, want %v", gotPost, tt.wantPost)
			}
		})
	}
}

func TestGetPostsOfGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		groupID     uint
		lastPostID  uint
		postsAmount uint
		mock        func(postsStorage *mock_posts.MockPostsStorage, groupID, lastPostID, postsAmount uint)
		wantPosts   []*domain.Post
		wantErr     bool
	}{
		{
			name:        "Test OK",
			groupID:     1,
			lastPostID:  0,
			postsAmount: 0,
			mock: func(postsStorage *mock_posts.MockPostsStorage, groupID, lastPostID, postsAmount uint) {
				postsStorage.EXPECT().GetPostsOfGroup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*domain.Post{{ID: 1}}, nil)
			},
			wantPosts: []*domain.Post{{ID: 1}},
			wantErr:   false,
		},
		{
			name:        "Test Error",
			groupID:     1,
			lastPostID:  0,
			postsAmount: 10,
			mock: func(postsStorage *mock_posts.MockPostsStorage, groupID, lastPostID, postsAmount uint) {
				postsStorage.EXPECT().GetPostsOfGroup(gomock.Any(), groupID, lastPostID, postsAmount).Return(nil, errors.ErrInternal)
			},
			wantPosts: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postsStorage := mock_posts.NewMockPostsStorage(ctrl)

			s := posts.NewPostsService(postsStorage, nil)

			tt.mock(postsStorage, tt.groupID, tt.lastPostID, tt.postsAmount)

			gotPosts, err := s.GetPostsOfGroup(context.Background(), tt.groupID, tt.lastPostID, tt.postsAmount)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPostsOfGroup() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotPosts, tt.wantPosts) {
				t.Errorf("GetPostsOfGroup() gotPosts = %v, want %v", gotPosts, tt.wantPosts)
			}
		})
	}
}

func TestGetGroupPostsBySubscriptionIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		subIDs      []uint
		lastPostID  uint
		postsAmount uint
		mock        func(postsStorage *mock_posts.MockPostsStorage, subIDs []uint, lastPostID, postsAmount uint)
		wantPosts   []*domain.Post
		wantErr     bool
	}{
		{
			name:        "Test OK",
			subIDs:      []uint{1, 2, 3},
			lastPostID:  0,
			postsAmount: 0,
			mock: func(postsStorage *mock_posts.MockPostsStorage, subIDs []uint, lastPostID, postsAmount uint) {
				postsStorage.EXPECT().GetGroupPostsBySubscriptionIDs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*domain.Post{{ID: 1}}, nil)
			},
			wantPosts: []*domain.Post{{ID: 1}},
			wantErr:   false,
		},
		{
			name:        "Test Error",
			subIDs:      []uint{1, 2, 3},
			lastPostID:  0,
			postsAmount: 10,
			mock: func(postsStorage *mock_posts.MockPostsStorage, subIDs []uint, lastPostID, postsAmount uint) {
				postsStorage.EXPECT().GetGroupPostsBySubscriptionIDs(gomock.Any(), subIDs, lastPostID, postsAmount).Return(nil, errors.ErrInternal)
			},
			wantPosts: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postsStorage := mock_posts.NewMockPostsStorage(ctrl)

			s := posts.NewPostsService(postsStorage, nil)

			tt.mock(postsStorage, tt.subIDs, tt.lastPostID, tt.postsAmount)

			gotPosts, err := s.GetGroupPostsBySubscriptionIDs(context.Background(), tt.subIDs, tt.lastPostID, tt.postsAmount)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetGroupPostsBySubscriptionIDs() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotPosts, tt.wantPosts) {
				t.Errorf("GetGroupPostsBySubscriptionIDs() gotPosts = %v, want %v", gotPosts, tt.wantPosts)
			}
		})
	}
}

func TestGetPostsByGroupSubIDsAndUserSubIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		groupSubIDs []uint
		userSubIDs  []uint
		lastPostID  uint
		postsAmount uint
		mock        func(postsStorage *mock_posts.MockPostsStorage, groupSubIDs, userSubIDs []uint, lastPostID, postsAmount uint)
		wantPosts   []*domain.Post
		wantErr     bool
	}{
		{
			name:        "Test OK",
			groupSubIDs: []uint{1, 2, 3},
			userSubIDs:  []uint{4, 5, 6},
			lastPostID:  0,
			postsAmount: 0,
			mock: func(postsStorage *mock_posts.MockPostsStorage, groupSubIDs, userSubIDs []uint, lastPostID, postsAmount uint) {
				postsStorage.EXPECT().GetPostsByGroupSubIDsAndUserSubIDs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*domain.Post{{ID: 1}}, nil)
			},
			wantPosts: []*domain.Post{{ID: 1}},
			wantErr:   false,
		},
		{
			name:        "Test Error",
			groupSubIDs: []uint{1, 2, 3},
			userSubIDs:  []uint{4, 5, 6},
			lastPostID:  0,
			postsAmount: 10,
			mock: func(postsStorage *mock_posts.MockPostsStorage, groupSubIDs, userSubIDs []uint, lastPostID, postsAmount uint) {
				postsStorage.EXPECT().GetPostsByGroupSubIDsAndUserSubIDs(gomock.Any(), groupSubIDs, userSubIDs, lastPostID, postsAmount).Return(nil, errors.ErrInternal)
			},
			wantPosts: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postsStorage := mock_posts.NewMockPostsStorage(ctrl)

			s := posts.NewPostsService(postsStorage, nil)

			tt.mock(postsStorage, tt.groupSubIDs, tt.userSubIDs, tt.lastPostID, tt.postsAmount)

			gotPosts, err := s.GetPostsByGroupSubIDsAndUserSubIDs(context.Background(), tt.groupSubIDs, tt.userSubIDs, tt.lastPostID, tt.postsAmount)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPostsByGroupSubIDsAndUserSubIDs() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotPosts, tt.wantPosts) {
				t.Errorf("GetPostsByGroupSubIDsAndUserSubIDs() gotPosts = %v, want %v", gotPosts, tt.wantPosts)
			}
		})
	}
}

func TestGetNewPosts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		lastPostID  uint
		postsAmount uint
		mock        func(postsStorage *mock_posts.MockPostsStorage, lastPostID, postsAmount uint)
		wantPosts   []*domain.Post
		wantErr     bool
	}{
		{
			name:        "Test OK",
			lastPostID:  0,
			postsAmount: 0,
			mock: func(postsStorage *mock_posts.MockPostsStorage, lastPostID, postsAmount uint) {
				postsStorage.EXPECT().GetNewPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*domain.Post{{ID: 1}}, nil)
			},
			wantPosts: []*domain.Post{{ID: 1}},
			wantErr:   false,
		},
		{
			name:        "Test Error",
			lastPostID:  0,
			postsAmount: 10,
			mock: func(postsStorage *mock_posts.MockPostsStorage, lastPostID, postsAmount uint) {
				postsStorage.EXPECT().GetNewPosts(gomock.Any(), lastPostID, postsAmount).Return(nil, errors.ErrInternal)
			},
			wantPosts: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			postsStorage := mock_posts.NewMockPostsStorage(ctrl)

			s := posts.NewPostsService(postsStorage, nil)

			tt.mock(postsStorage, tt.lastPostID, tt.postsAmount)

			gotPosts, err := s.GetNewPosts(context.Background(), tt.lastPostID, tt.postsAmount)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetNewPosts() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotPosts, tt.wantPosts) {
				t.Errorf("GetNewPosts() gotPosts = %v, want %v", gotPosts, tt.wantPosts)
			}
		})
	}
}

func TestGetGroupPostByPostID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostsStorage := mock_posts.NewMockPostsStorage(ctrl)

	tests := []struct {
		name    string
		postID  uint
		want    *domain.GroupPost
		wantErr error
		setup   func()
	}{
		{
			name:   "test case 1 - successful retrieval",
			postID: 1,
			want: &domain.GroupPost{
				PostID:  1,
				GroupID: 1,
			},
			wantErr: nil,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(&domain.Post{
					ID:       1,
					AuthorID: 1,
				}, nil)

				mockPostsStorage.EXPECT().GetGroupPostByPostID(gomock.Any(), uint(1)).Return(&domain.GroupPost{
					PostID:  1,
					GroupID: 1,
				}, nil)
			},
		},
		{
			name:   "test case 2",
			postID: 1,
			want: &domain.GroupPost{
				PostID:  1,
				GroupID: 1,
			},
			wantErr: errors.ErrNotFound,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(&domain.Post{
					ID:       1,
					AuthorID: 1,
				}, nil)

				mockPostsStorage.EXPECT().GetGroupPostByPostID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name:   "test case 3",
			postID: 1,
			want: &domain.GroupPost{
				PostID:  1,
				GroupID: 1,
			},
			wantErr: errors.ErrNotFound,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := posts.NewPostsService(mockPostsStorage, nil)

			got, err := s.GetGroupPostByPostID(context.Background(), tt.postID)
			assert.Equal(t, tt.wantErr, err)

			if tt.wantErr == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
