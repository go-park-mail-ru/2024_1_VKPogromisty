package rest_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"socio/domain"
	repository "socio/internal/repository/map"
	"socio/internal/rest"
	customtime "socio/pkg/time"
	"strings"
	"sync"
	"testing"
	"time"
)

var userStorage = repository.NewUsers(customtime.MockTimeProvider{}, &sync.Map{})
var sessionStorage, _ = repository.NewSessions(&sync.Map{})
var AuthHandler = rest.NewAuthHandler(customtime.MockTimeProvider{}, userStorage, sessionStorage)

type LoginTestCase struct {
	Method          string
	URL             string
	Body            string
	UserID          uint
	ShouldSetCookie bool
	Status          int
}

var LoginTestCases = map[string]LoginTestCase{
	"success": {
		Method:          "POST",
		URL:             "http://localhost:8080/api/v1/auth/login",
		Body:            `{"email":"petr09mitin@mail.ru", "password":"admin1"}`,
		UserID:          0,
		ShouldSetCookie: true,
		Status:          200,
	},
	"invalid json": {
		Method:          "POST",
		URL:             "http://localhost:8080/api/v1/auth/login",
		Body:            `{"email":""`,
		UserID:          0,
		ShouldSetCookie: false,
		Status:          400,
	},
	"no email": {
		Method:          "POST",
		URL:             "http://localhost:8080/api/v1/auth/login",
		Body:            `{"password":"admin1"}`,
		UserID:          0,
		ShouldSetCookie: false,
		Status:          401,
	},
	"no password": {
		Method:          "POST",
		URL:             "http://localhost:8080/api/v1/auth/login",
		Body:            `{"email":"petr09mitin@mail.ru"}`,
		UserID:          0,
		ShouldSetCookie: false,
		Status:          401,
	},
	"invalid email": {
		Method:          "POST",
		URL:             "http://localhost:8080/api/v1/auth/login",
		Body:            `{"email":"invalid_email", "password":"admin1"}`,
		UserID:          0,
		ShouldSetCookie: false,
		Status:          401,
	},
	"invalid password": {
		Method:          "POST",
		URL:             "http://localhost:8080/api/v1/auth/login",
		Body:            `{"email":"petr09mitin@mail.ru", "password":"invalid_password"}`,
		UserID:          0,
		ShouldSetCookie: false,
		Status:          401,
	},
}

func TestHandleLogin(t *testing.T) {
	date, _ := time.Parse(customtime.DateFormat, "1990-01-01")
	AuthHandler.Service.UserStorage.StoreUser(&domain.User{
		ID:        0,
		FirstName: "Petr",
		LastName:  "Mitin",
		Password:  "admin1",
		Email:     "petr09mitin@mail.ru",
		Avatar:    "default_avatar.png",
		DateOfBirth: customtime.CustomTime{
			Time: date,
		},
	})

	for name, tc := range LoginTestCases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(tc.Method, tc.URL, strings.NewReader(tc.Body))
			w := httptest.NewRecorder()
			AuthHandler.HandleLogin(w, req)

			if w.Code != tc.Status {
				t.Errorf("wrong StatusCode: got %d, expected %d", w.Code, tc.Status)
				return
			}

			resp := w.Result()
			defer resp.Body.Close()

			if !tc.ShouldSetCookie && len(resp.Cookies()) > 0 {
				t.Errorf("wrong Cookies: got %v, expected nil", resp.Cookies())
				return
			}

			if tc.ShouldSetCookie {
				if len(resp.Cookies()) == 0 {
					t.Errorf("wrong Cookies: got nil, expected %v", resp.Cookies())
					return
				}

				session := resp.Cookies()[0]

				if session.Name != "session_id" {
					t.Errorf("wrong Cookies: got %v, expected session_id", resp.Cookies())
					return
				}

				_, err := AuthHandler.Service.SessionStorage.GetUserIDBySession(session.Value)

				if err != nil {
					t.Error("wrong Cookies: cookie isn't set correctly in storage")
				}
			}
		})
	}
}

type LogoutTestCase struct {
	Method string
	URL    string
	Cookie *http.Cookie
	Status int
}

func TestHandleLogout(t *testing.T) {
	sessionID, _ := AuthHandler.Service.SessionStorage.CreateSession(0)

	var LogoutTestCases = map[string]LogoutTestCase{
		"success": {
			Method: "POST",
			URL:    "http://localhost:8080/api/v1/auth/logout",
			Cookie: &http.Cookie{
				Name:     "session_id",
				Value:    sessionID,
				MaxAge:   10 * 60 * 60,
				HttpOnly: true,
				Secure:   true,
				Path:     "/",
			},
			Status: 200,
		},
		"no cookie": {
			Method: "POST",
			URL:    "http://localhost:8080/api/v1/auth/logout",
			Cookie: &http.Cookie{},
			Status: 401,
		},
		"invalid cookie": {
			Method: "POST",
			URL:    "http://localhost:8080/api/v1/auth/logout",
			Cookie: &http.Cookie{
				Name:  "session_id",
				Value: "invalid_session_id",
			},
			Status: 401,
		},
	}

	for name, tc := range LogoutTestCases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(tc.Method, tc.URL, nil)
			req.AddCookie(tc.Cookie)

			w := httptest.NewRecorder()

			AuthHandler.HandleLogout(w, req)

			if w.Code != tc.Status {
				t.Errorf("wrong StatusCode: got %d, expected %d", w.Code, tc.Status)
				return
			}

			resp := w.Result()
			defer resp.Body.Close()

			if tc.Status == 200 {
				_, err := AuthHandler.Service.SessionStorage.GetUserIDBySession(sessionID)

				if err == nil {
					t.Error("session wasn't deleted")
				}
			}

		})
	}
}

