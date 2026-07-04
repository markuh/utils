package db

import (
	"context"
	"database/sql"
)

type PgxScanner interface {
	Scan(dest ...any) error
}

// IQuery interface for make db queires
type IQuery interface {
	Exec(sql string, arguments ...any) (sql.Result, error)
	ExecContext(ctx context.Context, sql string, arguments ...any) (sql.Result, error)
	Query(sql string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error)
	QueryRow(sql string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, sql string, args ...any) *sql.Row
	PrepareContext(ctx context.Context, sql string) (*sql.Stmt, error)
}

// TxRepository use for extends repository interfaces
type TxRepository interface {
	WithTx(ctx context.Context, handler func(ctx context.Context) error) error
	GetDb(ctx context.Context) IQuery
	GetNative() *sql.DB
	Lock(ctx context.Context, table, code string) (bool, error)
	Unlock(ctx context.Context, table, code string) (bool, error)
}
