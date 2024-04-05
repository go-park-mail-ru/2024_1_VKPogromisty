package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"socio/domain"
	"socio/errors"
	"socio/pkg/sanitizer"
)

const (
	sendChanSize                   = 256
	SendMessageAction   ChatAction = "SEND_MESSAGE"
	UpdateMessageAction ChatAction = "UPDATE_MESSAGE"
	DeleteMessageAction ChatAction = "DELETE_MESSAGE"
	SetOnlineAction     ChatAction = "SET_ONLINE"
	SetOfflineAction    ChatAction = "SET_OFFLINE"
)

type ChatAction string

type Action struct {
	Type     ChatAction      `json:"type"`
	Receiver uint            `json:"receiver"`
	Payload  json.RawMessage `json:"payload"`
}

type PersonalMessagesRepository interface {
	GetLastMessageID(ctx context.Context, senderID, receiverID uint) (lastMessageID uint, err error)
	GetMessagesByDialog(ctx context.Context, senderID, receiverID, lastMessageID uint) (messages []*domain.PersonalMessage, err error)
	GetDialogsByUserID(ctx context.Context, userID uint) (dialogs []*domain.Dialog, err error)
	StoreMessage(ctx context.Context, message *domain.PersonalMessage) (newMessage *domain.PersonalMessage, err error)
	UpdateMessage(ctx context.Context, message *domain.PersonalMessage) (updatedMessage *domain.PersonalMessage, err error)
	DeleteMessage(ctx context.Context, messageID uint) (err error)
}

type PubSubRepository interface {
	ReadActions(ctx context.Context, userID uint, ch chan *Action) (err error)
	WriteAction(ctx context.Context, action *Action) (err error)
}

// Client will: read Actions from redis and write Actions into Send, subscribe to corresponding redis channel
type Client struct {
	UserID               uint
	Send                 chan *Action
	PersonalMessagesRepo PersonalMessagesRepository
	PubSubRepository     PubSubRepository
	Sanitizer            *sanitizer.Sanitizer
}

func NewClient(userID uint, pubSubRepo PubSubRepository, messagesRepo PersonalMessagesRepository, sanitizer *sanitizer.Sanitizer) (client *Client, err error) {
	if err != nil {
		return
	}

	client = &Client{
		UserID:               userID,
		Send:                 make(chan *Action, sendChanSize),
		PubSubRepository:     pubSubRepo,
		PersonalMessagesRepo: messagesRepo,
		Sanitizer:            sanitizer,
	}

	return
}

func (c *Client) ReadPump(ctx context.Context) {
	actionsCh := make(chan *Action)
	defer close(actionsCh)

	go c.PubSubRepository.ReadActions(ctx, c.UserID, c.Send)
}

func (c *Client) HandleAction(ctx context.Context, action *Action) {
	switch action.Type {
	case SendMessageAction:
		payload := new(SendMessagePayload)
		err := json.NewDecoder(bytes.NewReader(action.Payload)).Decode(payload)
		if err != nil {
			return
		}
		c.handleSendMessageAction(ctx, action, payload)

	case UpdateMessageAction:
		payload := new(UpdateMessagePayload)
		json.NewDecoder(bytes.NewReader(action.Payload)).Decode(payload)
		c.handleUpdateMessageAction(ctx, action, payload)

	case DeleteMessageAction:
		payload := new(DeleteMessagePayload)
		json.NewDecoder(bytes.NewReader(action.Payload)).Decode(payload)
		c.handleDeleteMessageAction(ctx, action, payload.MessageID)
	}
}

func (c *Client) handleSendMessageAction(ctx context.Context, action *Action, message *SendMessagePayload) {
	msg := &domain.PersonalMessage{
		Content:    message.Content,
		SenderID:   c.UserID,
		ReceiverID: action.Receiver,
	}

	newMessage, err := c.PersonalMessagesRepo.StoreMessage(ctx, msg)
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		c.PubSubRepository.WriteAction(ctx, action)
		return
	}

	c.Sanitizer.SanitizePersonalMessage(newMessage)

	action.Payload, err = json.Marshal(newMessage)
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		c.PubSubRepository.WriteAction(ctx, action)
		return
	}

	c.PubSubRepository.WriteAction(ctx, action)
}

func (c *Client) handleUpdateMessageAction(ctx context.Context, action *Action, message *UpdateMessagePayload) {
	msg := &domain.PersonalMessage{
		ID:      message.MessageID,
		Content: message.Content,
	}

	newMessage, err := c.PersonalMessagesRepo.UpdateMessage(ctx, msg)
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		c.PubSubRepository.WriteAction(ctx, action)
		return
	}

	c.Sanitizer.SanitizePersonalMessage(newMessage)

	action.Payload, err = json.Marshal(newMessage)
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		c.PubSubRepository.WriteAction(ctx, action)
		return
	}

	c.PubSubRepository.WriteAction(ctx, action)
}

func (c *Client) handleDeleteMessageAction(ctx context.Context, action *Action, messageID uint) {
	err := c.PersonalMessagesRepo.DeleteMessage(ctx, messageID)
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		c.PubSubRepository.WriteAction(ctx, action)
		return
	}

	c.PubSubRepository.WriteAction(ctx, action)
}
