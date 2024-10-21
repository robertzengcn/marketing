package dto

type EmailServiceEntityDto struct {
	Id int64 `json:"id"`
	From     string `json:"from"`
	Password string    `json:"password"`
	Host     string   `json:"host"`
	Port     string     `json:"port"`
	Name    string     `json:"name"`
	Ssl    int8     `json:"ssl"`
}