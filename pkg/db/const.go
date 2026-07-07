package db

import (
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

var PgDialect = goqu.Dialect("postgres")
