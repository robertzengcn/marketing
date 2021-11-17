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

func (this *CampaignController) Prepare() {
    this.EnableXSRF = false
}

func (this *CampaignController) CreateCampaign() {
	l := logs.GetLogger()
	l.Println("this is a message of get create campaign")
	
	campaign_name := this.GetString("campaign_name")

	campaing_id,err:=models.CreateCampaign(campaign_name)

	if err != nil {
		this.ErrorJson(20211117161926,err.Error(),nil)
	}

	this.SuccessJson(campaing_id)
}

