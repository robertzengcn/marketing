package models

import (
	"encoding/json"
	"errors"
	"marketing/utils"
	"path/filepath"
	"runtime"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	// "errors"
	"sync"
)

type FetchEmail struct {
	Id int64  `orm:"pk;auto" json:"-"`
	Url string `orm:"size(150)" json:"url"`
	Email string `orm:"size(150)" json:"email"`
	Description string `orm:"size(300)" json:"description"`
	RunId int64 `orm:"column(taskrunid)"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
}
func (u *FetchEmail) TableName() string {
	return "fetchemail"
}
func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_", new(FetchEmail))
	// create table
	// orm.RunSyncdb("default", false, true)
}
///fetch email from task
func (u *FetchEmail)Fetchtaskemail(taskrunid int64)(error){
	var searchreqModel SearchRequest
	// searchreqModel
	seaReqArr,seaErr:=searchreqModel.Getrequestrunid(taskrunid)
	logs.Info(seaReqArr);
	if(seaErr!=nil){
		return seaErr
	}
	for _,sv:=range seaReqArr{
	var serplinkModel SerpLink
	
	serpList,_,serpLerr:=serplinkModel.GetlistbyReqid(sv.Id)
	if(serpLerr!=nil){
		// return serpLerr
		logs.Error(serpLerr)
		continue
	}
	if(len(serpList)<1){
		// return errors.New("not find any link to fetch email")
		logs.Error("not find any link to fetch email")
		continue
	}
	blacklistVar:=Blacklist{}
	var wg sync.WaitGroup
	for _, s := range serpList {
		topDomain,derror:=utils.Gettopdomain(s.Domain)
		if(derror!=nil){
		logs.Error(derror)
			continue
		}
		//check is the item in blacklist
		bres,_:=blacklistVar.Getone(topDomain)
		if(bres!=nil){//item in black list
			continue
		}
		//check whether the domain alreay have email in db
		sCount,_:=u.Fetchemaildomain(s.Domain)
		if(sCount>0){
			//domain alreay exist
			continue
		}
		// Increment the WaitGroup counter.
		wg.Add(1)
		go u.Sendquerycom(s.Link,taskrunid,&wg,true)

	}
	wg.Wait()
	}
	// logs.Info("fetch email complete")
	return nil
}
///send query email command
func (u *FetchEmail)Sendquerycom(url string,runid int64,wg *sync.WaitGroup,sendEmail bool)(error){
	// Decrement the counter when the goroutine completes.
	defer wg.Done()
	gHost, gherr := beego.AppConfig.String("emailscrape::host")
	if gherr != nil {
		// taskModel.Handletaskerror(&Result{Runid: runid, Output: "", Err: gherr})
		return gherr
	}
	gPort, gperr := beego.AppConfig.String("emailscrape::port")
	if gperr != nil {
		// taskModel.Handletaskerror(&Result{Runid: runid, Output: "", Err: gperr})
		return gperr
	}
	gUser, gerr := beego.AppConfig.String("emailscrape::user")
	if gerr != nil {
		// logs.Error(gerr)
		// taskModel.Handletaskerror(&Result{Runid: runid, Output: "", Err: gerr})
		return gerr
	}
	gPass, gperr := beego.AppConfig.String("emailscrape::pass")
	if gperr != nil {
		// taskModel.Handletaskerror(&Result{Runid: runid, Output: "", Err: gperr})
		return gperr
	}
	conn, cerr := utils.Connect(gHost+":"+gPort, gUser, gPass)
	if cerr != nil {
		// logs.Error(cerr)
		// taskModel.Handletaskerror(&Result{Runid: runid, Output: "", Err: cerr})
		return cerr 
	}
	savefilemd:=utils.Md5V2("/app/workspace/"+url)
	savefile:=savefilemd+".json"
	fetCommand:="Emailscrapy -u "+url+" -o "+savefile
	logs.Info(fetCommand)
	kout, kerr := conn.SendCommands(fetCommand)
	logs.Info(string(kout))
	if kerr != nil {
		logs.Error(kerr)
		// taskModel.Handletaskerror(&Result{Runid: runid, Output: string(kout), Err: kerr})
		return kerr
	}
	// logs.Info(kout)
	sftpClient, sftperr := conn.Createsfptclient()
	if sftperr != nil {
		logs.Error(sftperr)		
		return sftperr
	}
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	localFilepath := apppath + "/output/"+savefilemd+".json"
	// findfile:="/app/workspace/"+url+".json"
	derr := conn.Downloadfile(sftpClient, savefile, localFilepath)
	if derr != nil {
		logs.Error(derr)
		// u.Handletaskerror(&Result{Runid: runid, Output: string(kout), Err: derr})
		return derr
	}
	// serequestarr, rerr := u.Readfile(localFilepath)
	// if rerr != nil {
	// 	logs.Error(rerr)
	// 	u.Handletaskerror(&Result{Runid: runid, Output: string(kout), Err: rerr})
	// 	return
	// }
	byteValue, _ :=utils.ReadFile(localFilepath)
	var fetcharr []FetchEmail
	json.Unmarshal(byteValue, &fetcharr)
	logs.Info(fetcharr)
	if(len(fetcharr)<1){
		return errors.New("now find email")
	}
	mailModel:=MailLog{}
	emailser := EmailService{}
	for _,x:= range fetcharr {
		x.RunId=runid
		u.SaveEmail(x)
		if(sendEmail){//require send email
			mbools,mErr:=mailModel.Checkemailsend(x.Email,runid)
			if(mErr!=nil){
				logs.Error(mErr)
			}
			if(mbools){
				continue
			}
			emailser.Sendemailtask(&x, runid)

		}
	}
	
	return nil
}
///save email to local
func (u *FetchEmail)SaveEmail(fetchemail FetchEmail)(int64,error){
	o := orm.NewOrm()
	// fetchObj := FetchEmail{Email: fetchemail.Email,RunId:fetchemail.RunId}
	logs.Info(fetchemail)
	fetchemailM:=FetchEmail{}
	qs := o.QueryTable(u)
	err :=qs.Filter("email",fetchemail.Email).Filter("taskrunid", fetchemail.RunId).One(&fetchemailM)
	logs.Error(err)
	if(err == orm.ErrNoRows){
		id, aerr := o.Insert(&fetchemail)
		if(aerr!=nil){
			return 0,aerr
		}
		return id, nil
	}
	return 0,nil
	// if created, id, err := o.ReadOrCreate(&fetchObj, "email","taskrunid"); err == nil {
	// 	if created {
	// 		return id,err
	// 	} else {
	// 		return id,err
	// 	}
	// }
	// return 0,errors.New("not found")
}
///get all email by task run id
func (u *FetchEmail)Fetchallemail(taskrunid int64)([]*FetchEmail,int64,error){
	o := orm.NewOrm()
	var fetmails []*FetchEmail
	qs := o.QueryTable(u)
	num, err :=qs.Filter("taskrunid", taskrunid).All(&fetmails)
	return fetmails,num,err
}
///check whethe the domain already have email relative
func (u *FetchEmail)Fetchemaildomain(domain string)(int64,error){
	
	o := orm.NewOrm()
	qs := o.QueryTable(u)

	return qs.Filter("url__contains", domain).Count()
}
