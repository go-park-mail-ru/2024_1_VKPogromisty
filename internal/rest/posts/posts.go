package rest

import (
	defJSON "encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"socio/domain"
	"socio/errors"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
	"socio/usecase/posts"
	"strconv"
	"strings"

	postspb "socio/internal/grpc/post/proto"
	uspb "socio/internal/grpc/user/proto"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	UserIDQueryParam      = "userId"
	LastPostIDQueryParam  = "lastPostId"
	PostsAmountQueryParam = "postsAmount"
	LastLikeIDQueryParam  = "lastLikeId"
	BatchSize             = 1 << 23
)

type ListUserPostsResponse struct {
	Posts  []*domain.Post `json:"posts"`
	Author *domain.User   `json:"author"`
}

type PostsHandler struct {
	PostsClient postspb.PostClient
	UserClient  uspb.UserClient
}

func NewPostsHandler(postsClient postspb.PostClient, userClient uspb.UserClient) (handler *PostsHandler) {
	handler = &PostsHandler{
		PostsClient: postsClient,
		UserClient:  userClient,
	}
	return
}

func (h *PostsHandler) uploadAvatar(r *http.Request, fh *multipart.FileHeader) (string, error) {
	fileName := uuid.NewString() + filepath.Ext(fh.Filename)
	stream, err := h.PostsClient.Upload(r.Context())
	if err != nil {
		return "", err
	}

	file, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, BatchSize)
	batchNumber := 1

	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		chunk := buf[:num]

		err = stream.Send(&postspb.UploadRequest{
			FileName: fileName,
			Chunk:    chunk,
		})

		if err != nil {
			return "", err
		}
		batchNumber += 1
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}

	return res.FileName, nil
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

	post, err := h.PostsClient.GetPostByID(r.Context(), &postspb.GetPostByIDRequest{
		PostId: uint64(postIDData),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
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

	posts, err := h.PostsClient.GetUserPosts(r.Context(), &postspb.GetUserPostsRequest{
		UserId:      uint64(input.UserID),
		LastPostId:  uint64(input.LastPostID),
		PostsAmount: uint64(input.PostsAmount),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	author, err := h.UserClient.GetByID(r.Context(), &uspb.GetByIDRequest{
		UserId: uint64(input.UserID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	response := ListUserPostsResponse{
		Posts:  postspb.ToPosts(posts),
		Author: uspb.ToUser(author.User),
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

	posts, err := h.PostsClient.GetUserFriendsPosts(r.Context(), &postspb.GetUserFriendsPostsRequest{
		UserId:      uint64(userID),
		LastPostId:  uint64(input.LastPostID),
		PostsAmount: uint64(input.PostsAmount),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	postsWithAuthors := make([]*domain.PostWithAuthor, 0, len(posts.Posts))
	for _, post := range posts.Posts {
		author, err := h.UserClient.GetByID(r.Context(), &uspb.GetByIDRequest{
			UserId: post.AuthorId,
		})
		if err != nil {
			json.ServeGRPCStatus(r.Context(), w, err)
			return
		}

		postsWithAuthors = append(postsWithAuthors, &domain.PostWithAuthor{
			Post:   postspb.ToPost(post),
			Author: uspb.ToUser(author.User),
		})
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

	postInput.Content = strings.TrimSpace(r.PostFormValue("content"))

	for _, fh := range r.MultipartForm.File {
		for _, f := range fh {
			fileName, err := h.uploadAvatar(r, f)
			if err != nil {
				json.ServeJSONError(r.Context(), w, err)
				return
			}
			postInput.Attachments = append(postInput.Attachments, fileName)
		}
	}

	postData, err := h.PostsClient.CreatePost(r.Context(), &postspb.CreatePostRequest{
		AuthorId:    uint64(postInput.AuthorID),
		Content:     postInput.Content,
		Attachments: postInput.Attachments,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	post := postspb.ToPost(postData.Post)

	author, err := h.UserClient.GetByID(r.Context(), &uspb.GetByIDRequest{
		UserId: uint64(post.AuthorID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	postWithAuthor := &domain.PostWithAuthor{
		Post:   post,
		Author: uspb.ToUser(author.User),
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

	updatedPost, err := h.PostsClient.UpdatePost(r.Context(), &postspb.UpdatePostRequest{
		PostId:  uint64(input.PostID),
		Content: input.Content,
		UserId:  uint64(userID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, postspb.ToPost(updatedPost.Post), http.StatusOK)

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

	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	_, err = h.PostsClient.DeletePost(r.Context(), &postspb.DeletePostRequest{
		PostId: uint64(input.PostID),
		UserId: uint64(userID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PostsHandler) HandleGetLikedPosts(w http.ResponseWriter, r *http.Request) {
	authorizedUserID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	lastLikeIDData := r.URL.Query().Get(LastLikeIDQueryParam)
	var lastLikeID uint64

	if lastLikeIDData == "" {
		lastLikeID = 0
	} else {
		lastLikeID, err = strconv.ParseUint(lastLikeIDData, 0, 0)
		if err != nil {
			json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
			return
		}
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

	postsRes, err := h.PostsClient.GetLikedPosts(r.Context(), &postspb.GetLikedPostsRequest{
		UserId:      uint64(authorizedUserID),
		LastLikeId:  lastLikeID,
		PostsAmount: postsAmount,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	likedPosts := postspb.ToLikesWithPosts(postsRes)

	var likedPostsWithUsers []posts.LikeWithPostAndUser

	for _, post := range likedPosts {
		var likedPostWithUser posts.LikeWithPostAndUser

		author, err := h.UserClient.GetByID(r.Context(), &uspb.GetByIDRequest{
			UserId: uint64(post.Like.UserID),
		})
		if err != nil {
			json.ServeGRPCStatus(r.Context(), w, err)
			return
		}

		likedPostWithUser.Like = post.Like
		likedPostWithUser.Post = post.Post
		likedPostWithUser.User = uspb.ToUser(author.User)

		likedPostsWithUsers = append(likedPostsWithUsers, likedPostWithUser)
	}

	json.ServeJSONBody(r.Context(), w, likedPostsWithUsers, http.StatusOK)
}

func (h *PostsHandler) HandleLikePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input *domain.PostLike

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

	res, err := h.PostsClient.LikePost(r.Context(), &postspb.LikePostRequest{
		PostId: uint64(input.PostID),
		UserId: uint64(userID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, res, http.StatusCreated)
}

func (h *PostsHandler) HandleUnlikePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input *domain.PostLike

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

	_, err = h.PostsClient.UnlikePost(r.Context(), &postspb.UnlikePostRequest{
		PostId: uint64(input.PostID),
		UserId: uint64(userID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
