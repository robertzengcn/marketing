package models

import (
	"crypto/tls"
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	// "net/mail"
	"math/rand"
	"net/smtp"
	"strconv"
	"strings"
	"time"
	"marketing/utils"
	//"errors"
)

type EmailService struct {
	Id       int64     `orm:"pk;auto"`
	From     string    `orm:"size(250)" valid:"Required"`
	Password string    `orm:"size(250)" valid:"Required"`
	Host     string    `orm:"size(150)" valid:"Required"`
	Port     string    `orm:"size(5)" valid:"Required"`
	Campaign *Campaign `orm:"rel(fk);on_delete(do_nothing);column(campaign_id)"`
	Name     string    `orm:"size(250);description(the name of mailservice)"`
	Ssl      int8       `orm:"size(1);default(1);column(ssl);description(whether to use ssl)"`
	// SenderName string  `orm:"size(250);"`
	Status   int       `orm:"size(1);default(1);description(this mean status of the mailservice)"`
	Usetime  time.Time `orm:"null;type(datetime)"`
	AccountId   *Account   `orm:"rel(fk);on_delete(do_nothing);column(account_id)" valid:"Required"`
	Created  time.Time      `orm:"null;auto_now_add;type(datetime)"`
	Updated  time.Time      `orm:"null;auto_now;type(datetime)"` 
}

///defined table name
func (u *EmailService) TableName() string {
	return "email_service"
}
func init() {
	orm.RegisterModelWithPrefix("mk_", new(EmailService))
}

///get smtp connect client
func (u *EmailService) GetsmtpAuth(emailService *EmailService) smtp.Auth {
	return smtp.PlainAuth("", emailService.From, emailService.Password, emailService.Host)
}

///send email
func (u *EmailService) Sendemailtsl(emailService *EmailService, toList []string, subject string, body string) error {

	toHeader := strings.Join(toList, ",")
	msg := []byte("From: " + emailService.From + "\n" +
		"To: " + toHeader + "\n" + // use toHeader
		"Subject: " + subject + "\n\n" +
		body)
	logs.Info(string(msg))
	// Connect to the SMTP Server
	servername := emailService.Host + ":" + emailService.Port

	auth := u.GetsmtpAuth(emailService)
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         emailService.Host,
	}
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		logs.Error(err)
		return err
	}
	c, err := smtp.NewClient(conn, emailService.Host)
	if err != nil {
		logs.Error(err)
		return err
	}
	// Auth
	if err = c.Auth(auth); err != nil {
		logs.Error(err)
		return err
	}
	// seerr:=smtp.SendMail(emailService.Host+":"+emailService.Port, auth, emailService.From, toList,msg)
	// logs.Error(seerr)
	// if(seerr!=nil){
	// 	return seerr
	// }
	// To && From
	if err = c.Mail(emailService.From); err != nil {
		logs.Error(err)
		return err
	}
	logs.Info(toList)
	for _, v := range toList {	
		sendemail:=strings.TrimSpace(v)
		if(!utils.ValidEmail(sendemail)){
			logs.Error("error email:"+sendemail)
			continue
		}
		if err = c.Rcpt(strings.TrimSpace(v)); err != nil {
			logs.Error(err)
			return err
		}
	}
	// Data
	w, err := c.Data()
	u.Updatesendtime(emailService.Id)
	if err != nil {
		logs.Error(err)
		return err
	}
	
	_, err = w.Write([]byte(msg))
	//update send time log
	

	if err != nil {
		//disable email account
		u.Disableemail(emailService.Id)
		logs.Error(err)
		return err
	}


	err = w.Close()
	if err != nil {
		logs.Error(err)
		return err
	}
	c.Quit()
	
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
	if created, id, err := o.ReadOrCreate(&emser, "from", "password", "host", "port", "campaign_id", "name"); err == nil {
		if created {
			// fmt.Println("New Insert an object. Id:", id)
			logs.Info("New Insert an object. Id:" + strconv.FormatInt(id, 10))
		} else {
			// fmt.Println("Get an object. Id:", id)
			logs.Info("Get an object. Id:" + strconv.FormatInt(id, 10))
		}
		return id, err
	} else {
		return 0, err
	}
	// return 0,errors.New("unkown error")
}

