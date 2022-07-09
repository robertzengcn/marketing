package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
	"github.com/google/uuid"
	// "errors"
)

type TaskRun struct {
	Id      int64     `orm:"pk;auto"`
	Task *Task `orm:"rel(fk);" valid:"Required;"`
	Logid string `orm:"size(50)"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
}
func (u *TaskRun) TableName() string {
	return "task_run"
}

func init() {
	
	orm.RegisterModelWithPrefix("mk_",new(TaskRun))	
}
///create a task run record
func (u *TaskRun)CreateRun(taskid int64)(int64,error){
	taskModel:=Task{}
	taskVar,taskErr:=taskModel.GetOne(taskid)
	if(taskErr!=nil){
		return 0,taskErr
	}
	o := orm.NewOrm()
	taskrun := new(TaskRun)
	taskrun.Task=taskVar
	logid:= uuid.New()
	
	taskrun.Logid=logid.String()
	
	taskrunid,terr:=o.Insert(taskrun)
	if(terr!=nil){
		return 0,terr
	}
	return taskrunid,nil
}
///get one task run
func(u *TaskRun)GetOne(taskrunid int64)(*TaskRun,error){
	o := orm.NewOrm()
	taskrun := TaskRun{Id: taskrunid}
	err := o.Read(&taskrun)
	if(err!=nil){
		return nil, err	
	}else{
		return &taskrun,nil
	}
}