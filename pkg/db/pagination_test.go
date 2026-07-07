package db

import (
	"strings"
	"testing"

	"github.com/doug-martin/goqu/v9"
)

func TestSortingApplyNullsLastForNullableDateFields(t *testing.T) {
	dialect := goqu.Dialect("postgres")
	table := TableConfig[any]{
		SortingFields: map[string]bool{
			"started_at": true,
		},
	}

	q := dialect.From("tasks").Select("id")
	q = Sorting{Sort: "started_at", SortOrder: SortOrderDesc}.Apply(table, q)

	sql, _, err := q.Prepared(true).ToSQL()
	if err != nil {
		t.Fatalf("ToSQL: %v", err)
	}
	if !strings.Contains(sql, "NULLS LAST") {
		t.Fatalf("expected NULLS LAST in SQL, got: %s", sql)
	}
	if !strings.Contains(sql, "started_at") || !strings.Contains(strings.ToUpper(sql), "DESC") {
		t.Fatalf("expected DESC sort by started_at, got: %s", sql)
	}
}
