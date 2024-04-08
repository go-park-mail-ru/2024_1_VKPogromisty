package rest

import (
	"context"
	defJSON "encoding/json"
	"net/http"
	"socio/domain"
	"socio/errors"
	"socio/internal/rest/middleware"
	"socio/pkg/json"
	"socio/pkg/requestcontext"
	"socio/pkg/sanitizer"
	"socio/usecase/chat"
	"socio/usecase/csrf"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait        = 60 * time.Second
	writeWait       = 10 * time.Second
	pingPeriod      = 1 * time.Second
	maxMessageSize  = 10000
	readBufferSize  = 4096
	writeBufferSize = 4096
	newline         = '\n'
)

type ChatServer struct {
	Service *chat.Service
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  readBufferSize,
	WriteBufferSize: writeBufferSize,
	CheckOrigin:     middleware.CheckOrigin,
}

func NewChatServer(pubSubRepo chat.PubSubRepository, messagesRepo chat.PersonalMessagesRepository, sanitizer *sanitizer.Sanitizer) (chatServer *ChatServer) {
	return &ChatServer{
		Service: chat.NewChatService(pubSubRepo, messagesRepo, sanitizer),
	}
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

	json.ServeJSONBody(r.Context(), w, dialogs)
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
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/chat/messages/ [get]
func (c *ChatServer) HandleGetMessagesByDialog(w http.ResponseWriter, r *http.Request) {
	userID, err := requestcontext.GetUserID(r.Context())
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	peerIDData := r.URL.Query().Get("peerId")
	if peerIDData == "" {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	peerID, err := strconv.ParseUint(peerIDData, 0, 0)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrInvalidData)
		return
	}

	lastMessageIDData := r.URL.Query().Get("lastMessageId")
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

	messages, err := c.Service.GetMessagesByDialog(r.Context(), userID, uint(peerID), uint(lastMessageID))
	if err != nil {
		json.ServeJSONError(r.Context(), w, err)
		return
	}

	if messages == nil {
		messages = make([]*domain.PersonalMessage, 0)
	}

	json.ServeJSONBody(r.Context(), w, messages)

}

// ServeWS godoc
//
//	@Summary		serve websocket connection
//	@Description	Serve websocket connection. You can send actions to connection following simple structure:
//	@Description
//	@Description	{
//	@Description	"type": ActionType,
//	@Description	"receiver": uint,
//	@Description	"payload": interface{}
//	@Description	}
//	@Description
//	@Description	ActionType is a string with one of following values: "SEND_MESSAGE", "UPDATE_MESSAGE", "DELETE_MESSAGE"
//	@Description
//	@Description	If "type" = "SEND_MESSAGE", then payload should be {"content": string}
//	@Description	If "type" = "UPDATE_MESSAGE", then payload should be {"messageId": uint, "content": string}
//	@Description	If "type" = "DELETE_MESSAGE", then payload should be {"messageId": uint}
//	@Description
//	@Description	In response clients, subscribed to corresponding channel, will get same structure back:
//	@Description	{
//	@Description	"type": ActionType,
//	@Description	"receiver": uint,
//	@Description	"payload": interface{}
//	@Description	}
//	@Description
//	@Description	"payload" can be:
//	@Description	PersonalMessage if "type" = "SEND_MESSAGE"
//	@Description	PersonalMessage if "type" = "UPDATE_MESSAGE"
//	@Description	Absent if "type" = "DELETE_MESSAGE"
//	@Description	{"error": string} if error happened at any point of query processing
//	@Description
//
//	@Tags			chat
//	@license.name	Apache 2.0
//	@ID				chat/serve_ws
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//	@Param			X-CSRF-Token	header	string	true	"CSRF token"
//
//	@Produce		json
//	@Success		200
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/chat/ [get]
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

	go c.listenWrite(r.Context(), conn, client)
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
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
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
		err = defJSON.Unmarshal(jsonMessage, action)
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

		err = csrf.NewCSRFService().Check(sessionID, userID, action.CSRFToken)
		if err != nil {
			return
		}

		client.HandleAction(ctx, action)
	}
}

func (c *ChatServer) listenWrite(ctx context.Context, conn *websocket.Conn, client *chat.Client) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		err := conn.Close()
		if err != nil {
			return
		}

		err = c.Service.Unregister(client.UserID)
		if err != nil {
			return
		}
	}()

	for {
		select {
		case message, ok := <-client.Send:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			messageData, err := defJSON.Marshal(message)
			if err != nil {
				return
			}

			w.Write(messageData)

			n := len(client.Send)
			for i := 0; i < n; i++ {
				messageData, err = defJSON.Marshal(<-client.Send)
				if err != nil {
					return
				}
				w.Write([]byte{newline})
				w.Write(messageData)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
