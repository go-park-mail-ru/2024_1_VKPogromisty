package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"socio/domain"
	"socio/errors"
	"sync"
	"time"
)

const (
	sendChanSize                        = 256
	tickerInterval                      = 5 * time.Minute
	SendMessageAction        ChatAction = "SEND_MESSAGE"
	UpdateMessageAction      ChatAction = "UPDATE_MESSAGE"
	DeleteMessageAction      ChatAction = "DELETE_MESSAGE"
	SendStickerMessageAction ChatAction = "SEND_STICKER_MESSAGE"
)

type PersonalMessagesRepository interface {
	GetLastMessageID(ctx context.Context, senderID, receiverID uint) (lastMessageID uint, err error)
	GetMessagesByDialog(ctx context.Context, senderID, receiverID, lastMessageID, messagesAmount uint) (messages []*domain.PersonalMessage, err error)
	GetDialogsByUserID(ctx context.Context, userID uint) (dialogs []*domain.Dialog, err error)
	StoreMessage(ctx context.Context, message *domain.PersonalMessage) (newMessage *domain.PersonalMessage, err error)
	UpdateMessage(ctx context.Context, msg *domain.PersonalMessage, attachmentsToDelete []string) (updatedMsg *domain.PersonalMessage, err error)
	DeleteMessage(ctx context.Context, messageID uint) (err error)
	GetStickerByID(ctx context.Context, stickerID uint) (sticker *domain.Sticker, err error)
	GetStickersByAuthorID(ctx context.Context, authorID uint) (stickers []*domain.Sticker, err error)
	GetAllStickers(ctx context.Context) (stickers []*domain.Sticker, err error)
	StoreSticker(ctx context.Context, sticker *domain.Sticker) (newSticker *domain.Sticker, err error)
	DeleteSticker(ctx context.Context, stickerID uint) (err error)
	StoreStickerMessage(ctx context.Context, senderID, receiverID, stickerID uint) (newStickerMessage *domain.PersonalMessage, err error)
}

type PubSubRepository interface {
	ReadActions(ctx context.Context, userID uint, ch chan *Action) (err error)
	WriteAction(ctx context.Context, action *Action) (err error)
}

type ChatAction string

type Action struct {
	Type      ChatAction      `json:"type"`
	Receiver  uint            `json:"receiver"`
	CSRFToken string          `json:"csrfToken"`
	Payload   json.RawMessage `json:"payload"`
}

type SendMessagePayload struct {
	Content     string   `json:"content"`
	Attachments []string `json:"attachments"`
}

type UpdateMessagePayload struct {
	MessageID           uint     `json:"messageId"`
	Content             string   `json:"content"`
	AttachmentsToDelete []string `json:"attachmentsToDelete"`
}

type DeleteMessagePayload struct {
	MessageID uint `json:"messageId"`
}

type SendStickerMessagePayload struct {
	StickerID uint `json:"stickerId"`
}

type Client struct {
	UserID                    uint
	Send                      chan *Action
	ChatService               *Service
	UnsentAttachmentReceivers *sync.Map
}

func NewClient(userID uint, chatService *Service) (client *Client, err error) {
	if err != nil {
		return
	}

	client = &Client{
		UserID:                    userID,
		Send:                      make(chan *Action, sendChanSize),
		ChatService:               chatService,
		UnsentAttachmentReceivers: &sync.Map{},
	}

	return
}

func (c *Client) ReadPump(ctx context.Context) {
	go func() {
		defer c.ClearUnsentAttachments(ctx)

		ticker := time.NewTicker(tickerInterval)
		defer ticker.Stop()

		for range ticker.C {
			c.ClearUnsentAttachments(ctx)
		}
	}()

	go c.ChatService.PubSubRepository.ReadActions(ctx, c.UserID, c.Send)
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
		err := json.NewDecoder(bytes.NewReader(action.Payload)).Decode(payload)
		if err != nil {
			return
		}
		c.handleUpdateMessageAction(ctx, action, payload)

	case DeleteMessageAction:
		payload := new(DeleteMessagePayload)
		err := json.NewDecoder(bytes.NewReader(action.Payload)).Decode(payload)
		if err != nil {
			return
		}
		c.handleDeleteMessageAction(ctx, action, payload.MessageID)

	case SendStickerMessageAction:
		payload := new(SendStickerMessagePayload)
		err := json.NewDecoder(bytes.NewReader(action.Payload)).Decode(payload)
		if err != nil {
			return
		}
		c.handleSendStickerMessageAction(ctx, action, payload)
	}
}

