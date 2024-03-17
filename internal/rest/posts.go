package rest

import (
	"net/http"
	"socio/pkg/json"
	"socio/usecase/posts"
)

type PostsHandler struct {
	Service *posts.Service
}

func NewPostsHandler(postsStorage posts.PostsStorage, usersStorage posts.UsersStorage) (handler *PostsHandler) {
	handler = &PostsHandler{
		Service: posts.NewPostsService(postsStorage, usersStorage),
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
	postsWithAuthors, err := api.Service.ListPosts()
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	json.ServeJSONBody(w, map[string][]posts.PostWithAuthor{"posts": postsWithAuthors})
}
