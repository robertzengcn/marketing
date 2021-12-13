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
	EmailTpl   *EmailTpl	`orm:"rel(fk);on_delete(do_nothing)"`
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
func (u *Campaign)CreateCampaign(username string) (id int64, err error) {
	o := orm.NewOrm()
	var us Campaign
	us.CampaignName = username
	id, err = o.Insert(&us)
		return id,err
}
/// show all campaign
func (u *Campaign)ListCampaign(start int,limitNum int)([]Campaign,error){
	o := orm.NewOrm()
	var cam []Campaign
	count, e := o.QueryTable(new(Campaign)).Limit(limitNum,start).All(&cam, "Campaign_id", "Campaign_name")
	if e != nil {
		return nil, e
	}
	if count <= 0 {
		return nil, errors.New("nothing found")
	}
	return cam, nil
}