///get email service by campaign id
func (u *EmailService) GetEsbycam(campaignId int64) (*EmailService, error) {
	var ess EmailService
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	err:=qs.Filter("campaign_id", campaignId).Filter("status",1).OrderBy("usetime").One(&ess, "Id", "From", "Password", "Host", "Port")
	if err !=nil{
		return nil, err	
	}
	return  &ess,nil	
}

///update email send time
func (u *EmailService) Updatesendtime(sid int64) (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	currentTime := time.Now()
	return qs.Filter("id", sid).Update(orm.Params{
		"usetime": currentTime.Format("2006-01-02 15:04:05"),
	})
}
//disable email account
func (u *EmailService) Disableemail(sid int64) (error) {
	o := orm.NewOrm()
	emailser:=EmailService{Id:sid}
	if o.Read(&emailser) == nil {
		emailser.Status=0
		if _, err := o.Update(&emailser); 
		err != nil {
			return err
		}else{
			return nil
		}
	}
	return errors.New("not found")
}

///get one email service by id
func (u *EmailService) GetOne(serId int64) (*EmailService, error) {
	o := orm.NewOrm()
	emailSeModel := EmailService{Id: serId}
	err := o.Read(&emailSeModel)
	if err != nil {
		return nil, err
	} else {
		return &emailSeModel, nil
	}
}

