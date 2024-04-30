package rest

import (
	"net/http"
	"socio/errors"
	pgpb "socio/internal/grpc/public_group/proto"
	"socio/internal/rest/uploaders"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type PublicGroupHandler struct {
	PublicGroupClient pgpb.PublicGroupClient
}

func NewPublicGroupHandler(publicGroupClient pgpb.PublicGroupClient) (h *PublicGroupHandler) {
	return &PublicGroupHandler{
		PublicGroupClient: publicGroupClient,
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

	json.ServeJSONBody(r.Context(), w, pgpb.ToSubscription(res.GetSubscription()), http.StatusNoContent)
}

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
