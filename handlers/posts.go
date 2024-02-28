package handlers

import (
	"net/http"
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
	w.Write([]byte("posts"))
}
