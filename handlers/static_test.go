package handlers_test

import (
	"net/http/httptest"
	"socio/handlers"
	"testing"
)

type ServeStaticTestCase struct {
	Method           string
	URL              string
	Body             string
	ExpectedFileName string
	Status           int
}

var StaticHandler = &handlers.StaticHandler{}

var ServeStaticTestCases = map[string]ServeStaticTestCase{
	"success": {
		Method: "GET",
		URL:    "http://localhost:8080/api/v1/static/default_avatar.png",
		Body:   "",
		Status: 200,
	},
	"no path": {
		Method: "GET",
		URL:    "http://localhost:8080/api/v1/static/",
		Body:   "",
		Status: 400,
	},
	"invalid path": {
		Method: "GET",
		URL:    "http://localhost:8080/api/v1/static/invalid_path.png",
		Body:   "",
		Status: 404,
	},
}

func TestHandleServeStatic(t *testing.T) {
	for name, tc := range ServeStaticTestCases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(tc.Method, tc.URL, nil)
			w := httptest.NewRecorder()
			StaticHandler.HandleServeStatic(w, req)

			if w.Code != tc.Status {
				t.Errorf("wrong StatusCode: got %d, expected %d", w.Code, tc.Status)
			}

			resp := w.Result()
			defer resp.Body.Close()
		})
	}
}
