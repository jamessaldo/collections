package dto

import "math"

type Pagination struct {
	Page      int           `json:"page"`
	PageSize  int           `json:"page_size"`
	TotalData int64         `json:"total_data"`
	TotalPage int           `json:"total_page"`
	NextPage  int           `json:"next_page"`
	PrevPage  int           `json:"prev_page"`
	HasNext   bool          `json:"has_next"`
	HasPrev   bool          `json:"has_prev"`
	Data      []interface{} `json:"data"`
}

// create a function to generate paginate data
func Paginate(page, pageSize int, totalData int64, data []interface{}) Pagination {
	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))

	nextPage := -1
	if page < totalPage {
		nextPage = page + 1
	}

	prevPage := -1
	if page > 1 && page <= totalPage {
		prevPage = page - 1
	} else if page > totalPage {
		prevPage = totalPage
	}

	hasNext := page < totalPage
	hasPrev := page > 1
	return Pagination{
		Page:      page,
		PageSize:  pageSize,
		TotalData: totalData,
		TotalPage: totalPage,
		NextPage:  nextPage,
		PrevPage:  prevPage,
		HasNext:   hasNext,
		HasPrev:   hasPrev,
		Data:      data,
	}
}
