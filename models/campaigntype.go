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

//get campaign type by id
func (u *CampaignType) GetCampaignTypeById(id int32) (CampaignType, error) {
	o := orm.NewOrm()
	var camType CampaignType
	err := o.QueryTable(new(CampaignType)).Filter("campaign_type_id", id).One(&camType)
	return camType, err
}