package chat

import (
	"encoding/json"
	"socio/errors"
	"sync"
)

// Service will: register and unregister new clients by userID, establish pub/sub connection to redis
type Service struct {
	Clients          *sync.Map
	PubSubRepository PubSubRepository
	MessagesRepo     PersonalMessagesRepository
}

type Action struct {
	Type     string `json:"type"`
	Receiver uint   `json:"receiver"`
	Payload  json.RawMessage
}

type SendMessagePayload struct {
	Content string `json:"content"`
}

type UpdateMessagePayload struct {
	MessageID uint   `json:"messageId"`
	Content   string `json:"content"`
}

type DeleteMessagePayload struct {
	MessageID uint `json:"messageId"`
}

func NewChatService(pubSubRepo PubSubRepository, messagesRepo PersonalMessagesRepository) (chatService *Service) {
	return &Service{
		Clients:          &sync.Map{},
		PubSubRepository: pubSubRepo,
		MessagesRepo:     messagesRepo,
	}
}

func (s *Service) Register(userID uint) (c *Client, err error) {
	cData, ok := s.Clients.Load(userID)
	if ok {
		c = cData.(*Client)
		return
	}

	c, err = NewClient(userID, s.PubSubRepository, s.MessagesRepo)
	if err != nil {
		return
	}

	go c.ReadPump()

	s.Clients.Store(userID, c)
	return
}

func (s *Service) Unregister(userID uint) (err error) {
	_, ok := s.Clients.LoadAndDelete(userID)
	if !ok {
		return errors.ErrNotFound
	}

	return
}
