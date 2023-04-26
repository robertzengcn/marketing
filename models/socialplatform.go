package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	// "marketing/utils"
	"strings"
)
var DefaultSocialPlatform *SocialPlatform
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
//find social platform by name
func (u *SocialPlatform)FindSocialPlatformByName(name string) (SocialPlatform, error) {
	//convert name to lower case
	name=strings.ToLower(name)
	o := orm.NewOrm()
	var socialplatform SocialPlatform
	qs := o.QueryTable(new(SocialPlatform))
	qerr:=qs.Filter("name", name).One(&socialplatform, "Id")
	if(qerr!=nil){
		return socialplatform,qerr
	}
	return socialplatform,nil
}
//get social platform by id
func (u *SocialPlatform)GetSocialPlatformById(id int64) (SocialPlatform, error) {
	o := orm.NewOrm()
	var socialplatform SocialPlatform
	qs := o.QueryTable(new(SocialPlatform))
	qerr:=qs.Filter("id", id).One(&socialplatform, "Id","name")
	if(qerr!=nil){
		return socialplatform,qerr
	}
	return socialplatform,nil
}
