package posts_test

import (
	"reflect"
	"socio/domain"
	repository "socio/internal/repository/map"
	"socio/posts"
	"socio/utils"
	"sync"
	"testing"
	"time"
)

func TestListPosts(t *testing.T) {
	creationDate, _ := time.Parse(utils.DateFormat, "2000-01-01")

	postsStorage := repository.NewPosts(utils.MockTimeProvider{}, &sync.Map{})
	userStorage := repository.NewUsers(utils.MockTimeProvider{}, &sync.Map{})
	postsService := posts.NewPostsService(postsStorage, userStorage)

	author, _ := userStorage.GetUserByID(0)

	post1 := domain.Post{
		ID:           0,
		AuthorID:     0,
		Text:         "Заснял такие вот красивые деревья)",
		Attachments:  []string{"tree1.jpeg", "tree2.jpeg", "tree3.jpeg"},
		CreationDate: utils.CustomTime{Time: creationDate},
	}

	post2 := domain.Post{
		ID:           1,
		AuthorID:     0,
		Text:         "Озеро недалеко от моего домика в Швейцарии. Красота!",
		Attachments:  []string{"lake.jpeg"},
		CreationDate: utils.CustomTime{Time: creationDate},
	}

	post3 := domain.Post{
		ID:           2,
		AuthorID:     0,
		Text:         "Moя подруга - очень хороший фотограф",
		Attachments:  []string{"camera.jpeg"},
		CreationDate: utils.CustomTime{Time: creationDate},
	}

	post4 := domain.Post{
		ID:           3,
		AuthorID:     0,
		Text:         "Мост в бесконечность",
		Attachments:  []string{"bridge.jpeg"},
		CreationDate: utils.CustomTime{Time: creationDate},
	}

	post5 := domain.Post{
		ID:           4,
		AuthorID:     0,
		Text:         "Белые розы, белые розы... Не совсем белые, но все равно прекрасно)",
		Attachments:  []string{"rose.jpeg"},
		CreationDate: utils.CustomTime{Time: creationDate},
	}

	tests := []struct {
		name string
		want []posts.PostWithAuthor
	}{
		{"List posts", []posts.PostWithAuthor{
			{Post: post1, Author: *author},
			{Post: post2, Author: *author},
			{Post: post3, Author: *author},
			{Post: post4, Author: *author},
			{Post: post5, Author: *author},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := postsService.ListPosts()
			if err != nil {
				t.Errorf("ListPosts() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListPosts() = %v, want %v", got, tt.want)
			}
		})
	}
}
