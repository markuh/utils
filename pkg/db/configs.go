package db

type TableConfig struct {
	Name          string
	Alias         string
	InsertFields  []any
	SelectFields  []any
	UpdatedFields []any
	SortingFields map[string]bool // true for nulls last
}
