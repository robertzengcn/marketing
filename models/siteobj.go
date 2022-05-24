package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"errors"
)
var DefaultSiteObj *SiteObj
type SiteObj struct {
	SiteId      int64     `orm:"pk;auto"`
	Campaign *Campaign  `orm:"rel(fk);on_delete(do_nothing)"`	  
	SiteRecord time.Time `orm:"auto_now;type(datetime)"`
	Email string `orm:"size(300)"`
	Url string   `orm:"size(400)"`
}
func (u *SiteObj) TableName() string {
	return "site_obj"
}

func init() {
orm.RegisterModelWithPrefix("mk_", new(SiteObj))
}

/// add site
func (u *SiteObj) AddSite(campaign *Campaign,email string,url string)(int64,error){
	o := orm.NewOrm()
	siteObj := SiteObj{Campaign: campaign,Email:email,Url:url}
	// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
	if created, id, err := o.ReadOrCreate(&siteObj, "campaign_id","email","url"); err == nil {
		if created {
			return id,err
		} else {
			return id,err
		}
	}
	return 0,errors.New("not found")
}



