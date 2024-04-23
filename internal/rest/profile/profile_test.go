package rest

// import (
// 	"bytes"
// 	"context"
// 	"net/http"
// 	"net/http/httptest"
// 	"socio/domain"
// 	"socio/errors"
// 	mock_user "socio/mocks/usecase/user"
// 	"socio/pkg/requestcontext"
// 	"socio/pkg/sanitizer"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/gorilla/mux"
// 	"github.com/microcosm-cc/bluemonday"
// 	"github.com/stretchr/testify/assert"
// )

// type fields struct {
// 	UserStorage    *mock_user.MockUserStorage
// 	SessionStorage *mock_user.MockSessionStorage
// 	Sanitizer      *sanitizer.Sanitizer
// }

// func TestHandleGetSubscriptions(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	tests := []struct {
// 		name           string
// 		expectedStatus int
// 		userID         string
// 		ctx            context.Context
// 		setupMocks     func(*fields)
// 	}{
// 		{
// 			name:           "valid context",
// 			expectedStatus: http.StatusOK,
// 			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().GetUserByIDWithSubsInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.User{
// 					ID: 2,
// 				}, true, true, nil)
// 			},
// 		},
// 		{
// 			name:           "valid userID",
// 			expectedStatus: http.StatusOK,
// 			userID:         "2",
// 			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().GetUserByIDWithSubsInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.User{
// 					ID: 2,
// 				}, true, true, nil)
// 			},
// 		},
// 		{
// 			name:           "invalid userID",
// 			expectedStatus: http.StatusBadRequest,
// 			userID:         "invalid",
// 			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
// 			setupMocks: func(fields *fields) {

// 			},
// 		},
// 		{
// 			name:           "invalid context",
// 			expectedStatus: http.StatusBadRequest,
// 			ctx:            context.Background(),
// 			setupMocks: func(fields *fields) {
// 			},
// 		},
// 		{
// 			name:           "error getting subscriptions",
// 			expectedStatus: http.StatusNotFound,
// 			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2)),
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().GetUserByIDWithSubsInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.User{
// 					ID: 2,
// 				}, true, true, errors.ErrNotFound)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			fields := &fields{
// 				UserStorage:    mock_user.NewMockUserStorage(ctrl),
// 				SessionStorage: mock_user.NewMockSessionStorage(ctrl),
// 				Sanitizer:      sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
// 			}

// 			tt.setupMocks(fields)

// 			handler := NewProfileHandler(fields.UserStorage)

// 			req, err := http.NewRequest("GET", "/get-subs", bytes.NewBuffer(nil))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Set("Content-Type", "application/json")

// 			req = req.WithContext(tt.ctx)

// 			req = mux.SetURLVars(req, map[string]string{
// 				"userID": tt.userID,
// 			})

// 			rr := httptest.NewRecorder()
// 			handler.HandleGetProfile(rr, req)

// 			assert.Equal(t, tt.expectedStatus, rr.Code)
// 		})
// 	}
// }

// func TestHandleDeleteSubscriptions(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	validCtx := context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2))
// 	validCtx = context.WithValue(validCtx, requestcontext.SessionIDKey, "session_id")

// 	ctxNoSession := context.WithValue(context.Background(), requestcontext.UserIDKey, uint(2))

// 	tests := []struct {
// 		name           string
// 		expectedStatus int
// 		ctx            context.Context
// 		setupMocks     func(*fields)
// 	}{
// 		{
// 			name:           "valid context",
// 			expectedStatus: http.StatusNoContent,
// 			ctx:            validCtx,
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(nil)
// 				fields.SessionStorage.EXPECT().DeleteSession(gomock.Any(), gomock.Any()).Return(nil)
// 			},
// 		},
// 		{
// 			name:           "invalid context",
// 			expectedStatus: http.StatusBadRequest,
// 			ctx:            context.Background(),
// 			setupMocks: func(fields *fields) {
// 			},
// 		},
// 		{
// 			name:           "invalid context",
// 			expectedStatus: http.StatusBadRequest,
// 			ctx:            ctxNoSession,
// 			setupMocks: func(fields *fields) {
// 			},
// 		},
// 		{
// 			name:           "error",
// 			expectedStatus: http.StatusNotFound,
// 			ctx:            validCtx,
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(errors.ErrNotFound)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			fields := &fields{
// 				UserStorage:    mock_user.NewMockUserStorage(ctrl),
// 				SessionStorage: mock_user.NewMockSessionStorage(ctrl),
// 				Sanitizer:      sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
// 			}

// 			tt.setupMocks(fields)

// 			handler := NewProfileHandler(fields.UserStorage)

// 			req, err := http.NewRequest("GET", "/get-subs", bytes.NewBuffer(nil))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Set("Content-Type", "application/json")

// 			req = req.WithContext(tt.ctx)

// 			rr := httptest.NewRecorder()
// 			handler.HandleDeleteProfile(rr, req)

// 			assert.Equal(t, tt.expectedStatus, rr.Code)
// 		})
// 	}
// }
