package repository

import (
	"bytes"
	"encoding/json"
	"socio/usecase/chat"

	"github.com/gomodule/redigo/redis"
)

type ChatPubSub struct {
	pool *redis.Pool
}

func NewChatPubSub(pool *redis.Pool) (chatPubSub *ChatPubSub) {
	return &ChatPubSub{
		pool: pool,
	}
}

func (c *ChatPubSub) ReadActions(userID uint, ch chan *chat.Action) (err error) {
	conn := c.pool.Get()
	defer func() {
		err = conn.Close()
		if err != nil {
			return
		}
	}()

	psc := redis.PubSubConn{Conn: conn}

	err = psc.Subscribe(userID)
	if err != nil {
		return
	}

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			action := new(chat.Action)
			err = json.NewDecoder(bytes.NewReader(v.Data)).Decode(action)
			if err != nil {
				return
			}

			ch <- action
		}
	}
}

func (c *ChatPubSub) WriteAction(action *chat.Action) (err error) {
	conn := c.pool.Get()
	defer func() {
		err = conn.Close()
		if err != nil {
			return
		}
	}()

	data, err := json.Marshal(action)
	if err != nil {
		return
	}

	_, err = conn.Do("PUBLISH", action.Receiver, data)
	if err != nil {
		return
	}

	return
}
