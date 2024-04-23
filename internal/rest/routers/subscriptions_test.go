package routers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	routers "socio/internal/rest/routers"
	mock_auth "socio/mocks/usecase/auth"
	mock_subscriptions "socio/mocks/usecase/subscriptions"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMountSubscriptionsRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	subStorage := mock_subscriptions.NewMockSubscriptionsStorage(ctrl)
	userStorage := mock_subscriptions.NewMockUserStorage(ctrl)
	sessionStorage := mock_auth.NewMockSessionStorage(ctrl)

	router := mux.NewRouter()
	routers.MountSubscriptionsRouter(router, subStorage, userStorage, sessionStorage)

	// Test if the routes are correctly mounted
	testCases := []struct {
		method string
		path   string
	}{
		{"POST", "/subscriptions/"},
		{"OPTIONS", "/subscriptions/"},
		{"DELETE", "/subscriptions/"},
		{"OPTIONS", "/subscriptions/"},
		{"GET", "/subscriptions/subscribers"},
		{"OPTIONS", "/subscriptions/subscribers"},
		{"GET", "/subscriptions/subscriptions"},
		{"OPTIONS", "/subscriptions/subscriptions"},
		{"GET", "/subscriptions/friends"},
		{"OPTIONS", "/subscriptions/friends"},
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
