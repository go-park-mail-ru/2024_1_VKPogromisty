package chat

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/pkg/sanitizer"
	"sync"
)

// Service will: register and unregister new clients by userID, establish pub/sub connection to redis
type Service struct {
	Clients          *sync.Map
	PubSubRepository PubSubRepository
	MessagesRepo     PersonalMessagesRepository
	Sanitizer        *sanitizer.Sanitizer
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

func NewChatService(pubSubRepo PubSubRepository, messagesRepo PersonalMessagesRepository, sanitizer *sanitizer.Sanitizer) (chatService *Service) {
	return &Service{
		Clients:          &sync.Map{},
		PubSubRepository: pubSubRepo,
		MessagesRepo:     messagesRepo,
		Sanitizer:        sanitizer,
	}
}

func (s *Service) Register(ctx context.Context, userID uint) (c *Client, err error) {
	cData, ok := s.Clients.Load(userID)
	if ok {
		c = cData.(*Client)
		return
	}

	c, err = NewClient(userID, s.PubSubRepository, s.MessagesRepo, s.Sanitizer)
	if err != nil {
		return
	}

	go c.ReadPump(ctx)

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

func (s *Service) GetMessagesByDialog(ctx context.Context, userID, peerID, lastMessageID uint) (messages []*domain.PersonalMessage, err error) {
	if lastMessageID == 0 {
		lastMessageID, err = s.MessagesRepo.GetLastMessageID(ctx, userID, peerID)
		if err != nil {
			return
		}
		lastMessageID++
	}

	messages, err = s.MessagesRepo.GetMessagesByDialog(ctx, userID, peerID, lastMessageID)
	if err != nil {
		return
	}

	for _, message := range messages {
		s.Sanitizer.SanitizePersonalMessage(message)
	}

	return
}

func (s *Service) GetDialogsByUserID(ctx context.Context, userID uint) (dialogs []*domain.Dialog, err error) {
	dialogs, err = s.MessagesRepo.GetDialogsByUserID(ctx, userID)
	if err != nil {
		return
	}

	for _, dialog := range dialogs {
		s.Sanitizer.SanitizeDialog(dialog)
	}

	return
}