type RegistrateUserTestCase struct {
	Method    string
	URL       string
	BodyValue map[string]string
	Avatar    string
	Expected  string
	Status    int
}

var RegistrateUserTestCases = map[string]RegistrateUserTestCase{
	"success": {
		Method: "POST",
		URL:    "http://localhost:8080/api/v1/auth/registration",
		BodyValue: map[string]string{
			"email":          "petr01mitin@gmail.com",
			"password":       "admin1",
			"repeatPassword": "admin1",
			"firstName":      "Petr",
			"lastName":       "Mitin",
			"dateOfBirth":    "1990-01-01",
		},
		Avatar:   "default_avatar.png",
		Expected: `{"body":{"user":{"userId":3,"firstName":"Petr","lastName":"Mitin","email":"petr01mitin@gmail.com","registrationDate":"2021-01-01T00:00:00Z","avatar":"default_avatar.png","dateOfBirth":"1990-01-01T00:00:00Z"}}}`,
		Status:   201,
	},
	"invalid body": {
		Method: "POST",
		URL:    "http://localhost:8080/api/v1/auth/registration",
		BodyValue: map[string]string{
			"email":          "petr01mitin",
			"password":       "admin1",
			"repeatPassword": "admin1",
			"firstName":      "Petr",
			"lastName":       "Mitin",
			"dateOfBirth":    "1990-01-01",
		},
		Avatar:   "default_avatar.png",
		Expected: `{"error":"invalid email"}`,
		Status:   400,
	},
}

func TestHandleRegistration(t *testing.T) {
	for name, tc := range RegistrateUserTestCases {
		t.Run(name, func(t *testing.T) {
			var requestBody bytes.Buffer

			writer := multipart.NewWriter(&requestBody)
			for key, value := range tc.BodyValue {
				writer.WriteField(key, value)
			}
			writer.Close()

			req := httptest.NewRequest(tc.Method, tc.URL, &requestBody)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()
			AuthHandler.HandleRegistration(w, req)

			if w.Code != tc.Status {
				t.Errorf("wrong StatusCode: got %d, expected %d", w.Code, tc.Status)
				return
			}

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()

			bodyStr := string(body)
			if bodyStr != tc.Expected {
				t.Errorf("wrong Response: \ngot %+v, \nexpected %+v", bodyStr, tc.Expected)
			}

			if tc.Status == 201 {
				_, err := AuthHandler.Service.UserStorage.GetUserByEmail(tc.BodyValue["email"])
				if err != nil {
					t.Error("user wasn't created")
					return
				}

				if len(resp.Cookies()) == 0 {
					t.Error("session wasn't created")
					return
				}

				session := resp.Cookies()[0]
				_, err = AuthHandler.Service.IsAuthorized(session)
				if session.Name != "session_id" || err != nil {
					t.Error("session wasn't created")
				}
			}
		})
	}
}

func TestCheckIsAuthorized(t *testing.T) {
	var userStorage = repository.NewUsers(customtime.MockTimeProvider{}, &sync.Map{})
	var sessionStorage, _ = repository.NewSessions(&sync.Map{})
	var authHandler = rest.NewAuthHandler(customtime.MockTimeProvider{}, userStorage, sessionStorage)

	sessionID, _ := authHandler.Service.SessionStorage.CreateSession(0)

	tests := []struct {
		name     string
		cookie   *http.Cookie
		wantBody string
	}{
		{
			name:     "Valid session",
			cookie:   &http.Cookie{Name: "session_id", Value: sessionID},
			wantBody: `{"body":{"isAuthorized":true}}`,
		},
		{
			name:     "Invalid session",
			cookie:   &http.Cookie{Name: "session_id", Value: "invalidSessionValue"},
			wantBody: `{"body":{"isAuthorized":false}}`,
		},
		{
			name:     "No session",
			cookie:   nil,
			wantBody: `{"body":{"isAuthorized":false}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			rr := httptest.NewRecorder()

			authHandler.CheckIsAuthorized(rr, req)

			if body := rr.Body.String(); body != tt.wantBody {
				t.Errorf("handler returned wrong body: got %v want %v", body, tt.wantBody)
			}
		})
	}
}
