package chat

import (
	"bytes"
	"encoding/json"
	"socio/domain"
	"socio/errors"
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
	GetLastMessageID(senderID, receiverID uint) (lastMessageID uint, err error)
	GetMessagesByDialog(senderID, receiverID, lastMessageID uint) (messages []*domain.PersonalMessage, err error)
	GetDialogsByUserID(userID uint) (dialogs []*Dialog, err error)
	StoreMessage(message *domain.PersonalMessage) (newMessage *domain.PersonalMessage, err error)
	UpdateMessage(message *domain.PersonalMessage) (updatedMessage *domain.PersonalMessage, err error)
	DeleteMessage(messageID uint) (err error)
}

type PubSubRepository interface {
	ReadActions(userID uint, ch chan *Action) (err error)
	WriteAction(action *Action) (err error)
}

// Client will: read Actions from redis and write Actions into Send, subscribe to corresponding redis channel
type Client struct {
	UserID               uint
	Send                 chan *Action
	PersonalMessagesRepo PersonalMessagesRepository
	PubSubRepository     PubSubRepository
}

func NewClient(userID uint, pubSubRepo PubSubRepository, messagesRepo PersonalMessagesRepository) (client *Client, err error) {
	if err != nil {
		return
	}

	client = &Client{
		UserID:               userID,
		Send:                 make(chan *Action, sendChanSize),
		PubSubRepository:     pubSubRepo,
		PersonalMessagesRepo: messagesRepo,
	}

	return
}

func (c *Client) ReadPump() {
	actionsCh := make(chan *Action)
	defer close(actionsCh)

	go c.PubSubRepository.ReadActions(c.UserID, c.Send)
}

func (c *Client) HandleAction(action *Action) {
	switch action.Type {
	case SendMessageAction:
		payload := new(SendMessagePayload)
		err := json.NewDecoder(bytes.NewReader(action.Payload)).Decode(payload)
		if err != nil {
			return
		}
		c.handleSendMessageAction(action, payload)

	case UpdateMessageAction:
		payload := new(UpdateMessagePayload)
		json.NewDecoder(bytes.NewReader(action.Payload)).Decode(payload)
		c.handleUpdateMessageAction(action, payload)

	case DeleteMessageAction:
		payload := new(DeleteMessagePayload)
		json.NewDecoder(bytes.NewReader(action.Payload)).Decode(payload)
		c.handleDeleteMessageAction(action, payload.MessageID)
	}
}

func (c *Client) handleSendMessageAction(action *Action, message *SendMessagePayload) (err error) {
	msg := &domain.PersonalMessage{
		Content:    message.Content,
		SenderID:   c.UserID,
		ReceiverID: action.Receiver,
	}

	newMessage, err := c.PersonalMessagesRepo.StoreMessage(msg)
	if err != nil {
		action.Payload = errors.MarshalError(err)
		c.PubSubRepository.WriteAction(action)
		return
	}

	action.Payload, err = json.Marshal(newMessage)
	if err != nil {
		action.Payload = errors.MarshalError(err)
		c.PubSubRepository.WriteAction(action)
		return
	}

	c.PubSubRepository.WriteAction(action)

	return
}

func (c *Client) handleUpdateMessageAction(action *Action, message *UpdateMessagePayload) {
	msg := &domain.PersonalMessage{
		ID:      message.MessageID,
		Content: message.Content,
	}

	newMessage, err := c.PersonalMessagesRepo.UpdateMessage(msg)
	if err != nil {
		action.Payload = errors.MarshalError(err)
		c.PubSubRepository.WriteAction(action)
		return
	}

	action.Payload, err = json.Marshal(newMessage)
	if err != nil {
		action.Payload = errors.MarshalError(err)
		c.PubSubRepository.WriteAction(action)
		return
	}

	c.PubSubRepository.WriteAction(action)
}

func (c *Client) handleDeleteMessageAction(action *Action, messageID uint) {
	err := c.PersonalMessagesRepo.DeleteMessage(messageID)
	if err != nil {
		action.Payload = errors.MarshalError(err)
		c.PubSubRepository.WriteAction(action)
		return
	}

	c.PubSubRepository.WriteAction(action)
}
