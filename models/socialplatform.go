package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	// "marketing/utils"
	"strings"
)
var DefaultSocialPlatform *SocialPlatform
type SocialPlatform struct {
	Id int64 `orm:"pk;auto" json:"id"`
	Name string `orm:"size(100)" json:"name"`
	Url string `orm:"size(1000)" json:"url"`
	LoginUrl string `orm:"size(1000)" json:"login_url"`
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
	qerr:=qs.Filter("id", id).One(&socialplatform, "id","name","url","login_url")
	if(qerr!=nil){
		return socialplatform,qerr
	}
	return socialplatform,nil
}
//list social platform
func (u *SocialPlatform)Listsocialplatform()([]SocialPlatform,error){
	o := orm.NewOrm()
	var socialplatforms []SocialPlatform
	qs := o.QueryTable(new(SocialPlatform))
	_,err:=qs.All(&socialplatforms)
	if(err!=nil){
		return socialplatforms,err
	}
	return socialplatforms,nil
}
