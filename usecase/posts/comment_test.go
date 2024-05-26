package posts_test

import (
	"context"
	"socio/domain"
	"socio/errors"
	mock_posts "socio/mocks/usecase/posts"
	customtime "socio/pkg/time"
	"socio/usecase/posts"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetCommentsByPostID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostsStorage := mock_posts.NewMockPostsStorage(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name    string
		postID  uint
		want    []*domain.Comment
		wantErr bool
		setup   func()
	}{
		{
			name:   "test case 1 - successful retrieval",
			postID: 1,
			want: []*domain.Comment{
				{
					ID:         1,
					PostID:     1,
					AuthorID:   1,
					Content:    "Sanitized comment",
					LikedByIDs: []uint64{2},
					CreatedAt:  customtime.CustomTime{Time: tp.Now()},
					UpdatedAt:  customtime.CustomTime{Time: tp.Now()},
				},
			},
			wantErr: false,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(&domain.Post{
					ID:       1,
					AuthorID: 1,
					Content:  "Sanitized post",
				}, nil)

				mockPostsStorage.EXPECT().GetCommentsByPostID(gomock.Any(), uint(1)).Return([]*domain.Comment{
					{
						ID:         1,
						PostID:     1,
						AuthorID:   1,
						Content:    "Sanitized comment",
						LikedByIDs: []uint64{2},
						CreatedAt:  customtime.CustomTime{Time: tp.Now()},
						UpdatedAt:  customtime.CustomTime{Time: tp.Now()},
					},
				}, nil)
			},
		},
		{
			name:    "test case 2",
			postID:  1,
			want:    nil,
			wantErr: true,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(&domain.Post{
					ID:       1,
					AuthorID: 1,
					Content:  "Sanitized post",
				}, nil)

				mockPostsStorage.EXPECT().GetCommentsByPostID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name:    "test case 3",
			postID:  1,
			want:    nil,
			wantErr: true,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := posts.NewPostsService(mockPostsStorage, nil)

			got, err := s.GetCommentsByPostID(context.Background(), tt.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommentsByPostID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCreateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostsStorage := mock_posts.NewMockPostsStorage(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name    string
		comment *domain.Comment
		want    *domain.Comment
		wantErr bool
		setup   func()
	}{
		{
			name: "test case 1 - successful creation",
			comment: &domain.Comment{
				PostID:  1,
				Content: "Sanitized comment",
			},
			want: &domain.Comment{
				ID:         1,
				PostID:     1,
				AuthorID:   1,
				Content:    "Sanitized comment",
				LikedByIDs: []uint64{2},
				CreatedAt:  customtime.CustomTime{Time: tp.Now()},
				UpdatedAt:  customtime.CustomTime{Time: tp.Now()},
			},
			wantErr: false,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(&domain.Post{
					ID:       1,
					AuthorID: 1,
				}, nil)

				mockPostsStorage.EXPECT().StoreComment(gomock.Any(), gomock.Any()).Return(&domain.Comment{
					ID:         1,
					PostID:     1,
					AuthorID:   1,
					Content:    "Sanitized comment",
					LikedByIDs: []uint64{2},
					CreatedAt:  customtime.CustomTime{Time: tp.Now()},
					UpdatedAt:  customtime.CustomTime{Time: tp.Now()},
				}, nil)
			},
		},
		{
			name: "test case 2",
			comment: &domain.Comment{
				PostID:  1,
				Content: "Sanitized comment",
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(&domain.Post{
					ID:       1,
					AuthorID: 1,
				}, nil)

				mockPostsStorage.EXPECT().StoreComment(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name: "test case 3",
			comment: &domain.Comment{
				PostID:  1,
				Content: "Sanitized comment",
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockPostsStorage.EXPECT().GetPostByID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name: "test case 3",
			comment: &domain.Comment{
				PostID:  1,
				Content: "",
			},
			want:    nil,
			wantErr: true,
			setup: func() {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := posts.NewPostsService(mockPostsStorage, nil)

			got, err := s.CreateComment(context.Background(), tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUpdateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostsStorage := mock_posts.NewMockPostsStorage(ctrl)

	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name    string
		comment *domain.Comment
		want    *domain.Comment
		wantErr bool
		setup   func()
	}{
		{
			name: "test case 1 - successful update",
			comment: &domain.Comment{
				ID:       1,
				AuthorID: 1,
				Content:  "Sanitized comment",
			},
			want: &domain.Comment{
				ID:         1,
				AuthorID:   1,
				Content:    "Sanitized comment",
				PostID:     1,
				LikedByIDs: []uint64{2},
				CreatedAt:  customtime.CustomTime{Time: tp.Now()},
				UpdatedAt:  customtime.CustomTime{Time: tp.Now()},
			},
			wantErr: false,
			setup: func() {
				mockPostsStorage.EXPECT().GetCommentByID(gomock.Any(), uint(1)).Return(&domain.Comment{
					ID:       1,
					AuthorID: 1,
					Content:  "Old comment",
				}, nil)

				mockPostsStorage.EXPECT().UpdateComment(gomock.Any(), gomock.Any()).Return(&domain.Comment{
					ID:         1,
					AuthorID:   1,
					Content:    "Sanitized comment",
					PostID:     1,
					LikedByIDs: []uint64{2},
					CreatedAt:  customtime.CustomTime{Time: tp.Now()},
					UpdatedAt:  customtime.CustomTime{Time: tp.Now()},
				}, nil)
			},
		},
		{
			name: "test case 2",
			comment: &domain.Comment{
				ID:       1,
				AuthorID: 1,
				Content:  "Sanitized comment",
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockPostsStorage.EXPECT().GetCommentByID(gomock.Any(), uint(1)).Return(&domain.Comment{
					ID:       1,
					AuthorID: 1,
					Content:  "Old comment",
				}, nil)

				mockPostsStorage.EXPECT().UpdateComment(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name: "test case 3",
			comment: &domain.Comment{
				ID:       1,
				AuthorID: 1,
				Content:  "",
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockPostsStorage.EXPECT().GetCommentByID(gomock.Any(), uint(1)).Return(&domain.Comment{
					ID:       1,
					AuthorID: 1,
					Content:  "Old comment",
				}, nil)
			},
		},
		{
			name: "test case 4",
			comment: &domain.Comment{
				ID:       1,
				AuthorID: 1,
				Content:  "Sanitized comment",
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockPostsStorage.EXPECT().GetCommentByID(gomock.Any(), uint(1)).Return(&domain.Comment{
					ID:       1,
					AuthorID: 2,
					Content:  "Old comment",
				}, nil)
			},
		},
		{
			name: "test case 5",
			comment: &domain.Comment{
				ID:       1,
				AuthorID: 1,
				Content:  "Sanitized comment",
			},
			want:    nil,
			wantErr: true,
			setup: func() {
				mockPostsStorage.EXPECT().GetCommentByID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := posts.NewPostsService(mockPostsStorage, nil)

			got, err := s.UpdateComment(context.Background(), tt.comment)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDeleteComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostsStorage := mock_posts.NewMockPostsStorage(ctrl)

	tests := []struct {
		name    string
		comment *domain.Comment
		wantErr error
		setup   func()
	}{
		{
			name: "test case 1 - successful deletion",
			comment: &domain.Comment{
				ID:       1,
				AuthorID: 1,
			},
			wantErr: nil,
			setup: func() {
				mockPostsStorage.EXPECT().GetCommentByID(gomock.Any(), uint(1)).Return(&domain.Comment{
					ID:       1,
					AuthorID: 1,
				}, nil)

				mockPostsStorage.EXPECT().DeleteComment(gomock.Any(), uint(1)).Return(nil)
			},
		},
		{
			name: "test case 2 - forbidden deletion",
			comment: &domain.Comment{
				ID:       1,
				AuthorID: 2,
			},
			wantErr: errors.ErrForbidden,
			setup: func() {
				mockPostsStorage.EXPECT().GetCommentByID(gomock.Any(), uint(1)).Return(&domain.Comment{
					ID:       1,
					AuthorID: 1,
				}, nil)
			},
		},
		{
			name: "test case 3",
			comment: &domain.Comment{
				ID:       1,
				AuthorID: 1,
			},
			wantErr: errors.ErrInternal,
			setup: func() {
				mockPostsStorage.EXPECT().GetCommentByID(gomock.Any(), uint(1)).Return(&domain.Comment{
					ID:       1,
					AuthorID: 1,
				}, nil)

				mockPostsStorage.EXPECT().DeleteComment(gomock.Any(), uint(1)).Return(errors.ErrInternal)
			},
		},
		{
			name: "test case 4",
			comment: &domain.Comment{
				ID:       1,
				AuthorID: 1,
			},
			wantErr: errors.ErrInternal,
			setup: func() {
				mockPostsStorage.EXPECT().GetCommentByID(gomock.Any(), uint(1)).Return(nil, errors.ErrInternal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := posts.NewPostsService(mockPostsStorage, nil)

			err := s.DeleteComment(context.Background(), tt.comment)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestLikeComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostsStorage := mock_posts.NewMockPostsStorage(ctrl)

	tests := []struct {
		name        string
		commentLike *domain.CommentLike
		want        *domain.CommentLike
		wantErr     error
		setup       func()
	}{
		{
			name: "test case 1 - successful like",
			commentLike: &domain.CommentLike{
				UserID:    1,
				CommentID: 1,
			},
			want: &domain.CommentLike{
				ID:        1,
				UserID:    1,
				CommentID: 1,
			},
			wantErr: nil,
			setup: func() {
				mockPostsStorage.EXPECT().GetCommentByID(gomock.Any(), uint(1)).Return(&domain.Comment{
					ID:       1,
					AuthorID: 1,
				}, nil)

				mockPostsStorage.EXPECT().StoreCommentLike(gomock.Any(), gomock.Any()).Return(&domain.CommentLike{
					ID:        1,
					UserID:    1,
					CommentID: 1,
				}, nil)
			},
		},
		{
			name: "test case 2",
			commentLike: &domain.CommentLike{
				UserID:    1,
				CommentID: 1,
			},
			want:    nil,
			wantErr: errors.ErrNotFound,
			setup: func() {
				mockPostsStorage.EXPECT().GetCommentByID(gomock.Any(), uint(1)).Return(&domain.Comment{
					ID:       1,
					AuthorID: 1,
				}, nil)

				mockPostsStorage.EXPECT().StoreCommentLike(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
		},
		{
			name: "test case 3",
			commentLike: &domain.CommentLike{
				UserID:    1,
				CommentID: 1,
			},
			want:    nil,
			wantErr: errors.ErrNotFound,
			setup: func() {
				mockPostsStorage.EXPECT().GetCommentByID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := posts.NewPostsService(mockPostsStorage, nil)

			got, err := s.LikeComment(context.Background(), tt.commentLike)
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUnlikeComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostsStorage := mock_posts.NewMockPostsStorage(ctrl)

	tests := []struct {
		name    string
		like    *domain.CommentLike
		wantErr error
		setup   func()
	}{
		{
			name: "test case 1 - successful deletion",
			like: &domain.CommentLike{
				UserID:    1,
				CommentID: 1,
			},
			wantErr: nil,
			setup: func() {
				mockPostsStorage.EXPECT().DeleteCommentLike(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "test case 2",
			like: &domain.CommentLike{
				UserID:    1,
				CommentID: 1,
			},
			wantErr: errors.ErrInternal,
			setup: func() {
				mockPostsStorage.EXPECT().DeleteCommentLike(gomock.Any(), gomock.Any()).Return(errors.ErrInternal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := posts.NewPostsService(mockPostsStorage, nil)

			err := s.UnlikeComment(context.Background(), tt.like)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
