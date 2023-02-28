package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

type SocialProfile struct {
	Id int64 `orm:"pk;auto"`
	Name string `orm:"size(100)"`
	PhoneNumber string `orm:"size(100)"`
	Email string `orm:"size(100)"`
}

// set engine as INNODB
func (u *SocialProfile) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(SocialProfile))
}

//create social SocialProfile Data
func (u *SocialProfile)CreateSocialProfile(name string, phoneNumber string,email string) (int64, error) {
	o := orm.NewOrm()
	socialProfile := SocialProfile{Name: name, PhoneNumber: phoneNumber,Email:email}
	id, err := o.Insert(&socialProfile)
	return id, err
}

