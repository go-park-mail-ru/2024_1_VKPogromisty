package repository

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/pkg/contextlogger"

	"github.com/jackc/pgx/v4"
)

const (
	GetStickerByIDQuery = `
	SELECT id,
		author_id,
		name,
		file_name,
		created_at,
		updated_at
	FROM public.sticker
	WHERE id = $1;
	`
	GetStickersByAuthorIDQuery = `
	SELECT id,
		author_id,
		name,
		file_name,
		created_at,
		updated_at
	FROM public.sticker
	WHERE author_id = $1;
	`
	GetAllStickersQuery = `
	SELECT id,
		author_id,
		name,
		file_name,
		created_at,
		updated_at
	FROM public.sticker;
	`
	StoreStickerQuery = `
	INSERT INTO public.sticker (author_id, name, file_name)
	VALUES ($1, $2, $3)
	RETURNING id,
		author_id, 
		name, 
		file_name,
		created_at,
		updated_at;
	`
	DeleteStickerQuery = `
	DELETE FROM public.sticker
	WHERE id = $1;
	`
	StoreStickerMessageQuery = `
	INSERT INTO public.personal_message (sender_id, receiver_id, sticker_id)
	VALUES ($1, $2, $3)
	RETURNING id,
		sender_id,
		receiver_id,
		sticker_id,
		created_at, 
		updated_at;
	`
)

func (pm *PersonalMessages) GetStickerByID(ctx context.Context, stickerID uint) (sticker *domain.Sticker, err error) {
	sticker = new(domain.Sticker)

	contextlogger.LogSQL(ctx, GetStickerByIDQuery, stickerID)

	err = pm.db.QueryRow(context.Background(), GetStickerByIDQuery, stickerID).Scan(
		&sticker.ID,
		&sticker.AuthorID,
		&sticker.Name,
		&sticker.FileName,
		&sticker.CreatedAt.Time,
		&sticker.UpdatedAt.Time,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
		}

		return
	}

	return
}

func (pm *PersonalMessages) GetStickersByAuthorID(ctx context.Context, authorID uint) (stickers []*domain.Sticker, err error) {
	stickers = make([]*domain.Sticker, 0)

	contextlogger.LogSQL(ctx, GetStickersByAuthorIDQuery, authorID)

	rows, err := pm.db.Query(context.Background(), GetStickersByAuthorIDQuery, authorID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		sticker := new(domain.Sticker)

		err = rows.Scan(
			&sticker.ID,
			&sticker.AuthorID,
			&sticker.Name,
			&sticker.FileName,
			&sticker.CreatedAt.Time,
			&sticker.UpdatedAt.Time,
		)
		if err != nil {
			return
		}

		stickers = append(stickers, sticker)
	}

	return
}

func (pm *PersonalMessages) GetAllStickers(ctx context.Context) (stickers []*domain.Sticker, err error) {
	stickers = make([]*domain.Sticker, 0)

	contextlogger.LogSQL(ctx, GetAllStickersQuery)

	rows, err := pm.db.Query(context.Background(), GetAllStickersQuery)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		sticker := new(domain.Sticker)

		err = rows.Scan(
			&sticker.ID,
			&sticker.AuthorID,
			&sticker.Name,
			&sticker.FileName,
			&sticker.CreatedAt.Time,
			&sticker.UpdatedAt.Time,
		)
		if err != nil {
			return
		}

		stickers = append(stickers, sticker)
	}

	return
}

func (pm *PersonalMessages) StoreSticker(ctx context.Context, sticker *domain.Sticker) (newSticker *domain.Sticker, err error) {
	newSticker = new(domain.Sticker)

	contextlogger.LogSQL(ctx, StoreStickerQuery, sticker.AuthorID, sticker.Name, sticker.FileName)

	err = pm.db.QueryRow(context.Background(), StoreStickerQuery, sticker.AuthorID, sticker.Name, sticker.FileName).Scan(
		&newSticker.ID,
		&newSticker.AuthorID,
		&newSticker.Name,
		&newSticker.FileName,
		&newSticker.CreatedAt.Time,
		&newSticker.UpdatedAt.Time,
	)
	if err != nil {
		return
	}

	return
}

func (pm *PersonalMessages) DeleteSticker(ctx context.Context, stickerID uint) (err error) {
	contextlogger.LogSQL(ctx, DeleteStickerQuery, stickerID)

	_, err = pm.db.Exec(context.Background(), DeleteStickerQuery, stickerID)
	if err != nil {
		return
	}

	return
}

func (pm *PersonalMessages) StoreStickerMessage(ctx context.Context, senderID, receiverID, stickerID uint) (newStickerMessage *domain.PersonalMessage, err error) {
	newStickerMessage = new(domain.PersonalMessage)
	sticker := new(domain.Sticker)

	contextlogger.LogSQL(ctx, StoreStickerMessageQuery, senderID, receiverID, stickerID)

	err = pm.db.QueryRow(context.Background(), StoreStickerMessageQuery, senderID, receiverID, stickerID).Scan(
		&newStickerMessage.ID,
		&newStickerMessage.SenderID,
		&newStickerMessage.ReceiverID,
		&sticker.ID,
		&newStickerMessage.CreatedAt.Time,
		&newStickerMessage.UpdatedAt.Time,
	)
	if err != nil {
		return
	}

	sticker, err = pm.GetStickerByID(ctx, stickerID)
	if err != nil {
		return
	}

	newStickerMessage.Sticker = sticker

	return
}
