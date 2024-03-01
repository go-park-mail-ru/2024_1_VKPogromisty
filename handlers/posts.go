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

// HandleListPosts godoc
//
//	@Summary		list all posts
//	@Description	list posts to authorized user
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Produce		json
//	@Success		200	{object}	utils.JSONResponse{body=services.ListPostsResponse}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/ [get]
func (api *PostsHandler) HandleListPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := api.service.ListPosts()
	if err != nil {
		utils.ServeJSONError(w, err)
		return
	}

	utils.ServeJSONBody(w, map[string][]services.PostWithAuthor{"posts": posts})
}
