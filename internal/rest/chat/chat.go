package rest

import (
	"context"
	"mime/multipart"
	"net/http"
	"socio/domain"
	"socio/errors"
	"socio/internal/rest/middleware"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
	customtime "socio/pkg/time"
	"socio/usecase/chat"
	"socio/usecase/csrf"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	easyjson "github.com/mailru/easyjson"
)

const (
	pongWait                 = 60 * time.Second
	writeWait                = 10 * time.Second
	pingPeriod               = 1 * time.Second
	maxMessageSize           = 10000
	readBufferSize           = 4096
	writeBufferSize          = 4096
	newline                  = '\n'
	PeerIDQueryParam         = "peerId"
	LastMessageIDQueryParam  = "lastMessageId"
	MessagesAmountQueryParam = "messagesAmount"
)

type ChatServer struct {
	Service ChatService
	wsConns *sync.Map
}

type ChatService interface {
	CreateSticker(ctx context.Context, sticker *domain.Sticker, image *multipart.FileHeader) (newSticker *domain.Sticker, err error)
	CreateUnsentMessageAttachments(ctx context.Context, attachs *domain.UnsentMessageAttachment, fhs []*multipart.FileHeader) (filenames []string, err error)
	DeleteSticker(ctx context.Context, stickerID uint, userID uint) (err error)
	DeleteUnsentMessageAttachment(ctx context.Context, attach *domain.UnsentMessageAttachment) (err error)
	DeleteUnsentMessageAttachments(ctx context.Context, attach *domain.UnsentMessageAttachment) (err error)
	GetAllStickers(ctx context.Context) (stickers []*domain.Sticker, err error)
	GetClient(ctx context.Context, userID uint) (c *chat.Client, err error)
	GetDialogsByUserID(ctx context.Context, userID uint) (dialogs []*domain.Dialog, err error)
	GetMessagesByDialog(ctx context.Context, userID uint, peerID uint, lastMessageID uint, messagesAmount uint) (messages []*domain.PersonalMessage, err error)
	GetStickersByAuthorID(ctx context.Context, authorID uint) (stickers []*domain.Sticker, err error)
	GetUnsentMessageAttachments(ctx context.Context, attach *domain.UnsentMessageAttachment) (fileNames []string, err error)
	Register(ctx context.Context, userID uint) (c *chat.Client, err error)
	Unregister(userID uint) (err error)
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  readBufferSize,
	WriteBufferSize: writeBufferSize,
	CheckOrigin:     middleware.CheckOrigin,
}

func NewChatServer(service ChatService) (chatServer *ChatServer) {
	return &ChatServer{
		Service: service,
		wsConns: &sync.Map{},
	}
}

func (c *ChatServer) getWSConns(userID uint) (conns []*websocket.Conn, ok bool) {
	untypedConns, ok := c.wsConns.Load(userID)
	if !ok {
		conns = nil
		return
	}

	conns, ok = untypedConns.([]*websocket.Conn)
	if !ok {
		return
	}

	return
}

// HandleGetDialogs godoc
//
//	@Summary		get user dialogs
//	@Description	get user dialogs
//	@Tags			chat
//	@license.name	Apache 2.0
//	@ID				chat/get_dialogs
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=[]domain.Dialog}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/chat/dialogs/ [get]
func (c *ChatServer) HandleGetDialogs(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	dialogs, err := c.Service.GetDialogsByUserID(r.Context(), userID)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, dialogs, http.StatusOK)
}

