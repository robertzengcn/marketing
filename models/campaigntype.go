package models

import (
	"github.com/beego/beego/v2/client/orm"
)

//defined CampaignType struct
type CampaignType struct {
	CampaignTypeId int32 `orm:"pk;auto"`
	CampaignTypeName string `orm:"size(100)"`
}

func (u *CampaignType) TableName() string {
	return "campaign_type"
}

func (u *CampaignType) TableEngine() string {
	return "INNODB"
}

func init() {
	orm.RegisterModelWithPrefix("mk_", new(CampaignType))
}