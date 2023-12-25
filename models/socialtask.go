package models
import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)
type SocialTask struct {
	Id int64 `orm:"pk;auto"`
	TaskName string `orm:"size(200);column(task_name)"`
	Campaign *Campaign `orm:"rel(fk);on_delete(do_nothing);column(campaign_id)"`
	Type string `orm:"size(100)"`
	Disable int `orm:"default(0)"`
	UpdateTime  time.Time `orm:"null;auto_now;type(datetime)"`
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
func (u *SocialTask)CreateSocialTask(campaignId int64, taskType string,taskName string) (int64, error) {
	o := orm.NewOrm()
	socialTask := SocialTask{Campaign: &Campaign{CampaignId: campaignId},Type:taskType,TaskName: taskName}
	id, err := o.Insert(&socialTask)
	return id, err
}
//get social task list by campaignId
func (u *SocialTask)GetSocialTaskList(campaignId int64,page int64,size int64,search string) ([]*SocialTask,int64,error) {
	o := orm.NewOrm()
	var socialTask []*SocialTask
	querySet := o.QueryTable(new(SocialTask)).Filter("disable", 0).Filter("campaign_id", campaignId)
	if(len(search)>0){
		querySet=querySet.Filter("task_name__icontains",search)
	}
	_,err:=querySet.Limit(size,page).All(&socialTask,"id","task_name","type","campaign_id")
	
	// get the search number
	num,numerr:=u.GetSocialTaskNum(campaignId,search)
	if(numerr!=nil){
		return nil,0,numerr
	}
	campaigModel:=Campaign{}
	for _,soci:= range socialTask{
		campaign,cerr:=campaigModel.FindCambyid(soci.Campaign.CampaignId)
		if(cerr!=nil){
			continue
		}
		soci.Campaign=campaign
	}
	return socialTask, num,err
}

//get social task number by campaignId and seach conditon
func (u *SocialTask)GetSocialTaskNum(campaignId int64,search string) (int64,error) {
	o := orm.NewOrm()
	
	querySet := o.QueryTable(new(SocialTask)).Filter("disable", 0).Filter("campaign_id", campaignId)
	if(len(search)>0){
		querySet=querySet.Filter("task_name__icontains",search)
	}
	num,err:=querySet.Count()
	return num,err
}

//get social task by id
func (u *SocialTask)GetSocialTaskById(id int64) (*SocialTask, error) {
	o := orm.NewOrm()
	var socialTask SocialTask
	err := o.QueryTable(new(SocialTask)).Filter("id", id).Filter("disable", 0).One(&socialTask)
	if(socialTask.Id>0){
		//get Campaign info
		campaignModel:=Campaign{}
		camVar,camErr:=campaignModel.FindCambyid(socialTask.Campaign.CampaignId)
		if(camErr==nil){
			socialTask.Campaign.CampaignName=camVar.CampaignName
			socialTask.Campaign.Tags=camVar.Tags
			socialTask.Campaign.Types=camVar.Types
			socialTask.Campaign.Disable=camVar.Disable
		}
	}
	return &socialTask, err
}
//get social task by taskid and campaignId
func (u *SocialTask)GetsocialcamIdTid(taskid int64, campaignId int64)(*SocialTask, error){
	o := orm.NewOrm()
	var socialTask SocialTask
	err := o.QueryTable(new(SocialTask)).Filter("id", taskid).Filter("campaign_id", campaignId).One(&socialTask)
	return &socialTask,err
}