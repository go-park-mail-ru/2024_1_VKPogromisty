package rest

import (
	"net/http"
	"socio/domain"
	"socio/errors"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
	"socio/usecase/posts"
	"strconv"
	"strings"

	postspb "socio/internal/grpc/post/proto"
	pgpb "socio/internal/grpc/public_group/proto"
	uspb "socio/internal/grpc/user/proto"
	"socio/internal/rest/uploaders"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
)

const (
	UserIDQueryParam      = "userId"
	LastPostIDQueryParam  = "lastPostId"
	PostsAmountQueryParam = "postsAmount"
	LastLikeIDQueryParam  = "lastLikeId"
	BatchSize             = 1 << 23
)

//easyjson:json
type LikePostInput struct {
	PostID uint `json:"postId"`
}

//easyjson:json
type UnlikePostInput struct {
	PostID uint `json:"postId"`
}

//easyjson:json
type ListUserPostsResponse struct {
	Posts  []*domain.Post `json:"posts"`
	Author *domain.User   `json:"author"`
}

//easyjson:json
type CreateCommentInput struct {
	PostID  uint   `json:"postId"`
	Content string `json:"content"`
}

//easyjson:json
type UpdateCommentInput struct {
	CommentID uint   `json:"commentId"`
	Content   string `json:"content"`
}

//easyjson:json
type DeleteCommentInput struct {
	CommentID uint `json:"commentId"`
}

//easyjson:json
type LikeCommentInput struct {
	CommentID uint `json:"commentId"`
}

//easyjson:json
type UnlikeCommentInput struct {
	CommentID uint `json:"commentId"`
}

type PostsHandler struct {
	PostsClient       postspb.PostClient
	UserClient        uspb.UserClient
	PublicGroupClient pgpb.PublicGroupClient
}

