package repository

import (
	"context"
	"fmt"
	"socio/domain"
	"strconv"
	"strings"

	"github.com/gomodule/redigo/redis"
)

func GetUnsentMessageAttachmentKey(attach *domain.UnsentMessageAttachment) string {
	return "unsent_msg_attach_" + fmt.Sprint(attach.SenderID) + ":" + fmt.Sprint(attach.ReceiverID)
}

func ParseUnsentMessageAttachmentKey(key string) (senderID, receiverID uint) {
	ids := strings.Split(strings.TrimPrefix(key, "unsent_msg_attach_"), ":")
	senderIDData, err := strconv.ParseUint(ids[0], 10, 64)
	if err != nil {
		return
	}

	receiverIDData, err := strconv.ParseUint(ids[1], 10, 64)
	if err != nil {
		return
	}

	senderID = uint(senderIDData)
	receiverID = uint(receiverIDData)

	return
}

type UnsentMessageAttachments struct {
	pool *redis.Pool
}

func NewUnsentMessageAttachments(pool *redis.Pool) (unsetMessageAttachments *UnsentMessageAttachments) {
	return &UnsentMessageAttachments{
		pool: pool,
	}
}

func (u *UnsentMessageAttachments) Store(ctx context.Context, attach *domain.UnsentMessageAttachment) (err error) {
	c := u.pool.Get()
	defer c.Close()

	_, err = c.Do("RPUSH", GetUnsentMessageAttachmentKey(attach), attach.FileName)
	if err != nil {
		return
	}

	return
}

func (u *UnsentMessageAttachments) GetAll(ctx context.Context, attach *domain.UnsentMessageAttachment) (fileNames []string, err error) {
	c := u.pool.Get()
	defer c.Close()

	fileNames, err = redis.Strings(c.Do("LRANGE", GetUnsentMessageAttachmentKey(attach), 0, -1))
	if err != nil {
		return
	}

	return
}

func (u *UnsentMessageAttachments) DeleteAll(ctx context.Context, attach *domain.UnsentMessageAttachment) (err error) {
	c := u.pool.Get()
	defer c.Close()

	_, err = c.Do("DEL", GetUnsentMessageAttachmentKey(attach))
	if err != nil {
		return
	}

	return
}

func (u *UnsentMessageAttachments) Delete(ctx context.Context, attach *domain.UnsentMessageAttachment) (err error) {
	c := u.pool.Get()
	defer c.Close()

	_, err = c.Do("LREM", GetUnsentMessageAttachmentKey(attach), 0, attach.FileName)
	if err != nil {
		return
	}

	return
}
