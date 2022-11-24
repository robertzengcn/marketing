package mysqlinit

import (
	"marketing/utils"
	// "fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/config"
	_ "github.com/go-sql-driver/mysql"
	"time"
)
func init() {
	mysql_user, uerr := config.String("mysql_user")
	if(uerr!=nil){
		utils.PanicFunc(uerr)
	}
	mysql_pass, perr := config.String("mysql_pass")
	if(perr!=nil){
		utils.PanicFunc(perr)	
	}
	mysql_host, hoerr := config.String("mysql_host")
	if(hoerr!=nil){
		utils.PanicFunc(hoerr)	
	}
	mysql_port, poerr := config.String("mysql_port")
	if(poerr!=nil){
		utils.PanicFunc(poerr)
	}
	mysql_dbname, dbname_err := config.String("mysql_dbname")
	if(dbname_err!=nil){
		utils.PanicFunc(dbname_err)
	}
	orm.RegisterDriver("mysql", orm.DRMySQL)
	 mysqlconurl:=mysql_user+":"+mysql_pass+"@tcp("+mysql_host+":"+mysql_port+")/"+mysql_dbname+"?charset=utf8&parseTime=True&loc=Local"
	
	orm.RegisterDataBase("default", "mysql", mysqlconurl)
	
	// register model
	orm.RunSyncdb("default", false, true)
	orm.Debug,_ = config.Bool("dbdebug")
	logDB,dbErr := orm.GetDB("default")
	if(dbErr!=nil){
		utils.PanicFunc(dbErr)
	}else{
		logDB.SetConnMaxLifetime(30 *time.Second)
	}
}
