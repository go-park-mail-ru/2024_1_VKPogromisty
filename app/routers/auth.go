package routers

import (
	repository "socio/internal/repository/map"
	"socio/internal/rest"
	"socio/utils"

	"github.com/gorilla/mux"
)

func MountAuthRouter(rootRouter *mux.Router, userStorage *repository.Users, sessionStorage *repository.Sessions) {
	r := rootRouter.PathPrefix("/auth").Subrouter()
	h := rest.NewAuthHandler(utils.RealTimeProvider{}, userStorage, sessionStorage)

	r.HandleFunc("/login", h.HandleLogin).Methods("POST", "OPTIONS")
	r.HandleFunc("/signup", h.HandleRegistration).Methods("POST", "OPTIONS")
	r.HandleFunc("/logout", h.HandleLogout).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/is-authorized", h.CheckIsAuthorized).Methods("GET", "OPTIONS")
}