func (c *Client) handleSendMessageAction(ctx context.Context, action *Action, message *SendMessagePayload) {
	attachments, err := c.ChatService.UnsentMessageAttachmentsStorage.GetAll(ctx, &domain.UnsentMessageAttachment{
		SenderID:   c.UserID,
		ReceiverID: action.Receiver,
	})
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		err = c.ChatService.PubSubRepository.WriteAction(ctx, action)
		if err != nil {
			return
		}

		return
	}

	msg := &domain.PersonalMessage{
		Content:     message.Content,
		SenderID:    c.UserID,
		ReceiverID:  action.Receiver,
		Attachments: attachments,
	}

	c.ChatService.Sanitizer.SanitizePersonalMessage(msg)

	if len(msg.Content) == 0 && len(attachments) == 0 {
		action.Payload, err = errors.MarshalError(errors.ErrInvalidData)
		if err != nil {
			return
		}

		err = c.ChatService.PubSubRepository.WriteAction(ctx, action)
		if err != nil {
			return
		}

		return
	}

	newMessage, err := c.ChatService.MessagesRepo.StoreMessage(ctx, msg)
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		err = c.ChatService.PubSubRepository.WriteAction(ctx, action)
		if err != nil {
			return
		}

		return
	}

	err = c.ChatService.UnsentMessageAttachmentsStorage.DeleteAll(ctx, &domain.UnsentMessageAttachment{
		SenderID:   c.UserID,
		ReceiverID: action.Receiver,
	})
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		err = c.ChatService.PubSubRepository.WriteAction(ctx, action)
		if err != nil {
			return
		}

		return
	}

	c.ChatService.Sanitizer.SanitizePersonalMessage(newMessage)

	action.Payload, err = json.Marshal(newMessage)
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		err = c.ChatService.PubSubRepository.WriteAction(ctx, action)
		if err != nil {
			return
		}

		return
	}

	err = c.ChatService.PubSubRepository.WriteAction(ctx, action)
	if err != nil {
		return
	}
}

func (c *Client) handleUpdateMessageAction(ctx context.Context, action *Action, message *UpdateMessagePayload) {
	attachments, err := c.ChatService.UnsentMessageAttachmentsStorage.GetAll(ctx, &domain.UnsentMessageAttachment{
		SenderID:   c.UserID,
		ReceiverID: action.Receiver,
	})
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		err = c.ChatService.PubSubRepository.WriteAction(ctx, action)
		if err != nil {
			return
		}

		return
	}

	msg := &domain.PersonalMessage{
		ID:          message.MessageID,
		Content:     message.Content,
		Attachments: attachments,
	}

	if len(message.Content) == 0 && len(attachments) == len(message.AttachmentsToDelete) {
		action.Payload, err = errors.MarshalError(errors.ErrInvalidData)
		if err != nil {
			return
		}

		c.ChatService.PubSubRepository.WriteAction(ctx, action)
		return
	}

	newMessage, err := c.ChatService.MessagesRepo.UpdateMessage(ctx, msg, message.AttachmentsToDelete)
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		c.ChatService.PubSubRepository.WriteAction(ctx, action)
		return
	}

	for _, attach := range message.AttachmentsToDelete {
		err = c.ChatService.MessageAttachmentStorage.Delete(attach)
		if err != nil {
			action.Payload, err = errors.MarshalError(err)
			if err != nil {
				return
			}

			c.ChatService.PubSubRepository.WriteAction(ctx, action)
			return
		}
	}

	err = c.ChatService.UnsentMessageAttachmentsStorage.DeleteAll(ctx, &domain.UnsentMessageAttachment{
		SenderID:   c.UserID,
		ReceiverID: action.Receiver,
	})
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		err = c.ChatService.PubSubRepository.WriteAction(ctx, action)
		if err != nil {
			return
		}

		return
	}

	c.ChatService.Sanitizer.SanitizePersonalMessage(newMessage)

	action.Payload, err = json.Marshal(newMessage)
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		c.ChatService.PubSubRepository.WriteAction(ctx, action)
		return
	}

	c.ChatService.PubSubRepository.WriteAction(ctx, action)
}

func (c *Client) handleDeleteMessageAction(ctx context.Context, action *Action, messageID uint) {
	err := c.ChatService.MessagesRepo.DeleteMessage(ctx, messageID)
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		c.ChatService.PubSubRepository.WriteAction(ctx, action)
		return
	}

	c.ChatService.PubSubRepository.WriteAction(ctx, action)
}

func (c *Client) handleSendStickerMessageAction(ctx context.Context, action *Action, message *SendStickerMessagePayload) {
	newStickerMessage, err := c.ChatService.MessagesRepo.StoreStickerMessage(ctx, c.UserID, action.Receiver, message.StickerID)
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		c.ChatService.PubSubRepository.WriteAction(ctx, action)
		return
	}

	action.Payload, err = json.Marshal(newStickerMessage)
	if err != nil {
		action.Payload, err = errors.MarshalError(err)
		if err != nil {
			return
		}

		c.ChatService.PubSubRepository.WriteAction(ctx, action)
		return
	}

	c.ChatService.PubSubRepository.WriteAction(ctx, action)
}

func (c *Client) ClearUnsentAttachments(ctx context.Context) {
	c.UnsentAttachmentReceivers.Range(func(key, value interface{}) bool {
		receiverID := key.(uint)

		attachs, err := c.ChatService.UnsentMessageAttachmentsStorage.GetAll(ctx, &domain.UnsentMessageAttachment{
			SenderID:   c.UserID,
			ReceiverID: receiverID,
		})
		if err != nil {
			return false
		}

		for _, fileName := range attachs {
			err = c.ChatService.MessageAttachmentStorage.Delete(fileName)
			if err != nil {
				return false
			}
		}

		err = c.ChatService.UnsentMessageAttachmentsStorage.DeleteAll(ctx, &domain.UnsentMessageAttachment{
			SenderID:   c.UserID,
			ReceiverID: receiverID,
		})

		return err == nil
	})
}
