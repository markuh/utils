package db

import (
	"strings"

	"github.com/doug-martin/goqu/v9"
)

const (
	SortOrderAsc  = "ASC"
	SortOrderDesc = "DESC"
)

type Pagination struct {
	All     *bool
	Page    int
	PerPage int
}

func (p Pagination) Apply(q *goqu.SelectDataset) *goqu.SelectDataset {
	if p.All != nil && *p.All {
		return q
	}
	if off := p.GetOffset(); off > 0 {
		q = q.Offset(uint(off))
	}
	if limit := p.GetLimit(); limit > 0 {
		q = q.Limit(uint(limit))
	}
	return q
}

func (p Pagination) GetPage() int {
	if p.Page <= 0 {
		return 0
	}
	return p.Page
}

func (p Pagination) GetPerPage() int {
	if p.PerPage <= 0 || p.PerPage > 100 {
		return 100
	}
	return p.PerPage
}

func (p Pagination) GetLimit() int {
	return p.GetPerPage()
}

func (p Pagination) GetOffset() int {
	return p.GetPage() * p.GetPerPage()
}

type Sorting struct {
	Sort      string
	SortOrder string
}

func (s Sorting) Apply(tableConfig TableConfig[any], q *goqu.SelectDataset) *goqu.SelectDataset {
	return s.ApplyColumn(tableConfig, s.Sort, q)
}

func (s Sorting) ApplyColumn(tableConfig TableConfig[any], column string, q *goqu.SelectDataset) *goqu.SelectDataset {
	if s.Sort == "" || column == "" {
		return q
	}

	nullsLast, ok := tableConfig.SortingFields[s.Sort]
	if !ok {
		return q
	}

	orderBy := goqu.I(column).Asc()
	if s.GetSortOrder() == SortOrderDesc {
		orderBy = goqu.I(column).Desc()
	}
	if nullsLast {
		orderBy = orderBy.NullsLast()
	}
	return q.Order(orderBy)
}

func (s Sorting) GetSortOrder() string {
	if strings.ToUpper(s.SortOrder) == SortOrderDesc {
		return SortOrderDesc
	}
	return SortOrderAsc
}

func (s Sorting) GetOrderBy(tableConfig TableConfig[any]) (string, string) {
	sortCol := s.Sort
	if sortCol == "" {
		return tableConfig.DefaultSort, SortOrderAsc
	}
	if _, ok := tableConfig.SortingFields[sortCol]; !ok {
		return tableConfig.DefaultSort, SortOrderAsc
	}
	if s.GetSortOrder() == SortOrderDesc {
		return sortCol, SortOrderDesc
	}
	return sortCol, SortOrderAsc
}
