package dto
// import(
// 	"time"
// )
type ProxyDto struct {
	Id int64 `json:"id"`
	Host        string    `json:"host"`
	Port        string    `json:"port"`
	User        string    `json:"user"`
	Pass        string    `json:"pass"`
	Protocol    string    `json:"protocol"`
	CountryCode string    `json:"country_code"`
	Addtime     string    `json:"addtime"`
}

type ProxyRespDto struct {
	Total int64 `json:"total"`
	Records  []ProxyDto `json:"records"`
}

// type ProxyDetailDto struct {
// 	Id int64 `json:"id"`
// 	Host        string    `json:"host"`
// 	Port        string    `json:"port"`
// 	User        string    `json:"user"`
// 	Pass        string    `json:"pass"`
// 	Protocol    string    `json:"protocol"`
// 	CountryCode string    `json:"country_code"`
// }
type SaveProxyDto struct{
	Id int64 `json:"id"`
}
type ImportProxyDto struct {
	
	Host        string    `json:"host"`
	Port        string    `json:"port"`
	User        string    `json:"user"`
	Pass        string    `json:"pass"`
	Protocol    string    `json:"protocol"`
}
