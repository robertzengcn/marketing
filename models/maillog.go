package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"

)
var DefaultMailLog *MailLog
type MailLog struct {
	LogId      int64     `orm:"pk;auto"`
	Campaign *Campaign  `orm:"rel(fk);on_delete(do_nothing)"`	
	Subject string `orm:"size(150);column(mail_subject)"`
	Content string `orm:"type(text);column(mail_content)"`
	Receiver string `orm:"size(150);column(mail_receiver)"`  
	TplRecord time.Time `orm:"auto_now;type(datetime)"`
	TaskrunId *TaskRun  `orm:"rel(fk);on_delete(do_nothing);column(taskrun_id)"`
	EmailService *EmailService 	 `orm:"rel(fk);on_delete(do_nothing);column(emailser_id)"` 
}
func (u *MailLog) TableName() string {
	return "mail_log"
}

func init() {
orm.RegisterModelWithPrefix("mk_", new(MailLog))
}

///find one result from mail log by email and campaign
func (u *MailLog)Getemailcam(email string,campaignId int64)(int64,error){
	o := orm.NewOrm()
	mailModel:=MailLog{}
	qs := o.QueryTable(&mailModel)
	return qs.Filter("mail_receiver", email).Filter("campaign_id", campaignId).Count()	
}
///add mail log
func (u *MailLog)Addmaillog(maillog MailLog)(int64,error){
	o := orm.NewOrm()
	id, err := o.Insert(&maillog)
	if err==nil{
		return 0,err
	}
	return id,err
}
///check a email whether has been send before
func (u *MailLog)Checkemailsend(email string,taskrunId int64)(bool, error){
	taskrunModel:=TaskRun{}
	task,terr:=taskrunModel.Gettaskbyrun(taskrunId)
	if(terr!=nil){
		return false,terr
	}
	logNum, _ := u.Getemailcam(email, task.CampaignId.CampaignId)
		if logNum > 0 { //mail already exist in log
			return true,nil
		}else{
			return false,nil
		}

}

