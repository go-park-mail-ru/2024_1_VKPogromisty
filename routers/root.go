package routers

import (
	"net/http"
	"socio/handlers"

	"github.com/gorilla/mux"
)

func SetUpCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func NewRootRouter() (rootRouter *mux.Router) {
	rootRouter = mux.NewRouter().PathPrefix("/api/v1/").Subrouter()
	rootRouter.Use(SetUpCORS)

	// need auth handler in post router to check if user is authenticated, will be removed when db is added
	authHandler := handlers.NewAuthHandler()
	MountAuthRouter(rootRouter, authHandler)

	MountPostsRouter(rootRouter, authHandler)
	MountStaticRouter(rootRouter)

	return
}
