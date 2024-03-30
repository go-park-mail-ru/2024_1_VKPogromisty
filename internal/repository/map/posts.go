package repository

import (
	"socio/domain"
	customtime "socio/pkg/time"
	"sync"
	"time"
)

type Posts struct {
	Posts  *sync.Map
	NextID uint
}

func NewPosts(tp customtime.TimeProvider, posts *sync.Map) (postsStorage *Posts) {
	postsStorage = &Posts{}
	postsStorage.Posts = posts

	creationDate, _ := time.Parse(customtime.DateFormat, "2000-01-01")

	postsStorage.NextID = 5

	postsStorage.Posts.Store(0, &domain.Post{
		ID:          0,
		AuthorID:    0,
		Content:     "Заснял такие вот красивые деревья)",
		Attachments: []string{"tree1.jpeg", "tree2.jpeg", "tree3.jpeg"},
		CreatedAt:   customtime.CustomTime{Time: creationDate},
		UpdatedAt:   customtime.CustomTime{Time: creationDate},
	})

	postsStorage.Posts.Store(1, &domain.Post{
		ID:          1,
		AuthorID:    0,
		Content:     "Озеро недалеко от моего домика в Швейцарии. Красота!",
		Attachments: []string{"lake.jpeg"},
		CreatedAt:   customtime.CustomTime{Time: creationDate},
		UpdatedAt:   customtime.CustomTime{Time: creationDate},
	})

	postsStorage.Posts.Store(2, &domain.Post{
		ID:          2,
		AuthorID:    0,
		Content:     "Moя подруга - очень хороший фотограф",
		Attachments: []string{"camera.jpeg"},
		CreatedAt:   customtime.CustomTime{Time: creationDate},
		UpdatedAt:   customtime.CustomTime{Time: creationDate},
	})

	postsStorage.Posts.Store(3, &domain.Post{
		ID:          3,
		AuthorID:    0,
		Content:     "Мост в бесконечность",
		Attachments: []string{"bridge.jpeg"},
		CreatedAt:   customtime.CustomTime{Time: creationDate},
		UpdatedAt:   customtime.CustomTime{Time: creationDate},
	})

	postsStorage.Posts.Store(4, &domain.Post{
		ID:          4,
		AuthorID:    0,
		Content:     "Белые розы, белые розы... Не совсем белые, но все равно прекрасно)",
		Attachments: []string{"rose.jpeg"},
		CreatedAt:   customtime.CustomTime{Time: creationDate},
		UpdatedAt:   customtime.CustomTime{Time: creationDate},
	})

	return
}

func (s *Posts) GetAll() (posts []*domain.Post, err error) {
	s.Posts.Range(func(key, value interface{}) bool {
		posts = append(posts, value.(*domain.Post))
		return true
	})

	return
}
