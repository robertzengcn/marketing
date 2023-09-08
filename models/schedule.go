package models

import (
	"errors"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	// "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

type Schedule struct {
	Id int64  `orm:"pk;auto" json:"-"`
	Name string `orm:"size(200)" json:"name"`
	CampaignId *Campaign `orm:"rel(fk);on_delete(do_nothing);column(campaign_id)"`
	Status bool `orm:"null" json:"status"`
	Cycle string `orm:"null" json:"cycle"`
	Created time.Time `orm:"null;auto_now_add;type(datetime)"`
}

func (u *Schedule) TableName() string {
	return "schedule"
}

func init() {
	orm.RegisterModelWithPrefix("mk_", new(Schedule))
}
///create schedule
func (u *Schedule) CreateSchedule(scheduleVar Schedule) (int64, error) {
	valid := validation.Validation{}
	b, verr := valid.Valid(&scheduleVar)
	if verr != nil {
		// handle error
		return 0, verr
	}
	if !b {
		// validation does not pass
		var errMessage string
		for _, err := range valid.Errors {
			errMessage += err.Key + ":" + err.Message
		}
		return 0, errors.New(errMessage)
	}
	o := orm.NewOrm()
	id, err := o.Insert(&scheduleVar)
	if err != nil {
		return 0, err
	}
	return id, err
}
///find schedule by cycle
func (u *Schedule)Findonebycyc(cycle string)(*Schedule,error){
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	scheduleVar:=Schedule{}
	err :=qs.Filter("cycle",cycle).Filter("status",1).One(&scheduleVar)
	if(err!=nil){
		return nil,err
	}
	return &scheduleVar,nil

}
//find all schedule by cycle
func (u *Schedule)Findallschedule(cyc string)([]Schedule,error){
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	var scheduleList []Schedule
	_, err := qs.Filter("cycle",cyc).Filter("status",1).All(&scheduleList)
	if err != nil {
		return nil, err
	}
	return scheduleList, err
}


///create task in schedule
func (u *Schedule)Createtask(scheduleId int64)(int64,error){
	scheduleModel:=Schedule{}
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	err :=qs.Filter("id",scheduleId).Filter("status",1).One(&scheduleModel)
	if(err!=nil){
		return 0,err
	}
	campaignModel:=Campaign{}
	campaignVar,campaignErr:=campaignModel.FindCambyid(scheduleModel.CampaignId.CampaignId)
	if(campaignErr!=nil){
		return 0,campaignErr
	}
	currentTime := time.Now()
	taskstatusModel := TaskStatus{Id: 1}

	taskVar:=Task{
		TaskName: "schedule create task "+currentTime.Format("2006-01-02 15:04:05"),
		TaskStatus: &taskstatusModel,
		CampaignId: campaignVar,
	}
	tags :=strings.Split(campaignVar.Tags,",")
	// logs.Info(tags)
	if(len(tags)<=0){
		return 0,errors.New("tag empty")
	}
	keywordModel:=Keyword{}
	keywordArr,kErr:=keywordModel.Getkeywordbytag(tags,5)
	if(kErr!=nil){
		return 0,kErr
	}
	if(len(keywordArr)<=0){
		return 0,errors.New("the campaign do not have keyword related")
	}
	// logs.Info(keywordArr)
	var keyContent string

	for _, v := range keywordArr {
		keyContent+=v.Keyword+"\n"
	}
	// logs.Info(keyContent)
	taskModel:=Task{}
	taskId,terr:=taskModel.Createtask(taskVar)
	if(terr!=nil){
		return 0,terr
	}
	taskOne,taskErr:=taskModel.GetOne(taskId)
	if(taskErr!=nil){
		return 0,taskErr
	}
	taskDetailModel:=TaskDetail{
		Task:taskOne, 
		Taskkeyword: keyContent,
	}
	_,tadetailErr:=taskDetailModel.Savetaskdetail(taskDetailModel)
	if(tadetailErr!=nil){
		return 0,tadetailErr
	}
	return taskId,nil
}
///show schedule list
func (u *Schedule)ScheduleList(limt int, offet int)([]Schedule,error){
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	var scheduleList []Schedule
	_, err := qs.Limit(limt, offet).All(&scheduleList)
	if err != nil {
		return nil, err
	}
	return scheduleList, err
}



