package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"

)
var DefaultEmailTpl *EmailTpl
type EmailTpl struct {
	TplId      int     `orm:"pk;auto"`
	TplTitle string  `orm:"size(250)"`
	TplContent string `orm:"type(text)"`
	TplRecord time.Time `orm:"auto_now;type(datetime)"`
}
func (u *EmailTpl) TableName() string {
	return "email_tpl"
}

func init() {
orm.RegisterModelWithPrefix("mk_", new(EmailTpl))
}



