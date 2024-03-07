package controllers

import (
	// beego "github.com/beego/beego/v2/server/web"
	// "strings"
	"marketing/models"

	// "github.com/beego/beego/v2/adapter/logs"
	// "github.com/beego/beego/v2/adapter/logs"
	// "github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/i18n"
	// "github.com/beego/beego/v2/core/logs"
	// "strings"
	"marketing/dto"
)

type SocialController struct {
	BaseController
	i18n.Locale
}

func (c *SocialController) ChildPrepare() {
	// l := logs.GetLogger()
	//      l.Println("22222")
}


type SocialEntity struct {
	CampaignId   int64    `json:"campaign_id"`
	CampaignName string   `json:"campaign_name"`
	Tags         string   `json:"tags"`
	Types        string   `json:"types"`
	Keyword      []string `json:"keyword"`
}
type SocialTaskEntity struct {
	Id         int64  `json:"id"`
	TaskName   string `json:"task_name"`
	CampaignId int64  `json:"campaign_id"`
	// CampaignName string `json:"campaign_name"`
	Tags string `json:"tag"`
	Type string `json:"type"`
	// Keywords     []string `json:"keywords"`
	ExtraTaskIfo models.TaskExtaInfo `json:"extra_task_info"`
}

type SocialTaskResponse struct {
	Records []SocialTaskEntity
	Total   int64
}

