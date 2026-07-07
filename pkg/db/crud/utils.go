package crud

import (
	"github.com/doug-martin/goqu/v9"

	"github.com/markuh/utils/pkg/db"
)

func (r *crudRepository[T, SearchFilter, GetFilter]) defaultSelectQuery() *goqu.SelectDataset {
	return db.PgDialect.
		Select(r.table.SelectFields...).
		From(r.table.Name).
		Prepared(true)
}
