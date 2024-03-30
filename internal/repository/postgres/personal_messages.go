package repository

import (
	"context"
	"socio/domain"
	"socio/errors"
	customtime "socio/pkg/time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	getMessagesByDialogQuery = `
	SELECT id,
		sender_id,
		receiver_id,
		content,
		created_at,
		updated_at
	FROM public.personal_message
	WHERE (
			sender_id = $1
			AND receiver_id = $2
		)
		OR (
			sender_id = $2
			AND receiver_id = $1
		)
	ORDER BY created_at DESC;
	`
	storePersonalMessageQuery = `
	INSERT INTO public.personal_message (sender_id, receiver_id, content)
	VALUES ($1, $2, $3)
	RETURNING id,
		sender_id,
		receiver_id,
		content,
		created_at,
		updated_at;
	`
	updatePersonalMessageQuery = `
	UPDATE public.personal_message
	SET content = $1
	WHERE id = $2
	RETURNING id,
		sender_id,
		receiver_id,
		content,
		created_at,
		updated_at;
	`
	deletePersonalMessageQuery = `
	DELETE FROM public.personal_message
	WHERE id = $1;
	`
)

type PersonalMessages struct {
	db *pgxpool.Pool
	TP customtime.TimeProvider
}

func NewPersonalMessages(db *pgxpool.Pool, tp customtime.TimeProvider) *PersonalMessages {
	return &PersonalMessages{
		db: db,
		TP: tp,
	}
}

func (pm *PersonalMessages) GetMessagesByDialog(senderID, receiverID uint) (messages []*domain.PersonalMessage, err error) {
	rows, err := pm.db.Query(context.Background(), getMessagesByDialogQuery, senderID, receiverID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		msg := new(domain.PersonalMessage)

		err = rows.Scan(
			&msg.ID,
			&msg.SenderID,
			&msg.ReceiverID,
			&msg.Content,
			&msg.CreatedAt.Time,
			&msg.UpdatedAt.Time,
		)
		if err != nil {
			return
		}

		messages = append(messages, msg)
	}

	return

}

func (pm *PersonalMessages) StorePersonalMessage(msg *domain.PersonalMessage) (newMsg *domain.PersonalMessage, err error) {
	newMsg = new(domain.PersonalMessage)
	err = pm.db.QueryRow(context.Background(), storePersonalMessageQuery,
		msg.SenderID,
		msg.ReceiverID,
		msg.Content,
	).Scan(
		&newMsg.ID,
		&newMsg.SenderID,
		&newMsg.ReceiverID,
		&newMsg.Content,
		&newMsg.CreatedAt.Time,
		&newMsg.UpdatedAt.Time,
	)
	if err != nil {
		return
	}

	return
}

func (pm *PersonalMessages) UpdatePersonalMessage(msg *domain.PersonalMessage) (updatedMsg *domain.PersonalMessage, err error) {
	updatedMsg = new(domain.PersonalMessage)
	err = pm.db.QueryRow(context.Background(), updatePersonalMessageQuery, msg.Content, msg.ID).Scan(
		&updatedMsg.ID,
		&updatedMsg.SenderID,
		&updatedMsg.ReceiverID,
		&updatedMsg.Content,
		&updatedMsg.CreatedAt.Time,
		&updatedMsg.UpdatedAt.Time,
	)
	if err != nil {
		return
	}

	return
}

func (pm *PersonalMessages) DeletePersonalMessage(msgID uint) (err error) {
	result, err := pm.db.Exec(context.Background(), deletePersonalMessageQuery, msgID)
	if err != nil {
		return
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return errors.ErrInvalidBody
	}
	if rowsAffected != 1 {
		return errors.ErrRowsAffected
	}

	return
}
