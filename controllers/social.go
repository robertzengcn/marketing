package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"marketing/models"
	// "github.com/beego/beego/v2/core/logs"
	"strings"
)

type SocialController struct {
	BaseController
	i18n.Locale
}

func (c *SocialController) ChildPrepare() {
	// l := logs.GetLogger()
	//      l.Println("22222")
}

// func (c *SocialController) CreateSocialAccount() {
// 	// l := logs.GetLogger()
// 	// l.Println("this is a message of get create social account")

// 	social_name := c.GetString("social_name")
// 	social_phone := c.GetString("social_phone")
// 	social_platformid,serr := c.GetInt64("social_type")
// 	if(serr!=nil){
// 		c.ErrorJson(202302211119,serr.Error(),nil)
// 	}
// 	//social_url := c.GetString("social_url")
// 	social_username := c.GetString("social_username")
// 	social_password := c.GetString("social_password")
// 	//social_token := c.GetString("social_token")
// 	social_email := c.GetString("social_email")
// 	campaign_id,_ := c.GetInt64("campaign_id")

// 	social_id,err:=models.DefaultSocialAccount.CreateSocialAccount(campaign_id,social_username,social_password,social_platformid,social_name,social_phone,social_email,)

// 	if err != nil {
// 		c.ErrorJson(20211117161926,err.Error(),nil)
// 	}

// 	c.SuccessJson(social_id)
// }
type SocialEntity struct {
	CampaignId   int64    `json:"campaign_id"`
	CampaignName string   `json:"campaign_name"`
	Tags         string   `json:"tags"`
	Types        string   `json:"types"`
	Keyword      []string `json:"keyword"`
}
type SocialTaskEntity struct {
	Id           int64  `json:"id"`
	TaskName	 string `json:"task_name"`
	CampaignId   int64  `json:"campaign_id"`
	// CampaignName string `json:"campaign_name"`
	Tags         string `json:"tag"`
	Type         string `json:"type"`
	// Keywords     []string `json:"keywords"`
	ExtraTaskIfo models.TaskExtaInfo `json:"extra_task_info"`
}

///list social campaign
func (c *SocialController) Listsocialcampaigin() {
	types,_ := c.GetInt32("types")
	start, sterr := c.GetInt("start", 0)
	if sterr != nil {
		c.ErrorJson(20230525093049, sterr.Error(), nil)
	}
	end, enderr := c.GetInt("end", 0)
	if enderr != nil {
		c.ErrorJson(20230525093053, sterr.Error(), nil)
	}
	campaignModel := models.Campaign{}
	campaigns, err := campaignModel.ListCambytype(types, start, end)
	if err != nil {
		c.ErrorJson(20230525093057, err.Error(), nil)
	}
	keywordModel := models.Keyword{}
	var socials []SocialEntity
	// extraInfoModel:=models.TaskExtaInfo{}
	// logs.Info(campaigns)
	for i := 0; i < len(campaigns); i++ {
		var soentity SocialEntity
		soentity.CampaignId = campaigns[i].CampaignId
		soentity.CampaignName = campaigns[i].CampaignName
		soentity.Tags = campaigns[i].Tags
		soentity.Types = campaigns[i].Types.CampaignTypeName
		tags := strings.Split(campaigns[i].Tags, ",")
		// extraInfoModel.Getextrainfotaskid(campaigns[i].)
		keywordArr, kErr := keywordModel.Getkeywordbytag(tags, 5)
		if kErr != nil {
			continue
		}
		if len(keywordArr) <= 0 {
			continue
		}
		// soentity.Keyword=[]string{}
		for _, s := range keywordArr {
			soentity.Keyword = append(soentity.Keyword, s.Keyword)
		}
		socials = append(socials, soentity)
	}
	c.SuccessJson(socials)

}

