package rest

import (
	defJSON "encoding/json"
	"net/http"
	"socio/domain"
	"socio/errors"
	postpb "socio/internal/grpc/post/proto"
	pgpb "socio/internal/grpc/public_group/proto"
	uspb "socio/internal/grpc/user/proto"
	"socio/internal/rest/uploaders"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
	"socio/usecase/posts"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type CreatePublicGroupAdminInput struct {
	UserID uint64 `json:"userId"`
}

type DeletePublicGroupAdminInput struct {
	UserID uint64 `json:"userId"`
}

type CheckIfUserIsAdminRes struct {
	IsAdmin bool `json:"isAdmin"`
}

type PublicGroupHandler struct {
	PublicGroupClient pgpb.PublicGroupClient
	PostClient        postpb.PostClient
	UserClient        uspb.UserClient
}

func NewPublicGroupHandler(publicGroupClient pgpb.PublicGroupClient, postClient postpb.PostClient, userClient uspb.UserClient) (h *PublicGroupHandler) {
	return &PublicGroupHandler{
		PublicGroupClient: publicGroupClient,
		PostClient:        postClient,
		UserClient:        userClient,
	}
}

// HandleGetByID godoc
//
//	@Summary		get public group by ID
//	@Description	get public group by ID
//	@Tags			groups
//	@license.name	Apache 2.0
//	@ID				groups/get
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			groupID	path	string	true	"Group ID"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=publicgroup.PublicGroupWithInfo}
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/groups/{groupID} [get]
func (h *PublicGroupHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	groupIDData := mux.Vars(r)["groupID"]
	if len(groupIDData) == 0 {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	var groupID uint64
	var err error

	groupID, err = strconv.ParseUint(groupIDData, 10, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	authorizedUserID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	res, err := h.PublicGroupClient.GetByID(r.Context(), &pgpb.GetByIDRequest{
		Id:     groupID,
		UserId: uint64(authorizedUserID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, pgpb.ToPublicGroupWithInfo(res.GetPublicGroup()), http.StatusOK)
}

// HandleSearchByName godoc
//
//	@Summary		search public groups by name
//	@Description	search public groups by name
//	@Tags			groups
//	@license.name	Apache 2.0
//	@ID				groups/search
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			query	query	string	true	"Search query"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=[]publicgroup.PublicGroupWithInfo}
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/groups/search [get]
func (h *PublicGroupHandler) HandleSearchByName(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if len(query) == 0 {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	authorizedUserID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	res, err := h.PublicGroupClient.SearchByName(r.Context(), &pgpb.SearchByNameRequest{
		Query:  query,
		UserId: uint64(authorizedUserID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, pgpb.ToPublicGroupsWithInfo(res.GetPublicGroups()), http.StatusOK)
}

// HandleCreate godoc
//
//	@Summary		create public group
//	@Description	create public group
//	@Tags			groups
//	@license.name	Apache 2.0
//	@ID				groups/create
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			name	formData	string	true	"Name of the group"
//	@Param			description	formData	string	true	"Description of the group"
//	@Param			avatar	formData	file	false	"Avatar of the group"
//
//	@Produce		json
//	@Success		201	{object}	json.JSONResponse{body=domain.PublicGroup}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/groups/ [post]
func (h *PublicGroupHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(4 * 1024 * 1024)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidBody)
		return
	}

	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrUnauthorized)
		return
	}

	var input pgpb.CreateRequest
	input.Name = strings.TrimSpace(r.PostFormValue("name"))
	input.Description = strings.TrimSpace(r.PostFormValue("description"))
	_, avatarFH, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	if avatarFH != nil {
		avatarFileName, err := uploaders.UploadPublicGroupAvatar(r, h.PublicGroupClient, avatarFH)
		if err != nil {
			json.ServeJSONError(r.Context(), w, err)
			return
		}

		input.Avatar = avatarFileName
	}

	res, err := h.PublicGroupClient.Create(r.Context(), &input)
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	_, err = h.UserClient.CreatePublicGroupAdmin(r.Context(), &uspb.CreatePublicGroupAdminRequest{
		UserId:        uint64(userID),
		PublicGroupId: uint64(res.GetPublicGroup().Id),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, pgpb.ToPublicGroup(res.GetPublicGroup()), http.StatusCreated)
}

// HandleUpdate godoc
//
//	@Summary		update public group
//	@Description	update public group
//	@Tags			groups
//	@license.name	Apache 2.0
//	@ID				groups/update
//	@Accept			mpfd
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			groupID	path	string	true	"Group ID"
//	@Param			name	formData	string	false	"Name of the group"
//	@Param			description	formData	string	false	"Description of the group"
//	@Param			avatar	formData	file	false	"Avatar of the group"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=domain.PublicGroup}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/groups/{groupID} [put]
func (h *PublicGroupHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(4 * 1024 * 1024)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidBody)
		return
	}

	groupIDData := mux.Vars(r)["groupID"]
	if len(groupIDData) == 0 {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	groupID, err := strconv.ParseUint(groupIDData, 10, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	var input pgpb.UpdateRequest
	input.Id = groupID
	input.Name = strings.TrimSpace(r.PostFormValue("name"))
	input.Description = strings.TrimSpace(r.PostFormValue("description"))
	_, avatarFH, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	if avatarFH != nil {
		avatarFileName, err := uploaders.UploadPublicGroupAvatar(r, h.PublicGroupClient, avatarFH)
		if err != nil {
			json.ServeJSONError(r.Context(), w, err)
			return
		}

		input.Avatar = avatarFileName
	}

	res, err := h.PublicGroupClient.Update(r.Context(), &input)
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, pgpb.ToPublicGroup(res.GetPublicGroup()), http.StatusOK)
}

// HandleDelete godoc
//
//	@Summary		delete public group
//	@Description	delete public group
//	@Tags			groups
//	@license.name	Apache 2.0
//	@ID				groups/delete
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			groupID	path	string	true	"Group ID"
//
//	@Produce		json
//	@Success		204	{object}	json.JSONResponse
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/groups/{groupID} [delete]
func (h *PublicGroupHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	groupIDData := mux.Vars(r)["groupID"]
	if len(groupIDData) == 0 {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	var groupID uint64
	var err error

	groupID, err = strconv.ParseUint(groupIDData, 10, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	_, err = h.PublicGroupClient.Delete(r.Context(), &pgpb.DeleteRequest{
		Id: groupID,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, nil, http.StatusNoContent)
}

// HandleGetSubscriptionByPublicGroupIDAndSubscriberID godoc
//
//	@Summary		get subscription by public group ID and subscriber ID
//	@Description	get subscription by public group ID and subscriber ID
//	@Tags			groups
//	@license.name	Apache 2.0
//	@ID				groups/is-sub
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			groupID	path	string	true	"Group ID"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=domain.PublicGroupSubscription}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/groups/{groupID}/is-sub [get]
func (h *PublicGroupHandler) HandleGetSubscriptionByPublicGroupIDAndSubscriberID(w http.ResponseWriter, r *http.Request) {
	groupIDData := mux.Vars(r)["groupID"]
	if len(groupIDData) == 0 {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	groupID, err := strconv.ParseUint(groupIDData, 10, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	authorizedUserID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	res, err := h.PublicGroupClient.GetSubscriptionByPublicGroupIDAndSubscriberID(r.Context(), &pgpb.GetSubscriptionByPublicGroupIDAndSubscriberIDRequest{
		PublicGroupId: groupID,
		SubscriberId:  uint64(authorizedUserID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, pgpb.ToSubscription(res.GetSubscription()), http.StatusOK)
}

// HandleGetBySubscriberID godoc
//
//	@Summary		get public groups by subscriber ID
//	@Description	get public groups by subscriber ID
//	@Tags			groups
//	@license.name	Apache 2.0
//	@ID				groups/by-sub
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			userID	path	string	false	"User ID"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=[]publicgroup.PublicGroup}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/groups/by-sub/{userID} [get]
func (h *PublicGroupHandler) HandleGetBySubscriberID(w http.ResponseWriter, r *http.Request) {
	userIDData := mux.Vars(r)["userID"]
	var userID uint64
	var err error

	authorizedUserID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	if len(userIDData) != 0 {
		userID, err = strconv.ParseUint(userIDData, 10, 0)
		if err != nil {
			json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
			return
		}
	} else {
		userID = uint64(authorizedUserID)
	}

	res, err := h.PublicGroupClient.GetBySubscriberID(r.Context(), &pgpb.GetBySubscriberIDRequest{
		SubscriberId: userID,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, pgpb.ToPublicGroups(res.GetPublicGroups()), http.StatusOK)
}

// HandleSubscribe godoc
//
//	@Summary		subscribe to public group
//	@Description	subscribe to public group
//	@Tags			groups
//	@license.name	Apache 2.0
//	@ID				groups/sub
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			groupID	path	string	true	"Group ID"
//
//	@Produce		json
//	@Success		201	{object}	json.JSONResponse
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/groups/{groupID}/sub [post]
func (h *PublicGroupHandler) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	groupIDData := mux.Vars(r)["groupID"]
	if len(groupIDData) == 0 {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	groupID, err := strconv.ParseUint(groupIDData, 10, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	authorizedUserID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	res, err := h.PublicGroupClient.Subscribe(r.Context(), &pgpb.SubscribeRequest{
		PublicGroupId: groupID,
		SubscriberId:  uint64(authorizedUserID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, pgpb.ToSubscription(res.GetSubscription()), http.StatusCreated)
}

// HandleUnsubscribe godoc
//
//	@Summary		unsubscribe from public group
//	@Description	unsubscribe from public group
//	@Tags			groups
//	@license.name	Apache 2.0
//	@ID				groups/unsub
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			groupID	path	string	true	"Group ID"
//
//	@Produce		json
//	@Success		204	{object}	json.JSONResponse
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/groups/{groupID}/unsub [post]
func (h *PublicGroupHandler) HandleUnsubscribe(w http.ResponseWriter, r *http.Request) {
	groupIDData := mux.Vars(r)["groupID"]
	if len(groupIDData) == 0 {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	groupID, err := strconv.ParseUint(groupIDData, 10, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	authorizedUserID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	_, err = h.PublicGroupClient.Unsubscribe(r.Context(), &pgpb.UnsubscribeRequest{
		PublicGroupId: groupID,
		SubscriberId:  uint64(authorizedUserID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, nil, http.StatusNoContent)
}

// HandleCreateGroupPost godoc
//
//	@Summary		create post in public group
//	@Description	create post in public group
//	@Tags			groups
//	@license.name	Apache 2.0
//	@ID				groups/posts/create
//	@Accept			mpfd
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			groupID	path	string	true	"Group ID"
//	@Param			content	formData	string	true	"Content of the post"
//	@Param			attachments	formData	file	false	"Attachments of the post"
//
//	@Produce		json
//	@Success		201	{object}	json.JSONResponse{body=domain.PostWithAuthorAndGroup}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/groups/{groupID}/posts/ [post]
func (h *PublicGroupHandler) HandleCreateGroupPost(w http.ResponseWriter, r *http.Request) {
	groupIDData := mux.Vars(r)["groupID"]
	if len(groupIDData) == 0 {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	groupID, err := strconv.ParseUint(groupIDData, 10, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	err = r.ParseMultipartForm(1000 << 20)
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
			fileName, err := uploaders.UploadPostAttachment(r, h.PostClient, f)
			if err != nil {
				json.ServeJSONError(r.Context(), w, err)
				return
			}
			postInput.Attachments = append(postInput.Attachments, fileName)
		}
	}

	postData, err := h.PostClient.CreatePost(r.Context(), &postpb.CreatePostRequest{
		AuthorId:    uint64(postInput.AuthorID),
		Content:     postInput.Content,
		Attachments: postInput.Attachments,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	post := postpb.ToPost(postData.Post)

	_, err = h.PostClient.CreateGroupPost(r.Context(), &postpb.CreateGroupPostRequest{
		GroupId: groupID,
		PostId:  uint64(post.ID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	group, err := h.PublicGroupClient.GetByID(r.Context(), &pgpb.GetByIDRequest{
		Id: groupID,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	postWithGroup := &domain.PostWithAuthorAndGroup{
		Post:  post,
		Group: pgpb.ToPublicGroup(group.GetPublicGroup().PublicGroup),
	}

	json.ServeJSONBody(r.Context(), w, postWithGroup, http.StatusCreated)
}

// HandleGetGroupPosts godoc
//
//	@Summary		get posts of public group
//	@Description	get posts of public group
//	@Tags			groups
//	@license.name	Apache 2.0
//	@ID				groups/posts
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			groupID	path	string	true	"Group ID"
//	@Param			lastPostId	query	string	false	"Last post ID"
//	@Param			postsAmount	query	string	false	"Posts amount"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=[]domain.Post}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/groups/{groupID}/posts/ [get]
func (h *PublicGroupHandler) HandleGetGroupPosts(w http.ResponseWriter, r *http.Request) {
	groupIDData := mux.Vars(r)["groupID"]
	if len(groupIDData) == 0 {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	groupID, err := strconv.ParseUint(groupIDData, 10, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	lastPostIDData := r.URL.Query().Get("lastPostId")
	var lastPostID uint64
	if len(lastPostIDData) != 0 {
		lastPostID, err = strconv.ParseUint(lastPostIDData, 10, 0)
		if err != nil {
			json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
			return
		}
	}

	postsAmountData := r.URL.Query().Get("postsAmount")
	var postsAmount uint64
	if len(postsAmountData) != 0 {
		postsAmount, err = strconv.ParseUint(postsAmountData, 10, 0)
		if err != nil {
			json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
			return
		}
	}

	res, err := h.PostClient.GetPostsOfGroup(r.Context(), &postpb.GetPostsOfGroupRequest{
		GroupId:     groupID,
		LastPostId:  lastPostID,
		PostsAmount: postsAmount,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	posts := postpb.ToPosts(res.GetPosts())

	json.ServeJSONBody(r.Context(), w, posts, http.StatusOK)
}

// HandleCreatePublicGroupAdmin godoc
//
//	@Summary		create public group admin
//	@Description	create public group admin
//	@Tags			groups
//	@license.name	Apache 2.0
//	@ID				groups/admin/create
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			groupID	path	string	true	"Group ID"
//
//	@Produce		json
//	@Success		201	{object}	json.JSONResponse
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/groups/{groupID}/admins/ [post]
func (h *PublicGroupHandler) HandleCreatePublicGroupAdmin(w http.ResponseWriter, r *http.Request) {
	groupIDData := mux.Vars(r)["groupID"]
	if len(groupIDData) == 0 {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	groupID, err := strconv.ParseUint(groupIDData, 10, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	var input CreatePublicGroupAdminInput

	decoder := defJSON.NewDecoder(r.Body)
	err = decoder.Decode(&input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	_, err = h.UserClient.CreatePublicGroupAdmin(r.Context(), &uspb.CreatePublicGroupAdminRequest{
		UserId:        input.UserID,
		PublicGroupId: groupID,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, nil, http.StatusCreated)
}

// HandleDeletePublicGroupAdmin godoc
//
//	@Summary		delete public group admin
//	@Description	delete public group admin
//	@Tags			groups
//	@license.name	Apache 2.0
//	@ID				groups/admin/delete
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			groupID	path	string	true	"Group ID"
//
//	@Produce		json
//	@Success		204	{object}	json.JSONResponse
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/groups/{groupID}/admins/ [delete]
func (h *PublicGroupHandler) HandleDeletePublicGroupAdmin(w http.ResponseWriter, r *http.Request) {
	groupIDData := mux.Vars(r)["groupID"]
	if len(groupIDData) == 0 {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	groupID, err := strconv.ParseUint(groupIDData, 10, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	var input DeletePublicGroupAdminInput

	decoder := defJSON.NewDecoder(r.Body)
	err = decoder.Decode(&input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	_, err = h.UserClient.DeletePublicGroupAdmin(r.Context(), &uspb.DeletePublicGroupAdminRequest{
		UserId:        input.UserID,
		PublicGroupId: groupID,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, nil, http.StatusNoContent)
}

// HandleGetAdminsByPublicGroupID godoc
//
//	@Summary		get admins by public group ID
//	@Description	get admins by public group ID
//	@Tags			groups
//	@license.name	Apache 2.0
//	@ID				groups/admins
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			groupID	path	string	true	"Group ID"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=[]domain.User}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/groups/{groupID}/admins/ [get]
func (h *PublicGroupHandler) HandleGetAdminsByPublicGroupID(w http.ResponseWriter, r *http.Request) {
	groupIDData := mux.Vars(r)["groupID"]
	if len(groupIDData) == 0 {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	groupID, err := strconv.ParseUint(groupIDData, 10, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	res, err := h.UserClient.GetAdminsByPublicGroupID(r.Context(), &uspb.GetAdminsByPublicGroupIDRequest{
		PublicGroupId: groupID,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, uspb.ToUsers(res.GetAdmins()), http.StatusOK)
}

// HandleCheckIfUserIsAdmin godoc
//
//	@Summary		check if user is admin
//	@Description	check if user is admin
//	@Tags			groups
//	@license.name	Apache 2.0
//	@ID				groups/admins/check
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			groupID	path	string	true	"Group ID"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=CheckIfUserIsAdminRes}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/groups/{groupID}/admins/check [get]
func (h *PublicGroupHandler) HandleCheckIfUserIsAdmin(w http.ResponseWriter, r *http.Request) {
	groupIDData := mux.Vars(r)["groupID"]
	if len(groupIDData) == 0 {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	groupID, err := strconv.ParseUint(groupIDData, 10, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
		return
	}

	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	res, err := h.UserClient.CheckIfUserIsAdmin(r.Context(), &uspb.CheckIfUserIsAdminRequest{
		UserId:        uint64(userID),
		PublicGroupId: groupID,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, &CheckIfUserIsAdminRes{
		IsAdmin: res.GetIsAdmin(),
	}, http.StatusOK)
}
