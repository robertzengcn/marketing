package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	"errors"
)
type TaskDetail struct {
	Id int64  `orm:"pk;auto"`
	TaskId *Task `orm:"rel(fk);" valid:"Required;"`
	Taskkeyword string `orm:"size(1000)"`
}
// set engineer as INNODB
func (u *TaskDetail) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_", new(TaskDetail))	
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
	id, err := o.Insert(&tadetail)
	if err != nil {
		return 0, err
	}
	return id,err
}