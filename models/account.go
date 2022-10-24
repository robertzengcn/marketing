package models

import (
	// "fmt"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"marketing/utils"
)

type Account struct {
	Id      int64     `orm:"pk;auto"`
	Name    string    `orm:"size(100)"`
	Password  string    `orm:"size(255)"`
	Email   string    `orm:"size(150)"`
	Created time.Time `orm:"null;auto_now_add;type(datetime)"`
	Updated time.Time `orm:"null;auto_now;type(datetime)"`
}

func (u *Account) TableName() string {
	return "account"
}

// 设置引擎为 INNODB
func (u *Account) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(Account))
	// create table
	// orm.RunSyncdb("default", false, true)
}

/**
* query user account data
 */
func IndexAllAccount(start int,number int) (accounts []Account, err error) {
	o := orm.NewOrm()

	var us []Account
	count, e := o.QueryTable(new(Account)).Limit(start,number).All(&us, "Name", "Email")
	if e != nil {
		return nil, e
	}
	if count <= 0 {
		return nil, errors.New("nothing found")
	}
	for _, element := range us {
		fmt.Printf("%+v\n", element)
	}
	// qs:=o.QueryTable(new(Account));
	// var result=qs.Filter("id", 1);
	// fmt.Printf(result)
	// fmt.Printf("%+v\n", us)
	return us, nil
}

//按照id查询
func Selectbyid(id int) (account Account, err error) {
	o := orm.NewOrm()
	var us Account
	errs := o.Raw("SELECT id, name,email,created,updated FROM gotest_account WHERE id = ?", id).QueryRow(&us)
	return us, errs
}

//query all rows
func SelectAccountlist() (accounts []Account, err error) {
	o := orm.NewOrm()
	var users []Account
	_, errs := o.QueryTable("gotest_account").All(&users)
	return users, errs
}
///add account
func AddAccount(username string, email string) (id int64, err error) {
	o := orm.NewOrm()
	var us Account
	us.Name = username
	us.Email = email
	id, err = o.Insert(&us)
	// if err == nil {
		return id,err
	// }
}
///check is account valid
func Validaccount(username string, pass string) (account Account, err error) {
	o := orm.NewOrm()
	
	qs := o.QueryTable(new(Account))
	if(utils.ValidEmail(username)){
		err =qs.Filter("email", username).Filter("password", pass).One(&account,"Id","Name","Email")
	}else{
	err =qs.Filter("name", username).Filter("password", pass).One(&account,"Id","Name","Email")
		}
	// account = Account{Name: username,Email:email}
	// err=o.Read(&account)
	
	return account,err	
}
