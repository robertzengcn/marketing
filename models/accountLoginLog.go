package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"

)
var DefaultAccountLoginLog *AccountLoginLog
type AccountLoginLog struct {
	Id      int64     `orm:"auto"`
	Account *Account  `orm:"rel(fk);on_delete(do_nothing)"`
	LoginTime time.Time `orm:"auto_now_add;type(datetime)"`
}
func (u *AccountLoginLog) TableName() string {
	return "account_login_log"
}
func (u *AccountLoginLog) TableEngine() string {
	return "MyISAM"
}
func init() {
orm.RegisterModelWithPrefix("mk_", new(AccountLoginLog))
}

func (u *AccountLoginLog) AccountLogin(account *Account) (id int64,err error){
	o := orm.NewOrm()
	var accountloginlog AccountLoginLog
	accountloginlog.Account=account
	id, err = o.Insert(&accountloginlog)
	return id,err	
}