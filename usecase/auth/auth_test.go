package auth_test

// import (
// 	"context"
// 	"net/http"
// 	"reflect"
// 	"socio/domain"
// 	"socio/errors"
// 	mock_auth "socio/mocks/usecase/auth"
// 	"socio/pkg/hash"
// 	"socio/pkg/sanitizer"
// 	customtime "socio/pkg/time"
// 	"socio/usecase/auth"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/microcosm-cc/bluemonday"
// )

// func TestService_RegistrateUser(t *testing.T) {
// 	timeProv := customtime.MockTimeProvider{}

// 	sanitizer := sanitizer.NewSanitizer(bluemonday.StripTagsPolicy())

// 	type fields struct {
// 		UserStorage    *mock_auth.MockUserStorage
// 		SessionStorage *mock_auth.MockSessionStorage
// 	}

// 	type args struct {
// 		ctx       context.Context
// 		userInput auth.RegistrationInput
// 	}

// 	tests := []struct {
// 		name        string
// 		prepareMock func(*fields)
// 		args        args
// 		wantUser    *domain.User
// 		wantSession *http.Cookie
// 		wantErr     bool
// 	}{
// 		{
// 			name: "success",
// 			args: args{
// 				ctx: context.Background(),
// 				userInput: auth.RegistrationInput{
// 					FirstName:      "John",
// 					LastName:       "Doe",
// 					Password:       "password",
// 					RepeatPassword: "password",
// 					Email:          "john@mail.ru",
// 					Avatar:         nil,
// 					DateOfBirth:    "2021-01-01",
// 				},
// 			},
// 			wantUser: &domain.User{
// 				FirstName: "John",
// 				LastName:  "Doe",
// 				Email:     "john@mail.ru",
// 				Password:  "password",
// 				Avatar:    "default_avatar.png",
// 				DateOfBirth: customtime.CustomTime{
// 					Time: timeProv.Now(),
// 				},
// 			},
// 			wantSession: &http.Cookie{
// 				Name:     "session_id",
// 				Value:    "session_id",
// 				MaxAge:   10 * 60 * 60,
// 				HttpOnly: true,
// 				Secure:   true,
// 				Path:     "/",
// 				SameSite: http.SameSiteNoneMode,
// 			},
// 			wantErr: false,
// 			prepareMock: func(f *fields) {
// 				f.UserStorage.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
// 				f.UserStorage.EXPECT().StoreUser(gomock.Any(), gomock.Any()).Return(nil)
// 				f.SessionStorage.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return("session_id", nil)
// 			},
// 		},
// 		{
// 			name: "invalid passwords",
// 			args: args{
// 				ctx: context.Background(),
// 				userInput: auth.RegistrationInput{
// 					FirstName:      "John",
// 					LastName:       "Doe",
// 					Password:       "password",
// 					RepeatPassword: "tyazhelo",
// 					Email:          "john@mail.ru",
// 					Avatar:         nil,
// 					DateOfBirth:    "2021-01-01",
// 				},
// 			},
// 			wantUser:    nil,
// 			wantSession: nil,
// 			wantErr:     true,
// 			prepareMock: func(f *fields) {
// 			},
// 		},
// 		{
// 			name: "invalid date",
// 			args: args{
// 				ctx: context.Background(),
// 				userInput: auth.RegistrationInput{
// 					FirstName:      "John",
// 					LastName:       "Doe",
// 					Password:       "password",
// 					RepeatPassword: "tyazhelo",
// 					Email:          "john@mail.ru",
// 					Avatar:         nil,
// 					DateOfBirth:    "",
// 				},
// 			},
// 			wantUser:    nil,
// 			wantSession: nil,
// 			wantErr:     true,
// 			prepareMock: func(f *fields) {
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()
// 			f := fields{
// 				UserStorage:    mock_auth.NewMockUserStorage(ctrl),
// 				SessionStorage: mock_auth.NewMockSessionStorage(ctrl),
// 			}

// 			if tt.prepareMock != nil {
// 				tt.prepareMock(&f)
// 			}

// 			s := auth.NewService(f.UserStorage, f.SessionStorage, sanitizer)

// 			gotUser, gotSession, err := s.RegistrateUser(tt.args.ctx, tt.args.userInput)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Service.RegistrateUser() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotUser, tt.wantUser) {
// 				t.Errorf("Service.RegistrateUser() gotUser = %v, want %v", gotUser, tt.wantUser)
// 			}
// 			if !reflect.DeepEqual(gotSession, tt.wantSession) {
// 				t.Errorf("Service.RegistrateUser() gotSession = %v, want %v", gotSession, tt.wantSession)
// 			}
// 		})
// 	}
// }

// func TestService_Login(t *testing.T) {
// 	timeProv := customtime.MockTimeProvider{}

// 	sanitizer := sanitizer.NewSanitizer(bluemonday.StripTagsPolicy())