///list social campaign use type
func (c *SocialController) Listsocialcampaigin() {
	types, _ := c.GetInt32("types")
	num, _ := c.GetInt("keyword_num") //the keyword number return in task
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

	uid := c.GetSession("uid")
	accountId := uid.(int64)
	// extraInfoModel:=models.TaskExtaInfo{}
	// logs.Info(campaigns)
	for i := 0; i < len(campaigns); i++ {
		var soentity SocialEntity
		soentity.CampaignId = campaigns[i].CampaignId
		soentity.CampaignName = campaigns[i].CampaignName
		// soentity.Tags = campaigns[i].Tags
		soentity.Types = campaigns[i].Types.CampaignTypeName
		var tags []string
		for _, s := range campaigns[i].Tags {
			tags = append(tags, s.TagName)
		}
		//tags := strings.Split(campaigns[i].Tags, ",")
		// extraInfoModel.Getextrainfotaskid(campaigns[i].)
		keywordArr, kErr := keywordModel.Getkeywordbytag(tags, num, accountId)
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
	campagid, _ := c.GetInt64("campaiginId")
	page, perr := c.GetInt64("page", 0)
	if perr != nil {
		c.ErrorJson(202312210933114, perr.Error(), nil)
	}
	size, serr := c.GetInt64("size", 10)
	if serr != nil {
		c.ErrorJson(202312210934118, serr.Error(), nil)
	}
	if campagid <= 0 {
		c.ErrorJson(202307070932112, "campaign id empty", nil)
	}
	searchVal := c.GetString("search")
	socialTaskmodel := models.SocialTask{}
	soTask, soNum, soterr := socialTaskmodel.GetSocialTaskList(campagid, page, size, searchVal)
	if soterr != nil {
		c.ErrorJson(20230526095299, soterr.Error(), nil)
	}
	soentityArr := []SocialTaskEntity{}
	// keywordModel := models.Keyword{}

	taskExtraModel := models.TaskExtaInfo{}

	for _, stv := range soTask {
		var soentity SocialTaskEntity
		soentity.Id = stv.Id
		soentity.CampaignId = stv.Campaign.CampaignId
		soentity.TaskName = stv.TaskName
		//find social task type
		socialtasktypeModel := models.SocialTaskType{}
		socialtasktype, _ := socialtasktypeModel.FindSocialTaskTypeById(stv.Type.TypeId)

		// soentity.Tags = stv.Campaign.Tags
		soentity.Type = socialtasktype.TypeName
		// tags := strings.Split(stv.Campaign.Tags, ",")
		// keywordArr, _ := keywordModel.Getkeywordbytag(tags, 5)
		// if kErr != nil {
		// 	continue
		// }
		// for _, s := range keywordArr {
		// 	soentity.Keywords = append(soentity.Keywords, s.Keyword)
		// }

		taskEx, taskerr := taskExtraModel.Getextrainfotaskid(stv.Id)
		if taskerr == nil {
			soentity.ExtraTaskIfo = taskEx
		}
		soentityArr = append(soentityArr, soentity)
	}
	sot := SocialTaskResponse{
		Records: soentityArr,
		Total:   soNum,
	}
	c.SuccessJson(sot)
}

func (c *SocialController) Getsocialtaskinfo() {
	taskId, _ := c.GetInt64("task_id")
	if taskId <= 0 {
		c.ErrorJson(202307100914155, "task id empty", nil)
	}
	socialTaskmodel := models.SocialTask{}
	stv, soerr := socialTaskmodel.GetSocialTaskById(taskId)
	if soerr != nil {
		c.ErrorJson(202307100918160, soerr.Error(), nil)
	}
	var soentity dto.SocialtaskDto
	taskExtraModel := models.TaskExtaInfo{}
	soentity.Id = stv.Id
	soentity.CampaignId = stv.Campaign.CampaignId
	soentity.TaskName = stv.TaskName
	//get social task type
	socialtasktypeModel := models.SocialTaskType{}
	socialtasktype, _ := socialtasktypeModel.FindSocialTaskTypeById(stv.Type.TypeId)
	// soentity.Tags = stv.Campaign.Tags
	soentity.TypeId = socialtasktype.TypeId
	soentity.TypeName = socialtasktype.TypeName
	//keyword list
	var keywordArr []string
	//loop keywords
	for _, s := range stv.Keywords {
		// logs.Info("the keywords is"+s.Keyword)
		keywordArr = append(keywordArr, s.Keyword)
	}
	soentity.Keywords = keywordArr

	var tags []string
	for _, s := range stv.Tags {
		tags = append(tags, s.TagName)
	}
	soentity.Tags = tags

	taskEx, taskerr := taskExtraModel.Getextrainfotaskid(stv.Id)
	if taskerr == nil {
		soentity.ExtraTaskIfo = taskEx
	}
	c.SuccessJson(soentity)
}

//get keywords in task
func (c *SocialController) Gettaskkeyword() {
	taskId, _ := c.GetInt64("task_id")
	num, _ := c.GetInt("task_id", 5)
	if taskId <= 0 {
		c.ErrorJson(202307110927188, "task id empty", nil)
	}
	uid := c.GetSession("uid")
	accountId := uid.(int64)
	socialTaskmodel := models.SocialTask{}
	stv, soerr := socialTaskmodel.GetSocialTaskById(taskId)
	if soerr != nil {
		c.ErrorJson(202309270711193, soerr.Error(), nil)
	}
	keywordModel := models.Keyword{}
	var tags []string
	for _, s := range stv.Campaign.Tags {
		tags = append(tags, s.TagName)
	}
	if len(tags) <= 0 {
		c.ErrorJson(202312300929195, "tags empty", nil)
	}
	// tags := strings.Split(stv.Campaign.Tags, ",")
	keywordArr, kErr := keywordModel.Getkeywordbytag(tags, num, accountId)
	if kErr != nil {
		c.ErrorJson(202307110935199, kErr.Error(), nil)
	}
	c.SuccessJson(keywordArr)
}

type SocialTaskResult struct {
	Id int64 `json:"id"`
}

func (c *SocialController) Savesocialtask() {
	campaignid, _ := c.GetInt64("campaign_id")
	if campaignid <= 0 {
		c.ErrorJson(202307120949207, "campaign id empty", nil)
	}
	campaignModel := models.Campaign{}
	camPaign, _ := campaignModel.FindCambyid(campaignid)
	if camPaign.CampaignId <= 0 {
		c.ErrorJson(202307121021212, "campaign not exist", nil)
	}

	uid := c.GetSession("uid")
	accountId := uid.(int64)
	if camPaign.AccountId.Id != accountId {
		c.ErrorJson(202401151106279, "campaign not belong to the account", nil)
	}
	taskname := c.GetString("task_name")
	if len(taskname) <= 0 {
		c.ErrorJson(202307120950211, "task name empty", nil)
	}
	tasktype, tasktyperr := c.GetInt64("type_id")
	if tasktyperr != nil {
		c.ErrorJson(202401050956, tasktyperr.Error(), nil)
	}
	//receive keywords list
	keywords := []string{}
	//receive tag list in post
	tags := []string{}
	inputValues, _ := c.Input()
	for k, v := range inputValues {
		if k == "tags[]" {
			if len(v) > 0 {
				tags = append(tags, v...)
			}
		} else if k == "keywords[]" {
			if len(v) > 0 {
				keywords = append(keywords, v...)
			}
		}
	}
	// logs.Info(tags)
	socialtaskId, _ := c.GetInt64("socialtask_id")
	socialtaskModel := models.SocialTask{}
	if socialtaskId > 0 {
		//update social task basic info
		_, serr := socialtaskModel.UpdateSocialTaskById(socialtaskId, taskname, tasktype, camPaign.CampaignId)
		if serr != nil {
			c.ErrorJson(202401121001300, serr.Error(), nil)
		}
		socialtaskModel.UpdateSocialTaskTags(socialtaskId, tags, accountId)
		//update keywords
		socialtaskModel.UpdateSocialTaskKeywords(socialtaskId, keywords, accountId)
	} else {

		sid, serr := socialtaskModel.CreateSocialTask(campaignid, tasktype, taskname)
		if serr != nil {
			c.ErrorJson(202307121023223, serr.Error(), nil)
		}
		//UPDATE social task tags
		socialtaskModel.UpdateSocialTaskTags(sid, tags, accountId)
		// logs.Info(keywords)
		//update social task keywords
		socialtaskModel.UpdateSocialTaskKeywords(sid, keywords, accountId)
		socialtaskId = sid

	}
	sotsk := SocialTaskResult{Id: socialtaskId}
	c.SuccessJson(sotsk)
	// result_task_id, _ := c.GetInt64("result_task_id")

	// if(result_task_id>0){
	// taskinfoModel := models.TaskExtaInfo{}
	// taskinfoModel.CreateTaskExtraInfo(sid, result_task_id)
	// }

}

//list social task type
func (c *SocialController) Listsocialtasktype() {
	socialtasktypeModel := models.SocialTaskType{}
	types, err := socialtasktypeModel.ListSocialTaskType()
	if err != nil {
		c.ErrorJson(202307121023223, err.Error(), nil)
	}
	c.SuccessJson(types)
}
