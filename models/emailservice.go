package models
import(
	"net/smtp"
	"strings"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	"errors"
)

type EmailService struct{
	Id int64 `orm:"pk;auto"`
	From string `orm:"size(250)" valid:"Required"`
	Password string `orm:"size(250)" valid:"Required"`
	Host string `orm:"size(150)" valid:"Required"`
	Port string `orm:"size(4)" valid:"Required"`	
}
func (u *EmailService) TableName() string {
	return "email_service"
}
func init() {
	orm.RegisterModelWithPrefix("mk_", new(EmailService))
}
///get smtp connect client
func(u *EmailService)GetsmtpAuth(emailService EmailService)(smtp.Auth){
	return smtp.PlainAuth("", emailService.From, emailService.Password, emailService.Host)
}
///send email
func(u *EmailService)Sendemail(emailService EmailService,toList []string, subject string, body string)(error){
	toHeader := strings.Join(toList, ",")
	msg := []byte("From: " + emailService.From + "\n" +
        "To: " + toHeader + "\n" + // use toHeader
        "Subject: "+subject+"\n\n" +
        body)
	auth:=u.GetsmtpAuth(emailService)	
	return smtp.SendMail(emailService.Host+":"+emailService.Port, auth, emailService.From, toList,msg)
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
	id, err := o.Insert(&emser)
	if err != nil {
		return 0, err
	}
	return id, err
}

