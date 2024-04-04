package repository

import (
	"context"
	"socio/errors"
	"sync"

	"github.com/google/uuid"
)

type Sessions struct {
	Sessions *sync.Map
}

func NewSessions(sessions *sync.Map) (sessionsStorage *Sessions, err error) {
	sessionsStorage = &Sessions{}
	sessionsStorage.Sessions = sessions

	return
}

func (s *Sessions) CreateSession(ctx context.Context, userID uint) (sessionID string, err error) {
	sessionID = uuid.NewString()
	s.Sessions.Store(sessionID, userID)

	return
}

func (s *Sessions) GetUserIDBySession(ctx context.Context, sessionID string) (userID uint, err error) {
	userIDData, ok := s.Sessions.Load(sessionID)
	if !ok {
		err = errors.ErrUnauthorized
		return
	}

	userID, ok = userIDData.(uint)
	if !ok {
		err = errors.ErrInternal
		return
	}

	return
}

func (s *Sessions) DeleteSession(ctx context.Context, sessionID string) (err error) {
	_, ok := s.Sessions.LoadAndDelete(sessionID)
	if !ok {
		err = errors.ErrNotFound
		return
	}

	return
}
