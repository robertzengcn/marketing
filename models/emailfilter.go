package models

import (
	"errors"
	"time"
"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"marketing/utils"
)

//this class is used for filter email
type EmailFilter struct {
	Id int64  `orm:"pk;auto"`
	Name string `orm:"size(50)"`
	Description string `orm:"type(text);column(description)"`
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
	valid := validation.Validation{}
	
	b, err := valid.Valid(filter)
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
	id, err := o.Insert(filter)
	return id,err
}

// GetEmailFilterById retrieves an EmailFilter by its Id and account id
func (u *EmailFilter)GetEmailFilterById(id int64,accountId int64) (*EmailFilter, error) {
	o := orm.NewOrm()
	filter := &EmailFilter{}
	err := o.QueryTable(u).Filter("id",id).Filter("account_id",accountId).One(filter,"id","name","description","created","updated")
	if err == orm.ErrNoRows {
		return nil, err
	}
	return filter, err
}

// UpdateEmailFilter updates an existing EmailFilter in the database
func (u *EmailFilter)UpdateEmailFilter(filter *EmailFilter) error {
	valid := validation.Validation{}
	
	b, err := valid.Valid(filter)
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
	o := orm.NewOrm()
	_, uerr := o.Update(filter)
	return uerr
}

// DeleteEmailFilter deletes an EmailFilter from the database
func (u *EmailFilter)DeleteEmailFilter(id int64,accountId int64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable(u).Filter("id",id).Filter("account_id",accountId).Delete()
	return err
}
//update filter detail to filter item, param is filter id and filter detail id array
func (u *EmailFilter)UpdateEmailFilterDetail(id int64,accountId int64,detialIds []int64) error {
	//get old filter detail with filter id
	fileterDetialModel:=EmailFilterDetail{}
	oldfilterdetail,_:=fileterDetialModel.GetEmailFilterDetailByFilterId(id,accountId)
	//delete old filter detail
	for _,olddetail:=range oldfilterdetail{
		//remove one if not exist in new detail array
		oexist := utils.ContainsType(detialIds, olddetail.Id)
		if(!oexist){
			fderr:=fileterDetialModel.DeleteEmailFilterDetail(olddetail.Id,accountId)
			if(fderr!=nil){
				return fderr
			}
		}
	}
	//add new filter detail
	for _,detailId:=range detialIds{
		//check if detail exist
		detail,_:=fileterDetialModel.GetEmailFilterDetailById(detailId,accountId)
		if(detail==nil){
			//add new detail
			emailFilterDetialEntity:=EmailFilterDetail{
				AccountId:&Account{Id:accountId},
				FilterId:&EmailFilter{Id:id},
				Content:detail.Content,
			}
			_,ferr:=fileterDetialModel.CreateEmailFilterDetail(&emailFilterDetialEntity)
			if(ferr!=nil){
				return ferr
			}
		}
	}
	return nil
}
//list email filter by account id
func (u *EmailFilter)ListEmailFilter(accountId int64,page int64, size int64, search string,orderby string) ([]*EmailFilter, error) {
	var emps []*EmailFilter
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	// orm.Debug = true
	cond := orm.NewCondition()
	// qs.Filter("account_id", accountId)
	cond = cond.And("account_id", accountId)
	//qs = qs.SetCond(cond1)
	if(len(search)>0){
		searchCond := orm.NewCondition()
		searchCond = searchCond.Or("name__contains", search)
		//cond =cond.AndCond(cond.Or("tpl_title__contains",search).Or("tpl_content__contains",search))
		cond = cond.AndCond(searchCond)
	}
	qs=qs.SetCond(cond)
	if(len(orderby)>0){
		qs=qs.OrderBy(orderby)
	}else{
		qs=qs.OrderBy("id")
	}
	_,err:=qs.Limit(size, page).All(&emps,"Id","Name","Description","Created","Updated")
	return emps,err
}
//count email filter by account id
func (u *EmailFilter)CountEmailFilter(accountId int64,search string) (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(u)
	// qs.Filter("account_id", accountId)
	cond := orm.NewCondition()
	// qs.Filter("account_id", accountId)
	cond = cond.And("account_id", accountId)
	//qs = qs.SetCond(cond1)
	if(len(search)>0){
		searchCond := orm.NewCondition()
		searchCond = searchCond.Or("name__contains", search)
		//cond =cond.AndCond(cond.Or("tpl_title__contains",search).Or("tpl_content__contains",search))
		cond = cond.AndCond(searchCond)
	}
	qs=qs.SetCond(cond)
	return qs.Count()
}