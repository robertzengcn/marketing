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

//get socail account relation with campaign use campaign Id
func (c *CampaignController) GetSocialAccount() {
	campaign_id,_ := c.GetInt64("campaign_id",0)
	
	socialaccounts,err:=models.DefaultCampaign.GetSocialAccount(campaign_id)
	if(err!=nil){
		c.ErrorJson(20211208153839,err.Error(),nil)

	}
	c.SuccessJson(socialaccounts)
}



