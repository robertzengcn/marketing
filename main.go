package main

import (
	beego "github.com/beego/beego/v2/server/web"
	_ "marketing/routers"
    // "fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/beego/beego/v2/server/web/session/redis"
	// "github.com/beego/beego/v2/task"
	"marketing/controllers"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/i18n"
)

func init() {
	mysql_user, uerr := config.String("mysql_user")
	if(uerr!=nil){
		innerFunc(uerr)
	}
	mysql_pass, perr := config.String("mysql_pass")
	if(perr!=nil){
		innerFunc(perr)	
	}
	mysql_host, hoerr := config.String("mysql_host")
	if(hoerr!=nil){
		innerFunc(hoerr)	
	}
	mysql_port, poerr := config.String("mysql_port")
	if(poerr!=nil){
		innerFunc(poerr)
	}
	mysql_dbname, dbname_err := config.String("mysql_dbname")
	if(dbname_err!=nil){
		innerFunc(dbname_err)
	}
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", mysql_user+":"+mysql_pass+"@tcp("+mysql_host+":"+mysql_port+")/"+mysql_dbname+"?charset=utf8&parseTime=True&loc=Local")
	
	// register model
	orm.RunSyncdb("default", false, true)
	orm.Debug,_ = config.Bool("dbdebug")
	// beego.BConfig.WebConfig.Session.SessionProvider = "redis"
	// beego.BConfig.WebConfig.Session.SessionProviderConfig = "redis:6379"
}
func innerFunc(errorObj error ) {
	panic(errorObj.Error())
}


func main() {
	beego.InsertFilter("/*", beego.BeforeExec, controllers.Filter_user)
	// utils.InitTask()
	// task.StartTask()
	// defer task.StopTask()
	beego.AddFuncMap("i18n", i18n.Tr)
	beego.Run()
}
