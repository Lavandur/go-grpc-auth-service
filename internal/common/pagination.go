package common

import "github.com/doug-martin/goqu/v9"

type Pagination struct {
	OrderBy *string `json:"orderBy"`
	Offset  *uint   `json:"offset"`
	Size    *uint   `json:"size"`
}

func GetPagination(
	query *goqu.SelectDataset,
	pagination *Pagination,
) *goqu.SelectDataset {

	if pagination != nil {
		if pagination.Offset != nil {
			query = query.Offset(*pagination.Offset)
		}
		if pagination.Size != nil {
			query = query.Limit(*pagination.Size)
		}
		if pagination.OrderBy != nil {
			query = query.Order(goqu.I(*pagination.OrderBy).Asc())
		}
	}

	return query
}
