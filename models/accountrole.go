package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type AccountRole struct {
	Id   int64  `orm:"pk;auto"`
	Name string `orm:"size(100)"`
}

func (u *AccountRole) TableName() string {
	return "account_role"
}

func (u *AccountRole) TableEngine() string {
	return "INNODB"
}

func init() {
	orm.RegisterModelWithPrefix("mk_", new(AccountRole))
}

const (
	Admin   = "admin"
	Manager = "manager"
	Normal  = "normal"
)



