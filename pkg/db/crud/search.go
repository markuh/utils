package crud

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"

	"github.com/markuh/utils/pkg/apperrors"
	"github.com/markuh/utils/pkg/db"
)

// Search default CRUD method for interface implementation.
func (r *crudRepository[T, SearchFilter, GetFilter]) Search(ctx context.Context, filter SearchFilter) ([]*T, error) {
	filters := r.searchFilter(filter)

	query, args, err := r.searchQuery().
		Where(goqu.And(filters...)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, apperrors.Wrap(err, "failed to build search table query")
	}

	rows, err := r.db.GetDb(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, "failed to search table")
	}
	defer rows.Close()

	result, err := db.ScanRows(rows, r.scanRow)
	if err != nil {
		return nil, apperrors.Wrap(err, "failed to scan table items")
	}

	return result, nil
}

// SearchWithTotal default CRUD method for interface implementation.
func (r *crudRepository[T, SearchFilter, GetFilter]) SearchWithTotal(ctx context.Context, filter SearchFilter) ([]*T, int64, error) {

	filters := r.searchFilter(filter)

	query, args, err := r.searchQuery().
		Select(goqu.COUNT("*")).
		Where(filters...).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, 0, apperrors.Wrap(err, "failed to build query")
	}

	var total int64
	if err := r.db.GetDb(ctx).QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return nil, 0, apperrors.Wrap(err, "failed to get table total")
	}

	result, err := r.Search(ctx, filter)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, "failed to search table")
	}

	return result, total, nil
}

// Get default CRUD method for interface implementation.
func (r *crudRepository[T, SearchFilter, GetFilter]) Get(ctx context.Context, filter GetFilter) (*T, error) {

	filters := r.getFilter(filter)

	query, args, err := r.getQuery().
		Where(goqu.And(filters...)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, apperrors.Wrap(err, "failed to build query")
	}

	row := r.db.GetDb(ctx).QueryRow(ctx, query, args...)
	item, err := r.scanRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, apperrors.Wrap(err, "failed to get table item")
	}

	return item, nil
}
