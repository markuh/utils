package crud

import (
	"context"

	"github.com/doug-martin/goqu/v9"
)

type Repository[T any, SearchFilter any, GetFilter any] interface {
	Search(ctx context.Context, filter SearchFilter) ([]*T, error)
	SearchWithTotal(ctx context.Context, filter SearchFilter) ([]*T, int64, error)
	Get(ctx context.Context, filter GetFilter) (*T, error)
	Create(ctx context.Context, item *T) (*T, error)
	CreateWithValues(ctx context.Context, item *T, vals [][]any) (*T, error)
	CreateWithParams(ctx context.Context, item *T, vals [][]any, postSaveFunc func(ctx context.Context, item *T) error) (*T, error)
	Update(ctx context.Context, item *T) (*T, error)
	UpdateWithValues(ctx context.Context, item *T, vals *goqu.Record) (*T, error)
	UpdateWithParams(ctx context.Context, item *T, vals *goqu.Record, updateFilter func(item *T) []goqu.Expression, preUpdateFunc func(ctx context.Context, item *T) error) (*T, error)
	Delete(ctx context.Context, item *T) error
	SoftDelete(ctx context.Context, item *T) error
}
