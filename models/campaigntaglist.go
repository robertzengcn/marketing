package models

import (
	"github.com/beego/beego/v2/client/orm"

)
type CampaignTagList struct {
	Id       int64     `orm:"pk;auto"`
	Campaign *Campaign `orm:"rel(fk)"`
	Tag      *Tag      `orm:"rel(fk)"`
}

func (u *CampaignTagList) TableName() string {
	return "campaign_tag_list"
}

func (u *CampaignTagList) TableEngine() string {
	return "INNODB"
}

func init() {
	// orm.RegisterModelWithPrefix("mk_", new(CampaignTagList))
}

//get campaign tag by campaign id
func (u *CampaignTagList) GetCampaignTagByCampaignId(campaignId int64) (campaignTag []*Tag, err error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("mk_tag.tag_id,mk_tag.tag_name").From("mk_campaign_tag_list").LeftJoin("mk_tag").On("mk_campaign_tag_list.tag_id = mk_tag.tag_id").Where("mk_campaign_tag_list.campaign_id = ?")
	sql := qb.String()
	//l.Println(sql)
	// 执行 SQL 语句
	o := orm.NewOrm()
	_, qerr := o.Raw(sql, campaignId).QueryRows(&campaignTag)
	if qerr != nil {
		return nil, qerr
	}
	return campaignTag, nil
}
