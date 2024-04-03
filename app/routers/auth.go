package routers

import (
	rest "socio/internal/rest/auth"
	"socio/usecase/auth"

	"github.com/gorilla/mux"
)

func MountAuthRouter(rootRouter *mux.Router, userStorage auth.UserStorage, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/auth").Subrouter()
	h := rest.NewAuthHandler(userStorage, sessionStorage)

	r.HandleFunc("/login", h.HandleLogin).Methods("POST", "OPTIONS")
	r.HandleFunc("/signup", h.HandleRegistration).Methods("POST", "OPTIONS")
	r.HandleFunc("/logout", h.HandleLogout).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/is-authorized", h.CheckIsAuthorized).Methods("GET", "OPTIONS")
}
