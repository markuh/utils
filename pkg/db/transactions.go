package db

import (
	"context"
	"hash/fnv"

	"database/sql"

	"github.com/markuh/utils/pkg/apperrors"
)

type txKey string

const (
	txCtxKey txKey = "txCtxKey"
)

// NewTxRepository wraps a repository for transactional execution
func NewTxRepository(db *sql.DB) TxRepository {
	return &txRepository{
		db: db,
	}
}

// txRepository
type txRepository struct {
	db *sql.DB
}

// getTx return tx and flag of it exists
func (t *txRepository) getTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txCtxKey).(*sql.Tx)
	return tx, ok
}

// WithTx create transaction and put in to context
func (t *txRepository) WithTx(ctx context.Context, handler func(ctx context.Context) error) error {

	var err error
	tx, txExists := t.getTx(ctx)

	if !txExists {
		// init TX
		tx, err = t.db.BeginTx(ctx, nil)
		if err != nil {
			return apperrors.Wrap(err, "can't init transaction")
		}
		defer func(ctx context.Context) { _ = tx.Rollback() }(ctx)

		// set TX to context
		ctx = context.WithValue(ctx, txCtxKey, tx)
	}

	if err := handler(ctx); err != nil {
		return apperrors.Wrap(err, "exec transaction error")
	}

	if !txExists {
		// finish only created transaction
		ctx = context.WithValue(ctx, txCtxKey, nil)

		// commit TX
		if err := tx.Commit(); err != nil {
			return apperrors.Wrap(err, "can't commit transaction")
		}
	}

	return nil
}

// GetDb return db connection for in tx queries
func (t *txRepository) GetDb(ctx context.Context) IQuery {
	tx, ok := ctx.Value(txCtxKey).(*sql.Tx)
	if !ok {
		return t.db
	}
	return tx
}

// GetNative return native db connection
func (t *txRepository) GetNative() *sql.DB {
	return t.db
}

// Lock logical lock
func (t *txRepository) Lock(ctx context.Context, table, code string) (bool, error) {
	var result bool
	tableHash := getInt16Hash(table)
	codeHash := getInt16Hash(code)

	query := `SELECT pg_try_advisory_lock($1, $2);`
	if err := t.GetDb(ctx).QueryRowContext(ctx, query, tableHash, codeHash).Scan(&result); err != nil {
		return result, apperrors.Wrap(err, "can't get lock")
	}
	return result, nil
}

// Unlock free logical lock
func (t *txRepository) Unlock(ctx context.Context, table, code string) (bool, error) {
	var result bool
	tableHash := getInt16Hash(table)
	codeHash := getInt16Hash(code)

	query := `SELECT pg_advisory_unlock($1, $2);`
	if err := t.GetDb(ctx).QueryRowContext(ctx, query, tableHash, codeHash).Scan(&result); err != nil {
		return result, apperrors.Wrap(err, "can't unlock")
	}
	return result, nil
}

func getInt16Hash(s string) int16 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int16(h.Sum32())
}