// 	type fields struct {
// 		UserStorage    *mock_auth.MockUserStorage
// 		SessionStorage *mock_auth.MockSessionStorage
// 	}

// 	type args struct {
// 		ctx        context.Context
// 		loginInput auth.LoginInput
// 	}

// 	tests := []struct {
// 		name        string
// 		args        args
// 		prepareMock func(*fields)
// 		wantUser    *domain.User
// 		wantSession *http.Cookie
// 		wantErr     bool
// 	}{
// 		{
// 			name: "success",
// 			args: args{
// 				ctx: context.Background(),
// 				loginInput: auth.LoginInput{
// 					Email:    "john@mail.ru",
// 					Password: "password",
// 				},
// 			},
// 			wantUser: &domain.User{
// 				FirstName: "John",
// 				LastName:  "Doe",
// 				Email:     "john@mail.ru",
// 				Password:  hash.HashPassword("password", []byte("salt")),
// 				Salt:      "salt",
// 				Avatar:    "default_avatar.png",
// 				DateOfBirth: customtime.CustomTime{
// 					Time: timeProv.Now(),
// 				},
// 			},
// 			wantSession: &http.Cookie{
// 				Name:     "session_id",
// 				Value:    "session_id",
// 				MaxAge:   10 * 60 * 60,
// 				HttpOnly: true,
// 				Secure:   true,
// 				Path:     "/",
// 				SameSite: http.SameSiteNoneMode,
// 			},
// 			wantErr: false,
// 			prepareMock: func(f *fields) {
// 				f.UserStorage.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(&domain.User{
// 					FirstName: "John",
// 					LastName:  "Doe",
// 					Email:     "john@mail.ru",
// 					Password:  hash.HashPassword("password", []byte("salt")),
// 					Salt:      "salt",
// 					Avatar:    "default_avatar.png",
// 					DateOfBirth: customtime.CustomTime{
// 						Time: timeProv.Now(),
// 					},
// 				}, nil)
// 				f.SessionStorage.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return("session_id", nil)
// 				f.UserStorage.EXPECT().RefreshSaltAndRehashPassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
// 			},
// 		},
// 		{
// 			name: "failed to get user by email",
// 			args: args{
// 				ctx: context.Background(),
// 				loginInput: auth.LoginInput{
// 					Email:    "john@mail.ru",
// 					Password: "password",
// 				},
// 			},
// 			wantUser:    nil,
// 			wantSession: nil,
// 			wantErr:     true,
// 			prepareMock: func(f *fields) {
// 				f.UserStorage.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
// 			},
// 		},
// 		{
// 			name: "invalid password",
// 			args: args{
// 				ctx: context.Background(),
// 				loginInput: auth.LoginInput{
// 					Email:    "john@mail.ru",
// 					Password: "",
// 				},
// 			},
// 			wantUser: &domain.User{
// 				FirstName: "John",
// 				LastName:  "Doe",
// 				Email:     "john@mail.ru",
// 				Password:  hash.HashPassword("password", []byte("salt")),
// 				Salt:      "salt",
// 				Avatar:    "default_avatar.png",
// 				DateOfBirth: customtime.CustomTime{
// 					Time: timeProv.Now(),
// 				},
// 			},
// 			wantSession: nil,
// 			wantErr:     true,
// 			prepareMock: func(f *fields) {
// 				f.UserStorage.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(&domain.User{
// 					FirstName: "John",
// 					LastName:  "Doe",
// 					Email:     "john@mail.ru",
// 					Password:  hash.HashPassword("password", []byte("salt")),
// 					Salt:      "salt",
// 					Avatar:    "default_avatar.png",
// 					DateOfBirth: customtime.CustomTime{
// 						Time: timeProv.Now(),
// 					},
// 				}, nil)
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()
// 			f := fields{
// 				UserStorage:    mock_auth.NewMockUserStorage(ctrl),
// 				SessionStorage: mock_auth.NewMockSessionStorage(ctrl),
// 			}

// 			if tt.prepareMock != nil {
// 				tt.prepareMock(&f)
// 			}

// 			s := auth.NewService(f.UserStorage, f.SessionStorage, sanitizer)

// 			gotUser, gotSession, err := s.Login(tt.args.ctx, tt.args.loginInput)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Service.Login() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotUser, tt.wantUser) {
// 				t.Errorf("Service.Login() gotUser = %v, want %v", gotUser, tt.wantUser)
// 			}
// 			if !reflect.DeepEqual(gotSession, tt.wantSession) {
// 				t.Errorf("Service.Login() gotSession = %v, want %v", gotSession, tt.wantSession)
// 			}
// 		})
// 	}
// }

// func TestService_Logout(t *testing.T) {
// 	type fields struct {
// 		SessionStorage *mock_auth.MockSessionStorage
// 	}

