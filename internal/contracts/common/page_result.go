package common

type PageResult[T any] struct {
	Data []T `json:"data"`
	Total int `json:"total"`
	Page int `json:"page"`
	PageSize int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

