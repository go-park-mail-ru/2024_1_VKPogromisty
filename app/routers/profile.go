package routers

import (
	"socio/internal/rest/middleware"
	rest "socio/internal/rest/profile"
	"socio/usecase/auth"
	"socio/usecase/profile"

	"github.com/gorilla/mux"
)

func MountProfileRouter(rootRouter *mux.Router, userStorage profile.UserStorage, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/profile").Subrouter()
	h := rest.NewProfileHandler(userStorage, sessionStorage)

	r.HandleFunc("/{userID}", h.HandleGetProfile).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.HandleGetProfile).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.HandleUpdateProfile).Methods("PUT", "OPTIONS")
	r.HandleFunc("/", h.HandleDeleteProfile).Methods("DELETE", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
}
