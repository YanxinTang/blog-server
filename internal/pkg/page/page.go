package page

import (
	"math"
)

type Pagination struct {
	Page    int `json:"page" form:"page,default=1"`
	PerPage int `json:"perpage"`
	Total   int `json:"total"`
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func (p *Pagination) Limit() int {
	return p.PerPage
}

func (p *Pagination) IsValid() bool {
	max := math.Ceil(float64(p.Total) / float64(p.PerPage))
	return float64(p.Page) <= max
}

func NewPagination() *Pagination {
	return &Pagination{
		PerPage: 10,
	}
}
