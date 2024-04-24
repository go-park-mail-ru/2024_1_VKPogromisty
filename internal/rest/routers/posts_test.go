package routers_test

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	routers "socio/internal/rest/routers"
// 	mock_auth "socio/mocks/usecase/auth"
// 	mock_posts "socio/mocks/usecase/posts"

// 	"github.com/golang/mock/gomock"
// 	"github.com/gorilla/mux"
// 	"github.com/stretchr/testify/assert"
// )

// func TestMountPostsRouter(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	postStorage := mock_posts.NewMockPostsStorage(ctrl)
// 	sessionStorage := mock_auth.NewMockSessionStorage(ctrl)

// 	router := mux.NewRouter()
// 	routers.MountPostsRouter(router, postStorage, sessionStorage)

// 	// Test if the routes are correctly mounted
// 	testCases := []struct {
// 		method string
// 		path   string
// 	}{
// 		{"GET", "/posts/"},
// 		{"OPTIONS", "/posts/"},
// 		{"GET", "/posts/friends"},
// 		{"OPTIONS", "/posts/friends"},
// 		{"POST", "/posts/"},
// 		{"OPTIONS", "/posts/"},
// 		{"PUT", "/posts/"},
// 		{"OPTIONS", "/posts/"},
// 		{"DELETE", "/posts/"},
// 		{"OPTIONS", "/posts/"},
// 	}

// 	for _, tc := range testCases {
// 		req, err := http.NewRequest(tc.method, tc.path, nil)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()
// 		router.ServeHTTP(rr, req)

// 		// We're just checking if the routes are mounted, so a 404 status code means the route is not mounted
// 		assert.NotEqual(t, http.StatusNotFound, rr.Code, "Route %s not mounted", tc.path)
// 	}
// }
