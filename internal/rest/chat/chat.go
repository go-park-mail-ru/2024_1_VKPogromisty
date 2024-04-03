package rest

import (
	"encoding/json"
	"net/http"
	"socio/internal/rest/middleware"
	"socio/usecase/chat"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	writeWait      = 10 * time.Second
	maxMessageSize = 10000
)

// ChatServer will: listen for ws connection messages, register and unregister clients based on state of ws connections
type ChatServer struct {
	Service *chat.Service
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     middleware.CheckOrigin, // FIX THIS
}

func NewChatServer(pubSubRepo chat.PubSubRepository, messagesRepo chat.PersonalMessagesRepository) (chatServer *ChatServer) {
	return &ChatServer{
		Service: chat.NewChatService(pubSubRepo, messagesRepo),
	}
}

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

	go c.ListenWrite(conn, client)
	go c.ListenRead(conn, client)
}

func (c *ChatServer) ListenRead(conn *websocket.Conn, client *chat.Client) {
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
		err = json.Unmarshal(jsonMessage, action)
		if err != nil {
			return
		}

		client.HandleAction(action)
	}
}

func (c *ChatServer) ListenWrite(conn *websocket.Conn, client *chat.Client) {
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

			messageData, err := json.Marshal(message)
			if err != nil {
				return
			}

			w.Write(messageData)

			n := len(client.Send)
			for i := 0; i < n; i++ {
				messageData, err = json.Marshal(<-client.Send)
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
