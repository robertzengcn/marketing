package models

import (
	// "fmt"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"marketing/utils"
	"github.com/beego/beego/v2/core/logs"
)

type Account struct {
	Id      int64     `orm:"pk;auto"`
	Name    string    `orm:"size(100)"`
	Password  string    `orm:"size(255)"`
	Email   string    `orm:"size(150)"`
	Roles  []*AccountRole `orm:"rel(m2m);rel_table(mk_accout_roles_list)"`
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
func (u *Account)IndexAllAccount(start int,number int) (accounts []Account, err error) {
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
// func (u *Account)Selectbyid(id int) (account Account, err error) {
// 	o := orm.NewOrm()
// 	var us Account
// 	errs := o.Raw("SELECT id, name,email,created,updated FROM gotest_account WHERE id = ?", id).QueryRow(&us)
// 	return us, errs
// }

//query all rows
func (u *Account)SelectAccountlist() (accounts []Account, err error) {
	o := orm.NewOrm()
	var users []Account
	_, errs := o.QueryTable("gotest_account").All(&users)
	return users, errs
}
///add account
func (u *Account)AddAccount(username string, email string,password string) (id int64, err error) {
	o := orm.NewOrm()
	var us Account
	us.Name = username
	us.Email = email
	us.Password=u.EncryptionPass(password)
	id, err = o.Insert(&us)
	// if err == nil {
		return id,err
	// }
}
///check is account valid
func (u *Account)Validaccount(username string, pass string) (Account, error) {
	o := orm.NewOrm()
	l := logs.GetLogger()
	epass:=u.EncryptionPass(pass)
	var account Account
	var err error
	//l.Println(epass)
	qs := o.QueryTable(new(Account))
	if(utils.ValidEmail(username)){
		err =qs.Filter("email", username).Filter("password", epass).One(&account,"Id","Name","Email")
	}else{
	err =qs.Filter("name", username).Filter("password", epass).One(&account,"Id","Name","Email")
		}
	// account = Account{Name: username,Email:email}
	// err=o.Read(&account)
	if(err==nil){
		//get account Roles
		accountrolesmodel:=AccountRolesList{}
		rules,erules:=accountrolesmodel.GetAccountRoleByAccountId(account.Id)
		if(erules!=nil){
		l.Println(erules)
		//loop roles
		
		}else{
			account.Roles=rules
			// l.Println(erules)
			return account, erules
		}
	}
	l.Println(account)
	return account,err	
}
///encryption user password
func (u *Account)EncryptionPass(pass string)(string){
	return utils.Md5V2(pass)
}
///get account by uid
func (u *Account)GetAccountbyid(uid int64) (account Account, err error) {
	o := orm.NewOrm()
	account = Account{Id: uid}
	err=o.Read(&account,"Id","Name","Email","Roles")
	return account,err	
}
//check account role is admin
func (u *Account)IsAdmin(uid int64) (bool) {
	o := orm.NewOrm()
	account := Account{Id: uid}
	err:=o.Read(&account,"Id","Name","Email","Roles")
	if(err!=nil){
		return false
	}
	for _, role := range account.Roles {
		if(role.Name=="admin"){
			return true
		}
	}
	return false
}



