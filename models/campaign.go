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
	Types string  `orm:"size(20);null"` //the type of campaign, email, social
	Disable int `orm:"default(0)"` //0: disabled, 1: enabled
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
func (u *Campaign)CreateCampaign(username string,tags string,types string) (id int64, err error) {
	o := orm.NewOrm()
	var us Campaign
	us.CampaignName = username
	us.Tags=tags
	us.Types=types
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
/// list campaign by type
func (u *Campaign)ListCambytype(types string,start int,limitNum int)([]Campaign,error){
	o := orm.NewOrm()
	var cam []Campaign
	count, e := o.QueryTable(new(Campaign)).Filter("types",types).Filter("disable",0).Limit(limitNum,start).All(&cam, "campaign_id", "campaign_name","tags","types")
	if e != nil {
		return nil, e
	}
	if count <= 0 {
		return nil, errors.New("nothing found")
	}
	return cam, nil
}
