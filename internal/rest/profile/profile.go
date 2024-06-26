package rest

import (
	"net/http"
	"socio/errors"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
	"socio/usecase/user"
	"strconv"
	"strings"

	uspb "socio/internal/grpc/user/proto"
	"socio/internal/rest/uploaders"

	"github.com/gorilla/mux"
)

type ProfileHandler struct {
	UserClient uspb.UserClient
}

func NewProfileHandler(userClient uspb.UserClient) (h *ProfileHandler) {
	return &ProfileHandler{
		UserClient: userClient,
	}
}

// HandleGetProfile godoc
//
//	@Summary		get user profile with subscriptions info
//	@Description	get user profile with subscriptions info
//	@Tags			profile
//	@license.name	Apache 2.0
//	@ID				profile/get
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			userID	path	string	false	"User ID, if empty - get authorized user profile"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=user.UserWithSubsInfo}
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/profile/{userID} [get]
func (h *ProfileHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
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

	userWithInfo, err := h.UserClient.GetByIDWithSubsInfo(r.Context(), &uspb.GetByIDWithSubsInfoRequest{
		UserId:           userID,
		AuthorizedUserId: uint64(authorizedUserID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, uspb.ToUserWithInfo(userWithInfo), http.StatusOK)
}

// HandleUpdateProfile godoc
//
//	@Summary		update user profile
//	@Description	update user profile
//	@Tags			profile
//	@license.name	Apache 2.0
//	@ID				profile/update
//	@Accept			mpfd
//
//	@Param			Cookie			header		string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//
//	@Param			firstName		formData	string	false	"First name"
//	@Param			lastName		formData	string	false	"Last name"
//	@Param			email			formData	string	false	"Email"
//	@Param			password		formData	string	false	"Password"
//	@Param			repeatPassword	formData	string	false	"Repeat password"
//	@Param			dateOfBirth		formData	string	false	"Date of birth"
//	@Param			avatar			formData	file	false	"Avatar"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=domain.User}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/profile/ [put]
func (h *ProfileHandler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(4 * 1024 * 1024)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidBody)
		return
	}

	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	var input user.UpdateUserInput
	input.ID = userID
	input.FirstName = strings.TrimSpace(r.PostFormValue("firstName"))
	input.LastName = strings.TrimSpace(r.PostFormValue("lastName"))
	input.Email = strings.TrimSpace(r.PostFormValue("email"))
	input.Password = r.PostFormValue("password")
	input.RepeatPassword = r.PostFormValue("repeatPassword")
	input.DateOfBirth = strings.TrimSpace(r.PostFormValue("dateOfBirth"))
	_, avatarFH, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	if avatarFH != nil {
		avatarFileName, err := uploaders.UploadAvatar(r, h.UserClient, avatarFH)
		if err != nil {
			json.ServeJSONError(r.Context(), w, err)
			return
		}

		input.Avatar = avatarFileName
	}

	var grpcInput = uspb.ToUpdateRequest(&input)

	updatedUser, err := h.UserClient.Update(r.Context(), grpcInput)
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, uspb.ToUser(updatedUser.User), http.StatusOK)
}

// HandleDeleteProfile godoc
//
//	@Summary		delete user profile
//	@Description	delete user profile
//	@Tags			profile
//	@license.name	Apache 2.0
//	@ID				profile/delete
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//
//	@Success		204
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/profile/ [delete]
func (h *ProfileHandler) HandleDeleteProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	_, err = h.UserClient.Delete(r.Context(), &uspb.DeleteRequest{
		UserId: uint64(userID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleSearchByName godoc
//
//	@Summary		search users by name
//	@Description	search users by name
//	@Tags			profile
//	@license.name	Apache 2.0
//	@ID				profile/search
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			query	query	string	true	"Search query"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=[]domain.User}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/profile/search [get]
func (h *ProfileHandler) HandleSearchByName(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	users, err := h.UserClient.SearchByName(r.Context(), &uspb.SearchByNameRequest{
		Query: query,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, uspb.ToUsers(users.Users), http.StatusOK)
}
