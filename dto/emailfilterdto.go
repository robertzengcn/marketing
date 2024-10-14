package dto

type EmailFilterEntityDto struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	FilterDetails []EmailFilterDetailDto `json:"filter_details"`
}
type EmailFilterDetailDto struct {
	Id int64 `json:"id"`
	Content string `json:"content"`
}