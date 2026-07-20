package crud

import (
	"context"

	"github.com/markuh/utils/pkg/apperrors"
	"github.com/markuh/utils/pkg/db"
)

// Create default CRUD method for interface implementation.
func (r *crudRepository[T, SearchFilter, GetFilter]) Create(ctx context.Context, item *T) (*T, error) {
	return r.CreateWithParams(ctx, item, r.table.InsertValues(item), r.postSaveFunc)
}

// CreateWithValues creates a new item with the given values.
func (r *crudRepository[T, SearchFilter, GetFilter]) CreateWithValues(ctx context.Context, item *T, vals [][]any) (*T, error) {
	return r.CreateWithParams(ctx, item, vals, r.postSaveFunc)
}

// CreateWithParams creates a new item with the given values and post save function.
func (r *crudRepository[T, SearchFilter, GetFilter]) CreateWithParams(
	ctx context.Context,
	item *T,
	vals [][]any,
	postSaveFunc func(ctx context.Context, item *T) error,
) (*T, error) {
	var newItem *T

	if err := r.db.WithTx(ctx, func(ctx context.Context) error {
		sql, args, err := db.PgDialect.
			Insert(r.table.Name).
			Cols(r.table.InsertFields...).
			Vals(vals...).
			Returning(r.table.SelectFields...).
			Prepared(true).
			ToSQL()
		if err != nil {
			return apperrors.Wrap(err, "failed to build insert item query")
		}

		newItem, err = r.scanRow(r.db.GetDb(ctx).QueryRow(ctx, sql, args...))
		if err != nil {
			return apperrors.Wrap(err, "failed to create item")
		}

		if r.postSaveFunc != nil {
			return r.postSaveFunc(ctx, newItem)
		}

		return nil

	}); err != nil {
		return nil, apperrors.Wrap(err, "failed to create new item")
	}

	return newItem, nil
}
