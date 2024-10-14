package models

import (
	//"errors"
	"errors"
	"time"
"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

type EmailFilterDetail struct {
	Id int64  `orm:"pk;auto"`
	FilterId *EmailFilter `orm:"rel(fk);column(filter_id)"`
	AccountId   *Account   `orm:"rel(fk);on_delete(do_nothing);column(account_id)" valid:"Required"` 
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
func (u *EmailFilterDetail)CreateEmailFilterDetail(e *EmailFilterDetail) (int64,error) {

	valid := validation.Validation{}
	
	b, err := valid.Valid(e)
    if err != nil {
		
       return 0,err
    }
	if !b {

		var errMessage string
	 // validation does not pass
	 for _, err := range valid.Errors {
		// log.Println(err.Key, err.Message)
		errMessage+=err.Key+":"+err.Message
		}
		return 0,errors.New(errMessage)
	}
	o := orm.NewOrm()
	id, ierr := o.Insert(e)
	return id, ierr
}

// GetEmailFilterDetailById retrieves EmailFilterDetail by Id. Returns error if Id doesn't exist
func (u *EmailFilterDetail)GetEmailFilterDetailById(id int64,accountId int64) (*EmailFilterDetail, error) {
	o := orm.NewOrm()
	e := &EmailFilterDetail{}
	err := o.QueryTable(u).Filter("Id", id).Filter("account_id", accountId).One(e)
	if err == orm.ErrNoRows {
		return nil, err
	}
	return e, nil
}

// UpdateEmailFilterDetail updates EmailFilterDetail by Id and returns error if the record to be updated doesn't exist
func (u *EmailFilterDetail)UpdateEmailFilterDetail(e *EmailFilterDetail) error {

	valid := validation.Validation{}
	
	b, err := valid.Valid(e)
    if err != nil {
		
       return err
    }
	if !b {

		var errMessage string
	 // validation does not pass
	 for _, err := range valid.Errors {
		// log.Println(err.Key, err.Message)
		errMessage+=err.Key+":"+err.Message
		}
		return errors.New(errMessage)
	}
	//check email filter belong to the user
	uef,uerr:=u.GetEmailFilterDetailById(e.Id,e.AccountId.Id)
	if uerr!=nil{
		return uerr
	}
	if(uef==nil){
		return errors.New("email filter detail not exist")
	}
	o := orm.NewOrm()
	_, uerrs := o.Update(e)
	return uerrs
}

// DeleteEmailFilterDetail deletes EmailFilterDetail by Id and returns error if the record to be deleted doesn't exist
func (u *EmailFilterDetail)DeleteEmailFilterDetail(id int64,accountId int64) error {
	o := orm.NewOrm()
	// e := &EmailFilterDetail{Id: id}
	_, err:=o.QueryTable(u).Filter("Id", id).Filter("account_id", accountId).Delete()
	// _, err := o.Delete(e)
	return err
}
//get filter detail by filter id
func (u *EmailFilterDetail)GetEmailFilterDetailByFilterId(filterId int64,accountId int64) ([]*EmailFilterDetail, error) {
	o := orm.NewOrm()
	var filterDetails []*EmailFilterDetail
	_, err := o.QueryTable(u).Filter("filter_id", filterId).Filter("account_id", accountId).All(&filterDetails)
	if err != nil {
		return nil, err
	}
	return filterDetails, nil
}
