package rest

import (
	"encoding/json"
	"net/http"
	"socio/domain"
	"socio/errors"
	repository "socio/internal/repository/map"
	"socio/usecase/auth"
	"socio/utils"
	"strings"
)

type AuthHandler struct {
	Service      *auth.Service
	TimeProvider utils.TimeProvider
}

func NewAuthHandler(tp utils.TimeProvider, userStorage *repository.Users, sessionStorage *repository.Sessions) (handler *AuthHandler) {
	handler = &AuthHandler{
		Service: auth.NewService(tp, userStorage, sessionStorage),
	}
	return
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
//	@Success		201	{object}	utils.JSONResponse{body=auth.User}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Header			200	{string}	Set-Cookie	"session_id=some_session_id; Path=/; Max-Age=36000; HttpOnly;"
//	@Router			/auth/signup/ [post]
func (api *AuthHandler) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(4 * 1024 * 1024)
	if err != nil {
		utils.ServeJSONError(w, errors.ErrInvalidBody)
		return
	}

	var regInput auth.RegistrationInput
	regInput.FirstName = strings.Trim(r.PostFormValue("firstName"), " \n\r\t")
	regInput.LastName = strings.Trim(r.PostFormValue("lastName"), " \n\r\t")
	regInput.Email = strings.Trim(r.PostFormValue("email"), " \n\r\t")
	regInput.Password = r.PostFormValue("password")
	regInput.RepeatPassword = r.PostFormValue("repeatPassword")
	regInput.DateOfBirth = strings.Trim(r.PostFormValue("dateOfBirth"), " \n\r\t")
	_, regInput.Avatar, err = r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		utils.ServeJSONError(w, err)
		return
	}

	user, session, err := api.Service.RegistrateUser(regInput)
	if err != nil {
		utils.ServeJSONError(w, err)
		return
	}

	http.SetCookie(w, session)
	w.WriteHeader(http.StatusCreated)
	utils.ServeJSONBody(w, map[string]*domain.User{"user": user})
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
//	@Success		200	{object}	utils.JSONResponse{body=auth.LoginResponse}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//
//	@Header			200	{string}	Set-Cookie	"session_id=some_session_id; Path=/; Max-Age=36000; HttpOnly;"
//
//	@Router			/auth/login/ [post]
func (api *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	loginInput := new(auth.LoginInput)
	err := decoder.Decode(loginInput)
	if err != nil {
		utils.ServeJSONError(w, errors.ErrJSONUnmarshalling)
		return
	}

	user, session, err := api.Service.Login(*loginInput)
	if err != nil {
		utils.ServeJSONError(w, err)
		return
	}

	http.SetCookie(w, session)
	utils.ServeJSONBody(w, map[string]any{"sessionId": session.Value, "user": user})
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
		utils.ServeJSONError(w, err)
		return
	}

	if err = api.Service.Logout(session); err != nil {
		utils.ServeJSONError(w, err)
		return
	}

	http.SetCookie(w, session)
}

// CheckIsAuthorized godoc
//
//	@Summary		check if user is authorized
//	@Tags			auth
//	@license.name	Apache 2.0
//	@ID				auth/is-authorized
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//
//	@Produce		json
//	@Success		200 {object}	utils.JSONResponse{body=auth.IsAuthorizedResponse}
//	@Failure		500
//
//	@Header			200	{string}	Set-Cookie	"session_id=some_session_id; Path=/; HttpOnly; Expires=Thu, 01 Jan 1970 00:00:00 GMT;"
//
//	@Router			/auth/is-authorized [get]
func (api *AuthHandler) CheckIsAuthorized(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		utils.ServeJSONBody(w, map[string]bool{"isAuthorized": false})
		return
	}

	if _, err := api.Service.IsAuthorized(session); err != nil {
		utils.ServeJSONBody(w, map[string]bool{"isAuthorized": false})
		return
	}

	utils.ServeJSONBody(w, map[string]bool{"isAuthorized": true})
}
