package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"marketing/models"
)

type TaskController struct {
	// beego.Controller
	BaseController
	i18n.Locale
}
///create task api
func (c *TaskController) CreateTask() {
	campaign_id, _ := c.GetInt64("campaign_id")
	emailtpl_id, _ := c.GetInt64("emailtpl_id")
	task_name:= c.GetString("task_name")

	//check campaign correct
	campaignModel := models.Campaign{CampaignId: campaign_id}
	campaignVar, cerr := campaignModel.FindCambyid(campaign_id)
	if cerr != nil {
		c.ErrorJson(20220615100923, cerr.Error(), nil)
	}
	//check email tpl id correct
	emailtplModel := models.EmailTpl{}
	emailtplVar, emerr := emailtplModel.GetOne(emailtpl_id)
	if emerr != nil {
		c.ErrorJson(20221615094023, emerr.Error(), nil)
	}
	taskstatusModel := models.TaskStatus{Id: 1}
	taskModel := models.Task{}
	taskVar := models.Task{
		TaskName: task_name,
		TaskStatus: &taskstatusModel,
		EmailTpl:   emailtplVar,
		CampaignId: campaignVar}

	taskid, taskerr := taskModel.Createtask(taskVar)
	if(taskerr!=nil){
		c.ErrorJson(20220615102240,taskerr.Error(),nil)
	}
	c.SuccessJson(taskid)
}
///update task status
func (c *TaskController) UpdateTaskstatus() {
	status_id, _ := c.GetInt64("status_id")
	task_id, _ := c.GetInt64("task_id")
	TaskModel:=models.Task{}
	
	terr:=TaskModel.Updatetaskstatus(task_id,status_id)	
	if(terr!=nil){
		c.ErrorJson(20220617155052,terr.Error(),nil)
	}
	c.SuccessJson(nil)
}
