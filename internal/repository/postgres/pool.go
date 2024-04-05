package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPool(connStr string) (pool *pgxpool.Pool, err error) {
	pgConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return
	}

	pool, err = pgxpool.ConnectConfig(context.Background(), pgConfig)
	if err != nil {
		return
	}

	return
}
