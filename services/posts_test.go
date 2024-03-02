package services_test

import (
	"reflect"
	"socio/services"
	"socio/utils"
	"testing"
	"time"
)

func TestListPosts(t *testing.T) {
	creationDate, _ := time.Parse(utils.DateFormat, "2000-01-01")
	date, _ := time.Parse(utils.DateFormat, "1990-01-01")
	author := services.User{
		ID:        0,
		FirstName: "Petr",
		LastName:  "Mitin",
		Email:     "petr09mitin@mail.ru",
		RegistrationDate: utils.CustomTime{
			Time: date,
		},
		Avatar: "default_avatar.png",
		DateOfBirth: utils.CustomTime{
			Time: date,
		},
	}

	postsService := services.NewPostsService()
	post1 := services.Post{
		ID:           0,
		AuthorID:     0,
		Text:         "Заснял такие вот красивые деревья)",
		Attachments:  []string{"tree1.jpeg", "tree2.jpeg", "tree3.jpeg"},
		CreationDate: utils.CustomTime{Time: creationDate},
	}

	post2 := services.Post{
		ID:           1,
		AuthorID:     1,
		Text:         "Озеро недалеко от моего домика в Швейцарии. Красота!",
		Attachments:  []string{"lake.jpeg"},
		CreationDate: utils.CustomTime{Time: creationDate},
	}

	post3 := services.Post{
		ID:           2,
		AuthorID:     1,
		Text:         "Moя подруга - очень хороший фотограф",
		Attachments:  []string{"camera.jpeg"},
		CreationDate: utils.CustomTime{Time: creationDate},
	}

	post4 := services.Post{
		ID:           3,
		AuthorID:     0,
		Text:         "Мост в бесконечность",
		Attachments:  []string{"bridge.jpeg"},
		CreationDate: utils.CustomTime{Time: creationDate},
	}

	post5 := services.Post{
		ID:           4,
		AuthorID:     0,
		Text:         "Белые розы, белые розы... Не совсем белые, но все равно прекрасно)",
		Attachments:  []string{"rose.jpeg"},
		CreationDate: utils.CustomTime{Time: creationDate},
	}

	tests := []struct {
		name string
		want []services.PostWithAuthor
	}{
		{"List posts", []services.PostWithAuthor{
			{Post: post1, Author: author},
			{Post: post2, Author: author},
			{Post: post3, Author: author},
			{Post: post4, Author: author},
			{Post: post5, Author: author},
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