func NewPostsHandler(postsClient postspb.PostClient, userClient uspb.UserClient, publicGroupClient pgpb.PublicGroupClient) (handler *PostsHandler) {
	handler = &PostsHandler{
		PostsClient:       postsClient,
		UserClient:        userClient,
		PublicGroupClient: publicGroupClient,
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
//	@Success		200	{object}	json.JSONResponse{body=domain.PostWithAuthorAndGroup}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/{id} [get]
func (h *PostsHandler) HandleGetPostByID(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrUnauthorized)
		return
	}

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

	postData, err := h.PostsClient.GetPostByID(r.Context(), &postspb.GetPostByIDRequest{
		PostId: uint64(postIDData),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	post := postspb.ToPost(postData.GetPost())

	authorData, err := h.UserClient.GetByID(r.Context(), &uspb.GetByIDRequest{
		UserId: uint64(post.AuthorID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	groupPostData, err := h.PostsClient.GetGroupPostByPostID(r.Context(), &postspb.GetGroupPostByPostIDRequest{
		PostId: uint64(post.ID),
	})
	if err != nil {
		_, code := errors.ParseGRPCError(err)
		if code != http.StatusNotFound {
			json.ServeGRPCStatus(r.Context(), w, err)
			return
		}
	}

	var group *domain.PublicGroup

	if err == nil {
		groupData, err := h.PublicGroupClient.GetByID(r.Context(), &pgpb.GetByIDRequest{
			Id:     groupPostData.GroupPost.GroupId,
			UserId: uint64(userID),
		})
		if err != nil {
			json.ServeJSONError(r.Context(), w, err)
			return
		}

		group = pgpb.ToPublicGroup(groupData.PublicGroup.PublicGroup)
	} else {
		group = nil
	}

	postWithData := &domain.PostWithAuthorAndGroup{
		Post:   post,
		Author: uspb.ToUser(authorData.GetUser()),
		Group:  group,
	}

	json.ServeJSONBody(r.Context(), w, postWithData, http.StatusOK)
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
		Posts:  postspb.ToPosts(posts.GetPosts()),
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
			fileName, err := uploaders.UploadPostAttachment(r, h.PostsClient, f)
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
//	@Accept			mpfd
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			postId	formData	uint	true	"ID of the post"
//	@Param			content	formData	string	true	"Content of the post"
//	@Param			attachmentsToDelete	formData	[]string	false	"Attachments to delete"
//	@Param			attachments	formData	[]file	false	"Attachments of the post"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=map[string]domain.Post}	"application/json"	"Attachments is always null!!!"
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/ [put]
func (h *PostsHandler) HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1000 << 20)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidBody)
		return
	}

	var input posts.PostUpdateInput

	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	postIDData, err := strconv.ParseUint(r.PostFormValue("postId"), 0, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	input.PostID = uint(postIDData)
	input.Content = strings.TrimSpace(r.PostFormValue("content"))
	input.AttachmentsToDelete = r.PostForm["attachmentsToDelete"]

	for _, f := range r.MultipartForm.File["attachments"] {
		fileName, err := uploaders.UploadPostAttachment(r, h.PostsClient, f)
		if err != nil {
			json.ServeJSONError(r.Context(), w, err)
			return
		}

		input.AttachmentsToAdd = append(input.AttachmentsToAdd, fileName)
	}

	updatedPost, err := h.PostsClient.UpdatePost(r.Context(), &postspb.UpdatePostRequest{
		PostId:              uint64(input.PostID),
		Content:             input.Content,
		UserId:              uint64(userID),
		AttachmentsToAdd:    input.AttachmentsToAdd,
		AttachmentsToDelete: input.AttachmentsToDelete,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, map[string]*domain.Post{
		"post": postspb.ToPost(updatedPost.Post),
	}, http.StatusOK)
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

	input := new(posts.DeletePostInput)

	err := easyjson.UnmarshalFromReader(r.Body, input)
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

// HandleGetLikedPosts godoc
//
//	@Summary		get liked posts
//	@Description	get posts that are authored by authorized user and liked by some people,
//	@Description	for every like it returns post and user that liked that post
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/get_liked_posts
//	@Accept			json
//
//	@Param			Cookie		header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			lastLikeId	query	uint	false	"ID of the last like, if 0 - get first likes"
//	@Param			postsAmount	query	uint	false	"Amount of liked posts to get, if 0 - get 20 liked posts"
//
//	@Produce		json
//	@Success		200	{object}	[]posts.LikeWithPostAndUser
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/liked [get]
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

// HandleLikePost godoc
//
//	@Summary		like post
//	@Description	like post
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/like
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			postId	body	uint	true	"ID of the post"
//
//	@Produce		json
//	@Success		201	{object}	domain.PostLike
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/like [post]
func (h *PostsHandler) HandleLikePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	input := new(LikePostInput)

	err := easyjson.UnmarshalFromReader(r.Body, input)
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

	json.ServeJSONBody(r.Context(), w, postspb.ToPostLike(res.Like), http.StatusCreated)
}

// HandleUnlikePost godoc
//
//	@Summary		unlike post
//	@Description	unlike post
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/unlike
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			postId	body	uint	true	"ID of the post"
//
//	@Produce		json
//	@Success		204	{object}	json.JSONResponse
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/unlike [delete]
func (h *PostsHandler) HandleUnlikePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	input := new(UnlikePostInput)

	err := easyjson.UnmarshalFromReader(r.Body, input)
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

// HandleGetGroupPostsBySubscriptions godoc
//
//	@Summary		get group posts by subscriptions
//	@Description	get group posts by subscriptions
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/get_group_posts_by_subscriptions
//	@Accept			json
//
//	@Param			Cookie		header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			lastPostId	query	uint	false	"ID of the last post, if 0 - get first posts"
//	@Param			postsAmount	query	uint	false	"Amount of posts to get, if 0 - get 20 posts"
//
//	@Produce		json
//	@Success		200	{object}	[]domain.PostWithAuthorAndGroup
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/groups [get]
func (h *PostsHandler) HandleGetGroupPostsBySubscriptions(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

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

	subscriptions, err := h.PublicGroupClient.GetBySubscriberID(r.Context(), &pgpb.GetBySubscriberIDRequest{
		SubscriberId: uint64(userID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	subIDs := make([]uint64, 0, len(subscriptions.GetPublicGroups()))
	groupsByID := map[uint64]*domain.PublicGroup{}
	for _, sub := range subscriptions.GetPublicGroups() {
		subIDs = append(subIDs, sub.Id)
		groupsByID[sub.Id] = pgpb.ToPublicGroup(sub)
	}

	postsRes, err := h.PostsClient.GetGroupPostsBySubscriptionIDs(r.Context(), &postspb.GetGroupPostsBySubscriptionIDsRequest{
		SubscriptionIds: subIDs,
		LastPostId:      lastPostID,
		PostsAmount:     postsAmount,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	posts := postspb.ToPosts(postsRes.GetPosts())
	res := make([]*domain.PostWithAuthorAndGroup, 0, len(posts))
	for _, post := range posts {
		postWithAuthorAndGroup := new(domain.PostWithAuthorAndGroup)
		author, err := h.UserClient.GetByID(r.Context(), &uspb.GetByIDRequest{
			UserId: uint64(post.AuthorID),
		})
		if err != nil {
			json.ServeGRPCStatus(r.Context(), w, err)
			return
		}

		postWithAuthorAndGroup.Author = uspb.ToUser(author.User)
		postWithAuthorAndGroup.Group = groupsByID[uint64(post.GroupID)]
		postWithAuthorAndGroup.Post = post

		res = append(res, postWithAuthorAndGroup)
	}

	json.ServeJSONBody(r.Context(), w, res, http.StatusOK)
}

// HandleGetPostsByGroupSubIDsAndUserSubIDs godoc
//
//		@Summary		get posts by group subscriptions and user subscriptions
//		@Description	get posts by group subscriptions and user subscriptions
//		@Tags			posts
//		@license.name	Apache 2.0
//		@ID				posts/get_posts_by_group_subscriptions_and_user_subscriptions
//		@Accept			json
//
//		@Param			Cookie		header	string	true	"session_id=some_session"
//		@Param			X-CSRF-Token	header	string	true	"CSRF token"
//		@Param			lastPostId	query	uint	false	"ID of the last post, if 0 - get first posts"
//		@Param			postsAmount	query	uint	false	"Amount of posts to get, if 0 - get 20 posts"
//
//		@Produce		json
//		@Success		200	{object}	[]domain.PostWithAuthorAndGroup
//		@Failure		400	{object}	errors.HTTPError
//		@Failure		401	{object}	errors.HTTPError
//		@Failure		403	{object}	errors.HTTPError
//	 @Failure		404	{object}	errors.HTTPError
//		@Failure		500	{object}	errors.HTTPError
//		@Router			/posts/all [get]
func (h *PostsHandler) HandleGetPostsByGroupSubIDsAndUserSubIDs(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

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

	groupSubs, err := h.PublicGroupClient.GetBySubscriberID(r.Context(), &pgpb.GetBySubscriberIDRequest{
		SubscriberId: uint64(userID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	groupSubIDs := make([]uint64, 0, len(groupSubs.GetPublicGroups()))
	groupsByID := map[uint64]*domain.PublicGroup{}
	for _, sub := range groupSubs.GetPublicGroups() {
		groupSubIDs = append(groupSubIDs, sub.Id)
		groupsByID[sub.Id] = pgpb.ToPublicGroup(sub)
	}

	userSubIDsRes, err := h.UserClient.GetSubscriptionIDs(r.Context(), &uspb.GetSubscriptionIDsRequest{
		UserId: uint64(userID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	userSubIDs := userSubIDsRes.GetSubscriptionIds()

	postsRes, err := h.PostsClient.GetPostsByGroupSubIDsAndUserSubIDs(r.Context(), &postspb.GetPostsByGroupSubIDsAndUserSubIDsRequest{
		GroupSubscriptionIds: groupSubIDs,
		UserSubscriptionIds:  userSubIDs,
		LastPostId:           lastPostID,
		PostsAmount:          postsAmount,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	posts := postspb.ToPosts(postsRes.GetPosts())
	res := make([]*domain.PostWithAuthorAndGroup, 0, len(posts))
	for _, post := range posts {
		postWithAuthorAndGroup := new(domain.PostWithAuthorAndGroup)
		author, err := h.UserClient.GetByID(r.Context(), &uspb.GetByIDRequest{
			UserId: uint64(post.AuthorID),
		})
		if err != nil {
			json.ServeGRPCStatus(r.Context(), w, err)
			return
		}

		postWithAuthorAndGroup.Author = uspb.ToUser(author.User)
		if post.GroupID != 0 {
			postWithAuthorAndGroup.Group = groupsByID[uint64(post.GroupID)]
		}
		postWithAuthorAndGroup.Post = post

		res = append(res, postWithAuthorAndGroup)
	}

	json.ServeJSONBody(r.Context(), w, res, http.StatusOK)
}

// HandleGetNewPosts godoc
//
//	@Summary		get new posts
//	@Description	get new posts
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/get_new_posts
//	@Accept			json
//
//	@Param			Cookie		header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			lastPostId	query	uint	false	"ID of the last post, if 0 - get first posts"
//	@Param			postsAmount	query	uint	false	"Amount of posts to get, if 0 - get 20 posts"
//
//	@Produce		json
//	@Success		200	{object}	[]domain.PostWithAuthorAndGroup
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/new [get]
func (h *PostsHandler) HandleGetNewPosts(w http.ResponseWriter, r *http.Request) {
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

	postsRes, err := h.PostsClient.GetNewPosts(r.Context(), &postspb.GetNewPostsRequest{
		LastPostId:  lastPostID,
		PostsAmount: postsAmount,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	posts := postspb.ToPosts(postsRes.GetPosts())
	res := make([]*domain.PostWithAuthorAndGroup, 0, len(posts))
	for _, post := range posts {
		postWithAuthorAndGroup := new(domain.PostWithAuthorAndGroup)
		author, err := h.UserClient.GetByID(r.Context(), &uspb.GetByIDRequest{
			UserId: uint64(post.AuthorID),
		})
		if err != nil {
			json.ServeGRPCStatus(r.Context(), w, err)
			return
		}

		postWithAuthorAndGroup.Author = uspb.ToUser(author.User)
		if post.GroupID != 0 {
			group, err := h.PublicGroupClient.GetByID(r.Context(), &pgpb.GetByIDRequest{
				Id: uint64(post.GroupID),
			})
			if err != nil {
				json.ServeGRPCStatus(r.Context(), w, err)
				return
			}
			postWithAuthorAndGroup.Group = pgpb.ToPublicGroup(group.PublicGroup.PublicGroup)
		}

		postWithAuthorAndGroup.Post = post

		res = append(res, postWithAuthorAndGroup)
	}

	json.ServeJSONBody(r.Context(), w, res, http.StatusOK)
}

// HandleGetPostByID godoc
//
//	@Summary		get comments by post id
//	@Description	get comments by post id
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/get_comments_by_post_id
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			postID	path	uint	true	"ID of the post"
//
//	@Produce		json
//	@Success		200	{object}	[]domain.CommentWithAuthor
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/{postID}/comments [get]
func (h *PostsHandler) HandleGetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
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

	comments, err := h.PostsClient.GetCommentsByPostID(r.Context(), &postspb.GetCommentsByPostIDRequest{
		PostId: uint64(postIDData),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	commentsWithAuthors := make([]*domain.CommentWithAuthor, 0, len(comments.Comments))

	for _, comment := range postspb.ToComments(comments.GetComments()) {
		commentAuthorData, err := h.UserClient.GetByID(r.Context(), &uspb.GetByIDRequest{
			UserId: uint64(comment.AuthorID),
		})
		if err != nil {
			json.ServeGRPCStatus(r.Context(), w, err)
			return
		}

		commentsWithAuthors = append(commentsWithAuthors, &domain.CommentWithAuthor{
			Comment: comment,
			Author:  uspb.ToUser(commentAuthorData.GetUser()),
		})
	}

	json.ServeJSONBody(r.Context(), w, commentsWithAuthors, http.StatusOK)
}

// HandleCreateComment godoc
//
//	@Summary		create comment
//	@Description	create comment
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/create_comment
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			content	body	string	true	"Content of the comment"
//	@Param			postID	body	uint	true	"ID of the post"
//
//	@Produce		json
//	@Success		201	{object}	domain.CommentWithAuthor
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/comment [post]
func (h *PostsHandler) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	commentInput := new(CreateCommentInput)

	err = easyjson.UnmarshalFromReader(r.Body, commentInput)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	commentData, err := h.PostsClient.CreateComment(r.Context(), &postspb.CreateCommentRequest{
		AuthorId: uint64(userID),
		Content:  commentInput.Content,
		PostId:   uint64(commentInput.PostID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	comment := postspb.ToComment(commentData.Comment)

	author, err := h.UserClient.GetByID(r.Context(), &uspb.GetByIDRequest{
		UserId: uint64(comment.AuthorID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	commentWithAuthor := &domain.CommentWithAuthor{
		Comment: comment,
		Author:  uspb.ToUser(author.User),
	}

	json.ServeJSONBody(r.Context(), w, commentWithAuthor, http.StatusCreated)
}

// HandleUpdateComment godoc
//
//	@Summary		update comment
//	@Description	update comment
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/update_comment
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			commentID	body	uint	true	"ID of the comment"
//	@Param			content	body	string	true	"Content of the comment"
//
//	@Produce		json
//	@Success		200	{object}	domain.CommentWithAuthor
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/comment [put]
func (h *PostsHandler) HandleUpdateComment(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	commentInput := new(UpdateCommentInput)

	err = easyjson.UnmarshalFromReader(r.Body, commentInput)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	commentData, err := h.PostsClient.UpdateComment(r.Context(), &postspb.UpdateCommentRequest{
		CommentId: uint64(commentInput.CommentID),
		Content:   commentInput.Content,
		UserId:    uint64(userID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	comment := postspb.ToComment(commentData.Comment)

	author, err := h.UserClient.GetByID(r.Context(), &uspb.GetByIDRequest{
		UserId: uint64(comment.AuthorID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	commentWithAuthor := &domain.CommentWithAuthor{
		Comment: comment,
		Author:  uspb.ToUser(author.User),
	}

	json.ServeJSONBody(r.Context(), w, commentWithAuthor, http.StatusOK)
}

// HandleDeleteComment godoc
//
//	@Summary		delete comment
//	@Description	delete comment
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/delete_comment
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			commentId	body	uint	true	"ID of the comment"
//
//	@Produce		json
//	@Success		204	{object}	json.JSONResponse
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/comment [delete]
func (h *PostsHandler) HandleDeleteComment(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	commentInput := new(DeleteCommentInput)

	err = easyjson.UnmarshalFromReader(r.Body, commentInput)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	_, err = h.PostsClient.DeleteComment(r.Context(), &postspb.DeleteCommentRequest{
		CommentId: uint64(commentInput.CommentID),
		UserId:    uint64(userID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleLikeComment godoc
//
//	@Summary		like comment
//	@Description	like comment
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/like_comment
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			commentID	body	uint	true	"ID of the comment"
//
//	@Produce		json
//	@Success		201	{object}	domain.CommentLike
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/comments/like [post]
func (h *PostsHandler) HandleLikeComment(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	input := new(LikeCommentInput)

	err = easyjson.UnmarshalFromReader(r.Body, input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	res, err := h.PostsClient.LikeComment(r.Context(), &postspb.LikeCommentRequest{
		CommentId: uint64(input.CommentID),
		UserId:    uint64(userID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, postspb.ToCommentLike(res.Like), http.StatusCreated)
}

// HandleUnlikeComment godoc
//
//	@Summary		unlike comment
//	@Description	unlike comment
//	@Tags			posts
//	@license.name	Apache 2.0
//	@ID				posts/unlike_comment
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			commentID	body	uint	true	"ID of the comment"
//
//	@Produce		json
//	@Success		204	{object}	json.JSONResponse
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/posts/comments/unlike [delete]
func (h *PostsHandler) HandleUnlikeComment(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	input := new(UnlikeCommentInput)

	err = easyjson.UnmarshalFromReader(r.Body, input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	_, err = h.PostsClient.UnlikeComment(r.Context(), &postspb.UnlikeCommentRequest{
		CommentId: uint64(input.CommentID),
		UserId:    uint64(userID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
