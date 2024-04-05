package rest

import (
	defJSON "encoding/json"
	"fmt"
	"net/http"
	"socio/domain"
	"socio/errors"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
	"socio/pkg/sanitizer"
	"socio/usecase/posts"
	"strconv"
	"strings"
)

type ListUserPostsResponse struct {
	Posts  []*domain.Post `json:"posts"`
	Author *domain.User   `json:"author"`
}

type PostsHandler struct {
	Service *posts.Service
}

func NewPostsHandler(postsStorage posts.PostsStorage, usersStorage posts.UserStorage, sanitizer *sanitizer.Sanitizer) (handler *PostsHandler) {
	handler = &PostsHandler{
		Service: posts.NewPostsService(postsStorage, usersStorage, sanitizer),
	}
	return
}

// HandleGetUserPosts godoc
//
//	@Summary		get user posts
//	@Description	get user posts
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/get_user_posts
//	@Accept			json
//
//	@Param			Cookie		header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			userId		query	uint	true	"ID of the user"
//	@Param			lastPostId	query	uint	false	"ID of the last post, if 0 - get first posts"
//
//	@Produce		json
//	@Success		200	{object}	ListUserPostsResponse
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/ [get]
func (h *PostsHandler) HandleGetUserPosts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input posts.ListUserPostsInput

	userID, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	input.UserID = uint(userID)

	lastPostID, err := strconv.Atoi(r.URL.Query().Get("lastPostId"))
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	input.LastPostID = uint(lastPostID)

	posts, author, err := h.Service.GetUserPosts(r.Context(), input.UserID, input.LastPostID)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	response := ListUserPostsResponse{
		Posts:  posts,
		Author: author,
	}
	json.ServeJSONBody(r.Context(), w, response)
}

// HandleGetUserFriendsPosts godoc
//
//	@Summary		get user friends posts
//	@Description	get user friends posts
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/get_user_friends_posts
//	@Accept			json
//
//	@Param			Cookie		header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			lastPostId	query	uint	false	"ID of the last post"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=[]domain.PostWithAuthor}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/friends [get]
func (h *PostsHandler) HandleGetUserFriendsPosts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input posts.ListUserFriendsPostsInput

	lastPostID, err := strconv.Atoi(r.URL.Query().Get("lastPostId"))
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	input.LastPostID = uint(lastPostID)

	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	postsWithAuthors, err := h.Service.GetUserFriendsPosts(r.Context(), userID, input.LastPostID)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, postsWithAuthors)
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
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			content		formData	string	true	"Content of the post"
//	@Param			attachments	formData	file	false	"Attachments of the post"
//
//	@Produce		json
//	@Success		201	{object}	json.JSONResponse{body=domain.PostWithAuthor}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/ [post]
func (h *PostsHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1000 << 20)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidBody)
		return
	}

	var postInput posts.PostInput

	postInput.AuthorID, err = requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	postInput.Content = strings.Trim(r.PostFormValue("content"), " \n\r\t")

	for _, fileHeaders := range r.MultipartForm.File {
		postInput.Attachments = append(postInput.Attachments, fileHeaders...)
	}

	postWithAuthor, err := h.Service.CreatePost(r.Context(), postInput)
	if err != nil {
		fmt.Println(err)
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.ServeJSONBody(r.Context(), w, postWithAuthor)
}

// HandleUpdatePost godoc
//
//	@Summary		update post
//	@Description	update post by id
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/update
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			post_id	body	uint	true	"ID of the post"
//	@Param			content	body	string	true	"Content of the post"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=domain.Post}	"application/json"	"Attachments is always null!!!"
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/ [put]
func (h *PostsHandler) HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input posts.PostUpdateInput

	decoder := defJSON.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	updatedPost, err := h.Service.UpdatePost(r.Context(), userID, input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, updatedPost)

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
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			post_id	body	uint	true	"ID of the post"
//
//	@Produce		json
//	@Success		204	{object}	json.JSONResponse
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/ [delete]
func (h *PostsHandler) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input posts.DeletePostInput

	decoder := defJSON.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	err = h.Service.DeletePost(r.Context(), input.PostID)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
