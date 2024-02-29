package handlers

import (
	"encoding/json"
	"net/http"
	"socio/errors"
	"socio/services"
)

type PostsHandler struct {
	service *services.PostsService
}

func NewPostsHandler() (handler *PostsHandler) {
	handler = &PostsHandler{
		service: services.NewPostsService(),
	}
	return
}

func (api *PostsHandler) HandleListPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := api.service.ListPosts()
	if err != nil {
		errors.ServeHttpError(&w, err)
		return
	}

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		errors.ServeHttpError(&w, err)
		return
	}
}
