package handlers

import (
	"net/http"
	"socio/services"
	"socio/utils"
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
		utils.ServeJSONError(w, err)
		return
	}

	utils.ServeJSONBody(w, map[string][]services.PostWithAuthor{"posts": posts})
}
