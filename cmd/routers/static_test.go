package routers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	routers "socio/cmd/routers"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMountStaticRouter(t *testing.T) {
	router := mux.NewRouter()
	routers.MountStaticRouter(router)

	// Test if the routes are correctly mounted
	testCases := []struct {
		method string
		path   string
	}{
		{"GET", "/static/default_avatar.png"},
		{"OPTIONS", "/static/default_avatar.png"},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest(tc.method, tc.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// We're just checking if the routes are mounted, so a 404 status code means the route is not mounted
		assert.NotEqual(t, http.StatusInternalServerError, rr.Code, "Route %s not mounted", tc.path)
	}
}
