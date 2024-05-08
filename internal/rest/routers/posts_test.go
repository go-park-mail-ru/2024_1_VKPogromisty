package routers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	routers "socio/internal/rest/routers"
	mock_auth "socio/mocks/grpc/auth_grpc"
	mock_posts "socio/mocks/grpc/post_grpc"
	mock_public_group "socio/mocks/grpc/public_group_grpc"
	mock_user "socio/mocks/grpc/user_grpc"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMountPostsRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postsClient := mock_posts.NewMockPostClient(ctrl)
	authClient := mock_auth.NewMockAuthClient(ctrl)
	userClient := mock_user.NewMockUserClient(ctrl)
	publicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)

	router := mux.NewRouter()
	routers.MountPostsRouter(router, postsClient, userClient, publicGroupClient, authClient)

	// Test if the routes are correctly mounted
	testCases := []struct {
		method string
		path   string
	}{
		{"GET", "/posts/"},
		{"OPTIONS", "/posts/"},
		{"GET", "/posts/friends"},
		{"OPTIONS", "/posts/friends"},
		{"POST", "/posts/"},
		{"OPTIONS", "/posts/"},
		{"PUT", "/posts/"},
		{"OPTIONS", "/posts/"},
		{"DELETE", "/posts/"},
		{"OPTIONS", "/posts/"},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest(tc.method, tc.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// We're just checking if the routes are mounted, so a 404 status code means the route is not mounted
		assert.NotEqual(t, http.StatusNotFound, rr.Code, "Route %s not mounted", tc.path)
	}
}
