package model

import (
	"math"
)

type Pagination struct {
	Page    uint64 `json:"page" form:"page,default=1"`
	PerPage uint64 `json:"perpage"`
	Total   uint64 `json:"total"`
}

func (p *Pagination) IsValid() bool {
	max := math.Ceil(float64(p.Total) / float64(p.PerPage))
	return float64(p.Page) <= max
}

func NewPagination() Pagination {
	return Pagination{
		PerPage: 10,
	}
}