// HandleGetMessagesByDialog godoc
//
//	@Summary		get messages by dialog
//	@Description	get messages by dialog with pagination
//	@Tags			chat
//	@license.name	Apache 2.0
//	@ID				chat/get_messages
//	@Accept			json
//
//	@Param			Cookie			header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			peerId			query	uint	true	"ID of the peer"
//	@Param			lastMessageId	query	uint	false	"ID of the last message, if last messages needed, should be set to 0"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=[]domain.PersonalMessage}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/chat/messages/ [get]
func (c *ChatServer) HandleGetMessagesByDialog(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	peerIDData := r.URL.Query().Get(PeerIDQueryParam)
	if peerIDData == "" {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	peerID, err := strconv.ParseUint(peerIDData, 0, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	lastMessageIDData := r.URL.Query().Get(LastMessageIDQueryParam)
	var lastMessageID uint64
	if lastMessageIDData == "" {
		lastMessageID = 0
	} else {
		lastMessageID, err = strconv.ParseUint(lastMessageIDData, 0, 0)
		if err != nil {
			json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
			return
		}
	}

	messagesAmountData := r.URL.Query().Get(MessagesAmountQueryParam)
	var messagesAmount uint64

	if messagesAmountData == "" {
		messagesAmount = 0
	} else {
		messagesAmount, err = strconv.ParseUint(messagesAmountData, 0, 0)
		if err != nil {
			json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
			return
		}
	}

	messages, err := c.Service.GetMessagesByDialog(r.Context(), userID, uint(peerID), uint(lastMessageID), uint(messagesAmount))
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	if messages == nil {
		messages = make([]*domain.PersonalMessage, 0)
	}

	json.ServeJSONBody(r.Context(), w, messages, http.StatusOK)
}

// ServeWS godoc
//
//		@Summary		serve websocket connection
//		@Description	Serve websocket connection. You can send actions to connection following simple structure:
//		@Description
//		@Description	{
//		@Description	"type": ActionType,
//		@Description	"receiver": uint,
//		@Description	"csrfToken": string,
//		@Description	"payload": interface{}
//		@Description	}
//		@Description
//		@Description	ActionType is a string with one of following values: "SEND_MESSAGE", "UPDATE_MESSAGE", "DELETE_MESSAGE", "SEND_STICKER_MESSAGE"
//		@Description
//		@Description	If "type" = "SEND_MESSAGE", then payload should be {"content": string, "attachments": []string}
//		@Description	If "type" = "UPDATE_MESSAGE", then payload should be {"messageId": uint, "content": string, attachmentsToDelete: []string}
//		@Description	If "type" = "DELETE_MESSAGE", then payload should be {"messageId": uint}
//		@Description	If "type" = "SEND_STICKER_MESSAGE", then payload should be {"stickerId": uint}
//		@Description
//		@Description	In response clients, subscribed to corresponding channel, will get same structure back:
//		@Description	{
//		@Description	"type": ActionType,
//		@Description	"receiver": uint,
//	 	@Description	 "csrfToken": string,
//		@Description	"payload": interface{}
//		@Description	}
//		@Description
//		@Description	"payload" can be:
//		@Description	PersonalMessage if "type" = "SEND_MESSAGE"
//		@Description	PersonalMessage if "type" = "UPDATE_MESSAGE"
//		@Description	Absent if "type" = "DELETE_MESSAGE"
//		@Description	PersonalMessage if "type" = "SEND_STICKER_MESSAGE"
//		@Description	{"error": string} if error happened at any point of query processing
//		@Description
//
//		@Tags			chat
//		@license.name	Apache 2.0
//		@ID				chat/serve_ws
//		@Accept			json
//
//		@Param			Cookie	header	string	true	"session_id=some_session"
//
//		@Produce		json
//		@Success		200
//		@Failure		400	{object}	errors.HTTPError
//		@Failure		401	{object}	errors.HTTPError
//		@Failure		500	{object}	errors.HTTPError
//		@Router			/chat/ [get]
func (c *ChatServer) ServeWS(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	client, err := c.Service.Register(r.Context(), userID)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	conns, ok := c.wsConns.Load(userID)
	if !ok {
		conns = make([]*websocket.Conn, 0, 1)
		conns = append(conns.([]*websocket.Conn), conn)
		c.wsConns.Store(userID, conns)

		go c.listenWrite(r.Context(), client)
	} else {
		conns = append(conns.([]*websocket.Conn), conn)
		c.wsConns.Store(userID, conns)
	}

	go c.listenRead(r.Context(), conn, client)
}

func (c *ChatServer) listenRead(ctx context.Context, conn *websocket.Conn, client *chat.Client) {
	defer func() {
		err := conn.Close()
		if err != nil {
			return
		}
	}()

	conn.SetReadLimit(maxMessageSize)
	err := conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return
	}

	conn.SetPongHandler(func(string) error {
		err := conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}

		return nil
	})

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return
			}
			break
		}

		action := new(chat.Action)
		err = easyjson.Unmarshal(jsonMessage, action)
		if err != nil {
			return
		}

		sessionID, err := requestcontext.GetSessionID(ctx)
		if err != nil {
			return
		}

		userID, err := requestcontext.GetUserID(ctx)
		if err != nil {
			return
		}

		err = csrf.NewCSRFService(customtime.RealTimeProvider{}).Check(sessionID, userID, action.CSRFToken)
		if err != nil {
			return
		}

		client.HandleAction(ctx, action)
	}
}

