package rest_test

import (
	"io"
	"net/http/httptest"
	repository "socio/internal/repository/map"
	"socio/internal/rest"
	customtime "socio/pkg/time"
	"sync"
	"testing"
)

type ListPostsTestCase struct {
	Method   string
	URL      string
	Body     string
	Expected string
	Status   int
}

var postsStorage = repository.NewPosts(customtime.MockTimeProvider{}, &sync.Map{})
var usersStorage = repository.NewUsers(customtime.MockTimeProvider{}, &sync.Map{})
var PostsHandler = rest.NewPostsHandler(postsStorage, usersStorage)

var ListPostsTestCases = map[string]ListPostsTestCase{
	"success": {
		Method:   "GET",
		URL:      "http://localhost:8080/api/v1/posts/",
		Body:     "",
		Expected: `{"body":{"posts":[{"post":{"postId":0,"authorId":0,"content":"Заснял такие вот красивые деревья)","attachments":["tree1.jpeg","tree2.jpeg","tree3.jpeg"],"createdAt":"2000-01-01T00:00:00Z","updatedAt":"2000-01-01T00:00:00Z"},"author":{"userId":0,"firstName":"Petr","lastName":"Mitin","email":"petr09mitin@mail.ru","avatar":"default_avatar.png","dateOfBirth":"1990-01-01T00:00:00Z","createdAt":"2021-01-01T00:00:00Z","updatedAt":"2021-01-01T00:00:00Z"}},{"post":{"postId":1,"authorId":0,"content":"Озеро недалеко от моего домика в Швейцарии. Красота!","attachments":["lake.jpeg"],"createdAt":"2000-01-01T00:00:00Z","updatedAt":"2000-01-01T00:00:00Z"},"author":{"userId":0,"firstName":"Petr","lastName":"Mitin","email":"petr09mitin@mail.ru","avatar":"default_avatar.png","dateOfBirth":"1990-01-01T00:00:00Z","createdAt":"2021-01-01T00:00:00Z","updatedAt":"2021-01-01T00:00:00Z"}},{"post":{"postId":2,"authorId":0,"content":"Moя подруга - очень хороший фотограф","attachments":["camera.jpeg"],"createdAt":"2000-01-01T00:00:00Z","updatedAt":"2000-01-01T00:00:00Z"},"author":{"userId":0,"firstName":"Petr","lastName":"Mitin","email":"petr09mitin@mail.ru","avatar":"default_avatar.png","dateOfBirth":"1990-01-01T00:00:00Z","createdAt":"2021-01-01T00:00:00Z","updatedAt":"2021-01-01T00:00:00Z"}},{"post":{"postId":3,"authorId":0,"content":"Мост в бесконечность","attachments":["bridge.jpeg"],"createdAt":"2000-01-01T00:00:00Z","updatedAt":"2000-01-01T00:00:00Z"},"author":{"userId":0,"firstName":"Petr","lastName":"Mitin","email":"petr09mitin@mail.ru","avatar":"default_avatar.png","dateOfBirth":"1990-01-01T00:00:00Z","createdAt":"2021-01-01T00:00:00Z","updatedAt":"2021-01-01T00:00:00Z"}},{"post":{"postId":4,"authorId":0,"content":"Белые розы, белые розы... Не совсем белые, но все равно прекрасно)","attachments":["rose.jpeg"],"createdAt":"2000-01-01T00:00:00Z","updatedAt":"2000-01-01T00:00:00Z"},"author":{"userId":0,"firstName":"Petr","lastName":"Mitin","email":"petr09mitin@mail.ru","avatar":"default_avatar.png","dateOfBirth":"1990-01-01T00:00:00Z","createdAt":"2021-01-01T00:00:00Z","updatedAt":"2021-01-01T00:00:00Z"}}]}}`,
		Status:   200,
	},
}

func TestHandleListPosts(t *testing.T) {
	for name, tc := range ListPostsTestCases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(tc.Method, tc.URL, nil)
			w := httptest.NewRecorder()
			PostsHandler.HandleListPosts(w, req)

			if w.Code != tc.Status {
				t.Errorf("wrong StatusCode: got %d, expected %d", w.Code, tc.Status)
			}

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()

			bodyStr := string(body)
			if bodyStr != tc.Expected {
				t.Errorf("wrong Response: \ngot %+v, \nexpected %+v", bodyStr, tc.Expected)
			}
		})
	}
}
