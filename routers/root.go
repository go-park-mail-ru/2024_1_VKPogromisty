package routers

import (
	"net/http"
	"socio/handlers"
	"socio/utils"

	"github.com/gorilla/mux"
)

func SetUpCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", utils.ALLOWED_ORIGIN)
		w.Header().Set("Access-Control-Allow-Headers", utils.ALLOWED_HEADERS)
		w.Header().Set("Access-Control-Allow-Methods", utils.ALLOWED_METHODS)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func DisableCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		h.ServeHTTP(w, r)
	})
}

func NewRootRouter() (rootRouter *mux.Router) {
	rootRouter = mux.NewRouter().PathPrefix("/api/v1/").Subrouter()

	// need auth handler in post router to check if user is authenticated, will be removed when db is added
	authHandler := handlers.NewAuthHandler(utils.RealTimeProvider{})
	MountAuthRouter(rootRouter, authHandler)
	MountPostsRouter(rootRouter, authHandler)
	MountStaticRouter(rootRouter)

	rootRouter.Use(SetUpCORS)
	rootRouter.Use(DisableCache)

	return
}
