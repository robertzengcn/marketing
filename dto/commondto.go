package dto

type CommonResponse[T any] struct {
	Record T   `json:"records"`
	Total int64 `json:"num"`
}
type IdResponse struct {
	Id int64   `json:"id"`
}