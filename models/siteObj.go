package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"

)
var DefaultSiteObj *SiteObj
type SiteObj struct {
	SiteId      int64     `orm:"pk;auto"`
	Campaign *Campaign  `orm:"rel(fk);on_delete(do_nothing)"`	  
	TplRecord time.Time `orm:"auto_now;type(datetime)"`
}
func (u *SiteObj) TableName() string {
	return "site_obj"
}

func init() {
orm.RegisterModelWithPrefix("mk_", new(EmailTpl))
}