func (c *ChatServer) listenWrite(ctx context.Context, client *chat.Client) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()

		conns, ok := c.getWSConns(client.UserID)
		if !ok {
			return
		}

		for _, conn := range conns {
			err := conn.Close()
			if err != nil {
				return
			}
		}

		err := c.Service.Unregister(client.UserID)
		if err != nil {
			return
		}
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				return
			}

			messages := make([][]byte, 0, len(client.Send)+1)

			messageData, err := easyjson.Marshal(message)
			if err != nil {
				return
			}

			messages = append(messages, messageData)

			n := len(client.Send)
			for i := 0; i < n; i++ {
				messageData, err = easyjson.Marshal(<-client.Send)
				if err != nil {
					return
				}

				messages = append(messages, messageData)
			}

			conns, ok := c.getWSConns(client.UserID)
			if !ok {
				return
			}

			for _, conn := range conns {
				err := conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err != nil {
					return
				}

				if !ok {
					err := conn.WriteMessage(websocket.CloseMessage, []byte{})
					if err != nil {
						return
					}

					return
				}

				w, err := conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}

				for _, message := range messages {
					_, err := w.Write([]byte{newline})
					if err != nil {
						return
					}

					_, err = w.Write(message)
					if err != nil {
						return
					}
				}

				if err := w.Close(); err != nil {
					return
				}
			}

		case <-ticker.C:
			conns, ok := c.getWSConns(client.UserID)
			if !ok {
				return
			}

			for _, conn := range conns {
				err := conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err != nil {
					return
				}

				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}

// GetStickersByAuthorID godoc
//
//	@Summary		get stickers by author ID
//	@Description	get stickers by author ID
//	@Tags			chat
//	@license.name	Apache 2.0
//	@ID				chat/get_stickers
//	@Accept			json
//
//	@Param			authorID	path	uint	true	"ID of the author"
//	@Param			Cookie		header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=[]domain.Sticker}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/chat/stickers/{authorID} [get]
func (c *ChatServer) HandleGetStickersByAuthorID(w http.ResponseWriter, r *http.Request) {
	authorIDData, ok := mux.Vars(r)["authorID"]
	if !ok {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	authorID, err := strconv.ParseUint(authorIDData, 0, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	stickers, err := c.Service.GetStickersByAuthorID(r.Context(), uint(authorID))
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, stickers, http.StatusOK)
}

// GetAllStickers godoc
//
//	@Summary		get all stickers
//	@Description	get all stickers
//	@Tags			chat
//	@license.name	Apache 2.0
//	@ID				chat/get_all_stickers
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=[]domain.Sticker}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/chat/stickers/ [get]
func (c *ChatServer) HandleGetAllStickers(w http.ResponseWriter, r *http.Request) {
	stickers, err := c.Service.GetAllStickers(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, stickers, http.StatusOK)
}

// CreateSticker godoc
//
//	@Summary		create sticker
//	@Description	create sticker
//	@Tags			chat
//	@license.name	Apache 2.0
//	@ID				chat/create_sticker
//	@Accept			multipart/form-data
//
//	@Param			name	formData	string	true	"Name of the sticker"
//	@Param			image	formData	file	true	"Image of the sticker"
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//
//	@Produce		json
//	@Success		201	{object}	json.JSONResponse{body=domain.Sticker}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/chat/stickers/ [post]
func (c *ChatServer) HandleCreateSticker(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1000 << 20)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidBody)
		return
	}

	sticker := new(domain.Sticker)

	sticker.AuthorID, err = requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	sticker.Name = strings.TrimSpace(r.PostFormValue("name"))

	_, fh, err := r.FormFile("image")
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidBody)
		return
	}

	sticker, err = c.Service.CreateSticker(r.Context(), sticker, fh)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, sticker, http.StatusCreated)
}

