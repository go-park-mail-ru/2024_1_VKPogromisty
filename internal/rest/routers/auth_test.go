package routers_test

import (
	"net/http"
	"net/http/httptest"
	mock_auth "socio/mocks/grpc/auth_grpc"
	mock_user "socio/mocks/grpc/user_grpc"
	"testing"

	"socio/internal/rest/routers"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMountAuthRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userClient := mock_user.NewMockUserClient(ctrl)
	authClient := mock_auth.NewMockAuthClient(ctrl)

	router := mux.NewRouter()
	routers.MountAuthRouter(router, authClient, userClient)

	// Test if the routes are correctly mounted
	testCases := []struct {
		method string
		path   string
	}{
		{"POST", "/auth/login"},
		{"OPTIONS", "/auth/login"},
		{"POST", "/auth/signup"},
		{"OPTIONS", "/auth/signup"},
		{"DELETE", "/auth/logout"},
		{"OPTIONS", "/auth/logout"},
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
