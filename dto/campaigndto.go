package dto


type CampaignDto struct {
	Records []CampaignItemDto `json:"records"`
	Num	 int64             `json:"num"`
}

type CampaignItemDto struct {
	CampaignId      int64     `orm:"pk;auto"`
	CampaignName    string    `orm:"size(100)"`
	CampaignDescription string `orm:"type(text);null"`
	Tags []string `orm:"type(text);null"` //the tag use to fetch keyword
	Types string  `orm:"null;rel(fk)"` //the type of campaign, email, social
	Disable int `orm:"default(0)"` //0: disabled, 1: enabled
	AccountId int64 `orm:"rel(fk);on_delete(do_nothing);column(account_id)"`
}

