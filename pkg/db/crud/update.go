package crud

import (
	"context"

	"github.com/doug-martin/goqu/v9"

	"github.com/markuh/utils/pkg/apperrors"
	"github.com/markuh/utils/pkg/db"
)

// Update default CRUD method for interface implementation.
func (r *crudRepository[T, SearchFilter, GetFilter]) Update(ctx context.Context, item *T) (*T, error) {
	return r.UpdateWithParams(ctx, item, r.table.UpdateValues(item), r.updateFilter, r.preUpdateFunc)
}

// UpdateWithParams updates an item with the given values and update filter.
func (r *crudRepository[T, SearchFilter, GetFilter]) UpdateWithParams(
	ctx context.Context,
	item *T,
	vals *goqu.Record,
	updateFilter func(item *T) []goqu.Expression,
	preUpdateFunc func(ctx context.Context, item *T) error,
) (*T, error) {
	var updatedItem *T

	if err := r.db.WithTx(ctx, func(ctx context.Context) error {
		if preUpdateFunc != nil {
			if err := preUpdateFunc(ctx, item); err != nil {
				return apperrors.Wrap(err, "failed to pre update item")
			}
		}

		q := db.PgDialect.
			Update(r.table.Name).
			Set(vals).
			Returning(r.table.SelectFields...)

		if updateFilter != nil {
			filters := updateFilter(item)
			if len(filters) > 0 {
				q = q.Where(goqu.And(filters...))
			}
		}

		sql, args, err := q.Prepared(true).ToSQL()
		if err != nil {
			return apperrors.Wrap(err, "failed to build update item query")
		}

		updatedItem, err = r.scanRow(r.db.GetDb(ctx).QueryRow(ctx, sql, args...))
		if err != nil {
			return apperrors.Wrap(err, "failed to update item")
		}

		return nil
	}); err != nil {
		return nil, apperrors.Wrap(err, "failed to update item")
	}

	return updatedItem, nil
}
