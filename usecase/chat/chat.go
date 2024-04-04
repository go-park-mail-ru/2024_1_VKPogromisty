package chat

import (
	"fmt"
	"socio/domain"
	"socio/errors"
	"sync"
)

// Service will: register and unregister new clients by userID, establish pub/sub connection to redis
type Service struct {
	Clients          *sync.Map
	PubSubRepository PubSubRepository
	MessagesRepo     PersonalMessagesRepository
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

type Dialog struct {
	User1       *domain.User            `json:"user1"`
	User2       *domain.User            `json:"user2"`
	LastMessage *domain.PersonalMessage `json:"lastMessage"`
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

func (s *Service) GetMessagesByDialog(userID, peerID, lastMessageID uint) (messages []*domain.PersonalMessage, err error) {
	if lastMessageID == 0 {
		lastMessageID, err = s.MessagesRepo.GetLastMessageID(userID, peerID)
		if err != nil {
			fmt.Println("here")
			return
		}
		lastMessageID++
	}

	messages, err = s.MessagesRepo.GetMessagesByDialog(userID, peerID, lastMessageID)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetDialogsByUserID(userID uint) (dialogs []*Dialog, err error) {
	dialogs, err = s.MessagesRepo.GetDialogsByUserID(userID)
	if err != nil {
		return
	}

	return
}
