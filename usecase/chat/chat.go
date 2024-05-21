package chat

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"socio/domain"
	"socio/errors"
	"socio/pkg/sanitizer"
	"socio/pkg/static"
	"sync"

	"github.com/google/uuid"
)

const (
	defaultMessagesAmount = 20
)

// Service will: register and unregister new clients by userID, establish pub/sub connection to redis
type Service struct {
	Clients                         *sync.Map
	PubSubRepository                PubSubRepository
	MessagesRepo                    PersonalMessagesRepository
	UnsentMessageAttachmentsStorage UnsentMessageAttachmentsStorage
	MessageAttachmentStorage        MessageAttachmentStorage
	StickerStorage                  StickerStorage
	Sanitizer                       *sanitizer.Sanitizer
}

type StickerStorage interface {
	Store(fileName string, filePath string, contentType string) (err error)
	Delete(fileName string) (err error)
}

func NewChatService(pubSubRepo PubSubRepository, unsentMessageAttachmentsStorage UnsentMessageAttachmentsStorage, messagesRepo PersonalMessagesRepository, stickerStorage StickerStorage, messageAttachmentStorage MessageAttachmentStorage, sanitizer *sanitizer.Sanitizer) (chatService *Service) {
	return &Service{
		Clients:                         &sync.Map{},
		PubSubRepository:                pubSubRepo,
		UnsentMessageAttachmentsStorage: unsentMessageAttachmentsStorage,
		MessagesRepo:                    messagesRepo,
		StickerStorage:                  stickerStorage,
		MessageAttachmentStorage:        messageAttachmentStorage,
		Sanitizer:                       sanitizer,
	}
}

func (s *Service) Register(ctx context.Context, userID uint) (c *Client, err error) {
	cData, ok := s.Clients.Load(userID)
	if ok {
		c = cData.(*Client)
		return
	}

	c, err = NewClient(userID, s)
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

func (s *Service) GetClient(ctx context.Context, userID uint) (c *Client, err error) {
	cData, ok := s.Clients.Load(userID)
	if !ok {
		err = errors.ErrNotFound
		return
	}

	c = cData.(*Client)

	return
}

func (s *Service) GetMessagesByDialog(ctx context.Context, userID, peerID, lastMessageID, messagesAmount uint) (messages []*domain.PersonalMessage, err error) {
	if lastMessageID == 0 {
		lastMessageID, err = s.MessagesRepo.GetLastMessageID(ctx, userID, peerID)
		if err != nil {
			return
		}
		lastMessageID++
	}

	if messagesAmount == 0 {
		messagesAmount = defaultMessagesAmount
	}

	messages, err = s.MessagesRepo.GetMessagesByDialog(ctx, userID, peerID, lastMessageID, messagesAmount)
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

func (s *Service) GetStickersByAuthorID(ctx context.Context, authorID uint) (stickers []*domain.Sticker, err error) {
	stickers, err = s.MessagesRepo.GetStickersByAuthorID(ctx, authorID)
	if err != nil {
		return
	}

	for _, sticker := range stickers {
		s.Sanitizer.SanitizeSticker(sticker)
	}

	return
}

func (s *Service) GetAllStickers(ctx context.Context) (stickers []*domain.Sticker, err error) {
	stickers, err = s.MessagesRepo.GetAllStickers(ctx)
	if err != nil {
		return
	}

	for _, sticker := range stickers {
		s.Sanitizer.SanitizeSticker(sticker)
	}

	return
}

func (s *Service) CreateSticker(ctx context.Context, sticker *domain.Sticker, image *multipart.FileHeader) (newSticker *domain.Sticker, err error) {
	if sticker.Name == "" || sticker.AuthorID == 0 || image == nil {
		err = errors.ErrInvalidData
		return
	}

	fileName := uuid.NewString() + filepath.Ext(image.Filename)

	err = static.SaveFile(image, "./"+fileName)
	if err != nil {
		return
	}

	err = s.StickerStorage.Store(fileName, "./"+fileName, image.Header.Get("Content-Type"))
	if err != nil {
		return
	}

	sticker.FileName = fileName

	err = static.RemoveFile("./" + fileName)
	if err != nil {
		return
	}

	newSticker, err = s.MessagesRepo.StoreSticker(ctx, sticker)
	if err != nil {
		return
	}

	s.Sanitizer.SanitizeSticker(newSticker)
	return
}

func (s *Service) DeleteSticker(ctx context.Context, stickerID uint, userID uint) (err error) {
	sticker, err := s.MessagesRepo.GetStickerByID(ctx, stickerID)
	if err != nil {
		return
	}

	if sticker.AuthorID != userID {
		err = errors.ErrForbidden
		return
	}

	err = s.StickerStorage.Delete(sticker.FileName)
	if err != nil {
		return
	}

	err = s.MessagesRepo.DeleteSticker(ctx, stickerID)
	if err != nil {
		return
	}

	return
}
