package posts

import (
	"socio/domain"
	repository "socio/internal/repository/map"
	"sort"
)

type PostWithAuthor struct {
	Post   domain.Post `json:"post"`
	Author domain.User `json:"author"`
}

type Service struct {
	PostsStorage *repository.Posts
	UsersStorage *repository.Users
}

type ListPostsResponse struct {
	Posts []PostWithAuthor `json:"posts"`
}

func NewPostsService(postsStorage *repository.Posts, usersStorage *repository.Users) (postsService *Service) {
	postsService = &Service{
		PostsStorage: postsStorage,
		UsersStorage: usersStorage,
	}

	return
}

func (p *Service) AugmentPostsWithAuthors() (postsWithAuthors []PostWithAuthor, err error) {
	for _, post := range p.PostsStorage.GetAll() {
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
