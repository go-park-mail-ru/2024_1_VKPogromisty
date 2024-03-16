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
		ID:           0,
		AuthorID:     0,
		Text:         "Заснял такие вот красивые деревья)",
		Attachments:  []string{"tree1.jpeg", "tree2.jpeg", "tree3.jpeg"},
		CreationDate: customtime.CustomTime{Time: creationDate},
	})

	postsStorage.Posts.Store(1, &domain.Post{
		ID:           1,
		AuthorID:     0,
		Text:         "Озеро недалеко от моего домика в Швейцарии. Красота!",
		Attachments:  []string{"lake.jpeg"},
		CreationDate: customtime.CustomTime{Time: creationDate},
	})

	postsStorage.Posts.Store(2, &domain.Post{
		ID:           2,
		AuthorID:     0,
		Text:         "Moя подруга - очень хороший фотограф",
		Attachments:  []string{"camera.jpeg"},
		CreationDate: customtime.CustomTime{Time: creationDate},
	})

	postsStorage.Posts.Store(3, &domain.Post{
		ID:           3,
		AuthorID:     0,
		Text:         "Мост в бесконечность",
		Attachments:  []string{"bridge.jpeg"},
		CreationDate: customtime.CustomTime{Time: creationDate},
	})

	postsStorage.Posts.Store(4, &domain.Post{
		ID:           4,
		AuthorID:     0,
		Text:         "Белые розы, белые розы... Не совсем белые, но все равно прекрасно)",
		Attachments:  []string{"rose.jpeg"},
		CreationDate: customtime.CustomTime{Time: creationDate},
	})

	return
}

func (s *Posts) GetAll() (posts []*domain.Post) {
	s.Posts.Range(func(key, value interface{}) bool {
		posts = append(posts, value.(*domain.Post))
		return true
	})

	return
}
