package routers

import (
	"socio/handlers"

	"github.com/gorilla/mux"
)

func MountAuthRouter(rootRouter *mux.Router, h *handlers.AuthHandler) {
	r := rootRouter.PathPrefix("/auth").Subrouter()

	r.HandleFunc("/login", h.HandleLogin).Methods("POST", "OPTIONS")
	r.HandleFunc("/signup", h.HandleRegistration).Methods("POST", "OPTIONS")
	r.HandleFunc("/logout", h.HandleLogout).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/is-authorized", h.CheckIsAuthorized).Methods("GET", "OPTIONS")
}
