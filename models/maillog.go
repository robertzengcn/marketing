package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"

)
var DefaultMailLog *MailLog
type MailLog struct {
	LogId      int64     `orm:"pk;auto"`
	Campaign *Campaign  `orm:"rel(fk);on_delete(do_nothing)"`	  
	TplRecord time.Time `orm:"auto_now;type(datetime)"`
}
func (u *MailLog) TableName() string {
	return "mail_log"
}

func init() {
orm.RegisterModelWithPrefix("mk_", new(MailLog))
}



