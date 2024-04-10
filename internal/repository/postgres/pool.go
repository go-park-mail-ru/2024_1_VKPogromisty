package repository

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBPool interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Close()
}

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
