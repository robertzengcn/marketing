package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	_ "github.com/go-sql-driver/mysql"
)
var DefaultEmailTpl *EmailTpl
type EmailTpl struct {
	TplId      int64     `orm:"pk;auto"`
	TplTitle string  `orm:"size(250)" valid:"Required"`
	TplContent string `orm:"type(text)" valid:"Required"`
	CampaignId *Campaign `orm:"rel(fk);on_delete(do_nothing)"`
	TplRecord time.Time `orm:"auto_now;type(datetime)"`
	UpdateTime time.Time `orm:"auto_now;type(datetime)"`
	
}
func (u *EmailTpl) TableName() string {
	return "email_tpl"
}

func init() {
orm.RegisterModelWithPrefix("mk_", new(EmailTpl))
}
///find one email tpl by email tpl id
func (u *EmailTpl)GetOne(tplId int64)(*EmailTpl,error){
	o := orm.NewOrm()
	emailtpl := EmailTpl{TplId: tplId}
	err := o.Read(&emailtpl)
	if(err!=nil){
		return nil, err	
	}else{
		return &emailtpl,nil
	}
}
///create email tpl
func (u *EmailTpl)Createone(emailtpl EmailTpl)(int64,error){
	valid := validation.Validation{}
	
	b, err := valid.Valid(emailtpl)
    if err != nil {
		logs.Error(err)
       return 0,err
    }
	if !b {

		var errMessage string
	 // validation does not pass
	 for _, err := range valid.Errors {
		// log.Println(err.Key, err.Message)
		errMessage+=err.Key+":"+err.Message
		}
		return 0,errors.New(errMessage)
	}
	logs.Info("valid pass")
	o := orm.NewOrm()		
	id, err := o.Insert(&emailtpl)
	if(err!=nil){
		return 0,err
	}
	return id,err
}





