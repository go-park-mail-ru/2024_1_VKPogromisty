package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"socio/services"
	"time"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var regInput services.RegistrationInput
	regInput.FirstName = r.PostFormValue("firstName")
	regInput.LastName = r.PostFormValue("lastName")
	regInput.Email = r.PostFormValue("email")
	regInput.Password = r.PostFormValue("password")
	regInput.RepeatPassword = r.PostFormValue("repeatPassword")
	regInput.DateOfBirth = r.PostFormValue("dateOfBirth")
	_, regInput.Avatar, err = r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := api.service.RegistrateUser(regInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)

}

func (api *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	loginInput := new(services.LoginInput)
	err := decoder.Decode(loginInput)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionID, err := api.service.Login(*loginInput)
	if err != nil {
		log.Printf("login error: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	w.Write([]byte(sessionID))
}

func (api *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Error(w, "no session_id cookie", http.StatusUnauthorized)
		return
	}

	if err = api.service.Logout(sessionID.Value); err != nil {
		http.Error(w, "no session_id cookie", http.StatusUnauthorized)
		return
	}

	sessionID.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, sessionID)
}
