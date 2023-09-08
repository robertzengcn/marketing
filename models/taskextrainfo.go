package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

type TaskExtaInfo struct {
	Id      int64     `orm:"pk;auto"`
	TaskId  int64     `orm:"null;default(0);column(task_id);unique"`
	ResulttaskId  int64     `orm:"null;default(0);column(result_task_id)"`
}

// multiple fields unique key
// func (u *TaskExtaInfo) TableUnique() [][]string {
// 	return [][]string{
// 		[]string{"task_id"},
// 	}
// }

func (u *TaskExtaInfo) TableName() string {
	return "task_extra_info"
}
func (u *TaskExtaInfo) TableEngine() string {
	return "MyISAM"
}
func init() {
orm.RegisterModelWithPrefix("mk_", new(TaskExtaInfo))
}
//get task info by task id
func (u *TaskExtaInfo) Getextrainfotaskid(taskid int64) (TaskExtaInfo, error) {
	o := orm.NewOrm()
	var taskextra TaskExtaInfo
	err := o.QueryTable(new(TaskExtaInfo)).Filter("task_id", taskid).One(&taskextra)
	return taskextra, err
}
//create task extra info
func (u *TaskExtaInfo) CreateTaskExtraInfo(taskid int64,resultTaskid int64) (int64, error) {
	o := orm.NewOrm()
	taskextra := TaskExtaInfo{TaskId: taskid,ResulttaskId:resultTaskid}
	id, err := o.Insert(&taskextra)
	return id, err
}