///send email to target email
func (u *EmailService) Sendemailtask(fetchemail *FetchEmail, taskrunId int64) error {

	//get CampaignId
	taskrunModel := TaskRun{}
	taskrun, terr := taskrunModel.GetOne(taskrunId)
	if terr != nil {
		return terr
	}
	taskModel := Task{}
	task, taerr := taskModel.GetOne(taskrun.Task.Id)
	if taerr != nil {
		return taerr
	}

	//get all email Email tpl
	emailtplModel := EmailTpl{}
	emArr, emErr := emailtplModel.Getalltpl(task.CampaignId.CampaignId)
	if emErr != nil {
		return emErr
	}
	if(len(emArr)<1){
		return errors.New("email tpl empty")
	}
	//getmail account for send email
	seremail, sererr := u.GetEsbycam(task.CampaignId.CampaignId)
	if sererr != nil {
		return sererr
	}
	//get random email tpl
	rand.Seed(time.Now().Unix())

	chooseEm := emArr[rand.Intn(len(emArr))]
	toMail := make([]string, 1)
	toMail[0] = fetchemail.Email

	//replace email content
	// emailtplModel:=EmailTpl{}
	chooseEm, reErr := emailtplModel.Replacevar(chooseEm, fetchemail)
	if reErr != nil {
		return reErr
	}

	//send email
	serErr := u.Sendemailtsl(seremail, toMail, chooseEm.TplTitle, chooseEm.TplContent)
	if serErr != nil {
		return serErr
	}
	maillogModel := MailLog{Campaign: task.CampaignId,
		Subject:   chooseEm.TplTitle,
		Content:   chooseEm.TplContent,
		Receiver:  toMail[0],
		TaskrunId: taskrun,
		EmailService:seremail,
	}
	maillogModel.Addmaillog(maillogModel)
	return nil
}
//create keyword list from csv data result
func (u *EmailService)CreateRescsv(filepath string)([]EmailService,error){
	data,err:=utils.Csvfilehandle(filepath)	
	if(err!=nil){
		return nil,err
	}
	var EmailServiceArrs []EmailService
	CampaignModel:= Campaign{}
	for i, line := range data {
        if i > 0 { // omit header line
            var rec EmailService
            for j, field := range line {
                if j == 0 {
                    rec.From = strings.TrimSpace(field)					
                } else if j == 1 {
                    rec.Password=strings.TrimSpace(field)
                }else if j==2{
					rec.Name=strings.TrimSpace(field)
				}else if j==3{
					rec.Host=strings.TrimSpace(field)
				}else if j==4{
					rec.Port=strings.TrimSpace(field)
				}else if j==5{
					if(len(field)<1){
						logs.Error("campaign is empty")
					}
					campaignId,_:=strconv.ParseInt(strings.TrimSpace(field),10,64)
					campgn,cerr:=CampaignModel.FindCambyid(campaignId)
					if(cerr!=nil){
						logs.Error(cerr)
						break
					}
					rec.Campaign=campgn
				}else if j==6{
					rec.Status,_=strconv.Atoi(strings.TrimSpace(field))
				}
            }
            EmailServiceArrs = append(EmailServiceArrs, rec)
        }
    }
    return EmailServiceArrs,nil
}
//get email service by id and account id
func (u *EmailService)GetEmailServiceById(id int64,accountId int64) (*EmailService, error) {
	o := orm.NewOrm()
	e := &EmailService{}
	err := o.QueryTable(e).Filter("Id", id).Filter("account_id", accountId).One(e)
	if err == orm.ErrNoRows {
		return nil, err
	}
	return e, nil
}
//update email service by id and account id
func (u *EmailService)UpdateEmailService(e *EmailService) error {
	valid := validation.Validation{}
	
	b, err := valid.Valid(e)
	if err != nil {
		
	   return err
	}
	if !b {

		var errMessage string
	 // validation does not pass
	 for _, err := range valid.Errors {
		// log.Println(err.Key, err.Message)
		errMessage+=err.Key+":"+err.Message
		}
		return errors.New(errMessage)
	}
	//check email service belong to the user
	uef,uerr:=u.GetEmailServiceById(e.Id,e.AccountId.Id)
	if uerr!=nil{
		return uerr
	}
	if(uef==nil){
		return errors.New("email service not found")
	}
	o := orm.NewOrm()
	_, ierr := o.Update(e)
	return ierr
}
///get email service list by account id
func (u *EmailService)GetEmailServiceListByAccountId(accountId int64, page int64, size int64, search string,orderby string)([]*EmailService,error){
	
	var emps []*EmailService
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	// qs.Filter("account_id", accountId)
	cond := orm.NewCondition()
	// qs.Filter("account_id", accountId)
	cond = cond.And("account_id", accountId)
	//qs = qs.SetCond(cond1)
	if(len(search)>0){
		searchCond := orm.NewCondition()
		searchCond = searchCond.Or("name__contains", search).Or("from__contains", search).Or("host__contains", search)
		//cond =cond.AndCond(cond.Or("tpl_title__contains",search).Or("tpl_content__contains",search))
		cond = cond.AndCond(searchCond)
	}
	qs=qs.SetCond(cond)
	if(len(orderby)>0){
		qs=qs.OrderBy(orderby)
	}else{
		qs=qs.OrderBy("-id")
	}
	qs.Limit(size, page).All(&emps)
	
	return emps,nil
}
//delete email service by id and account id
func (u *EmailService)DeleteEmailService(id int64,accountId int64) error {
	//check email service belong to the user
	uef,uerr:=u.GetEmailServiceById(id,accountId)
	if uerr!=nil{
		return uerr
	}
	if(uef==nil){
		return errors.New("email service not found")
	}
	o := orm.NewOrm()
	_, err := o.Delete(&EmailService{Id: id})
	return err
}
//count email service by account id
func (u *EmailService)CountEmailService(accountId int64,search string) (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	// qs.Filter("account_id", accountId)
	cond := orm.NewCondition()
	// qs.Filter("account_id", accountId)
	cond = cond.And("account_id", accountId)
	//qs = qs.SetCond(cond1)
	if(len(search)>0){
		searchCond := orm.NewCondition()
		searchCond = searchCond.Or("name__contains", search).Or("from__contains", search).Or("host__contains", search)
		//cond =cond.AndCond(cond.Or("tpl_title__contains",search).Or("tpl_content__contains",search))
		cond = cond.AndCond(searchCond)
	}
	qs=qs.SetCond(cond)
	return qs.Count()
}