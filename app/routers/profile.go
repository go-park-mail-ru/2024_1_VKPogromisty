package routers

import (
	"socio/internal/rest/middleware"
	rest "socio/internal/rest/profile"
	"socio/pkg/sanitizer"
	customtime "socio/pkg/time"
	"socio/usecase/auth"
	"socio/usecase/csrf"
	"socio/usecase/profile"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
)

func MountProfileRouter(rootRouter *mux.Router, userStorage profile.UserStorage, sessionStorage auth.SessionStorage) {
	r := rootRouter.PathPrefix("/profile").Subrouter()

	sanitizer := sanitizer.NewSanitizer(bluemonday.UGCPolicy())

	h := rest.NewProfileHandler(userStorage, sessionStorage, sanitizer)

	r.HandleFunc("/{userID}", h.HandleGetProfile).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.HandleGetProfile).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.HandleUpdateProfile).Methods("PUT", "OPTIONS")
	r.HandleFunc("/", h.HandleDeleteProfile).Methods("DELETE", "OPTIONS")
	r.Use(middleware.CreateCheckIsAuthorizedMiddleware(sessionStorage))
	r.Use(middleware.CreateCSRFMiddleware(csrf.NewCSRFService(customtime.RealTimeProvider{})))
}
