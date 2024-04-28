package rest_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"socio/domain"
// 	"socio/errors"
// 	rest "socio/internal/rest/auth"
// 	mock_auth "socio/mocks/usecase/auth"
// 	"socio/pkg/hash"
// 	"socio/pkg/sanitizer"
// 	"socio/usecase/auth"
// 	"testing"

// 	"github.com/stretchr/testify/assert"

// 	"github.com/golang/mock/gomock"
// 	"github.com/microcosm-cc/bluemonday"
// )

// type fields struct {
// 	UserStorage    *mock_auth.MockUserStorage
// 	SessionStorage *mock_auth.MockSessionStorage
// 	Sanitizer      *sanitizer.Sanitizer
// }

// func TestHandleLogin(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	tests := []struct {
// 		name           string
// 		input          *auth.LoginInput
// 		expectedStatus int
// 		setupMocks     func(*fields)
// 	}{
// 		{
// 			name: "valid input",
// 			input: &auth.LoginInput{
// 				Email:    "john.doe@example.com",
// 				Password: "secret",
// 			},
// 			expectedStatus: http.StatusOK,
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(&domain.User{
// 					ID:       1,
// 					Email:    "john.doe@example.com",
// 					Password: hash.HashPassword("secret", []byte("salt")),
// 					Salt:     "salt",
// 				}, nil)
// 				fields.UserStorage.EXPECT().RefreshSaltAndRehashPassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
// 				fields.SessionStorage.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return("session_id", nil)
// 			},
// 		},
// 		{
// 			name: "no user",
// 			input: &auth.LoginInput{
// 				Email:    "john.doe@example.com",
// 				Password: "secret",
// 			},
// 			expectedStatus: http.StatusUnauthorized,
// 			setupMocks: func(fields *fields) {
// 				fields.UserStorage.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			fields := &fields{
// 				UserStorage:    mock_auth.NewMockUserStorage(ctrl),
// 				SessionStorage: mock_auth.NewMockSessionStorage(ctrl),
// 				Sanitizer:      sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
// 			}

// 			tt.setupMocks(fields)

// 			handler := rest.NewAuthHandler(fields.UserStorage, fields.SessionStorage, fields.Sanitizer)

// 			b, err := json.Marshal(tt.input)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(b))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Set("Content-Type", "application/json")

// 			rr := httptest.NewRecorder()
// 			handler.HandleLogin(rr, req)

// 			assert.Equal(t, tt.expectedStatus, rr.Code)
// 		})
// 	}
// }

// func TestHandleLogout(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	tests := []struct {
// 		name           string
// 		cookie         *http.Cookie
// 		expectedStatus int
// 		setupMocks     func(*fields)
// 	}{
// 		{
// 			name: "valid session",
// 			cookie: &http.Cookie{
// 				Name:  "session_id",
// 				Value: "valid_session_id",
// 			},
// 			expectedStatus: http.StatusOK,
// 			setupMocks: func(fields *fields) {
// 				fields.SessionStorage.EXPECT().DeleteSession(gomock.Any(), "valid_session_id").Return(nil)
// 			},
// 		},
// 		{
// 			name: "no session",
// 			cookie: &http.Cookie{
// 				Value: "",
// 			},
// 			expectedStatus: http.StatusUnauthorized,
// 			setupMocks:     func(fields *fields) {},
// 		},
// 		{
// 			name: "err",
// 			cookie: &http.Cookie{
// 				Name:  "session_id",
// 				Value: "valid_session_id",
// 			},
// 			expectedStatus: http.StatusUnauthorized,
// 			setupMocks: func(fields *fields) {
// 				fields.SessionStorage.EXPECT().DeleteSession(gomock.Any(), gomock.Any()).Return(errors.ErrInternal)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			fields := &fields{
// 				UserStorage:    mock_auth.NewMockUserStorage(ctrl),
// 				SessionStorage: mock_auth.NewMockSessionStorage(ctrl),
// 				Sanitizer:      sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
// 			}

// 			tt.setupMocks(fields)

// 			handler := rest.NewAuthHandler(fields.UserStorage, fields.SessionStorage, fields.Sanitizer)

// 			req, err := http.NewRequest("POST", "/logout", nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.AddCookie(tt.cookie)

// 			rr := httptest.NewRecorder()
// 			handler.HandleLogout(rr, req)

// 			assert.Equal(t, tt.expectedStatus, rr.Code)
// 		})
// 	}
// }
