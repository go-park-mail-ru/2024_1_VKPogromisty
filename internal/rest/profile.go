package rest

import (
	"net/http"
	"socio/errors"
	"socio/internal/rest/middleware"
	"socio/pkg/json"
	"socio/usecase/profile"
	"strconv"

	"github.com/gorilla/mux"
)

type ProfileHandler struct {
	Service *profile.Service
}

func NewProfileHandler(userStorage profile.UserStorage) (h *ProfileHandler) {
	return &ProfileHandler{
		Service: profile.NewProfileService(userStorage),
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
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=profile.UserWithSubsInfo}
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/profile/{userID} [get]
func (h *ProfileHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	userIDData := mux.Vars(r)["userID"]
	if len(userIDData) == 0 {
		json.ServeJSONError(w, errors.ErrInvalidSlug)
		return
	}

	userID, err := strconv.ParseUint(userIDData, 10, 0)
	if err != nil {
		json.ServeJSONError(w, errors.ErrInvalidSlug)
		return
	}

	authorizedUserID := r.Context().Value(middleware.UserIDKey).(uint)

	userWithInfo, err := h.Service.GetUserByIDWithSubsInfo(uint(userID), authorizedUserID)
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	json.ServeJSONBody(w, userWithInfo)
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
//
//	@Success		204
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/profile/ [delete]
func (h *ProfileHandler) HandleDeleteProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)

	err := h.Service.DeleteUser(userID)
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
