package routers_test

import (
	"net/http"
	"testing"

	"socio/internal/rest/routers"
	mock_auth "socio/mocks/grpc/auth_grpc"
	mock_user "socio/mocks/grpc/user_grpc"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMountProfileRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock clients
	mockUserClient := mock_user.NewMockUserClient(ctrl)
	mockAuthManager := mock_auth.NewMockAuthClient(ctrl)

	// Create a new router
	router := mux.NewRouter()

	// Mount the routes
	routers.MountProfileRouter(router, mockUserClient, mockAuthManager)

	// Define the routes to test
	routes := []struct {
		method string
		path   string
	}{
		{"GET", "/profile/search"},
		{"GET", "/profile/"},
		{"PUT", "/profile/"},
		{"DELETE", "/profile/"},
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
