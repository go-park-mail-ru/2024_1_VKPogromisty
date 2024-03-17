package posts

import (
	"socio/domain"
	"sort"
)

type PostWithAuthor struct {
	Post   domain.Post `json:"post"`
	Author domain.User `json:"author"`
}

type PostsStorage interface {
	GetAll() (posts []*domain.Post, err error)
}

type UsersStorage interface {
	GetUserByID(userID uint) (user *domain.User, err error)
}

type Service struct {
	PostsStorage PostsStorage
	UsersStorage UsersStorage
}

type ListPostsResponse struct {
	Posts []PostWithAuthor `json:"posts"`
}

func NewPostsService(postsStorage PostsStorage, usersStorage UsersStorage) (postsService *Service) {
	postsService = &Service{
		PostsStorage: postsStorage,
		UsersStorage: usersStorage,
	}

	return
}

func (p *Service) AugmentPostsWithAuthors() (postsWithAuthors []PostWithAuthor, err error) {
	posts, err := p.PostsStorage.GetAll()
	if err != nil {
		return
	}

	for _, post := range posts {
		author, userErr := p.UsersStorage.GetUserByID(post.AuthorID)
		if userErr != nil {
			err = userErr
			return
		}

		postsWithAuthors = append(postsWithAuthors, PostWithAuthor{
			Post:   *post,
			Author: *author,
		})
	}

	sort.Slice(postsWithAuthors, func(i, j int) bool {
		return postsWithAuthors[i].Post.ID < postsWithAuthors[j].Post.ID
	})

	return
}

func (p *Service) ListPosts() (postsWithAuthors []PostWithAuthor, err error) {
	postsWithAuthors, err = p.AugmentPostsWithAuthors()
	if err != nil {
		return
	}

	return
}
