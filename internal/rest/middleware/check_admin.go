package middleware

import (
	"context"
	"fmt"
	"net/http"
	"socio/errors"
	uspb "socio/internal/grpc/user/proto"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
)

func CreateCheckAdminMiddleware(userClient uspb.UserClient) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := requestcontext.GetUserID(r.Context())
			if err != nil {
				fmt.Println(userID)
				json.ServeJSONError(r.Context(), w, errors.ErrUnauthorized)
				return
			}

			res, err := userClient.GetAdminByUserID(r.Context(), &uspb.GetAdminByUserIDRequest{
				UserId: uint64(userID),
			})
			if err != nil {
				json.ServeGRPCStatus(r.Context(), w, errors.ErrForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), requestcontext.AdminIDKey, uint(res.Admin.Id))

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
