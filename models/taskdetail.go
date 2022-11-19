package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	"errors"
	"marketing/utils"
	"strconv"
)
type TaskDetail struct {
	Id int64  `orm:"pk;auto"`
	Task *Task `orm:"rel(fk);" valid:"Required;"`
	Taskkeyword string `orm:"type(text)"`
	TaskFilename string `orm:"size(250)"`
}
// set engineer as INNODB
func (u *TaskDetail) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_", new(TaskDetail))	
}
func (u *TaskDetail) TableName() string {
	return "task_detail"
}
///save task detail to database
func (u *TaskDetail) Savetaskdetail(tadetail TaskDetail)(int64, error){
	valid := validation.Validation{}
	b, verr := valid.Valid(&tadetail)
    if verr != nil {
        // handle error
		return 0,verr
    }
    if !b {
        // validation does not pass
        var errMessage string
        for _, err := range valid.Errors {
			errMessage+=err.Key+":"+err.Message 
        }
		return 0,errors.New(errMessage)
    }
	o := orm.NewOrm()
	tadetail.TaskFilename=utils.Md5V2(tadetail.Taskkeyword+strconv.FormatInt(tadetail.Id,10))
	id, err := o.Insert(&tadetail)
	if err != nil {
		return 0, err
	}
	return id,err
}
///get task detail by task id
func (u *TaskDetail)Gettaskdetail(taskid int64)(*TaskDetail,error){
	var taskdetail TaskDetail
	o := orm.NewOrm()
	err := o.QueryTable("mk_task_detail").Filter("task_id", taskid).One(&taskdetail)
	if(err!=nil){
		return nil,err
	}
	return &taskdetail,nil 
}