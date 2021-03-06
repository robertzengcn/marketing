package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"marketing/models"
	// "github.com/beego/beego/v2/core/logs"
	// "marketing/utils"
	"github.com/beego/i18n"
	// "fmt"
)

type CampaignController struct {
	// beego.Controller
	BaseController
	i18n.Locale
}

// func (c *CampaignController) Prepare() {
//     c.EnableXSRF = false
// }
func (c *CampaignController)ChildPrepare(){
	
}
// func (c *CampaignController)Welcome(){
// 	l := logs.GetLogger()
// 	c.Data["langTemplateKey"] = c.GetString("lang")
// 	l.Println(c.Data["langTemplateKey"])
// 	c.TplName = "welcome.tpl" 
// }

//create campaign
func (c *CampaignController) CreateCampaign() {
	// l := logs.GetLogger()
	// l.Println("this is a message of get create campaign")
	
	campaign_name := c.GetString("campaign_name")

	campaing_id,err:=models.DefaultCampaign.CreateCampaign(campaign_name)

	if err != nil {
		c.ErrorJson(20211117161926,err.Error(),nil)
	}

	c.SuccessJson(campaing_id)
}
//list campaign use request
func (c *CampaignController) ListCampaign() {
	start,_ := c.GetInt("start",0)
	num,_:= c.GetInt("number",10)
	
	campagins,err:=models.DefaultCampaign.ListCampaign(start,num)
	if(err!=nil){
		c.ErrorJson(20211208153839,err.Error(),nil)

	}
	c.SuccessJson(campagins)
}
// create site api
// func (c *CampaignController) Createsite(){
// 	site:= c.GetString("site")
// 	email:= c.GetString("email")
	
// 	campaignId,_:=c.GetInt64("campaigin_id",0)
// 	if(campaignId<=0){
// 		c.ErrorJson(20211216154049,"campaign id empty",nil)
// 	}
// 	if(!utils.ValidEmail(email)){
// 		// fmt.Println(c.Lang)
// 		c.Lang=c.BaseController.Lang
// 		c.ErrorJson(20211217152054,c.Tr("welcome"),nil)
// 	}
// 	camPaign,err:=models.DefaultCampaign.FindCambyid(campaignId)
// 	if(err!=nil){
// 		c.ErrorJson(20211216154653,err.Error(),nil)
// 	}
// 	siteId,siteErr:=models.DefaultSiteObj.AddSite(camPaign,email,site)
// 	if(siteErr!=nil){
// 		c.ErrorJson(20211216155058,siteErr.Error(),nil)
// 	}
// 	c.SuccessJson(siteId)

// }