// DeleteSticker godoc
//
//	@Summary		delete sticker
//	@Description	delete sticker
//	@Tags			chat
//	@license.name	Apache 2.0
//	@ID				chat/delete_sticker
//	@Accept			json
//
//	@Param			stickerID	path	uint	true	"ID of the sticker"
//	@Param			Cookie		header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		404	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/chat/stickers/{stickerID} [delete]
func (c *ChatServer) HandleDeleteSticker(w http.ResponseWriter, r *http.Request) {
	stickerIDData, ok := mux.Vars(r)["stickerID"]
	if !ok {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	stickerID, err := strconv.ParseUint(stickerIDData, 0, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	err = c.Service.DeleteSticker(r.Context(), uint(stickerID), userID)
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, nil, http.StatusNoContent)
}

// HandleGetUnsentMessageAttachments godoc
//
//	@Summary		get unsent message attachments
//	@Description	get unsent message attachments, returns array of filenames
//	@Tags			chat
//	@license.name	Apache 2.0
//	@ID				chat/get_unsent_message_attachments
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			receiverID	path	uint	true	"ID of the receiver"
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=[]string}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/chat/dialogs/{receiverID}/unsent-attachments/ [get]
func (c *ChatServer) HandleGetUnsentMessageAttachments(w http.ResponseWriter, r *http.Request) {
	senderID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	receiverIDData, ok := mux.Vars(r)["receiverID"]
	if !ok {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	receiverID, err := strconv.ParseUint(receiverIDData, 0, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	attachments, err := c.Service.GetUnsentMessageAttachments(r.Context(), &domain.UnsentMessageAttachment{
		SenderID:   senderID,
		ReceiverID: uint(receiverID),
	})
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, attachments, http.StatusOK)
}

// HandleDeleteUnsentMessageAttachments godoc
//
//	@Summary		delete unsent message attachments
//	@Description	delete unsent message attachments
//	@Tags			chat
//	@license.name	Apache 2.0
//	@ID				chat/delete_unsent_message_attachments
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			receiverID	path	uint	true	"ID of the receiver"
//
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/chat/dialogs/{receiverID}/unsent-attachments/ [delete]
func (c *ChatServer) HandleDeleteUnsentMessageAttachments(w http.ResponseWriter, r *http.Request) {
	senderID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	receiverIDData, ok := mux.Vars(r)["receiverID"]
	if !ok {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	receiverID, err := strconv.ParseUint(receiverIDData, 0, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	err = c.Service.DeleteUnsentMessageAttachments(r.Context(), &domain.UnsentMessageAttachment{
		SenderID:   senderID,
		ReceiverID: uint(receiverID),
	})
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, nil, http.StatusNoContent)
}

// HandleCreateUnsentMessageAttachments godoc
//
//	@Summary		create unsent message attachments
//	@Description	create unsent message attachments
//	@Tags			chat
//	@license.name	Apache 2.0
//	@ID				chat/create_unsent_message_attachments
//	@Accept			multipart/form-data
//
//	@Param			attachment	formData	file	true	"Attachment file"
//	@Param			Cookie		header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			receiverID	path	uint	true	"ID of the receiver"
//
//	@Produce		json
//	@Success		201	{object}	json.JSONResponse{body=[]string}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/chat/dialogs/{receiverID}/unsent-attachments/ [post]
func (c *ChatServer) HandleCreateUnsentMessageAttachments(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1000 << 20)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidBody)
		return
	}
	senderID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	receiverIDData, ok := mux.Vars(r)["receiverID"]
	if !ok {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	receiverID, err := strconv.ParseUint(receiverIDData, 0, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	filenames, err := c.Service.CreateUnsentMessageAttachments(r.Context(), &domain.UnsentMessageAttachment{
		SenderID:   senderID,
		ReceiverID: uint(receiverID),
	},
		r.MultipartForm.File["attachment"])
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, filenames, http.StatusCreated)
}

// HandleDeleteUnsentMessageAttachment godoc
//
//	@Summary		delete unsent message attachment
//	@Description	delete unsent message attachment
//	@Tags			chat
//	@license.name	Apache 2.0
//	@ID				chat/delete_unsent_message_attachment
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//	@Param			receiverID	path	uint	true	"ID of the receiver"
//	@Param			fileName	path	string	true	"Name of the file"
//
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		403	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/chat/dialogs/{receiverID}/unsent-attachments/{fileName} [delete]
func (c *ChatServer) HandleDeleteUnsentMessageAttachment(w http.ResponseWriter, r *http.Request) {
	senderID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	receiverIDData, ok := mux.Vars(r)["receiverID"]
	if !ok {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	receiverID, err := strconv.ParseUint(receiverIDData, 0, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	fileNameData, ok := mux.Vars(r)["fileName"]
	if !ok {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	err = c.Service.DeleteUnsentMessageAttachment(r.Context(), &domain.UnsentMessageAttachment{
		SenderID:   senderID,
		ReceiverID: uint(receiverID),
		FileName:   fileNameData,
	})
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, nil, http.StatusNoContent)
}
