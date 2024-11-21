package dto

type EmailFilterEntityDto struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	FilterDetails []EmailFilterDetailDto `json:"filter_details"`
	CreatedTime string `json:"created_time"`
}
type EmailFilterDetailDto struct {
	Id int64 `json:"id"`
	Content string `json:"content"`

}
