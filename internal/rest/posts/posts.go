package rest

import (
	defJSON "encoding/json"
	"net/http"
	"socio/domain"
	"socio/errors"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
	"socio/pkg/sanitizer"
	"socio/usecase/posts"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

const (
	UserIDQueryParam      = "userId"
	LastPostIDQueryParam  = "lastPostId"
	PostsAmountQueryParam = "postsAmount"
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

// HandleGetPostByID godoc
//
//	@Summary		get post by id
//
//	@Description	get post by id
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/get_post_by_id
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			postId	query	uint	true	"ID of the post"
//
//	@Produce		json
//	@Success		200	{object}	domain.Post
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/ [get]
func (h *PostsHandler) HandleGetPostByID(w http.ResponseWriter, r *http.Request) {
	postID, ok := mux.Vars(r)["postID"]
	if !ok {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	postIDData, err := strconv.Atoi(postID)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	post, err := h.Service.GetPostByID(r.Context(), uint(postIDData))
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, post, http.StatusOK)
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
//	@Param			postsAmount	query	uint	false	"Amount of posts to get, if 0 - get 20 posts"
//
//	@Produce		json
//	@Success		200	{object}	ListUserPostsResponse
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/ [get]
func (h *PostsHandler) HandleGetUserPosts(w http.ResponseWriter, r *http.Request) {
	var input posts.ListUserPostsInput

	userID, err := strconv.Atoi(r.URL.Query().Get(UserIDQueryParam))
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	input.UserID = uint(userID)

	lastPostIDData := r.URL.Query().Get(LastPostIDQueryParam)
	var lastPostID uint64

	if lastPostIDData == "" {
		lastPostID = 0
	} else {
		lastPostID, err = strconv.ParseUint(lastPostIDData, 0, 0)
		if err != nil {
			json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
			return
		}
	}

	input.LastPostID = uint(lastPostID)

	postsAmountData := r.URL.Query().Get(PostsAmountQueryParam)
	var postsAmount uint64

	if postsAmountData == "" {
		postsAmount = 0
	} else {
		postsAmount, err = strconv.ParseUint(postsAmountData, 0, 0)
		if err != nil {
			json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
			return
		}
	}

	input.PostsAmount = uint(postsAmount)

	posts, author, err := h.Service.GetUserPosts(r.Context(), input.UserID, input.LastPostID, input.PostsAmount)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	response := ListUserPostsResponse{
		Posts:  posts,
		Author: author,
	}
	json.ServeJSONBody(r.Context(), w, response, http.StatusOK)
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
//	@Param			postsAmount	query	uint	false	"Amount of posts to get, if 0 - get 20 posts"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=[]domain.PostWithAuthor}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/friends [get]
func (h *PostsHandler) HandleGetUserFriendsPosts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input posts.ListUserFriendsPostsInput

	lastPostIDData := r.URL.Query().Get(LastPostIDQueryParam)
	var lastPostID uint64
	var err error

	if lastPostIDData == "" {
		lastPostID = 0
	} else {
		lastPostID, err = strconv.ParseUint(lastPostIDData, 0, 0)
		if err != nil {
			json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
			return
		}
	}

	input.LastPostID = uint(lastPostID)

	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	postsAmountData := r.URL.Query().Get(PostsAmountQueryParam)
	var postsAmount uint64

	if postsAmountData == "" {
		postsAmount = 0
	} else {
		postsAmount, err = strconv.ParseUint(postsAmountData, 0, 0)
		if err != nil {
			json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
			return
		}
	}

	input.PostsAmount = uint(postsAmount)

	postsWithAuthors, err := h.Service.GetUserFriendsPosts(r.Context(), userID, input.LastPostID, input.PostsAmount)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, postsWithAuthors, http.StatusOK)
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
//	@Failure		403	{object}	errors.HTTPError
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
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, postWithAuthor, http.StatusCreated)
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

	json.ServeJSONBody(r.Context(), w, updatedPost, http.StatusOK)

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
//	@Failure		403	{object}	errors.HTTPError
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

	err = h.Service.DeletePost(r.Context(), input.PostID)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
