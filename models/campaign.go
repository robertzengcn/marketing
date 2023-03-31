package models

import (
	// "fmt"
	 "errors"
	// "fmt"
	  "github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	// "time"
)
var DefaultCampaign *Campaign
type Campaign struct {
	CampaignId      int64     `orm:"pk;auto"`
	CampaignName    string    `orm:"size(100)"`
	// EmailTpl   *EmailTpl	`orm:"rel(fk);on_delete(do_nothing)"`
	Tags string `orm:"type(text);null"` //the tag use to fetch keyword
}

func (u *Campaign) TableName() string {
	return "campaign"
}

// 设置引擎为 INNODB
func (u *Campaign) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(Campaign))
	// create table
	// orm.RunSyncdb("default", false, true)
}

///create campaign with name
func (u *Campaign)CreateCampaign(username string,tags string) (id int64, err error) {
	o := orm.NewOrm()
	var us Campaign
	us.CampaignName = username
	us.Tags=tags
	id, err = o.Insert(&us)
		return id,err
}
/// show all campaign
func (u *Campaign)ListCampaign(start int,limitNum int)([]Campaign,error){
	o := orm.NewOrm()
	var cam []Campaign
	count, e := o.QueryTable(new(Campaign)).Limit(limitNum,start).All(&cam, "campaign_id", "campaign_name","tags")
	if e != nil {
		return nil, e
	}
	if count <= 0 {
		return nil, errors.New("nothing found")
	}
	return cam, nil
}
/// find campaign by campaign id
func (u *Campaign)FindCambyid(id int64)(*Campaign,error){
	o := orm.NewOrm()
	campaign := Campaign{CampaignId: id}
	err := o.Read(&campaign)

	return &campaign,err

} 
///get social account relation with campaign use CampaignId
func  (u *Campaign)GetSocialAccount(id int64)([]SocialAccount,error){
	o := orm.NewOrm()
	var socialAccount []SocialAccount
	_, err := o.QueryTable(new(SocialAccount)).Filter("campaign_id", id).All(&socialAccount, "id", "campaign_id", "user_name", "pass_word", "socialplatform_id", "createtime")
	if err != nil {
		return nil, err
	}
	return socialAccount, nil
}