package repository

import (
	"socio/errors"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type Session struct {
	pool *redis.Pool
}

func NewSession(pool *redis.Pool) (s *Session) {
	return &Session{
		pool: pool,
	}
}

func (s *Session) CreateSession(userID uint) (sessionID string, err error) {
	c := s.pool.Get()
	defer c.Close()

	sessionID = uuid.NewString()
	_, err = c.Do("SET", sessionID, userID)
	if err != nil {
		return
	}

	return
}

func (s *Session) DeleteSession(sessionID string) (err error) {
	c := s.pool.Get()
	defer c.Close()

	_, err = c.Do("DEL", sessionID)
	if err != nil {
		err = errors.ErrNotFound
	}

	return
}

func (s *Session) GetUserIDBySession(sessionID string) (userID uint, err error) {
	c := s.pool.Get()
	defer c.Close()

	result, err := c.Do("GET", sessionID)
	if err != nil {
		err = errors.ErrUnauthorized
		return
	}

	userIDData, err := redis.Uint64(result, err)
	if err != nil {
		err = errors.ErrUnauthorized
		return
	}

	userID = uint(userIDData)

	return
}
