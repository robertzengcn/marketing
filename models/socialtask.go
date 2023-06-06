package models
import (
	"github.com/beego/beego/v2/client/orm"
	// "time"
)
type SocialTask struct {
	Id int64 `orm:"pk;auto"`
	Campaign *Campaign `orm:"rel(fk);on_delete(do_nothing);column(campaign_id)"`
	Type string `orm:"size(100)"`
	Disable int `orm:"default(0)"`
	// UpdateTime  time.Time `orm:"null;auto_now;type(datetime)"`
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
func (u *SocialTask)CreateSocialTask(campaignId int64, taskType string) (int64, error) {
	o := orm.NewOrm()
	socialTask := SocialTask{Campaign: &Campaign{CampaignId: campaignId},Type:taskType}
	id, err := o.Insert(&socialTask)
	return id, err
}
//get social task list
func (u *SocialTask)GetSocialTaskList() ([]*SocialTask, error) {
	o := orm.NewOrm()
	var socialTask []*SocialTask
	_, err := o.QueryTable("mk_social_task").Filter("disable", 0).All(&socialTask,"Id","type","campaign_id")
	campaigModel:=Campaign{}
	for _,soci:= range socialTask{
		campaign,cerr:=campaigModel.FindCambyid(soci.Campaign.CampaignId)
		if(cerr!=nil){
			continue
		}
		soci.Campaign=campaign
	}
	return socialTask, err
}