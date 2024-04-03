package rest

import (
	defJSON "encoding/json"
	"net/http"
	"socio/errors"
	"socio/internal/rest/middleware"
	"socio/pkg/json"
	"socio/usecase/chat"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	writeWait      = 10 * time.Second
	maxMessageSize = 10000
)

type ChatServer struct {
	Service *chat.Service
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     middleware.CheckOrigin,
}

func NewChatServer(pubSubRepo chat.PubSubRepository, messagesRepo chat.PersonalMessagesRepository) (chatServer *ChatServer) {
	return &ChatServer{
		Service: chat.NewChatService(pubSubRepo, messagesRepo),
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
//
//	@Produce		json
//	@Success		200	{object}	json.JSONResponse{body=[]chat.Dialog}
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/chat/dialogs/ [get]
func (c *ChatServer) HandleGetDialogs(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)

	dialogs, err := c.Service.GetDialogsByUserID(userID)
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	json.ServeJSONBody(w, dialogs)
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
	userID := r.Context().Value(middleware.UserIDKey).(uint)

	peerIDData := r.URL.Query().Get("peerId")
	if peerIDData == "" {
		json.ServeJSONError(w, errors.ErrInvalidData)
		return
	}

	peerID, err := strconv.ParseUint(peerIDData, 0, 0)
	if err != nil {
		json.ServeJSONError(w, errors.ErrInvalidData)
		return
	}

	lastMessageIDData := r.URL.Query().Get("lastMessageId")
	var lastMessageID uint64
	if lastMessageIDData == "" {
		lastMessageID = 0
	} else {
		lastMessageID, err = strconv.ParseUint(lastMessageIDData, 0, 0)
		if err != nil {
			json.ServeJSONError(w, errors.ErrInvalidData)
			return
		}
	}

	messages, err := c.Service.GetMessagesByDialog(userID, uint(peerID), uint(lastMessageID))
	if err != nil {
		json.ServeJSONError(w, err)
		return
	}

	json.ServeJSONBody(w, messages)

}

// ServeWS godoc
//
//	@Summary		serve websocket connection
//	@Description	serve websocket connection
//	@Tags			chat
//	@license.name	Apache 2.0
//	@ID				chat/serve_ws
//	@Accept			json
//
//	@Param			Cookie	header	string	true	"session_id=some_session"
//
//	@Produce		json
//	@Success		200
//	@Failure		400	{object}	errors.HTTPError
//	@Failure		401	{object}	errors.HTTPError
//	@Failure		500	{object}	errors.HTTPError
//	@Router			/chat/ [get]
func (c *ChatServer) ServeWS(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client, err := c.Service.Register(userID)
	if err != nil {
		return
	}

	go c.listenWrite(conn, client)
	go c.listenRead(conn, client)
}

func (c *ChatServer) listenRead(conn *websocket.Conn, client *chat.Client) {
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

		client.HandleAction(action)
	}
}

func (c *ChatServer) listenWrite(conn *websocket.Conn, client *chat.Client) {
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
				w.Write([]byte{'\n'})
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
