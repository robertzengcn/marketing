package models

import (
	 //"errors"
	 "time"
	  "github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

type EmailFilterDetail struct {
	Id int64  `orm:"pk;auto"`
	FilterId *EmailFilter `orm:"rel(fk);column(filter_id)"`
	Content string `orm:"size(500)"`
	Created  time.Time      `orm:"null;auto_now_add;type(datetime)"`
	Updated  time.Time      `orm:"null;auto_now;type(datetime)"`
}


func (u *EmailFilterDetail) TableName() string {
	return "email_filter_detail"
}

// set mysql enginer as INNODB
func (u *EmailFilterDetail) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(EmailFilterDetail))
}
// CreateEmailFilterDetail inserts a new EmailFilterDetail into the database
func (u *EmailFilterDetail)CreateEmailFilterDetail(e *EmailFilterDetail) error {
	o := orm.NewOrm()
	_, err := o.Insert(e)
	return err
}

// GetEmailFilterDetailById retrieves EmailFilterDetail by Id. Returns error if Id doesn't exist
func (u *EmailFilterDetail)GetEmailFilterDetailById(id int64) (*EmailFilterDetail, error) {
	o := orm.NewOrm()
	e := &EmailFilterDetail{Id: id}
	err := o.Read(e)
	if err == orm.ErrNoRows {
		return nil, err
	}
	return e, nil
}

// UpdateEmailFilterDetail updates EmailFilterDetail by Id and returns error if the record to be updated doesn't exist
func (u *EmailFilterDetail)UpdateEmailFilterDetail(e *EmailFilterDetail) error {
	o := orm.NewOrm()
	_, err := o.Update(e)
	return err
}

// DeleteEmailFilterDetail deletes EmailFilterDetail by Id and returns error if the record to be deleted doesn't exist
func (u *EmailFilterDetail)DeleteEmailFilterDetail(id int64) error {
	o := orm.NewOrm()
	e := &EmailFilterDetail{Id: id}
	_, err := o.Delete(e)
	return err
}
//get filter detail by filter id
func (u *EmailFilterDetail)GetEmailFilterDetailByFilterId(filterId int64) ([]*EmailFilterDetail, error) {
	o := orm.NewOrm()
	var filterDetails []*EmailFilterDetail
	_, err := o.QueryTable(new(EmailFilterDetail)).Filter("filter_id", filterId).All(&filterDetails)
	if err != nil {
		return nil, err
	}
	return filterDetails, nil
}