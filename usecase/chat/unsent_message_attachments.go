package chat

import (
	"context"
	"mime/multipart"
	"socio/domain"
	"socio/pkg/static"
)

type UnsentMessageAttachmentsStorage interface {
	Store(ctx context.Context, attach *domain.UnsentMessageAttachment) (err error)
	GetAll(ctx context.Context, attach *domain.UnsentMessageAttachment) (fileNames []string, err error)
	DeleteAll(ctx context.Context, attach *domain.UnsentMessageAttachment) (err error)
	Delete(ctx context.Context, attach *domain.UnsentMessageAttachment) (err error)
}

type MessageAttachmentStorage interface {
	Store(fileName string, filePath string, contentType string) (err error)
	Delete(fileName string) (err error)
}

func (c *Service) CreateUnsentMessageAttachments(ctx context.Context, attachs *domain.UnsentMessageAttachment, fhs []*multipart.FileHeader) (filenames []string, err error) {
	for _, fh := range fhs {
		fileName := static.GetUniqueFileName(fh.Filename)

		err = static.SaveFile(fh, "./"+fileName)
		if err != nil {
			return
		}

		err = c.MessageAttachmentStorage.Store(fileName, "./"+fileName, fh.Header.Get("Content-Type"))
		if err != nil {
			return
		}

		err = static.RemoveFile("./" + fileName)
		if err != nil {
			return
		}

		err = c.UnsentMessageAttachmentsStorage.Store(ctx, &domain.UnsentMessageAttachment{
			SenderID:   attachs.SenderID,
			ReceiverID: attachs.ReceiverID,
			FileName:   fileName,
		})
		if err != nil {
			return
		}

		fileName = c.Sanitizer.Sanitize(fileName)

		filenames = append(filenames, fileName)
	}

	return
}

func (c *Service) GetUnsentMessageAttachments(ctx context.Context, attach *domain.UnsentMessageAttachment) (fileNames []string, err error) {
	unsanitizedFileNames, err := c.UnsentMessageAttachmentsStorage.GetAll(ctx, attach)
	if err != nil {
		return
	}

	for _, fileName := range unsanitizedFileNames {
		fileNames = append(fileNames, c.Sanitizer.Sanitize(fileName))
	}

	return
}

func (c *Service) DeleteUnsentMessageAttachments(ctx context.Context, attach *domain.UnsentMessageAttachment) (err error) {
	fileNames, err := c.UnsentMessageAttachmentsStorage.GetAll(ctx, attach)
	if err != nil {
		return
	}

	for _, fileName := range fileNames {
		err = c.MessageAttachmentStorage.Delete(fileName)
		if err != nil {
			return
		}
	}

	err = c.UnsentMessageAttachmentsStorage.DeleteAll(ctx, attach)
	if err != nil {
		return
	}

	return
}

func (c *Service) DeleteUnsentMessageAttachment(ctx context.Context, attach *domain.UnsentMessageAttachment) (err error) {
	err = c.MessageAttachmentStorage.Delete(attach.FileName)
	if err != nil {
		return
	}

	err = c.UnsentMessageAttachmentsStorage.Delete(ctx, attach)
	if err != nil {
		return
	}

	return
}
