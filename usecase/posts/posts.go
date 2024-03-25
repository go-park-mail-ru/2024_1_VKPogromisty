package posts

import (
	"mime/multipart"
	"socio/domain"
	"socio/errors"
)

type PostWithAuthor struct {
	Post   domain.Post `json:"post"`
	Author domain.User `json:"author"`
}

type PostInput struct {
	Content     string                  `json:"content"`
	AuthorID    uint                    `json:"author_id"`
	Attachments []*multipart.FileHeader `json:"attachments"`
}

type UserStorage interface {
	GetUserByID(userID uint) (user *domain.User, err error)
}

type PostsStorage interface {
	GetUserPosts(userID uint, lastPostID uint) (posts []*domain.Post, err error)
	StorePost(post *domain.Post, attachments []*multipart.FileHeader) (newPost *domain.Post, err error)
	DeletePost(postID uint) (err error)
}

type Service struct {
	PostsStorage PostsStorage
	UserStorage  UserStorage
}

type ListPostsResponse struct {
	Posts []PostWithAuthor `json:"posts"`
}

func NewPostsService(postsStorage PostsStorage, userStorage UserStorage) (postsService *Service) {
	postsService = &Service{
		PostsStorage: postsStorage,
		UserStorage:  userStorage,
	}

	return
}

func (s *Service) GetUserPosts(userID uint, lastPostID uint) (posts []*domain.Post, author *domain.User, err error) {
	author, err = s.UserStorage.GetUserByID(userID)
	if err != nil {
		return
	}

	posts, err = s.PostsStorage.GetUserPosts(userID, lastPostID)
	return
}

func (s *Service) CreatePost(input PostInput) (postWithAuthor *PostWithAuthor, err error) {
	if len(input.Content) == 0 && len(input.Attachments) == 0 {
		err = errors.ErrInvalidBody
		return
	}

	author, err := s.UserStorage.GetUserByID(input.AuthorID)
	if err != nil {
		return
	}

	newPost, err := s.PostsStorage.StorePost(&domain.Post{AuthorID: input.AuthorID, Content: input.Content}, input.Attachments)
	if err != nil {
		return
	}

	postWithAuthor = &PostWithAuthor{
		Post:   *newPost,
		Author: *author,
	}

	return
}

func (s *Service) DeletePost(postID uint) (err error) {
	err = s.PostsStorage.DeletePost(postID)
	return
}
