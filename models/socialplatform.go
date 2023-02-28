package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

type SocialPlatform struct {
	Id int64 `orm:"pk;auto"`
	Name string `orm:"size(100)"`
}


// set engine as INNODB
func (u *SocialPlatform) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(SocialPlatform))
}
