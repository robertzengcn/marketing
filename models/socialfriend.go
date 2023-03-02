package models

import (
	"time"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

type SocialFriend struct {
	Id         int64  `orm:"pk;auto"`
	Name 	 string `orm:"size(100)"`
	Url string `orm:"size(300)"`
	Createtime time.Time `orm:"auto_now;type(datetime)"`
}



// func (u *SocialFriend) TableName() string {
// 	return "social_account"
// }

// set engine as INNODB
func (u *SocialFriend) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(SocialFriend))
}

