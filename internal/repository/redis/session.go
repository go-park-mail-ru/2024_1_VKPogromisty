package repository

import (
	"context"
	"socio/errors"
	"socio/pkg/contextlogger"

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

func (s *Session) CreateSession(ctx context.Context, userID uint) (sessionID string, err error) {
	c := s.pool.Get()
	defer c.Close()

	sessionID = uuid.NewString()

	contextlogger.LogRedisAction(ctx, "SET", "SESSION_ID", userID)

	_, err = c.Do("SET", sessionID, userID)
	if err != nil {
		return
	}

	return
}

func (s *Session) DeleteSession(ctx context.Context, sessionID string) (err error) {
	c := s.pool.Get()
	defer c.Close()

	contextlogger.LogRedisAction(ctx, "DEL", "SESSION_ID", nil)

	_, err = c.Do("DEL", sessionID)
	if err != nil {
		err = errors.ErrNotFound
	}

	return
}

func (s *Session) GetUserIDBySession(ctx context.Context, sessionID string) (userID uint, err error) {
	c := s.pool.Get()
	defer c.Close()

	contextlogger.LogRedisAction(ctx, "GET", "SESSION_ID", nil)

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
