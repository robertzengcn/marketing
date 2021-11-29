package models

import (
	// "fmt"
	// "errors"
	// "fmt"
	  "github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	// "time"
)

type Campaign struct {
	Campaign_id      int64     `orm:"pk;auto"`
	Campaign_name    string    `orm:"size(100)"`
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
func CreateCampaign(username string) (id int64, err error) {
	o := orm.NewOrm()
	var us Campaign
	us.Campaign_name = username
	id, err = o.Insert(&us)
		return id,err
}
