package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"marketing/models"	
	"github.com/beego/i18n"
)

type TaskController struct {
	// beego.Controller
	BaseController
	i18n.Locale
}

///create task api
func (c *TaskController) CreateTask() {
	campaign_id, _ := c.GetInt64("campaign_id")
	// emailtpl_id, _ := c.GetInt64("emailtpl_id")
	task_name := c.GetString("task_name")
	taskkeyword := c.GetString("task_keyword")
	if len(taskkeyword) <= 0 {
		c.ErrorJson(20220621094521, "task keyword empty", nil)
	}

	//check campaign correct
	campaignModel := models.Campaign{CampaignId: campaign_id}
	campaignVar, cerr := campaignModel.FindCambyid(campaign_id)
	if cerr != nil {
		c.ErrorJson(20220615100923, cerr.Error(), nil)
	}
	//check email tpl id correct
	// emailtplModel := models.EmailTpl{}
	// emailtplVar, emerr := emailtplModel.GetOne(emailtpl_id)
	// if emerr != nil {
	// 	c.ErrorJson(20221615094023, emerr.Error(), nil)
	// }
	taskstatusModel := models.TaskStatus{Id: 1}
	taskModel := models.Task{}
	taskVar := models.Task{
		TaskName:   task_name,
		TaskStatus: &taskstatusModel,
		// EmailTpl:   emailtplVar,
		CampaignId: campaignVar}

	taskid, taskerr := taskModel.Createtask(taskVar)
	if taskerr != nil {
		c.ErrorJson(20220615102240, taskerr.Error(), nil)
	}
	taskEntity,terr:=taskModel.GetOne(taskid)
	if(terr!=nil){
		c.ErrorJson(20220621100351,"task not exist",nil)
	}
	taskdetailVar:=models.TaskDetail{Task:taskEntity,Taskkeyword:taskkeyword }
	taskdetailModel:=models.TaskDetail{}
	taskdetailModel.Savetaskdetail(taskdetailVar)
	c.SuccessJson(taskid)
}

///update task status
func (c *TaskController) UpdateTaskstatus() {
	status_id, _ := c.GetInt64("status_id")
	task_id, _ := c.GetInt64("task_id")
	searchenger := c.GetString("searchenger","google")
	TaskModel := models.Task{}
	if(status_id==3){
		//start search on google
	go TaskModel.Starttask(task_id,searchenger)
	}
	// terr := TaskModel.Updatetaskstatus(task_id, status_id)
	// if terr != nil {
	// 	c.ErrorJson(20220617155052, terr.Error(), nil)
	// }
	
	
	c.SuccessJson(nil)
}
