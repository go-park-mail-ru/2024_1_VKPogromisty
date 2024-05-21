package repository

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/pkg/contextlogger"

	"github.com/jackc/pgx/v4"
)

const (
	getStickerByIDQuery = `
	SELECT id,
		author_id,
		name,
		file_name,
		created_at,
		updated_at
	FROM public.sticker
	WHERE id = $1;
	`
	getStickersByAuthorIDQuery = `
	SELECT id,
		author_id,
		name,
		file_name,
		created_at,
		updated_at
	FROM public.sticker
	WHERE author_id = $1;
	`
	getAllStickersQuery = `
	SELECT id,
		author_id,
		name,
		file_name,
		created_at,
		updated_at
	FROM public.sticker;
	`
	storeStickerQuery = `
	INSERT INTO public.sticker (author_id, name, file_name)
	VALUES ($1, $2, $3)
	RETURNING id,
		author_id, 
		name, 
		file_name,
		created_at,
		updated_at;
	`
	deleteStickerQuery = `
	DELETE FROM public.sticker
	WHERE id = $1;
	`
	storeStickerMessageQuery = `
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

	contextlogger.LogSQL(ctx, getStickerByIDQuery, stickerID)

	err = pm.db.QueryRow(context.Background(), getStickerByIDQuery, stickerID).Scan(
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

	contextlogger.LogSQL(ctx, getStickersByAuthorIDQuery, authorID)

	rows, err := pm.db.Query(context.Background(), getStickersByAuthorIDQuery, authorID)
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

	contextlogger.LogSQL(ctx, getAllStickersQuery)

	rows, err := pm.db.Query(context.Background(), getAllStickersQuery)
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

	contextlogger.LogSQL(ctx, storeStickerQuery, sticker.AuthorID, sticker.Name, sticker.FileName)

	err = pm.db.QueryRow(context.Background(), storeStickerQuery, sticker.AuthorID, sticker.Name, sticker.FileName).Scan(
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
	contextlogger.LogSQL(ctx, deleteStickerQuery, stickerID)

	_, err = pm.db.Exec(context.Background(), deleteStickerQuery, stickerID)
	if err != nil {
		return
	}

	return
}

func (pm *PersonalMessages) StoreStickerMessage(ctx context.Context, senderID, receiverID, stickerID uint) (newStickerMessage *domain.PersonalMessage, err error) {
	newStickerMessage = new(domain.PersonalMessage)
	sticker := new(domain.Sticker)

	contextlogger.LogSQL(ctx, storeStickerMessageQuery, senderID, receiverID, stickerID)

	err = pm.db.QueryRow(context.Background(), storeStickerMessageQuery, senderID, receiverID, stickerID).Scan(
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
