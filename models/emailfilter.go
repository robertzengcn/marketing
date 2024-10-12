package models

import (
	//"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"marketing/utils"
)

//this class is used for filter email
type EmailFilter struct {
	Id int64  `orm:"pk;auto"`
	Name string `orm:"size(50)"`
	AccountId   *Account   `orm:"rel(fk);on_delete(do_nothing);column(account_id)" valid:"Required"` 
	Created  time.Time      `orm:"null;auto_now_add;type(datetime)"`
	Updated  time.Time      `orm:"null;auto_now;type(datetime)"`
}

func (u *EmailFilter) TableName() string {
	return "email_filter"
}

// set mysql enginer as INNODB
func (u *EmailFilter) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(EmailFilter))
}

// CreateEmailFilter inserts a new EmailFilter into the database
func (u *EmailFilter)CreateEmailFilter(filter *EmailFilter) (int64,error) {
	o := orm.NewOrm()
	id, err := o.Insert(filter)
	return id,err
}

// GetEmailFilterById retrieves an EmailFilter by its Id and account id
func (u *EmailFilter)GetEmailFilterById(id int64,accountId int64) (*EmailFilter, error) {
	o := orm.NewOrm()
	filter := &EmailFilter{Id: id,AccountId:&Account{Id:accountId}}
	err := o.Read(filter,"Id","Name")
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return filter, err
}

// UpdateEmailFilter updates an existing EmailFilter in the database
func (u *EmailFilter)UpdateEmailFilter(filter *EmailFilter) error {
	o := orm.NewOrm()
	_, err := o.Update(filter)
	return err
}

// DeleteEmailFilter deletes an EmailFilter from the database
func (u *EmailFilter)DeleteEmailFilter(id int64) error {
	o := orm.NewOrm()
	_, err := o.Delete(&EmailFilter{Id: id})
	return err
}
//update filter detail to filter item, param is filter id and filter detail id array
func (u *EmailFilter)UpdateEmailFilterDetail(id int64,detialId[]int64) error {
	//get old filter detail with filter id
	fileterDetialModel:=EmailFilterDetail{}
	oldfilterdetail,_:=fileterDetialModel.GetEmailFilterDetailByFilterId(id)
	if(oldfilterdetail!=nil){
		//delete old filter detail
		for _,olddetail:=range oldfilterdetail{
			//remove one if not exist in new detail array
			oexist := utils.ContainsType(detialId, olddetail.Id)
			
		}
	}

}