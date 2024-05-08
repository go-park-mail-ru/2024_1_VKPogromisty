package routers_test

import (
	"net/http"
	"testing"

	"socio/internal/rest/routers"
	mock_auth "socio/mocks/grpc/auth_grpc"
	mock_post "socio/mocks/grpc/post_grpc"
	mock_public_group "socio/mocks/grpc/public_group_grpc"
	mock_user "socio/mocks/grpc/user_grpc"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMountPublicGroupRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock clients
	mockGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)
	mockPostClient := mock_post.NewMockPostClient(ctrl)
	mockUserClient := mock_user.NewMockUserClient(ctrl)
	mockAuthManager := mock_auth.NewMockAuthClient(ctrl)

	// Create a new router
	router := mux.NewRouter()

	// Mount the routes
	routers.MountPublicGroupRouter(router, mockGroupClient, mockPostClient, mockUserClient, mockAuthManager)

	// Define the routes to test
	routes := []struct {
		method string
		path   string
	}{
		{"GET", "/groups/search"},
		{"GET", "/groups/1"},
		{"GET", "/groups/by-sub/1"},
		{"GET", "/groups/1/is-sub"},
		{"POST", "/groups/1/sub"},
		{"POST", "/groups/1/unsub"},
		{"POST", "/groups/"},
		{"GET", "/groups/1/posts/"},
		{"GET", "/groups/1/admins/"},
		{"GET", "/groups/1/admins/check"},
		{"POST", "/groups/1/admins/"},
		{"DELETE", "/groups/1/admins/"},
		{"PUT", "/groups/1"},
		{"DELETE", "/groups/1"},
		{"POST", "/groups/1/posts/"},
	}

	// Test each route
	for _, route := range routes {
		req, err := http.NewRequest(route.method, route.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		match := mux.RouteMatch{MatchErr: http.ErrNotSupported}
		if router.Match(req, &match) {
			assert.NotNil(t, match.Handler, "No handler found for route %s", route.path)
		} else {
			t.Errorf("No match for %s %s: %v", route.method, route.path, match.MatchErr)
		}
	}
}
