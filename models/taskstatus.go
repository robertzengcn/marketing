package models

import (
	"github.com/beego/beego/v2/client/orm"
)
type TaskStatus struct {
	Id      int64     `orm:"pk;auto"`
	Name    string    `orm:"size(100)"`
}

func (u *TaskStatus) TableName() string {
	return "task_status"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(TaskStatus))
	// create table
	// orm.RunSyncdb("default", false, true)
}

///find one task status by task status id
func (u *TaskStatus)GetOne(statusId int64)(*TaskStatus,error){
	o := orm.NewOrm()
	taskstatus := TaskStatus{Id: statusId}
	err := o.Read(&taskstatus)
	if(err!=nil){
		return nil, err	
	}else{
		return &taskstatus,nil
	}
}

