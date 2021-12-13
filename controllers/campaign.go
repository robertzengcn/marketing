package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"marketing/models"
	"github.com/beego/beego/v2/core/logs"
)

type CampaignController struct {
	// beego.Controller
	BaseController
}

func (c *CampaignController) Prepare() {
    c.EnableXSRF = false
}

func (c *CampaignController) CreateCampaign() {
	l := logs.GetLogger()
	l.Println("this is a message of get create campaign")
	
	campaign_name := c.GetString("campaign_name")

	campaing_id,err:=models.DefaultCampaign.CreateCampaign(campaign_name)

	if err != nil {
		c.ErrorJson(20211117161926,err.Error(),nil)
	}

	c.SuccessJson(campaing_id)
}
///list campaign use request
func (c *CampaignController) ListCampaign() {
	start,_ := c.GetInt("start",0)
	num,_:= c.GetInt("number",10)
	
	campagins,err:=models.DefaultCampaign.ListCampaign(start,num)
	if(err!=nil){
		c.ErrorJson(20211208153839,err.Error(),nil)

	}
	c.SuccessJson(campagins)
}



