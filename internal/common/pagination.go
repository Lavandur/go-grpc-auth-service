package common

type Pagination struct {
	OrderBy string `json:"orderBy"`
	Offset  int    `json:"offset"`
	Size    int    `json:"size"`
}
