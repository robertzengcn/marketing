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
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	filename:=apppath+"/tests/20220509-threaded-results.json"
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
