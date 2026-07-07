package crud

import (
	"context"

	"github.com/doug-martin/goqu/v9"

	"github.com/markuh/utils/pkg/db"
)

func (r *crudRepository[T, SearchFilter, GetFilter]) WithDb(dbTx db.TxRepository) *crudRepository[T, SearchFilter, GetFilter] {
	r.db = dbTx
	return r
}

func (r *crudRepository[T, SearchFilter, GetFilter]) WithSearchQuery(query func() *goqu.SelectDataset) *crudRepository[T, SearchFilter, GetFilter] {
	r.searchQuery = query
	if r.getQuery == nil {
		r.getQuery = r.defaultSelectQuery
	}
	return r
}

func (r *crudRepository[T, SearchFilter, GetFilter]) WithGetQuery(query func() *goqu.SelectDataset) *crudRepository[T, SearchFilter, GetFilter] {
	r.getQuery = query
	return r
}

func (r *crudRepository[T, SearchFilter, GetFilter]) WithSearchFilter(filter func(filter SearchFilter) []goqu.Expression) *crudRepository[T, SearchFilter, GetFilter] {
	r.searchFilter = filter
	return r
}

func (r *crudRepository[T, SearchFilter, GetFilter]) WithGetFilter(filter func(filter GetFilter) []goqu.Expression) *crudRepository[T, SearchFilter, GetFilter] {
	r.getFilter = filter
	return r
}

// func (r *crudRepository[T, SearchFilter, GetFilter]) WithInsertValues(values func(item *T) [][]any) *crudRepository[T, SearchFilter, GetFilter] {
// 	r.insertValues = values
// 	return r
// }

func (r *crudRepository[T, SearchFilter, GetFilter]) WithPostSaveFunc(postSaveFunc func(ctx context.Context, item *T) error) *crudRepository[T, SearchFilter, GetFilter] {
	r.postSaveFunc = postSaveFunc
	return r
}

// func (r *crudRepository[T, SearchFilter, GetFilter]) WithUpdateValues(values func(item *T) *goqu.Record) *crudRepository[T, SearchFilter, GetFilter] {
// 	r.updateValues = values
// 	return r
// }

func (r *crudRepository[T, SearchFilter, GetFilter]) WithUpdateFilter(filter func(item *T) []goqu.Expression) *crudRepository[T, SearchFilter, GetFilter] {
	r.updateFilter = filter
	return r
}

func (r *crudRepository[T, SearchFilter, GetFilter]) WithPreUpdateFunc(preUpdateFunc func(ctx context.Context, item *T) error) *crudRepository[T, SearchFilter, GetFilter] {
	r.preUpdateFunc = preUpdateFunc
	return r
}

func (r *crudRepository[T, SearchFilter, GetFilter]) WithSoftDeleteValues(values func(item *T) *goqu.Record) *crudRepository[T, SearchFilter, GetFilter] {
	r.softDeleteValues = values
	return r
}

func (r *crudRepository[T, SearchFilter, GetFilter]) WithSoftDeleteFilter(filter func(item *T) []goqu.Expression) *crudRepository[T, SearchFilter, GetFilter] {
	r.softDeleteFilter = filter
	return r
}
