package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"strings"
	"time"
)

type SocialTask struct {
	Id         int64           `orm:"pk;auto"`
	TaskName   string          `orm:"size(200);column(task_name)"`
	Campaign   *Campaign       `orm:"rel(fk);on_delete(do_nothing);column(campaign_id)"`
	Type       *SocialTaskType `orm:"rel(fk);on_delete(do_nothing);column(type_id)"`
	Tags       []*Tag          `orm:"-"`
	Keywords   []*Keyword      `orm:"-"`
	Disable    int             `orm:"default(0)"`
	UpdateTime time.Time       `orm:"null;auto_now;type(datetime)"`
}

func (u *SocialTask) TableName() string {
	return "social_task"
}

// set engineer as INNODB
func (u *SocialTask) TableEngine() string {
	return "INNODB"
}
func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_", new(SocialTask))
}

//create social task
func (u *SocialTask) CreateSocialTask(campaignId int64, taskType int64, taskName string) (int64, error) {
	//find social type entity
	socialtasktypeModel := SocialTaskType{}
	socialtasktype, err := socialtasktypeModel.FindSocialTaskTypeById(taskType)
	if err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	socialTask := SocialTask{Campaign: &Campaign{CampaignId: campaignId},
		Type: &SocialTaskType{
			TypeId: socialtasktype.TypeId},
		TaskName: taskName}
	id, err := o.Insert(&socialTask)
	return id, err
}

//update social task
func (u *SocialTask) UpdateSocialTaskById(id int64, taskName string, taskType int64, campaignId int64) (int64, error) {
	o := orm.NewOrm()
	socialTask := SocialTask{Id: id, TaskName: taskName, Type: &SocialTaskType{TypeId: taskType}, Campaign: &Campaign{CampaignId: campaignId}}
	id, err := o.Update(&socialTask)
	return id, err
}

//get social task list by campaignId
func (u *SocialTask) GetSocialTaskList(campaignId int64, page int64, size int64, search string) ([]*SocialTask, int64, error) {
	o := orm.NewOrm()
	var socialTask []*SocialTask
	querySet := o.QueryTable(new(SocialTask)).Filter("disable", 0).Filter("campaign_id", campaignId)
	if len(search) > 0 {
		querySet = querySet.Filter("task_name__icontains", search)
	}
	_, err := querySet.Limit(size, page).All(&socialTask, "id", "task_name", "type", "campaign_id")

	// get the search number
	num, numerr := u.GetSocialTaskNum(campaignId, search)
	if numerr != nil {
		return nil, 0, numerr
	}
	campaigModel := Campaign{}
	for _, soci := range socialTask {
		campaign, cerr := campaigModel.FindCambyid(soci.Campaign.CampaignId)
		if cerr != nil {
			continue
		}
		soci.Campaign = campaign
	}
	return socialTask, num, err
}

//get social task number by campaignId and seach conditon
func (u *SocialTask) GetSocialTaskNum(campaignId int64, search string) (int64, error) {
	o := orm.NewOrm()

	querySet := o.QueryTable(new(SocialTask)).Filter("disable", 0).Filter("campaign_id", campaignId)
	if len(search) > 0 {
		querySet = querySet.Filter("task_name__icontains", search)
	}
	num, err := querySet.Count()
	return num, err
}

