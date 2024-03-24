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
