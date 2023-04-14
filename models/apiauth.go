package models

import (
	"marketing/utils"
	"time"
	"strings"
	"github.com/beego/beego/v2/client/orm"
	// "github.com/beego/beego/v2/core/logs"
)
var DefaultApiauth *Apiauth

type Apiauth struct {
	Id         int64  `orm:"pk;auto"`
	UserName   string `orm:"size(100)"`
	Password   string `orm:"size(100)"`
	Updated time.Time `orm:"null;auto_now;type(datetime)"`
}

// set engine as INNODB
func (u *Apiauth) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(Apiauth))
}

//get api data id by username and password
func (u *Apiauth)GetApiAuth(username string,password string) (int64, error) {
	pass:=utils.Md5V2(password)
	// logs.Info(username)
	// logs.Info(pass)
	o := orm.NewOrm()
	var apiauth Apiauth
	err := o.QueryTable(new(Apiauth)).Filter("user_name", strings.TrimSpace(username)).Filter("password", strings.TrimSpace(pass)).One(&apiauth,"id")
	// logs.Error(err)
	if err != nil {
		return 0, err
	}
	// logs.Info(apiauth.Id)
	return apiauth.Id, nil
}