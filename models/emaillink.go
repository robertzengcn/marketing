package models

import (
	  "github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
    "io/ioutil"
    "os"
)

type EmailLink struct {
	Id int64  `orm:"pk;auto" json:"-"`
	Url string `orm:"size(100)" json:"url"`
	Email string `orm:"size(150)" json:"email"`
	Description string `orm:"size(550)" json:"description"`
	Campaign *Campaign `orm:"rel(fk);on_delete(do_nothing)"`
}

func (u *EmailLink) TableName() string {
	return "email"
}

func (u *EmailLink) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(EmailLink))
}
///save email to db
func (u *EmailLink)SaveEmaildb(email EmailLink)(int64,error){
	o := orm.NewOrm()
	id, err := o.Insert(email)
	if(err!=nil){
		return 0,err
	}
	return id,nil
}
///read email obj from json file
func (u *EmailLink)ReademailFile(filepath string)([]EmailLink,error){
	jsonFile, err := os.Open(filepath)
	if err != nil {
        return nil,err
    }	
	defer jsonFile.Close()
	// read our opened xmlFile as a byte array.
    byteValue, _ := ioutil.ReadAll(jsonFile)

    // we initialize our Users array
    var emailarr []EmailLink
	json.Unmarshal(byteValue, &emailarr)
	return 	emailarr,nil
}





