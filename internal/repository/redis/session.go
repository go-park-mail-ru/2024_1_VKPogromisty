package repository

import (
	"socio/errors"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type Session struct {
	c redis.Conn
}

func NewSession(c redis.Conn) (s *Session) {
	return &Session{
		c: c,
	}
}

func (s *Session) CreateSession(userID uint) (sessionID string, err error) {
	sessionID = uuid.NewString()
	_, err = s.c.Do("SET", sessionID, userID)
	if err != nil {
		return
	}

	return
}

func (s *Session) DeleteSession(sessionID string) (err error) {
	_, err = s.c.Do("DEL", sessionID)
	if err != nil {
		err = errors.ErrNotFound
	}

	return
}

func (s *Session) GetUserIDBySession(sessionID string) (userID uint, err error) {
	userIDData, err := redis.Uint64(s.c.Do("GET", sessionID))
	if err != nil {
		err = errors.ErrUnauthorized
	}

	userID = uint(userIDData)

	return
}
