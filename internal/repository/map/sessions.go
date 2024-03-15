package repository

import (
	"socio/errors"
	"sync"

	"github.com/google/uuid"
)

type Sessions struct {
	Sessions *sync.Map
}

func NewSessions(sessions *sync.Map) (sessionsStorage *Sessions) {
	sessionsStorage = &Sessions{}
	sessionsStorage.Sessions = sessions

	return
}

func (s *Sessions) CreateSession(userID uint) (sessionID string) {
	sessionID = uuid.NewString()
	s.Sessions.Store(sessionID, userID)

	return
}

func (s *Sessions) GetUserIDBySession(sessionID string) (userID uint, err error) {
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

func (s *Sessions) DeleteSession(sessionID string) (err error) {
	_, ok := s.Sessions.LoadAndDelete(sessionID)
	if !ok {
		err = errors.ErrNotFound
		return
	}

	return
}
