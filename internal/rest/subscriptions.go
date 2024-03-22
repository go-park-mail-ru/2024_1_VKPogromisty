package rest

import (
	defJSON "encoding/json"
	"net/http"
	"socio/domain"
	"socio/errors"
	"socio/internal/rest/middleware"
	"socio/pkg/json"
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

func (api *SubscriptionsHandler) HandleSubscription(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	input := new(SubscriptionInput)
	err := defJSON.NewDecoder(r.Body).Decode(input)
	if err != nil {
		json.ServeJSONError(w, errors.ErrJSONUnmarshalling)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(uint)

	subscription, err := api.Service.Subscribe(&domain.Subscription{SubscriberID: userID, SubscribedToID: input.SubscribedToID})
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	json.ServeJSONBody(w, map[string]*domain.Subscription{"subscription": subscription})
}

func (api *SubscriptionsHandler) HandleUnsubscription(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	input := new(SubscriptionInput)
	err := defJSON.NewDecoder(r.Body).Decode(input)
	if err != nil {
		json.ServeJSONError(w, errors.ErrJSONUnmarshalling)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(uint)

	err = api.Service.Unsubscribe(&domain.Subscription{SubscriberID: userID, SubscribedToID: input.SubscribedToID})
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *SubscriptionsHandler) HandleGetSubscriptions(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID := r.Context().Value(middleware.UserIDKey).(uint)

	subscriptions, err := api.Service.GetSubscriptions(userID)
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	json.ServeJSONBody(w, map[string][]*domain.User{"subscriptions": subscriptions})
}

func (api *SubscriptionsHandler) HandleGetSubscribers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID := r.Context().Value(middleware.UserIDKey).(uint)

	subscribers, err := api.Service.GetSubscribers(userID)
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	json.ServeJSONBody(w, map[string][]*domain.User{"subscribers": subscribers})
}

func (api *SubscriptionsHandler) HandleGetFriends(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID := r.Context().Value(middleware.UserIDKey).(uint)

	friends, err := api.Service.GetFriends(userID)
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	json.ServeJSONBody(w, map[string][]*domain.User{"friends": friends})
}
