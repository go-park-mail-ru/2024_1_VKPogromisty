package rest

import (
	"net/http"
	"socio/errors"
	"socio/internal/rest/middleware"
	"socio/pkg/json"
	"socio/usecase/profile"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type ProfileHandler struct {
	Service *profile.Service
}

func NewProfileHandler(userStorage profile.UserStorage, sessionStorage profile.SessionStorage) (h *ProfileHandler) {
	return &ProfileHandler{
		Service: profile.NewProfileService(userStorage, sessionStorage),
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
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/profile/ [put]
func (h *ProfileHandler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(4 * 1024 * 1024)
	if err != nil {
		json.ServeJSONError(w, errors.ErrInvalidBody)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(uint)

	var input profile.UpdateUserInput
	input.ID = userID
	input.FirstName = strings.Trim(r.PostFormValue("firstName"), " \n\r\t")
	input.LastName = strings.Trim(r.PostFormValue("lastName"), " \n\r\t")
	input.Email = strings.Trim(r.PostFormValue("email"), " \n\r\t")
	input.Password = r.PostFormValue("password")
	input.RepeatPassword = r.PostFormValue("repeatPassword")
	input.DateOfBirth = strings.Trim(r.PostFormValue("dateOfBirth"), " \n\r\t")
	_, input.Avatar, err = r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		json.ServeJSONError(w, err)
		return
	}

	updatedUser, err := h.Service.UpdateUser(input)
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	json.ServeJSONBody(w, updatedUser)
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
	sessionID := r.Context().Value(middleware.SessionIDKey).(string)

	err := h.Service.DeleteUser(userID, sessionID)
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