///get social task list
func (c *SocialController) Getsocialtasklist() {
	campagid,_ := c.GetInt64("campaiginId")
	if(campagid<=0){
		c.ErrorJson(202307070932112,"campaign id empty",nil)
	}
	socialTaskmodel := models.SocialTask{}
	soTask, soterr := socialTaskmodel.GetSocialTaskList(campagid)
	if soterr != nil {
		c.ErrorJson(20230526095299, soterr.Error(), nil)
	}
	soentityArr := []SocialTaskEntity{}
	// keywordModel := models.Keyword{}

	taskExtraModel:=models.TaskExtaInfo{}
	
	for _, stv := range soTask {
		var soentity SocialTaskEntity
		soentity.Id = stv.Id
		soentity.CampaignId = stv.Campaign.CampaignId
		soentity.TaskName=stv.TaskName
		
		soentity.Tags = stv.Campaign.Tags
		soentity.Type = stv.Type
		// tags := strings.Split(stv.Campaign.Tags, ",")
		// keywordArr, _ := keywordModel.Getkeywordbytag(tags, 5)
		// if kErr != nil {
		// 	continue
		// }
		// for _, s := range keywordArr {
		// 	soentity.Keywords = append(soentity.Keywords, s.Keyword)
		// }

		taskEx,taskerr:=taskExtraModel.Getextrainfotaskid(stv.Id)
		if(taskerr==nil){
			soentity.ExtraTaskIfo=taskEx
		}
		soentityArr = append(soentityArr, soentity)

	}
	c.SuccessJson(soentityArr)
}

func (c *SocialController) Getsocialtaskinfo() {
	taskId,_ := c.GetInt64("task_id")
	if(taskId<=0){
		c.ErrorJson(202307100914155,"task id empty",nil)
	}
	socialTaskmodel := models.SocialTask{}
	stv,soerr:=socialTaskmodel.GetSocialTaskById(taskId)
	if(soerr!=nil){
		c.ErrorJson(202307100918160,soerr.Error(),nil)
	}
	var soentity SocialTaskEntity
	taskExtraModel:=models.TaskExtaInfo{}
		soentity.Id = stv.Id
		soentity.CampaignId = stv.Campaign.CampaignId
		soentity.TaskName=stv.TaskName	
		soentity.Tags = stv.Campaign.Tags
		soentity.Type = stv.Type
		
	// tags := strings.Split(stv.Campaign.Tags, ",")
	// 	keywordArr, _ := keywordModel.Getkeywordbytag(tags, 5)
	// 	// if kErr != nil {
	// 	// 	continue
	// 	// }
	// 	for _, s := range keywordArr {
	// 		soentity.Keywords = append(soentity.Keywords, s.Keyword)
	// 	}

		taskEx,taskerr:=taskExtraModel.Getextrainfotaskid(stv.Id)
		if(taskerr==nil){
			soentity.ExtraTaskIfo=taskEx
		}
	c.SuccessJson(soentity)
}
func (c *SocialController) Gettaskkeyword(){
	taskId,_ := c.GetInt64("task_id")
	if(taskId<=0){
		c.ErrorJson(202307110927188,"task id empty",nil)
	}
	socialTaskmodel := models.SocialTask{}
	stv,soerr:=socialTaskmodel.GetSocialTaskById(taskId)
	if(soerr!=nil){
		c.ErrorJson(202309270711193,soerr.Error(),nil)
	}
	keywordModel := models.Keyword{}
	tags := strings.Split(stv.Campaign.Tags, ",")
	keywordArr, kErr := keywordModel.Getkeywordbytag(tags, 5)
	 if kErr != nil {
		c.ErrorJson(202307110935199,kErr.Error(),nil)
	 }
	 c.SuccessJson(keywordArr)
}
type SocialTaskResult struct {
	Id int64 `json:"id"`
}
func (c *SocialController) Savesocialtask(){
	campaignid,_ := c.GetInt64("campaign_id")
	if(campaignid<=0){
		c.ErrorJson(202307120949207,"campaign id empty",nil)
	}
	campaignModel:=models.Campaign{}
	camPaign,_:=campaignModel.FindCambyid(campaignid)
	if(camPaign.CampaignId<=0){
		c.ErrorJson(202307121021212,"campaign not exist",nil)
	}
	taskname:=c.GetString("task_name")
	if(len(taskname)<=0){
		c.ErrorJson(202307120950211,"task name empty",nil)
	}
	tasktype:=c.GetString("task_type")
	result_task_id,_:=c.GetInt64("result_task_id")
	socialtaskModel:=models.SocialTask{}
	sid,serr:=socialtaskModel.CreateSocialTask(campaignid,tasktype,taskname)
	if(serr!=nil){
		c.ErrorJson(202307121023223,serr.Error(),nil)
	}
	taskinfoModel:=models.TaskExtaInfo{}
	taskinfoModel.CreateTaskExtraInfo(sid,result_task_id)
	sotsk:=SocialTaskResult{Id:sid}
	c.SuccessJson(sotsk)
}


