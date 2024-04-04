package repository

import (
	"context"
	"socio/domain"
	"socio/errors"
	customtime "socio/pkg/time"
	"socio/usecase/chat"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	messagesLimit            = 20
	getMessagesByDialogQuery = `
	SELECT id,
		sender_id,
		receiver_id,
		content,
		created_at,
		updated_at
	FROM public.personal_message
	WHERE (
			(
				sender_id = $1
				AND receiver_id = $2
			)
			OR (
				sender_id = $2
				AND receiver_id = $1
			)
		)
		AND id < $3
	ORDER BY created_at DESC
	LIMIT $4;
	`
	getLastMessageIDQuery = `
	SELECT COALESCE(MAX(id), 0) AS last_message_id
	FROM public.personal_message
	WHERE (
			(
				sender_id = $1
				AND receiver_id = $2
			)
			OR (
				sender_id = $2
				AND receiver_id = $1
			)
		);
	`
	getDialogsByUserIDQuery = `
	SELECT u1.id,
		u1.first_name,
		u1.last_name,
		u1.email,
		u1.avatar,
		u1.date_of_birth,
		u1.created_at,
		u1.updated_at,
		u2.id,
		u2.first_name,
		u2.last_name,
		u2.email,
		u2.avatar,
		u2.date_of_birth,
		u2.created_at,
		u2.updated_at,
		pm1.id,
		pm1.sender_id,
		pm1.receiver_id,
		pm1.content,
		pm1.created_at,
		pm1.updated_at
	FROM public.user AS u1
		JOIN public.personal_message AS pm1 ON u1.id = pm1.sender_id
		JOIN public.user AS u2 ON pm1.receiver_id = u2.id
		LEFT JOIN public.personal_message AS pm2 ON (
			(
				pm1.sender_id = pm2.sender_id
				AND pm1.receiver_id = pm2.receiver_id
			)
			OR (
				pm1.sender_id = pm2.receiver_id
				AND pm1.receiver_id = pm2.sender_id
			)
		)
		AND pm1.created_at < pm2.created_at
	WHERE pm2.id IS NULL
		AND (
			pm1.sender_id = $1
			OR pm1.receiver_id = $1
		)
	ORDER BY pm1.created_at DESC;
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

func (pm *PersonalMessages) GetLastMessageID(senderID, receiverID uint) (lastMessageID uint, err error) {
	err = pm.db.QueryRow(context.Background(), getLastMessageIDQuery, senderID, receiverID).Scan(&lastMessageID)
	if err != nil {
		return
	}

	return
}

func (pm *PersonalMessages) GetMessagesByDialog(senderID, receiverID, lastMessageID uint) (messages []*domain.PersonalMessage, err error) {
	rows, err := pm.db.Query(context.Background(), getMessagesByDialogQuery, senderID, receiverID, lastMessageID, messagesLimit)
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

func (pm *PersonalMessages) GetDialogsByUserID(userID uint) (dialogs []*chat.Dialog, err error) {
	rows, err := pm.db.Query(context.Background(), getDialogsByUserIDQuery, userID)
	if err != nil {
		return
	}
	defer rows.Close()

	dialogs = make([]*chat.Dialog, 0)

	for rows.Next() {
		dialog := new(chat.Dialog)
		user1 := new(domain.User)
		user2 := new(domain.User)
		lastMessage := new(domain.PersonalMessage)

		err = rows.Scan(
			&user1.ID,
			&user1.FirstName,
			&user1.LastName,
			&user1.Email,
			&user1.Avatar,
			&user1.DateOfBirth.Time,
			&user1.CreatedAt.Time,
			&user1.UpdatedAt.Time,
			&user2.ID,
			&user2.FirstName,
			&user2.LastName,
			&user2.Email,
			&user2.Avatar,
			&user2.DateOfBirth.Time,
			&user2.CreatedAt.Time,
			&user2.UpdatedAt.Time,
			&lastMessage.ID,
			&lastMessage.SenderID,
			&lastMessage.ReceiverID,
			&lastMessage.Content,
			&lastMessage.CreatedAt.Time,
			&lastMessage.UpdatedAt.Time,
		)
		if err != nil {
			return
		}

		dialog.User1 = user1
		dialog.User2 = user2
		dialog.LastMessage = lastMessage

		dialogs = append(dialogs, dialog)
	}

	return
}

func (pm *PersonalMessages) StoreMessage(msg *domain.PersonalMessage) (newMsg *domain.PersonalMessage, err error) {
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

func (pm *PersonalMessages) UpdateMessage(msg *domain.PersonalMessage) (updatedMsg *domain.PersonalMessage, err error) {
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

func (pm *PersonalMessages) DeleteMessage(msgID uint) (err error) {
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
