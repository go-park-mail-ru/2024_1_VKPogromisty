package routers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func MountPostsRouter(rootRouter *mux.Router) {
	r := rootRouter.PathPrefix("/posts").Subrouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("posts")
	})
}
