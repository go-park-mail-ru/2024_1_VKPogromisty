package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"socio/pkg/contextlogger"
	"socio/pkg/requestcontext"
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

func (c *ChatPubSub) ReadActions(ctx context.Context, userID uint, ch chan *chat.Action) (err error) {
	conn := c.pool.Get()
	defer func() {
		err = conn.Close()
		if err != nil {
			return
		}
	}()

	psc := redis.PubSubConn{Conn: conn}

	contextlogger.LogRedisAction(ctx, "SUBSCRIBE", "userID", userID)

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

			contextlogger.LogRedisAction(ctx, "RECEIVE", action.Receiver, action)

			ch <- action
		}
	}
}

func (c *ChatPubSub) WriteAction(ctx context.Context, action *chat.Action) (err error) {
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

	contextlogger.LogRedisAction(ctx, "PUBLISH", action.Receiver, action)

	_, err = conn.Do("PUBLISH", action.Receiver, data)
	if err != nil {
		return
	}

	senderID, err := requestcontext.GetUserID(ctx)
	if err != nil {
		return
	}

	contextlogger.LogRedisAction(ctx, "PUBLISH", senderID, action)

	_, err = conn.Do("PUBLISH", senderID, data)
	if err != nil {
		return
	}

	return
}
