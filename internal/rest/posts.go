package rest

import (
	"fmt"
	"net/http"
	"socio/errors"
	"socio/internal/rest/middleware"
	"socio/pkg/json"
	"socio/usecase/posts"
	"strings"
)

type PostsHandler struct {
	Service *posts.Service
}

func NewPostsHandler(postsStorage posts.PostsStorage, usersStorage posts.UserStorage) (handler *PostsHandler) {
	handler = &PostsHandler{
		Service: posts.NewPostsService(postsStorage, usersStorage),
	}
	return
}

// HandleCreatePost godoc
//
//	@Summary		create post
//	@Description	create post with attachments
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/create
//	@Accept			mpfd
//
//	@Param			Cookie		header		string	true	"session_id=some_session"
//	@Param			content		formData	string	true	"Content of the post"
//	@Param			attachments	formData	file	false	"Attachments of the post"
//
//	@Produce		json
//	@Success		201	{object}	json.JSONResponse{body=posts.PostWithAuthor}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/ [post]
func (h *PostsHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		json.ServeJSONError(w, errors.ErrInvalidBody)
		return
	}

	var postInput posts.PostInput

	postInput.AuthorID = r.Context().Value(middleware.UserIDKey).(uint)
	postInput.Content = strings.Trim(r.PostFormValue("content"), " \n\r\t")

	for _, fileHeaders := range r.MultipartForm.File {
		postInput.Attachments = append(postInput.Attachments, fileHeaders...)
	}

	postWithAuthor, err := h.Service.CreatePost(postInput)
	if err != nil {
		fmt.Println(err)
		json.ServeJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.ServeJSONBody(w, postWithAuthor)
}
