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
//list schedule api
func (c *ScheduleController) ListSchedule(){
	limitVar,limitSerr:=c.GetInt("limit",25)
	if(limitSerr!=nil){
		c.ErrorJson(20230306144560,"get limit error",nil)
	}
	offsetVar,offsetSerr:=c.GetInt("offset",0)
	if(offsetSerr!=nil){
		c.ErrorJson(20230306144665,"get offset error",nil)
	}
	if(limitVar<1){
		limitVar=25
	}
	if(offsetVar<1){
		offsetVar=0
	}
	scheduleModel:=models.Schedule{}
	scheduleList,err:=scheduleModel.ScheduleList(limitVar,offsetVar)
	if(err!=nil){
		c.ErrorJson(20230306144775,"list schedule error",nil)
	}
	c.SuccessJson(scheduleList)
}
