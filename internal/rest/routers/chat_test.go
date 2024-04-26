package routers_test

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"socio/internal/rest/routers"
// 	mock_auth "socio/mocks/usecase/auth"
// 	mock_chat "socio/mocks/usecase/chat"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/gorilla/mux"
// 	"github.com/stretchr/testify/assert"
// )

// func TestMountChatRouter(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	pubSubRepo := mock_chat.NewMockPubSubRepository(ctrl)
// 	messagesRepo := mock_chat.NewMockPersonalMessagesRepository(ctrl)
// 	sessionStorage := mock_auth.NewMockSessionStorage(ctrl)

// 	router := mux.NewRouter()
// 	routers.MountChatRouter(router, pubSubRepo, messagesRepo, sessionStorage)

// 	// Test if the routes are correctly mounted
// 	testCases := []struct {
// 		method string
// 		path   string
// 	}{
// 		{"GET", "/chat/ws/"},
// 		{"OPTIONS", "/chat/ws/"},
// 		{"GET", "/chat/dialogs"},
// 		{"OPTIONS", "/chat/dialogs"},
// 		{"GET", "/chat/messages"},
// 		{"OPTIONS", "/chat/messages"},
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
