package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)
type Blacklist struct {
	Id int64      `orm:"pk;auto"`
	Domain string `orm:"size(255)"`
}
func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(Blacklist))
}
///get one from blacklist
func(u *Blacklist)Getone(domain string)(*Blacklist,error){
	o := orm.NewOrm()
	// blacklist := Blacklist{}
	var blacklist Blacklist
	qs := o.QueryTable(&blacklist)
	err := qs.Filter("domain", strings.TrimSpace(domain)).One(&blacklist) 
	if(err!=nil){
		return nil, err
	}
	return &blacklist,nil
}





