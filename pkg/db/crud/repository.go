package crud

import (
	"context"

	"github.com/doug-martin/goqu/v9"

	"github.com/markuh/utils/pkg/db"
)

type crudRepository[T any, SearchFilter any, GetFilter any] struct {
	db           db.TxRepository
	table        db.TableConfig[T]
	scanRow      func(row db.PgxScanner) (*T, error)
	searchQuery  func() *goqu.SelectDataset
	searchFilter func(filter SearchFilter) []goqu.Expression
	getQuery     func() *goqu.SelectDataset
	getFilter    func(filter GetFilter) []goqu.Expression

	preSaveFunc  func(ctx context.Context, item *T) (*T, error)
	postSaveFunc func(ctx context.Context, item *T) error

	updateFilter   func(item *T) []goqu.Expression
	preUpdateFunc  func(ctx context.Context, item *T) error
	postUpdateFunc func(ctx context.Context, item *T) (*T, error)

	softDeleteValues func(item *T) *goqu.Record
	softDeleteFilter func(item *T) []goqu.Expression

	preDeleteFunc  func(ctx context.Context, item *T) (*T, error)
	postDeleteFunc func(ctx context.Context, item *T) (*T, error)
}

func NewCrudRepository[T any, SearchFilter any, GetFilter any](
	db db.TxRepository,
	table db.TableConfig[T],
	scanRow func(row db.PgxScanner) (*T, error),
) *crudRepository[T, SearchFilter, GetFilter] {

	r := &crudRepository[T, SearchFilter, GetFilter]{
		db:      db,
		table:   table,
		scanRow: scanRow,
	}
	r.searchQuery = r.defaultSelectQuery
	r.getQuery = r.defaultSelectQuery
	return r
}