// 	type args struct {
// 		ctx     context.Context
// 		session *http.Cookie
// 	}

// 	tests := []struct {
// 		name        string
// 		args        args
// 		prepareMock func(*fields)
// 		wantErr     bool
// 	}{
// 		{
// 			name: "success",
// 			args: args{
// 				ctx: context.Background(),
// 				session: &http.Cookie{
// 					Name:     "session_id",
// 					Value:    "session_id",
// 					MaxAge:   10 * 60 * 60,
// 					HttpOnly: true,
// 					Secure:   true,
// 					Path:     "/",
// 					SameSite: http.SameSiteNoneMode,
// 				},
// 			},
// 			wantErr: false,
// 			prepareMock: func(f *fields) {
// 				f.SessionStorage.EXPECT().DeleteSession(gomock.Any(), gomock.Any()).Return(nil)
// 			},
// 		},
// 		{
// 			name: "failed to delete session",
// 			args: args{
// 				ctx: context.Background(),
// 				session: &http.Cookie{
// 					Name:     "session_id",
// 					Value:    "session_id",
// 					MaxAge:   10 * 60 * 60,
// 					HttpOnly: true,
// 					Secure:   true,
// 					Path:     "/",
// 					SameSite: http.SameSiteNoneMode,
// 				},
// 			},
// 			wantErr: true,
// 			prepareMock: func(f *fields) {
// 				f.SessionStorage.EXPECT().DeleteSession(gomock.Any(), gomock.Any()).Return(errors.ErrNotFound)
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()
// 			f := fields{
// 				SessionStorage: mock_auth.NewMockSessionStorage(ctrl),
// 			}

// 			if tt.prepareMock != nil {
// 				tt.prepareMock(&f)
// 			}

// 			s := auth.NewService(nil, f.SessionStorage, nil)

// 			if err := s.Logout(tt.args.ctx, tt.args.session); (err != nil) != tt.wantErr {
// 				t.Errorf("Service.Logout() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestService_IsAuthorized(t *testing.T) {
// 	type fields struct {
// 		SessionStorage *mock_auth.MockSessionStorage
// 	}

// 	type args struct {
// 		ctx     context.Context
// 		session *http.Cookie
// 	}

// 	tests := []struct {
// 		name        string
// 		args        args
// 		prepareMock func(*fields)
// 		want        uint
// 		wantErr     bool
// 	}{
// 		{
// 			name: "success",
// 			args: args{
// 				ctx: context.Background(),
// 				session: &http.Cookie{
// 					Name:     "session_id",
// 					Value:    "session_id",
// 					MaxAge:   10 * 60 * 60,
// 					HttpOnly: true,
// 					Secure:   true,
// 					Path:     "/",
// 					SameSite: http.SameSiteNoneMode,
// 				},
// 			},
// 			want:    1,
// 			wantErr: false,
// 			prepareMock: func(f *fields) {
// 				f.SessionStorage.EXPECT().GetUserIDBySession(gomock.Any(), gomock.Any()).Return(uint(1), nil)
// 			},
// 		},
// 		{
// 			name: "failed to get user id by session",
// 			args: args{
// 				ctx: context.Background(),
// 				session: &http.Cookie{
// 					Name:     "session_id",
// 					Value:    "session_id",
// 					MaxAge:   10 * 60 * 60,
// 					HttpOnly: true,
// 					Secure:   true,
// 					Path:     "/",
// 					SameSite: http.SameSiteNoneMode,
// 				},
// 			},
// 			want:    0,
// 			wantErr: true,
// 			prepareMock: func(f *fields) {
// 				f.SessionStorage.EXPECT().GetUserIDBySession(gomock.Any(), gomock.Any()).Return(uint(0), errors.ErrNotFound)
// 			},
// 		},
// 		{
// 			name: "invalid session",
// 			args: args{
// 				ctx: context.Background(),
// 				session: &http.Cookie{
// 					Name:     "",
// 					Value:    "",
// 					MaxAge:   0,
// 					HttpOnly: true,
// 					Secure:   true,
// 					Path:     "/",
// 					SameSite: http.SameSiteNoneMode,
// 				},
// 			},
// 			want:    0,
// 			wantErr: true,
// 			prepareMock: func(f *fields) {
// 				f.SessionStorage.EXPECT().GetUserIDBySession(gomock.Any(), gomock.Any()).Return(uint(0), nil)
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()
// 			f := fields{
// 				SessionStorage: mock_auth.NewMockSessionStorage(ctrl),
// 			}

// 			if tt.prepareMock != nil {
// 				tt.prepareMock(&f)
// 			}

// 			s := auth.NewService(nil, f.SessionStorage, nil)

// 			got, err := s.IsAuthorized(tt.args.ctx, tt.args.session)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Service.IsAuthorized() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("Service.IsAuthorized() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
