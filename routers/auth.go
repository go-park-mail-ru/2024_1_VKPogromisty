package routers

import (
	"socio/handlers"

	"github.com/gorilla/mux"
)

func MountAuthRouter(rootRouter *mux.Router) {
	r := rootRouter.PathPrefix("/auth").Subrouter()
	h := handlers.NewAuthHandler()

	r.HandleFunc("/login", h.HandleLogin).Methods("POST")
	r.HandleFunc("/signup", h.HandleRegistration).Methods("POST")
	r.HandleFunc("/logout", h.HandleLogout).Methods("DELETE")
}
