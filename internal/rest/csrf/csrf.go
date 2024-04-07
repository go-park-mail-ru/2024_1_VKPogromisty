package rest

import (
	"net/http"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
	customtime "socio/pkg/time"
	"socio/usecase/csrf"
	"time"
)

type CSRFHandler struct {
	CSRFService  *csrf.CSRFService
	TimeProvider customtime.TimeProvider
}

type CSRFTokenResponse struct {
	CSRFToken string `json:"csrfToken"`
}

func NewCSRFHandler(timeProvider customtime.TimeProvider) *CSRFHandler {
	return &CSRFHandler{
		CSRFService:  csrf.NewCSRFService(customtime.RealTimeProvider{}),
		TimeProvider: timeProvider,
	}
}

// GetCSRFToken godoc
//
//	@Summary		Get CSRF token
//	@Description	Get CSRF token
//	@Tags			csrf
//	@ID				csrf/get_csrf_token
//	@Accept			json
//	@Produce		json
//
//	@Success		200	{object}	json.JSONResponse{body=CSRFTokenResponse}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/csrf/ [get]
func (api *CSRFHandler) GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	sessionID, err := requestcontext.GetSessionID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	tokenExpTime := api.TimeProvider.Now().Add(30 * time.Minute).Unix()

	token, err := api.CSRFService.Create(sessionID, userID, tokenExpTime)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, map[string]string{"csrfToken": token})
}
