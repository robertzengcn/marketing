package models

import (
	"errors"
	"time"
	"net/url"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)
var DefaultEmailTpl *EmailTpl
type EmailTpl struct {
	TplId      int64     `orm:"pk;auto"`
	TplTitle string  `orm:"size(250)" valid:"Required"`
	TplContent string `orm:"type(text)" valid:"Required"`
	CampaignId *Campaign `orm:"rel(fk);on_delete(do_nothing);column(campaign_id)"`
	TplRecord time.Time `orm:"auto_now;type(datetime)"`
	UpdateTime time.Time `orm:"auto_now;type(datetime)"`
	Status int `orm:"size(1);default(1);description(this mean status of the email tpl)"`
	AccountId   *Account   `orm:"rel(fk);on_delete(do_nothing);column(account_id)" valid:"Required"` 
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
///get all email tpl by campaign
func (u *EmailTpl)Getalltpl(campaignId int64)([]*EmailTpl,error){
	var emps []*EmailTpl
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	qs.Filter("campaign_id", campaignId).Filter("status",1).All(&emps)
	return emps,nil
}
///replace email content
func (u *EmailTpl)Replacevar(et *EmailTpl, femail *FetchEmail)(*EmailTpl,error){
	url, err := url.Parse(femail.Url)
	if err != nil {
        return nil,err
    }
	now := time.Now()
	hostname := strings.TrimPrefix(url.Hostname(), "www.")
	et.TplTitle=strings.Replace(et.TplTitle,"{$host}",hostname,-1)
	et.TplTitle=strings.Replace(et.TplTitle,"{$receiver_email}",femail.Email,-1)
	et.TplTitle=strings.Replace(et.TplTitle,"{$send_time}",now.Format(time.ANSIC),-1)
	et.TplTitle=strings.Replace(et.TplTitle,"{$url}",femail.Url,-1)
	et.TplTitle=strings.Replace(et.TplTitle,"{$sender}",femail.Url,-1)

	et.TplContent=strings.Replace(et.TplContent,"{$host}",hostname,-1)
	et.TplContent=strings.Replace(et.TplContent,"{$receiver_email}",femail.Email,-1)
	et.TplContent=strings.Replace(et.TplContent,"{$send_time}",now.Format(time.ANSIC),-1)
	et.TplContent=strings.Replace(et.TplContent,"{$url}",femail.Url,-1)
	et.TplContent=strings.Replace(et.TplContent,"{$description}",femail.Description,-1)
	return et,nil
}
///get email template by id and account id
func (u *EmailTpl)GetEmailTplByIdAndAccountId(tplId int64,accountId int64)(*EmailTpl,error){
	o := orm.NewOrm()
	emailtpl := EmailTpl{TplId: tplId,AccountId:&Account{Id:accountId}}
	err := o.Read(&emailtpl,"TplId","AccountId")
	if(err!=nil){
		return nil, err	
	}else{
		return &emailtpl,nil
	}
}
///get email template list by account id
func (u *EmailTpl)GetEmailTplListByAccountId(accountId int64, page int64, size int64, search string,orderby string)([]*EmailTpl,error){
	
	var emps []*EmailTpl
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	// qs.Filter("account_id", accountId)
	cond := orm.NewCondition()
	// qs.Filter("account_id", accountId)
	cond = cond.And("account_id", accountId)
	//qs = qs.SetCond(cond1)
	if(len(search)>0){
		searchCond := orm.NewCondition()
		searchCond = searchCond.Or("tpl_title__contains", search).Or("tpl_content__contains", search)
		//cond =cond.AndCond(cond.Or("tpl_title__contains",search).Or("tpl_content__contains",search))
		cond = cond.AndCond(searchCond)
	}
	qs=qs.SetCond(cond)
	if(len(orderby)>0){
		qs=qs.OrderBy(orderby)
	}else{
		qs=qs.OrderBy("-tpl_id")
	}
	qs.Limit(size, page).All(&emps)
	
	return emps,nil
}
///update email template by id
func (u *EmailTpl)UpdateEmailTplById(emailtpl EmailTpl)(int64,error){
	//find item by id
	et,err:=u.GetOne(emailtpl.TplId)
	if(err!=nil){
		return 0,err
	}
	//update email item
	et.TplTitle=emailtpl.TplTitle
	et.TplContent=emailtpl.TplContent
	et.CampaignId=emailtpl.CampaignId
	et.Status=emailtpl.Status
	o := orm.NewOrm()
	num, err := o.Update(et)
	if(err!=nil){
		return 0,err
	}
return num,nil
}


///update email template by id and account id
func (u *EmailTpl)UpdateEmailTplByIdAndAccountId(emailtpl EmailTpl)(int64,error){
	//find item by id and account id
	et,err:=u.GetEmailTplByIdAndAccountId(emailtpl.TplId,emailtpl.AccountId.Id)
	if(err!=nil){
		return 0,err
	}
	//update email item
	return u.UpdateEmailTplById(*et)
}
///delete email template by id and account id
func (u *EmailTpl)DeleteEmailTplByIdAndAccountId(tplId int64,accountId int64)(int64,error){
	o := orm.NewOrm()
	num, err := o.QueryTable(u).Filter("TplId", tplId).Filter("AccountId",accountId).Delete()
	if(err!=nil){
		return 0,err
	}
	return num,nil
}







