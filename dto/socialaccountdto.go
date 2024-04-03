package dto

type SocialAccountDto struct {
	Id int64 `json:"id"`
	SocialType string `json:"social_type"`
	SocialTypeId int64 `json:"social_type_id"`
	User string `json:"user"`
	Password string `json:"pass"`
	Status int8 `json:"status"`
	// Proxy ProxyDto `json:"proxy"`
}
type SocialProxyDto struct {
	Id int64 `json:"id"`
	Url string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SocialAccountDetail struct {
	Id int64 `json:"id"`
	SocialType string `json:"social_type"`
	SocialTypeId int64 `json:"social_type_id"`
	SocialTypeUrl string `json:"social_type_url"`
	User   string  `json:"user"`
	Pass   string  `json:"pass"`
	Status int8 `json:"status"`
	Name string `json:"name"`
	PhoneNumber string `json:"phone"`
	Email string `json:"email"`
	Proxy  SocialProxyDto `json:"proxy"`
}
type SocialAccountListresp struct {
	Total int64 `json:"total"`
	Records  []SocialAccountDto `json:"records"`
}
type SocialAccountSaveresp struct {
	Id int64 `json:"id"`
}