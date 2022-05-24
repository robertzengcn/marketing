package models

import (
	"github.com/beego/beego/v2/client/orm"
)
type TaskStatus struct {
	Id      int64     `orm:"pk;auto"`
	Name    string    `orm:"size(100)"`
}

func (u *TaskStatus) TableName() string {
	return "taskstatus"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(TaskStatus))
	// create table
	// orm.RunSyncdb("default", false, true)
}

