package models

import (
	"errors"
	"net/smtp"
	"strings"
	"time"
	"strconv"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

type EmailService struct{
	Id int64 `orm:"pk;auto"`
	From string `orm:"size(250)" valid:"Required"`
	Password string `orm:"size(250)" valid:"Required"`
	Host string `orm:"size(150)" valid:"Required"`
	Port string `orm:"size(4)" valid:"Required"`	
	Campaign *Campaign `orm:"rel(fk);on_delete(do_nothing);column(campaign_id)"`
	Name string `orm:"size(250);description(the name of mailservice)"`
	Status int `orm:"size(1);default(1);description(this mean status of the mailservice)"`
	Usetime time.Time `orm:"null;type(datetime)"`
}

///defined table name
func (u *EmailService) TableName() string {
	return "email_service"
}
func init() {
	orm.RegisterModelWithPrefix("mk_", new(EmailService))
}
///get smtp connect client
func(u *EmailService)GetsmtpAuth(emailService *EmailService)(smtp.Auth){
	return smtp.PlainAuth("", emailService.From, emailService.Password, emailService.Host)
}
///send email
func(u *EmailService)Sendemail(emailService *EmailService,toList []string, subject string, body string)(error){
	toHeader := strings.Join(toList, ",")
	msg := []byte("From: " + emailService.From + "\n" +
        "To: " + toHeader + "\n" + // use toHeader
        "Subject: "+subject+"\n\n" +
        body)
	auth:=u.GetsmtpAuth(emailService)	
	seerr:=smtp.SendMail(emailService.Host+":"+emailService.Port, auth, emailService.From, toList,msg)
	if(seerr!=nil){
		return seerr
	}
	u.Updatesendtime(emailService.Id)
	return nil
}

///create email service
func (u *EmailService) Createemailser(emser EmailService) (int64, error) {
	valid := validation.Validation{}
	b, verr := valid.Valid(&emser)
	if verr != nil {
		// handle error
		return 0, verr
	}
	if !b {
		// validation does not pass
		var errMessage string
		for _, err := range valid.Errors {
			errMessage += err.Key + ":" + err.Message
		}
		return 0, errors.New(errMessage)
	}
	o := orm.NewOrm()
	// id, err := o.Insert(&emser)
	// if err != nil {
	// 	return 0, err
	// }
	// return id, err
	if created, id, err := o.ReadOrCreate(&emser, "from","password","host","port","campaign_id","name"); 
	err == nil {
		if created {
			// fmt.Println("New Insert an object. Id:", id)
			logs.Info("New Insert an object. Id:"+strconv.FormatInt(id,10))
			} else {
			// fmt.Println("Get an object. Id:", id)
			logs.Info("Get an object. Id:"+strconv.FormatInt(id,10))
			}
		return id,err	
	}else{
		return 0,err	
	}
	// return 0,errors.New("unkown error")
}
///get email service by campaign id
func (u *EmailService) GetEsbycam(campaignId int64)(*EmailService,error){
	var ess []EmailService
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	_,mailerr:=qs.Filter("campaign_id", campaignId).OrderBy("usetime asc").Limit(1).All(&ess, "Id", "From","Password","Host","Port")
	if(mailerr!=nil){
		return nil,mailerr
	}
	return &ess[0],nil
}
///update email send time
func (u *EmailService) Updatesendtime(sid int64)(int64, error){
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	currentTime := time.Now()
	return qs.Filter("id", sid).Update(orm.Params{
		"usetime":currentTime.Format("2006.01.02 15:04:05"),
	})
}
///get one email service by id
func (u *EmailService)GetOne(serId int64)(*EmailService,error){
	o := orm.NewOrm()
	emailSeModel := EmailService{Id: serId}
	err := o.Read(&emailSeModel)
	if(err!=nil){
		return nil, err	
	}else{
		return &emailSeModel,nil
	}
}

