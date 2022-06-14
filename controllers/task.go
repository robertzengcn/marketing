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

func (c *TaskController) CreateTask() {
	campaign_id,_ := c.GetInt("campaign_id")
	emailtpl_id,_ := c.GetInt("emailtpl_id")
	taskModel=models.Task{
		TaskStatus: 1,
		EmailTpl: emailtpl_id,
		CampaignId: campaign_id}

}


