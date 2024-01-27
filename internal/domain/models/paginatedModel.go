package models

type PaginatedModel[T any] struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
	Data   []T   `json:"data"`
}
