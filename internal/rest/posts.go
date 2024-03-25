package rest

import (
	defJSON "encoding/json"
	"fmt"
	"net/http"
	"socio/errors"
	"socio/internal/rest/middleware"
	"socio/pkg/json"
	"socio/usecase/posts"
	"strings"
)

type DeletePostInput struct {
	PostID uint `json:"post_id"`
}

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

// HandleDeletePost godoc
//
//	@Summary		delete post
//	@Description	delete post by id
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/delete
//	@Accept			json
//
//	@Param			Cookie		header		string	true	"session_id=some_session"
//	@Param			post_id		body		uint	true	"ID of the post"
//
//	@Produce		json
//	@Success		204	{object}	json.JSONResponse
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/ [delete]
func (h *PostsHandler) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input DeletePostInput

	decoder := defJSON.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		json.ServeJSONError(w, errors.ErrJSONUnmarshalling)
		return
	}

	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	err = h.Service.DeletePost(input.PostID)
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
