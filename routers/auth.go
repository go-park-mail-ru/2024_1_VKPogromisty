package routers

import (
	"socio/handlers"

	"github.com/gorilla/mux"
)

func MountAuthRouter(rootRouter *mux.Router, h *handlers.AuthHandler) {
	r := rootRouter.PathPrefix("/api/v1/auth").Subrouter()

	r.HandleFunc("/login", h.HandleLogin).Methods("POST")
	r.HandleFunc("/signup", h.HandleRegistration).Methods("POST")
	r.HandleFunc("/logout", h.HandleLogout).Methods("DELETE")
}
