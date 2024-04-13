package posts

import (
	"context"
	"mime/multipart"
	"socio/domain"
	"socio/errors"
	"socio/pkg/sanitizer"
)

const (
	defaultPostsAmount = 20
)

type PostInput struct {
	Content     string                  `json:"content"`
	AuthorID    uint                    `json:"authorId"`
	Attachments []*multipart.FileHeader `json:"attachments"`
}

type ListUserPostsInput struct {
	UserID      uint `json:"userId"`
	LastPostID  uint `json:"lastPostId"`
	PostsAmount uint `json:"postsAmount"`
}

type ListUserFriendsPostsInput struct {
	LastPostID  uint `json:"lastPostId"`
	PostsAmount uint `json:"postsAmount"`
}

type PostUpdateInput struct {
	PostID  uint   `json:"postId"`
	Content string `json:"content"`
}

type DeletePostInput struct {
	PostID uint `json:"postId"`
}

type UserStorage interface {
	GetUserByID(ctx context.Context, userID uint) (user *domain.User, err error)
}

type PostsStorage interface {
	GetPostByID(ctx context.Context, postID uint) (post *domain.Post, err error)
	GetUserPosts(ctx context.Context, userID uint, lastPostID uint, postsAmount uint) (posts []*domain.Post, err error)
	GetUserFriendsPosts(ctx context.Context, userID uint, lastPostID uint, postsAmount uint) (posts []domain.PostWithAuthor, err error)
	StorePost(ctx context.Context, post *domain.Post, attachments []*multipart.FileHeader) (newPost *domain.Post, err error)
	UpdatePost(ctx context.Context, post *domain.Post) (updatedPost *domain.Post, err error)
	DeletePost(ctx context.Context, postID uint) (err error)
}

type Service struct {
	PostsStorage PostsStorage
	UserStorage  UserStorage
	Sanitizer    *sanitizer.Sanitizer
}

type ListPostsResponse struct {
	Posts []domain.PostWithAuthor `json:"posts"`
}

func NewPostsService(postsStorage PostsStorage, userStorage UserStorage, sanitizer *sanitizer.Sanitizer) (postsService *Service) {
	postsService = &Service{
		PostsStorage: postsStorage,
		UserStorage:  userStorage,
		Sanitizer:    sanitizer,
	}

	return
}

func (s *Service) GetPostByID(ctx context.Context, postID uint) (post *domain.Post, err error) {
	post, err = s.PostsStorage.GetPostByID(ctx, postID)
	if err != nil {
		return
	}

	s.Sanitizer.SanitizePost(post)

	return
}

func (s *Service) GetUserPosts(ctx context.Context, userID uint, lastPostID uint, postsAmount uint) (posts []*domain.Post, author *domain.User, err error) {
	author, err = s.UserStorage.GetUserByID(ctx, userID)
	if err != nil {
		return
	}

	if postsAmount == 0 {
		postsAmount = defaultPostsAmount
	}

	posts, err = s.PostsStorage.GetUserPosts(ctx, userID, lastPostID, postsAmount)
	if err != nil {
		return
	}

	for _, post := range posts {
		s.Sanitizer.SanitizePost(post)
	}

	return
}

func (s *Service) GetUserFriendsPosts(ctx context.Context, userID uint, lastPostID uint, postsAmount uint) (posts []domain.PostWithAuthor, err error) {
	if postsAmount == 0 {
		postsAmount = defaultPostsAmount
	}

	posts, err = s.PostsStorage.GetUserFriendsPosts(ctx, userID, lastPostID, postsAmount)
	if err != nil {
		return
	}

	for _, post := range posts {
		s.Sanitizer.SanitizePostWithAuthor(&post)
	}

	return
}

func (s *Service) CreatePost(ctx context.Context, input PostInput) (postWithAuthor domain.PostWithAuthor, err error) {
	if len(input.Content) == 0 && len(input.Attachments) == 0 {
		err = errors.ErrInvalidBody
		return
	}

	author, err := s.UserStorage.GetUserByID(ctx, input.AuthorID)
	if err != nil {
		return
	}

	newPost, err := s.PostsStorage.StorePost(ctx, &domain.Post{AuthorID: input.AuthorID, Content: input.Content}, input.Attachments)
	if err != nil {
		return
	}

	postWithAuthor = domain.PostWithAuthor{
		Post:   newPost,
		Author: author,
	}

	s.Sanitizer.SanitizePostWithAuthor(&postWithAuthor)

	return
}

func (s *Service) UpdatePost(ctx context.Context, userID uint, input PostUpdateInput) (post *domain.Post, err error) {
	oldPost, err := s.PostsStorage.GetPostByID(ctx, input.PostID)
	if err != nil {
		return
	}

	if oldPost.AuthorID != userID {
		err = errors.ErrForbidden
		return
	}

	if len(input.Content) == 0 && len(oldPost.Attachments) == 0 {
		err = errors.ErrInvalidBody
		return
	}

	oldPost.Content = input.Content

	post, err = s.PostsStorage.UpdatePost(ctx, oldPost)
	if err != nil {
		return
	}

	s.Sanitizer.SanitizePost(post)

	return
}

func (s *Service) DeletePost(ctx context.Context, postID uint) (err error) {
	err = s.PostsStorage.DeletePost(ctx, postID)
	if err != nil {
		return
	}

	return
}
