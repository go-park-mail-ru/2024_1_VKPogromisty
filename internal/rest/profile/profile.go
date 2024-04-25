package rest

import (
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"socio/errors"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
	"socio/usecase/user"
	"strconv"
	"strings"

	uspb "socio/internal/grpc/user/proto"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	BatchSize = 1 << 23
)

type ProfileHandler struct {
	UserClient uspb.UserClient
}

func NewProfileHandler(userClient uspb.UserClient) (h *ProfileHandler) {
	return &ProfileHandler{
		UserClient: userClient,
	}
}

func (h *ProfileHandler) uploadAvatar(r *http.Request, avatarFH *multipart.FileHeader) (string, error) {
	fileName := uuid.NewString() + filepath.Ext(avatarFH.Filename)
	stream, err := h.UserClient.Upload(r.Context())
	if err != nil {
		return "", err
	}

	file, err := avatarFH.Open()
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

		err = stream.Send(&uspb.UploadRequest{
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
//	@Success		200	{object}	json.JSONResponse{body=profile.UserWithSubsInfo}
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
	input.FirstName = strings.Trim(r.PostFormValue("firstName"), " \n\r\t")
	input.LastName = strings.Trim(r.PostFormValue("lastName"), " \n\r\t")
	input.Email = strings.Trim(r.PostFormValue("email"), " \n\r\t")
	input.Password = r.PostFormValue("password")
	input.RepeatPassword = r.PostFormValue("repeatPassword")
	input.DateOfBirth = strings.Trim(r.PostFormValue("dateOfBirth"), " \n\r\t")
	_, avatarFH, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	if avatarFH != nil {
		avatarFileName, err := h.uploadAvatar(r, avatarFH)
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
