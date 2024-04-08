package routers

import (
	rest "socio/internal/rest/auth"
	"socio/pkg/sanitizer"
	"socio/usecase/auth"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
)

func MountAuthRouter(rootRouter *mux.Router, userStorage auth.UserStorage, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/auth").Subrouter()

	sanitizer := sanitizer.NewSanitizer(bluemonday.UGCPolicy())

	h := rest.NewAuthHandler(userStorage, sessionStorage, sanitizer)

	r.HandleFunc("/login", h.HandleLogin).Methods("POST", "OPTIONS")
	r.HandleFunc("/signup", h.HandleRegistration).Methods("POST", "OPTIONS")
	r.HandleFunc("/logout", h.HandleLogout).Methods("DELETE", "OPTIONS")
}
