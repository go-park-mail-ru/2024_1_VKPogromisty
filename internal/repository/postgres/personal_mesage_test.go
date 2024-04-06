package repository_test

import (
	repository "socio/internal/repository/postgres"
	customtime "socio/pkg/time"
	"testing"

	"github.com/chrisyxlee/pgxpoolmock"
)

func TestGetLastMessage(t *testing.T) {
	pool := pgxpoolmock.NewMockPgxIface()
	repo := repository.NewPersonalMessages(pool, customtime.MockTimeProvider{})
	repo.DeleteMessage(nil, 2)
}
