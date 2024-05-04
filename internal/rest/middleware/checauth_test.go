package middleware

import (
	"net/http"
	"net/http/httptest"
	"socio/errors"
	authpb "socio/internal/grpc/auth/proto"
	mock_auth "socio/mocks/grpc/auth_grpc"
	"socio/pkg/requestcontext"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestCreateCheckIsAuthorizedMiddleware(t *testing.T) {
	testCases := []struct {
		name           string
		cookie         *http.Cookie
		userID         uint
		expectedStatus int
		prepareMocks   func(authClient *mock_auth.MockAuthClient)
	}{
		{
			name:           "valid session",
			cookie:         &http.Cookie{Name: "session_id", Value: "testSessionID"},
			userID:         1,
			expectedStatus: http.StatusOK,
			prepareMocks: func(sessStorage *mock_auth.MockAuthClient) {
				sessStorage.EXPECT().ValidateSession(gomock.Any(), gomock.Any()).Return(
					&authpb.ValidateSessionResponse{UserId: 1}, nil,
				)
			},
		},
		{
			name:           "no cookie",
			cookie:         &http.Cookie{Value: ""},
			userID:         0,
			expectedStatus: http.StatusUnauthorized,
			prepareMocks:   func(sessStorage *mock_auth.MockAuthClient) {},
		},
		{
			name:           "error getting user ID",
			cookie:         &http.Cookie{Name: "session_id", Value: "testSessionID"},
			userID:         0,
			expectedStatus: http.StatusUnauthorized,
			prepareMocks: func(sessStorage *mock_auth.MockAuthClient) {
				sessStorage.EXPECT().ValidateSession(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrUnauthorized.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthClient := mock_auth.NewMockAuthClient(ctrl)

			tc.prepareMocks(mockAuthClient)

			handler := CreateCheckIsAuthorizedMiddleware(mockAuthClient)

			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.AddCookie(tc.cookie)

			rr := httptest.NewRecorder()
			handlerFunc := handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				if ctx.Value(requestcontext.UserIDKey) != tc.userID {
					t.Errorf("context does not contain correct user ID: got %v want %v", ctx.Value(requestcontext.UserIDKey), tc.userID)
				}
				if ctx.Value(requestcontext.SessionIDKey) != tc.cookie.Value {
					t.Errorf("context does not contain correct session ID: got %v want %v", ctx.Value(requestcontext.SessionIDKey), tc.cookie.Value)
				}
			}))

			handlerFunc.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}
		})
	}
}