//get social task by id
func (u *SocialTask) GetSocialTaskById(id int64) (*SocialTask, error) {
	o := orm.NewOrm()
	var socialTask SocialTask
	err := o.QueryTable(new(SocialTask)).Filter("id", id).Filter("disable", 0).One(&socialTask)
	if socialTask.Id == 0 {
		return nil, err
	}
	//get Keywords

	socialtasklistModel := SocialtaskKeywordList{}
	karr, kerr := socialtasklistModel.GetKeywordIdBySocialTaskId(socialTask.Id)
	if kerr == nil {
		socialTask.Keywords = karr
	}

	//get Tags
	socialtasktaglistModel := SocialtaskTagList{}
	tarr, terr := socialtasktaglistModel.GetTagIdBySocialTaskId(socialTask.Id)
	if terr == nil {
		socialTask.Tags = tarr
	}
	//get Campaign info
	campaignModel := Campaign{}
	camVar, camErr := campaignModel.FindCambyid(socialTask.Campaign.CampaignId)
	if camErr == nil {
		socialTask.Campaign.CampaignName = camVar.CampaignName
		socialTask.Campaign.Tags = camVar.Tags
		socialTask.Campaign.Types = camVar.Types
		socialTask.Campaign.Disable = camVar.Disable
	}
	//get social task type
	socialtasktypeModel := SocialTaskType{}
	socialtasktype, socialtasktypeErr := socialtasktypeModel.FindSocialTaskTypeById(socialTask.Type.TypeId)
	if socialtasktypeErr == nil {
		socialTask.Type.TypeName = socialtasktype.TypeName
	}
	return &socialTask, err
}

//get social task by taskid and campaignId
func (u *SocialTask) GetsocialcamIdTid(taskid int64, campaignId int64) (*SocialTask, error) {
	o := orm.NewOrm()
	var socialTask SocialTask
	err := o.QueryTable(new(SocialTask)).Filter("id", taskid).Filter("campaign_id", campaignId).One(&socialTask)
	return &socialTask, err
}

//update social task
func (u *SocialTask) UpdateSocialTask(taskid int64, campaignId int64, taskName string, taskType int64, tags []string, keywords []string) (int64, error) {
	//get account id by campaignId
	campaignModel := Campaign{}
	campaign, cerr := campaignModel.FindCambyid(campaignId)
	if cerr != nil {
		return 0, cerr
	}
	o := orm.NewOrm()
	socialTask := SocialTask{Id: taskid, TaskName: taskName, Campaign: &Campaign{CampaignId: campaignId},
		Type: &SocialTaskType{
			TypeId: taskType}}
	id, err := o.Update(&socialTask)
	if err != nil {
		return 0, err
	}
	// update tags
	u.UpdateSocialTaskTags(taskid, tags, campaign.AccountId.Id)

	//update keywords
	u.UpdateSocialTaskKeywords(taskid, keywords, campaign.AccountId.Id)
	return id, err
}

//update social task keywords
func (u *SocialTask) UpdateSocialTaskKeywords(taskid int64, keywords []string, accountId int64) error {
	//save keywords
	keywordModel := Keyword{}
	var kintarr []int64
	logs.Info(keywords)
	for _, keyword := range keywords {
		if len(strings.TrimSpace(keyword)) > 0 {
			kid, kerr := keywordModel.Savekeyworddb(Keyword{Keyword: keyword}, 0, accountId)
			logs.Info(kid)
			if kerr != nil {
				return kerr
			} else {
				kintarr = append(kintarr, kid)
			}
		}
	}
	logs.Info(kintarr)
	//update keywords
	socialtaskkeywordlistModel := SocialtaskKeywordList{}
	kerr := socialtaskkeywordlistModel.UpdateKeywordsToSocialTask(taskid, kintarr)
	if kerr != nil {
		return kerr
	}
	return nil
}

//update social task tag
func (u *SocialTask) UpdateSocialTaskTags(taskid int64, tags []string, accountId int64) error {
	//defined tag id list
	var tagIdList []int64
	//defined social tag model
	tagModel := Tag{}
	for _, tag := range tags {
		if len(strings.TrimSpace(tag)) > 0 {
			tagId, terr := tagModel.AddTagsByString(tag, accountId)
			logs.Info(tagId)
			if terr == nil {
				tagIdList = append(tagIdList, tagId)
			}
		}
	}

	//update tags
	socialtasktaglistModel := SocialtaskTagList{}

	_, err := socialtasktaglistModel.UpdateSocialTaskTags(taskid, tagIdList)
	return err
}
