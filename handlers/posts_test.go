package handlers_test

import (
	"io"
	"net/http/httptest"
	"socio/handlers"
	"testing"
)

type ListPostsTestCase struct {
	Method   string
	URL      string
	Body     string
	Expected string
	Status   int
}

var PostsHandler = handlers.NewPostsHandler()

var ListPostsTestCases = map[string]ListPostsTestCase{
	"success": {
		Method:   "GET",
		URL:      "http://localhost:8080/api/v1/posts/",
		Body:     "",
		Expected: `{"body":{"posts":[{"post":{"postId":0,"authorId":0,"text":"Заснял такие вот красивые деревья)","attachments":["tree1.jpeg","tree2.jpeg","tree3.jpeg"],"creationDate":"2000-01-01T00:00:00Z"},"author":{"userId":0,"firstName":"Petr","lastName":"Mitin","email":"petr09mitin@mail.ru","registrationDate":"1990-01-01T00:00:00Z","avatar":"default_avatar.png","dateOfBirth":"1990-01-01T00:00:00Z"}},{"post":{"postId":1,"authorId":1,"text":"Озеро недалеко от моего домика в Швейцарии. Красота!","attachments":["lake.jpeg"],"creationDate":"2000-01-01T00:00:00Z"},"author":{"userId":0,"firstName":"Petr","lastName":"Mitin","email":"petr09mitin@mail.ru","registrationDate":"1990-01-01T00:00:00Z","avatar":"default_avatar.png","dateOfBirth":"1990-01-01T00:00:00Z"}},{"post":{"postId":2,"authorId":1,"text":"Moя подруга - очень хороший фотограф","attachments":["camera.jpeg"],"creationDate":"2000-01-01T00:00:00Z"},"author":{"userId":0,"firstName":"Petr","lastName":"Mitin","email":"petr09mitin@mail.ru","registrationDate":"1990-01-01T00:00:00Z","avatar":"default_avatar.png","dateOfBirth":"1990-01-01T00:00:00Z"}},{"post":{"postId":3,"authorId":0,"text":"Мост в бесконечность","attachments":["bridge.jpeg"],"creationDate":"2000-01-01T00:00:00Z"},"author":{"userId":0,"firstName":"Petr","lastName":"Mitin","email":"petr09mitin@mail.ru","registrationDate":"1990-01-01T00:00:00Z","avatar":"default_avatar.png","dateOfBirth":"1990-01-01T00:00:00Z"}},{"post":{"postId":4,"authorId":0,"text":"Белые розы, белые розы... Не совсем белые, но все равно прекрасно)","attachments":["rose.jpeg"],"creationDate":"2000-01-01T00:00:00Z"},"author":{"userId":0,"firstName":"Petr","lastName":"Mitin","email":"petr09mitin@mail.ru","registrationDate":"1990-01-01T00:00:00Z","avatar":"default_avatar.png","dateOfBirth":"1990-01-01T00:00:00Z"}}]}}`,
		Status:   200,
	},
}

func TestListPosts(t *testing.T) {
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
