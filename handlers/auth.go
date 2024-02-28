package handlers

import (
	"encoding/json"
	"net/http"
	"socio/errors"
	"socio/services"
	"strings"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler() (handler *AuthHandler) {
	handler = &AuthHandler{
		service: services.NewAuthService(),
	}
	return
}

func (api *AuthHandler) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(4 * 1024 * 1024)
	if err != nil {
		errors.ServeHttpError(&w, errors.ErrInvalidBody)
		return
	}

	var regInput services.RegistrationInput
	regInput.FirstName = strings.Trim(r.PostFormValue("firstName"), " \n\r\t")
	regInput.LastName = strings.Trim(r.PostFormValue("lastName"), " \n\r\t")
	regInput.Email = strings.Trim(r.PostFormValue("email"), " \n\r\t")
	regInput.Password = r.PostFormValue("password")
	regInput.RepeatPassword = r.PostFormValue("repeatPassword")
	regInput.DateOfBirth = strings.Trim(r.PostFormValue("dateOfBirth"), " \n\r\t")
	_, regInput.Avatar, err = r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		errors.ServeHttpError(&w, errors.ErrInvalidBody)
		return
	}

	user, session, err := api.service.RegistrateUser(regInput)
	if err != nil {
		errors.ServeHttpError(&w, err)
		return
	}

	http.SetCookie(w, session)
	json.NewEncoder(w).Encode(user)
}

func (api *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	loginInput := new(services.LoginInput)
	err := decoder.Decode(loginInput)
	if err != nil {
		errors.ServeHttpError(&w, errors.ErrJSONUnmarshalling)
		return
	}

	session, err := api.service.Login(*loginInput)
	if err != nil {
		errors.ServeHttpError(&w, err)
		return
	}

	http.SetCookie(w, session)
	w.Write([]byte(session.Value))
}

func (api *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		errors.ServeHttpError(&w, err)
		return
	}

	if err = api.service.Logout(session); err != nil {
		errors.ServeHttpError(&w, err)
		return
	}

	http.SetCookie(w, session)
}

func (api *AuthHandler) CheckIsAuthorized(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			errors.ServeHttpError(&w, err)
			return
		}

		if err := api.service.IsAuthorized(session); err != nil {
			errors.ServeHttpError(&w, err)
			return
		}

		h.ServeHTTP(w, r)
	})
}
