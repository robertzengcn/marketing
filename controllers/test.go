package controllers

import (
	"marketing/models"
	"github.com/beego/i18n"
	"path/filepath"
	"runtime"
	"fmt"
	"marketing/utils"
	"github.com/beego/beego/v2/core/logs"
)

type TestController struct {
	BaseController
	i18n.Locale
}
///test save search request
func(c *TestController) Savesearchrequest(){
	filenameReq:=c.GetString("filename")
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	filename:=apppath+"/tests/"+filenameReq
	res,rerr:=models.DefaultTask.Readfile(filename)
	if(rerr!=nil){
		panic(rerr.Error())	
	}
	searchreqModel:=models.SearchRequest{}
	serr:=searchreqModel.Savesrlist(res,1)
	if(serr!=nil){
		fmt.Println(serr)
	}
	c.Ctx.WriteString("success")
}
func(c *TestController) Callemailscrape(){
	
}

func (c *TestController) Testsendemail(){
	ser_id,serErr := c.GetInt64("server_id")
	if(serErr!=nil){
		c.ErrorJson(20220816100967, serErr.Error(), nil)
	}
	toemail := c.GetString("to_email")
	if(len(toemail)<1){
		c.ErrorJson(20220815102320, "get to email error", nil)
	}
	if(!utils.ValidEmail(toemail)){
		c.ErrorJson(20220816102775, "to email format error", nil)
	}
	title := c.GetString("title")
	if(len(title)<1){
		c.ErrorJson(20220816103182, "email title error", nil)
	}
	content := c.GetString("content")
	if(len(content)<1){
		c.ErrorJson(20220816103286, "email content error", nil)
	}

	urls := c.GetString("url")
	if(len(urls)<1){
		c.ErrorJson(20220826101991, "email url error", nil)
	}
	description := c.GetString("description")
	if(len(description)<1){
		c.ErrorJson(20220826103495, "description error", nil)
	}

	emailModel:=models.EmailService{}
	emailSer,emailerr:=emailModel.GetOne(ser_id)
	if(emailerr!=nil){
		c.ErrorJson(20220816101672,emailerr.Error(),nil)		
	}
	fetchemail:=models.FetchEmail{Url:urls,Email:toemail,Description:description }
	var toList []string
	toList = append(toList, toemail)
	emailtplmodel:=models.EmailTpl{}
	emailtpl:=models.EmailTpl{TplTitle:title,TplContent:content  }
	emailtpls,emailtplerr:=emailtplmodel.Replacevar(&emailtpl,&fetchemail)
	if(emailtplerr!=nil){
		c.ErrorJson(202208261036110,emailtplerr.Error(),nil)
	}

	sendErr:=emailModel.Sendemailtsl(emailSer,toList,emailtpls.TplTitle,emailtpls.TplContent)
	if(sendErr!=nil){
		c.ErrorJson(20220816103399,sendErr.Error(),nil)		
	}

	c.SuccessJson(nil)
}
func (c *TestController) Checkemailsend(){
	email := c.GetString("email")
	taskrun_id, _ := c.GetInt64("taskrunid")

	mailModel:=models.MailLog{}
	mbools,mErr:=mailModel.Checkemailsend(email,taskrun_id)
	if(mErr!=nil){
		logs.Error(mErr)
	}
	logs.Info(mbools)
	c.SuccessJson(nil)
}
func (c *TestController) Getadultkeyword(){
	keywordModel:=models.Keyword{}
	err:=keywordModel.Getsexkeyword()
	if(err!=nil){
		logs.Error(err)
		c.ErrorJson(202211170947,"error",err)
	}
	c.SuccessJson(nil)
}
///create task by schedule
func (c *TestController)Createtasksched() {
	scheduleid, _ := c.GetInt64("scheduleid")
	scheduleModel:=models.Schedule{}
	sId,sErr:=scheduleModel.Createtask(scheduleid)
	if(sErr!=nil){
		logs.Error(sErr)
		c.ErrorJson(202211170947,sErr.Error(),sErr)
	}
	c.SuccessJson(sId)
}
func(c *TestController)CreatedayTask(){
	tp := c.GetString("type")
	sin := c.GetString("seachenginer","google")
	scheduleModel:=models.Schedule{}
		schVar, schErr:=scheduleModel.Findonebycyc(tp)
		if(schErr!=nil){
			logs.Error(schErr)	
			c.ErrorJson(202211191514126,schErr.Error(),nil)
		}
		staId,staerr:=scheduleModel.Createtask(schVar.Id)
		if(staerr!=nil){
			logs.Error(staerr)	
			c.ErrorJson(202211191514131,staerr.Error(),nil)
		}
		TaskModel:=models.Task{}
		
		go TaskModel.Starttask(staId,sin)
		
		c.SuccessJson(staId)
}

func(c *TestController)Getkeywordbytag(){
	var testarr []string
	testarr = append(testarr, "adult_site")
	keywordModel:=models.Keyword{}
	kArr,kErr:=keywordModel.Getkeywordbytag(testarr,5)
	if(kErr!=nil){
		c.ErrorJson(202211201058144,kErr.Error(),nil)
	}
	c.SuccessJson(kArr)
}
func(c *TestController)Getkeywordapi(){
	keywordModel:=models.Keyword{}
	_,kerr:=keywordModel.Getkeywordapi()
	if(kerr!=nil){
		c.ErrorJson(202212011410152,kerr.Error(),nil)
	}
	c.SuccessJson(nil)
}
//load csv file api
func(c *TestController)Loadkeywordapi(){
	//campaign_id, _ := c.GetInt64("campaign_id")	
	file := c.GetString("filepath")
	keywordModel:=models.Keyword{}
	keywordlist,kerr:=keywordModel.CreateRescsv(file)
	if(kerr!=nil){
		c.ErrorJson(202303011019,kerr.Error(),nil)
	}
	for _,v:=range keywordlist{
		// v=campaign_id
		_,kerr:=keywordModel.Savekeyworddb(v,v.Tag)
		if(kerr!=nil){
		logs.Error(kerr)
	}
	}
	c.SuccessJson(nil)

}
func(c *TestController)GetProxylist(){
	proxyModel:=models.Proxy{}
	go proxyModel.Handleproxy()
	// if(proxyerr!=nil){
	// 	c.ErrorJson(202212011410152,proxyerr.Error(),nil)
	// }
	c.SuccessJson(nil)
}

func(c *TestController)UpdatemulProxy(){
	proxymodel:=models.Proxy{}
		perr:=proxymodel.Updateproxy()
		if(perr!=nil){
			logs.Error(perr)
			c.ErrorJson(202303191057194,perr.Error(),nil)
		}
		c.SuccessJson(nil)
}
