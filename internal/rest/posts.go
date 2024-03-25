package rest

import (
	defJSON "encoding/json"
	"fmt"
	"net/http"
	"socio/domain"
	"socio/errors"
	"socio/internal/rest/middleware"
	"socio/pkg/json"
	"socio/usecase/posts"
	"strings"
)

type ListUserPostsResponse struct {
	Posts  []*domain.Post `json:"posts"`
	Author *domain.User   `json:"author"`
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

// HandleGetUserPosts godoc
//
//	@Summary		get user posts
//	@Description	get user posts
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/get_user_posts
//	@Accept			json
//
//	@Param			Cookie		header		string	true	"session_id=some_session"
//	@Param			user_id		body		uint	true	"ID of the user"
//	@Param			last_post_id	body		uint	false	"ID of the last post"
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

	decoder := defJSON.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		json.ServeJSONError(w, errors.ErrJSONUnmarshalling)
		return
	}

	posts, author, err := h.Service.GetUserPosts(input.UserID, input.LastPostID)
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	response := ListUserPostsResponse{
		Posts:  posts,
		Author: author,
	}
	json.ServeJSONBody(w, response)
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
//	@Param			Cookie		header		string	true	"session_id=some_session"
//	@Param			user_id		body		uint	true	"ID of the user"
//	@Param			last_post_id	body		uint	false	"ID of the last post"
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

	decoder := defJSON.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		json.ServeJSONError(w, errors.ErrJSONUnmarshalling)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(uint)

	postsWithAuthors, err := h.Service.GetUserFriendsPosts(userID, input.LastPostID)
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	json.ServeJSONBody(w, postsWithAuthors)
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
//	@Success		201	{object}	json.JSONResponse{body=domain.PostWithAuthor}
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

// HandleUpdatePost godoc
//
//	@Summary		update post
//	@Description	update post by id
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/update
//	@Accept			json
//
//	@Param			Cookie		header		string	true	"session_id=some_session"
//	@Param			post_id		body		uint	true	"ID of the post"
//	@Param			content		body		string	true	"Content of the post"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=domain.Post}
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
		json.ServeJSONError(w, errors.ErrJSONUnmarshalling)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(uint)

	updatedPost, err := h.Service.UpdatePost(userID, input)
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	json.ServeJSONBody(w, updatedPost)

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

	var input posts.DeletePostInput

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
