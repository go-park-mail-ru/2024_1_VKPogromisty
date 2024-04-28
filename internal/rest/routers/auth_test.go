package routers_test

//
// import (
// 	"net/http"
// 	"net/http/httptest"
// 	mock_auth "socio/mocks/usecase/auth"
// 	"testing"

// 	"socio/internal/rest/routers"

// 	"github.com/golang/mock/gomock"
// 	"github.com/gorilla/mux"
// 	"github.com/stretchr/testify/assert"
// )

// func TestMountAuthRouter(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	userStorage := mock_auth.NewMockUserStorage(ctrl)
// 	sessionStorage := mock_auth.NewMockSessionStorage(ctrl)

// 	router := mux.NewRouter()
// 	routers.MountAuthRouter(router, userStorage, sessionStorage)

// 	// Test if the routes are correctly mounted
// 	testCases := []struct {
// 		method string
// 		path   string
// 	}{
// 		{"POST", "/auth/login"},
// 		{"OPTIONS", "/auth/login"},
// 		{"POST", "/auth/signup"},
// 		{"OPTIONS", "/auth/signup"},
// 		{"DELETE", "/auth/logout"},
// 		{"OPTIONS", "/auth/logout"},
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
