package routers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"socio/internal/rest/routers"
	mock_auth "socio/mocks/grpc/auth_grpc"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMountCSRFRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authClient := mock_auth.NewMockAuthClient(ctrl)

	router := mux.NewRouter()
	routers.MountCSRFRouter(router, authClient)

	// Test if the routes are correctly mounted
	testCases := []struct {
		method string
		path   string
	}{
		{"GET", "/csrf/"},
		{"OPTIONS", "/csrf/"},
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
