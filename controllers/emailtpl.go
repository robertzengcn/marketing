package controllers

import (
	"marketing/models"
	// "errors"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/i18n"
)

type EmailtplController struct {
	// beego.Controller
	BaseController
	i18n.Locale
}
type CreateEmailres struct {
	Id int64
}

func (c *EmailtplController) CreateEmailtpl() {
	c.Lang = c.BaseController.Lang
	email_title := c.GetString("email_title")
	email_content := c.GetString("email_content")
	campaign_id,cerr:= c.GetInt64("campaign_id")
	if(cerr!=nil){
		c.ErrorJson(20220618162124,cerr.Error(),nil)
	}
	CampaignModel:=models.Campaign{}
	campaiginVar,camerr:=CampaignModel.FindCambyid(campaign_id)
	if(camerr!=nil){
		c.ErrorJson(20220618162729,camerr.Error(),nil)
	}
	if(campaiginVar==nil){
		
		c.ErrorJson(20220618162832,c.Tr("invail_campaign_id"),nil)
	}
	logs.Info(campaiginVar)
	emailtplModel := models.EmailTpl{}
	emailVar := models.EmailTpl{TplTitle: email_title, 
		TplContent: email_content,CampaignId: campaiginVar}
	logs.Info(emailVar)	
	emailId, emailerr := emailtplModel.Createone(emailVar)
	if emailerr != nil {
		c.ErrorJson(20220617155924, emailerr.Error(), nil)
	}
	c.SuccessJson(CreateEmailres{Id: emailId})
}
