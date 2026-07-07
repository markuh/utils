package db

import "github.com/doug-martin/goqu/v9"

type TableConfig[T any] struct {
	Name          string
	Alias         string
	InsertFields  []any
	SelectFields  []any
	UpdatedFields []any
	DefaultSort   string
	SortingFields map[string]bool // true for nulls last
	InsertValues  func(item *T) [][]any
	UpdateValues  func(item *T) *goqu.Record
}
