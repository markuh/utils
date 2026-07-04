package db

import (
	"database/sql"
	"errors"

	"github.com/markuh/utils/pkg/apperrors"
)

type PgxScanner interface {
	Scan(dest ...any) error
}

func ScanRows[T any](rows *sql.Rows, scanner func(row PgxScanner) (*T, error)) ([]*T, error) {
	result := make([]*T, 0)
	for rows.Next() {
		item, err := scanner(rows)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, apperrors.Wrap(err, "can't scan row")
		}
		if item != nil {
			result = append(result, item)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, "can't get rows error")
	}

	return result, nil
}
