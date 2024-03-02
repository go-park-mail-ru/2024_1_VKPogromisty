package handlers_test

import (
	"net/http/httptest"
	"socio/handlers"
	"socio/utils"
	"testing"

	"github.com/gorilla/mux"
)

type ServeStaticTestCase struct {
	Method   string
	URL      string
	FileName string
	Status   int
}

var StaticHandler = &handlers.StaticHandler{}

var ServeStaticTestCases = map[string]ServeStaticTestCase{
	"success": {
		Method:   "GET",
		URL:      "http://localhost:8080/api/v1/static/default_avatar.png",
		FileName: "default_avatar.png",
		Status:   200,
	},
	"no path": {
		Method:   "GET",
		URL:      "http://localhost:8080/api/v1/static/",
		FileName: "",
		Status:   400,
	},
	"invalid path": {
		Method:   "GET",
		URL:      "http://localhost:8080/api/v1/static/invalid_path.png",
		FileName: "invalid_path.png",
		Status:   404,
	},
}

func TestHandleServeStatic(t *testing.T) {
	utils.StaticFilePath = "../static"

	for name, tc := range ServeStaticTestCases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(tc.Method, tc.URL, nil)
			w := httptest.NewRecorder()

			req = mux.SetURLVars(req, map[string]string{"fileName": tc.FileName})

			StaticHandler.HandleServeStatic(w, req)

			if w.Code != tc.Status {
				t.Errorf("wrong StatusCode: got %d, expected %d", w.Code, tc.Status)
			}

			resp := w.Result()
			defer resp.Body.Close()
		})
	}
}
