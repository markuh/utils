package crud

import (
	"context"
	"database/sql"
	"errors"

	"github.com/markuh/utils/pkg/apperrors"
)

// Delete default CRUD method for interface implementation.
func (r *crudRepository[T, SearchFilter, GetFilter]) Delete(ctx context.Context, item *T) error {
	// if err := r.db.WithTx(ctx, func(ctx context.Context) error {
	// 	// return r.db.Delete(ctx, item)
	// }); err != nil {
	// 	return apperrors.Wrap(err, "failed to delete item")
	// }

	return apperrors.New("not implemented")
}

// SoftDelete default CRUD method for interface implementation.
func (r *crudRepository[T, SearchFilter, GetFilter]) SoftDelete(ctx context.Context, item *T) error {
	deletedItem, err := r.UpdateWithParams(ctx, item, r.softDeleteValues(item), r.softDeleteFilter, nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return apperrors.Wrap(err, "failed to update item to soft delete")
	}

	if deletedItem == nil {
		return apperrors.New("failed to get soft deleted item")
	}

	return nil
}
