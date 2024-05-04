package middleware

import (
	"net/http"
	"strconv"

	"socio/errors"
	uspb "socio/internal/grpc/user/proto"
	"socio/pkg/json"
	"socio/pkg/requestcontext"

	"github.com/gorilla/mux"
)

func CreateCheckPublicGroupAdminMiddleware(userClient uspb.UserClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userID, err := requestcontext.GetUserID(ctx)
			if err != nil {
				json.ServeJSONError(r.Context(), w, errors.ErrUnauthorized)
				return
			}

			groupID, err := strconv.ParseInt(mux.Vars(r)["groupID"], 10, 64)
			if err != nil {
				json.ServeJSONError(r.Context(), w, errors.ErrInvalidSlug)
				return
			}

			isAdmin, err := userClient.CheckIfUserIsAdmin(ctx, &uspb.CheckIfUserIsAdminRequest{
				UserId:        uint64(userID),
				PublicGroupId: uint64(groupID),
			})
			if err != nil {
				json.ServeJSONError(r.Context(), w, errors.ErrForbidden)
				return
			}

			if !isAdmin.IsAdmin {
				json.ServeJSONError(r.Context(), w, errors.ErrForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
