package rest

import (
	defJSON "encoding/json"
	"net/http"
	"socio/domain"
	"socio/errors"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
	"socio/usecase/subscriptions"
)

type SubscriptionsHandler struct {
	Service *subscriptions.Service
}

type SubscriptionInput struct {
	SubscribedToID uint `json:"subscribedTo"`
}

func NewSubscriptionsHandler(subStorage subscriptions.SubscriptionsStorage, userStorage subscriptions.UserStorage) (handler *SubscriptionsHandler) {
	handler = &SubscriptionsHandler{
		Service: subscriptions.NewService(subStorage, userStorage),
	}
	return
}

// HandleSubscription godoc
//
//	@Summary		handle user's subscription flow
//	@Description	subscribe to user
//	@Tags			subscriptions
//	@license.name	Apache 2.0
//	@ID				subscriptions/subscribe
//	@Accept			json
//
//	@Param			subscribedTo	body	int		true	"Subscribed to ID"
//	@Param			Cookie			header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//
//	@Produce		json
//	@Success		201	{object}	json.JSONResponse{body=domain.Subscription}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/subscriptions/ [post]
func (api *SubscriptionsHandler) HandleSubscription(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	input := new(SubscriptionInput)
	err := defJSON.NewDecoder(r.Body).Decode(input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	subscription, err := api.Service.Subscribe(r.Context(), &domain.Subscription{SubscriberID: userID, SubscribedToID: input.SubscribedToID})
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, map[string]*domain.Subscription{"subscription": subscription})
}

// HandleUnsubscription godoc
//
//		@Summary		handle user's unsubscription flow
//		@Description	unsubscribe from user
//		@Tags			subscriptions
//		@license.name	Apache 2.0
//		@ID				subscriptions/unsubscribe
//		@Accept			json
//
//		@Param			subscribedTo	body	int		true	"User to unsubscribe from"
//		@Param			Cookie			header	string	true	"session_id=some_session"
//		@Param			X-CSRF-Token	header	string	true	"CSRF token"
//
//		@Produce		json
//		@Success		204
//		@Failure		400	{object}	errors.HTTPError
//		@Failure		401	{object}	errors.HTTPError
//		@Failure		403	{object}	errors.HTTPError
//		@Failure		404	{object}	errors.HTTPError
//	 @Failure		500	{object}	errors.HTTPError
//
// @Router			/subscriptions/ [delete]
func (api *SubscriptionsHandler) HandleUnsubscription(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	input := new(SubscriptionInput)
	err := defJSON.NewDecoder(r.Body).Decode(input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	err = api.Service.Unsubscribe(r.Context(), &domain.Subscription{SubscriberID: userID, SubscribedToID: input.SubscribedToID})
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleGetSubscriptions godoc
//
//	@Summary		get user's subscriptions
//	@Description	get user's subscriptions
//	@Tags			subscriptions
//	@license.name	Apache 2.0
//	@ID				subscriptions/subscriptions
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=subscriptions.GetSubscriptionsResponse}
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/subscriptions/subscriptions/ [get]
func (api *SubscriptionsHandler) HandleGetSubscriptions(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	subscriptions, err := api.Service.GetSubscriptions(r.Context(), userID)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, map[string][]*domain.User{"subscriptions": subscriptions})
}

// HandleGetSubscribers godoc
//
//	@Summary		get user's subscribers
//	@Description	get user's subscribers
//	@Tags			subscriptions
//	@license.name	Apache 2.0
//	@ID				subscriptions/subscribers
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=subscriptions.GetSubscribersResponse}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/subscriptions/subscribers/ [get]
func (api *SubscriptionsHandler) HandleGetSubscribers(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	subscribers, err := api.Service.GetSubscribers(r.Context(), userID)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, map[string][]*domain.User{"subscribers": subscribers})
}

// HandleGetFriends godoc
//
//	@Summary		get user's friends
//	@Description	get user's friends
//	@Tags			subscriptions
//	@license.name	Apache 2.0
//	@ID				subscriptions/friends
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=subscriptions.GetFriendsResponse}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/subscriptions/friends/ [get]
func (api *SubscriptionsHandler) HandleGetFriends(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	friends, err := api.Service.GetFriends(r.Context(), userID)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, map[string][]*domain.User{"friends": friends})
}
