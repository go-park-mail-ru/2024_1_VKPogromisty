package rest

import (
	defJSON "encoding/json"
	"fmt"
	"net/http"
	"socio/errors"
	uspb "socio/internal/grpc/user/proto"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
	"strconv"

	"github.com/gorilla/mux"
)

type AdminHandler struct {
	UserClient uspb.UserClient
}

type CreateAdminInput struct {
	UserID uint `json:"userId"`
}

type DeleteAdminInput struct {
	AdminID uint `json:"adminId"`
}

func NewAdminHandler(userClient uspb.UserClient) (h *AdminHandler) {
	return &AdminHandler{
		UserClient: userClient,
	}
}

func (h *AdminHandler) HandleGetAdminByUserID(w http.ResponseWriter, r *http.Request) {
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

	admin, err := h.UserClient.GetAdminByUserID(r.Context(), &uspb.GetAdminByUserIDRequest{
		UserId: uint64(userID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, uspb.ToAdmin(admin.Admin), http.StatusOK)
}

func (h *AdminHandler) HandleGetAdmins(w http.ResponseWriter, r *http.Request) {
	admins, err := h.UserClient.GetAdmins(r.Context(), &uspb.GetAdminsRequest{})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, uspb.ToAdmins(admins.Admins), http.StatusOK)
}

func (h *AdminHandler) HandleCreateAdmin(w http.ResponseWriter, r *http.Request) {
	input := new(CreateAdminInput)

	decoder := defJSON.NewDecoder(r.Body)
	err := decoder.Decode(input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	fmt.Println(input)

	admin, err := h.UserClient.CreateAdmin(r.Context(), &uspb.CreateAdminRequest{
		UserId: uint64(input.UserID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, uspb.ToAdmin(admin.Admin), http.StatusOK)
}

func (h *AdminHandler) HandleDeleteAdmin(w http.ResponseWriter, r *http.Request) {
	input := new(DeleteAdminInput)

	decoder := defJSON.NewDecoder(r.Body)
	err := decoder.Decode(input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	_, err = h.UserClient.DeleteAdmin(r.Context(), &uspb.DeleteAdminRequest{
		AdminId: uint64(input.AdminID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
