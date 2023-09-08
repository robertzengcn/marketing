package models

import (
	"marketing/utils"
	"time"
	"strings"
	"github.com/beego/beego/v2/client/orm"
	// "github.com/beego/beego/v2/core/logs"
	"strconv"
)
var DefaultApiauth *Apiauth

type Apiauth struct {
	Id         int64  `orm:"pk;auto"`
	UserName   string `orm:"size(100)"`
	Password   string `orm:"size(100)"`
	Updated time.Time `orm:"null;auto_now;type(datetime)"`
}

// set engine as INNODB
func (u *Apiauth) TableEngine() string {
	return "INNODB"
}

func init() {
	// set default database
	orm.RegisterModelWithPrefix("mk_",new(Apiauth))
}

//get api data id by username and password
func (u *Apiauth)GetApiAuth(username string,password string) (int64, error) {
	rediskey:="apiauth:"+utils.Md5V2(username+":"+password)
	pass:=utils.Md5V2(password)
	
	//check key exist in redis first
	redisvalue,_:=utils.GetStr(rediskey)
	expectval,_:=strconv.ParseInt(redisvalue,10,64)
	if(expectval>0){
		return expectval,nil
	}
	o := orm.NewOrm()
	var apiauth Apiauth
	err := o.QueryTable(new(Apiauth)).Filter("user_name", strings.TrimSpace(username)).Filter("password", strings.TrimSpace(pass)).One(&apiauth,"id")
	
	if err != nil {
		return 0, err
	}else{
		utils.SetStr(rediskey,strconv.FormatInt(apiauth.Id,10),time.Hour*1)
	}
	// logs.Info(apiauth.Id)
	return apiauth.Id, nil
}