package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxScanner interface {
	Scan(dest ...any) error
}

// IQuery interface for make db queires
type IQuery interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// TxRepository interface for extends repository interfaces
type TxRepository interface {
	WithTx(ctx context.Context, handler func(ctx context.Context) error) error
	GetDb(ctx context.Context) IQuery
	GetNative() *pgxpool.Pool
	Lock(ctx context.Context, table, code string) (bool, error)
	Unlock(ctx context.Context, table, code string) (bool, error)
}
