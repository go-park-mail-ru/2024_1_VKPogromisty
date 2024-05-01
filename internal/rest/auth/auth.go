package rest

import (
	defJSON "encoding/json"
	"net/http"
	"socio/domain"
	"socio/errors"
	authpb "socio/internal/grpc/auth/proto"
	uspb "socio/internal/grpc/user/proto"
	"socio/internal/rest/uploaders"
	"socio/pkg/json"
	customtime "socio/pkg/time"
	"socio/usecase/auth"
	"socio/usecase/user"
	"strings"
)

type AuthHandler struct {
	AuthClient   authpb.AuthClient
	UserClient   uspb.UserClient
	TimeProvider customtime.TimeProvider
}

func NewAuthHandler(authClient authpb.AuthClient, userClient uspb.UserClient, tp customtime.TimeProvider) (handler *AuthHandler) {
	handler = &AuthHandler{
		AuthClient:   authClient,
		UserClient:   userClient,
		TimeProvider: tp,
	}
	return
}

func newSessionCookie(sessionID string) *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		MaxAge:   10 * 60 * 60,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	}
}

func clearSessionCookie(cookie *http.Cookie) {
	cookie.MaxAge = 0
	cookie.Value = ""
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteNoneMode
}

// HandleRegistration godoc
//
//	@Summary		handle user's registration flow
//	@Description	registrate user by his data
//	@Tags			auth
//	@license.name	Apache 2.0
//	@ID				auth/signup
//	@Accept			mpfd
//
//	@Param			firstName		formData	string	true	"First name of the user"
//	@Param			lastName		formData	string	true	"Last name of the user"
//	@Param			email			formData	string	true	"Email of the user"
//	@Param			password		formData	string	true	"Password of the user"			minLength(6)
//	@Param			repeatPassword	formData	string	true	"Repeat password of the user"	minLength(6)
//	@Param			dateOfBirth		formData	string	true	"Date of birth of the user"		format(date)	example(2021-01-01)
//	@Param			avatar			formData	file	false	"Avatar of the user"
//
//	@Produce		json
//	@Success		201	{object}	json.JSONResponse{body=domain.User}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Header			200	{string}	Set-Cookie	"session_id=some_session_id; Path=/; Max-Age=36000; HttpOnly;"
//	@Router			/auth/signup/ [post]
func (api *AuthHandler) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(4 * 1024 * 1024)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidBody)
		return
	}

	var regInput user.CreateUserInput
	regInput.FirstName = strings.TrimSpace(r.PostFormValue("firstName"))
	regInput.LastName = strings.TrimSpace(r.PostFormValue("lastName"))
	regInput.Email = strings.TrimSpace(r.PostFormValue("email"))
	regInput.Password = r.PostFormValue("password")
	regInput.RepeatPassword = r.PostFormValue("repeatPassword")
	regInput.DateOfBirth = strings.TrimSpace(r.PostFormValue("dateOfBirth"))
	_, avatarFH, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	regInput.Avatar, err = uploaders.UploadAvatar(r, api.UserClient, avatarFH)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	_, err = api.UserClient.Create(r.Context(), uspb.ToCreateRequest(&regInput))
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	res, err := api.AuthClient.Login(r.Context(), &authpb.LoginRequest{
		Email:    regInput.Email,
		Password: regInput.Password,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	sessionCookie := newSessionCookie(res.SessionId)

	http.SetCookie(w, sessionCookie)
	json.ServeJSONBody(r.Context(), w, map[string]*domain.User{
		"user": authpb.ToUser(res.User),
	}, http.StatusCreated)
}

// HandleLogin godoc
//
//	@Summary		handle user's login
//	@Description	login user by email and password
//	@Tags			auth
//	@license.name	Apache 2.0
//	@ID				auth/login
//	@Accept			json
//
//	@Param			email		body	string	true	"Email of the user"
//	@Param			password	body	string	true	"Password of the user"
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=domain.User}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//
//	@Header			200	{string}	Set-Cookie	"session_id=some_session_id; Path=/; Max-Age=36000; HttpOnly;"
//
//	@Router			/auth/login/ [post]
func (api *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidBody)
		return
	}

	defer r.Body.Close()

	decoder := defJSON.NewDecoder(r.Body)

	loginInput := new(auth.LoginInput)
	err := decoder.Decode(loginInput)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	res, err := api.AuthClient.Login(r.Context(), &authpb.LoginRequest{
		Email:    loginInput.Email,
		Password: loginInput.Password,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	sessionCookie := newSessionCookie(res.SessionId)

	http.SetCookie(w, sessionCookie)
	json.ServeJSONBody(r.Context(), w, map[string]any{"user": authpb.ToUser(res.User)}, http.StatusOK)
}

// HandleLogout godoc
//
//	@Summary		handle user's logout
//	@Description	logout user that is authorized
//	@Tags			auth
//	@license.name	Apache 2.0
//	@ID				auth/logout
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//
//	@Produce		json
//	@Success		200
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//
//	@Header			200	{string}	Set-Cookie	"session_id=some_session_id; Path=/; HttpOnly; Expires=Thu, 01 Jan 1970 00:00:00 GMT;"
//
//	@Router			/auth/logout/ [delete]
func (api *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	_, err = api.AuthClient.Logout(r.Context(), &authpb.LogoutRequest{
		SessionId: session.Value,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	clearSessionCookie(session)

	http.SetCookie(w, session)
}
