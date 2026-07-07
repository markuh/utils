package crud

import (
	"context"

	"github.com/doug-martin/goqu/v9"

	"github.com/markuh/utils/pkg/db"
)

type Repository[T any, SearchFilter any, GetFilter any] interface {
	Search(ctx context.Context, filter SearchFilter) ([]*T, error)
	SearchWithTotal(ctx context.Context, filter SearchFilter) ([]*T, int64, error)
	Get(ctx context.Context, filter GetFilter) (*T, error)
	Create(ctx context.Context, item *T) (*T, error)
	Update(ctx context.Context, item *T) (*T, error)
	Delete(ctx context.Context, item *T) error
	SoftDelete(ctx context.Context, item *T) error
	CreateWithParams(ctx context.Context, item *T, vals [][]any, postSaveFunc func(ctx context.Context, item *T) error) (*T, error)
	UpdateWithParams(ctx context.Context, item *T, vals *goqu.Record, updateFilter func(item *T) []goqu.Expression, preUpdateFunc func(ctx context.Context, item *T) error) (*T, error)
}

type crudRepository[T any, SearchFilter any, GetFilter any] struct {
	db           db.TxRepository
	table        db.TableConfig[T]
	scanRow      func(row db.PgxScanner) (*T, error)
	searchQuery  func() *goqu.SelectDataset
	searchFilter func(filter SearchFilter) []goqu.Expression
	getQuery     func() *goqu.SelectDataset
	getFilter    func(filter GetFilter) []goqu.Expression

	postSaveFunc func(ctx context.Context, item *T) error

	updateFilter  func(item *T) []goqu.Expression
	preUpdateFunc func(ctx context.Context, item *T) error

	softDeleteValues func(item *T) *goqu.Record
	softDeleteFilter func(item *T) []goqu.Expression
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
