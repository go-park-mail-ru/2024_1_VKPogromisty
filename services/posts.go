package services

import (
	"socio/utils"
	"sync"
	"time"
)

type Post struct {
	ID           uint             `json:"postId"`
	AuthorID     uint             `json:"authorId"`
	Text         string           `json:"text"`
	Attachments  []string         `json:"attachments"`
	CreationDate utils.CustomTime `json:"creationDate,omitempty"`
}

type PostWithAuthor struct {
	Post   *Post
	Author *User
}

type PostsService struct {
	posts  sync.Map
	nextID uint
}

func NewPostsService() (postsService *PostsService) {
	postsService = &PostsService{
		posts:  sync.Map{},
		nextID: 5,
	}

	postsService.posts.Store(0, &Post{
		ID:           0,
		AuthorID:     0,
		Text:         "Заснял такие вот красивые деревья)",
		Attachments:  []string{"tree1.jpeg", "tree2.jpeg", "tree3.jpeg"},
		CreationDate: utils.CustomTime{Time: time.Now()},
	})

	postsService.posts.Store(1, &Post{
		ID:           1,
		AuthorID:     1,
		Text:         "Озеро недалеко от моего домика в Швейцарии. Красота!",
		Attachments:  []string{"lake.jpeg"},
		CreationDate: utils.CustomTime{Time: time.Now()},
	})

	postsService.posts.Store(2, &Post{
		ID:           2,
		AuthorID:     1,
		Text:         "Moя подруга - очень хороший фотограф",
		Attachments:  []string{"camera.jpeg"},
		CreationDate: utils.CustomTime{Time: time.Now()},
	})

	postsService.posts.Store(3, &Post{
		ID:           3,
		AuthorID:     0,
		Text:         "Мост в бесконечность",
		Attachments:  []string{"bridge.jpeg"},
		CreationDate: utils.CustomTime{Time: time.Now()},
	})

	postsService.posts.Store(4, &Post{
		ID:           3,
		AuthorID:     0,
		Text:         "Белые розы, белые розы... Не совсем белые, но все равно прекрасно)",
		Attachments:  []string{"rose.jpeg"},
		CreationDate: utils.CustomTime{Time: time.Now()},
	})

	return
}

func (p *PostsService) augmentPostsWithAuthors() (postsWithAuthors []PostWithAuthor, err error) {
	author := &User{
		ID:        0,
		FirstName: "Petr",
		LastName:  "Mitin",
		Email:     "petr09mitin@mail.ru",
		RegistrationDate: utils.CustomTime{
			Time: time.Now(),
		},
		Avatar: "default_avatar.png",
		DateOfBirth: utils.CustomTime{
			Time: time.Now(),
		},
	}

	p.posts.Range(func(key any, value any) bool {
		postsWithAuthors = append(postsWithAuthors, PostWithAuthor{
			Post:   value.(*Post),
			Author: author,
		})
		return true
	})

	return
}

func (p *PostsService) ListPosts() (postsWithAuthors []PostWithAuthor, err error) {
	var posts []*Post
	p.posts.Range(func(key any, value any) bool {
		posts = append(posts, value.(*Post))
		return true
	})

	postsWithAuthors, err = p.augmentPostsWithAuthors()
	if err != nil {
		return
	}

	return
}
