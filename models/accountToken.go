package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	guuid "github.com/google/uuid"
)
var DefaultAccountToken *AccountToken
type AccountToken struct{
	TokenId int64 `orm:"pk;auto"`
	TokenVal string
	Account   *Account  `orm:"rel(fk);on_delete(do_nothing)"`
	TokenExpired   time.Time `orm:"type(datetime)"`
}
func (u *AccountToken) TableEngine() string {
	return "MYISAM"
} 

func (u *AccountToken) TableName() string {
	return "account_token"
}
func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(AccountToken))
	// create table
	// orm.RunSyncdb("default", false, true)
}
///gen account token
func (u *AccountToken) GenAccounttoken(account *Account) (token string,err error){
	token=guuid.NewString()
	o := orm.NewOrm()
	var ac AccountToken
	ac.Account=account
	ac.TokenVal=token
	now := time.Now()
	ac.TokenExpired=now.AddDate(0, 0, 2)
	_, err = o.Insert(&ac)
	return token,err
}



