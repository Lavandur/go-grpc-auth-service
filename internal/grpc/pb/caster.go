package pb

import "auth-service/internal/common"

func (v *Pagination) ToModel() *common.Pagination {
	if v == nil {
		return nil
	}

	var pagination *common.Pagination
	if v.Offset != nil {
		offset := uint(*v.Offset)
		pagination.Offset = &offset
	}
	if v.Size != nil {
		limit := uint(*v.Size)
		pagination.Size = &limit
	}
	if v.OrderBy != nil {
		pagination.OrderBy = v.OrderBy
	}

	return pagination
}
