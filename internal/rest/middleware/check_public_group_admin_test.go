package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	uspb "socio/internal/grpc/user/proto"
	user_mock "socio/mocks/grpc/user_grpc"
	"socio/pkg/requestcontext"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreateCheckPublicGroupAdminMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := user_mock.NewMockUserClient(ctrl)

	tests := []struct {
		name           string
		userID         uint
		groupID        int
		isAdmin        bool
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Successful admin check",
			userID:         1,
			groupID:        1,
			isAdmin:        true,
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "User is not admin",
			userID:         1,
			groupID:        1,
			isAdmin:        false,
			mockError:      nil,
			expectedStatus: http.StatusForbidden,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("GET", "/", nil)
			ctx := context.Background()
			ctx = context.WithValue(ctx, requestcontext.UserIDKey, tt.userID)
			r = r.WithContext(ctx)

			r = mux.SetURLVars(r, map[string]string{
				"groupID": strconv.FormatInt(int64(tt.groupID), 10),
			})

			// Set up the response recorder
			rr := httptest.NewRecorder()

			// Set up the mock expectations
			mockUserClient.EXPECT().CheckIfUserIsAdmin(gomock.Any(), &uspb.CheckIfUserIsAdminRequest{
				UserId:        uint64(tt.userID),
				PublicGroupId: uint64(tt.groupID),
			}).Return(&uspb.CheckIfUserIsAdminResponse{
				IsAdmin: tt.isAdmin,
			}, tt.mockError)

			// Create the middleware
			middleware := CreateCheckPublicGroupAdminMiddleware(mockUserClient)

			// Create a next handler
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Call the middleware
			middleware(next).ServeHTTP(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
