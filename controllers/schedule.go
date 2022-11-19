package controllers

import (

	"marketing/models"	
	"github.com/beego/i18n"

)

type ScheduleController struct {
	// beego.Controller
	BaseController
	i18n.Locale
}

///create task api
func (c *ScheduleController) CreateSchedule() {
	schedulModel:=models.Schedule{}
	scheduleName := c.GetString("schedule_name")
	campaignId,camErr:=c.GetInt64("campaignId",0)
	if(camErr!=nil){
		c.ErrorJson(20221114103421,"get campaign id error",nil)
	}
	camStaut,camSerr:=c.GetInt("status",1)
	if(camSerr!=nil){
		c.ErrorJson(20221114103525,"get schedule status error",nil)
	}
	cycle := c.GetString("cycle")
	if(len(scheduleName)<1){
		c.ErrorJson(202211141028,"schedule cycle empty",nil)
	}
	CampaignModel:=models.Campaign{}
	CampaiginVar,CampErr:=CampaignModel.FindCambyid(campaignId)
	if(CampErr!=nil){
		c.ErrorJson(20221114104333,"find campaigin error",nil)
	}
	
	var x bool
	if(camStaut==1){
		x=true
	}else{
		x=false
	}
	ScheduleVar:=models.Schedule{
		Name: scheduleName,
		CampaignId: CampaiginVar,
		Status: x,
		Cycle: cycle,
	}
	sresReu,sresErr:=schedulModel.CreateSchedule(ScheduleVar)
	if(sresErr!=nil){
		c.ErrorJson(20221114110650,"create schedule error",nil)
	}
	c.SuccessJson(sresReu)
}