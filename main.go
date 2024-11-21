package main

import (
	beego "github.com/beego/beego/v2/server/web"
	_ "marketing/routers"
    //  "fmt"
	// "github.com/beego/beego/v2/client/orm"
	// _ "github.com/go-sql-driver/mysql"
	_ "github.com/beego/beego/v2/server/web/session/redis"
	// "github.com/beego/beego/v2/task"
	"marketing/controllers"
	// "github.com/beego/beego/v2/core/config"
	"github.com/beego/i18n"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/filter/cors"
	// "time"
	"marketing/job"
	"github.com/beego/beego/v2/task"
	_ "marketing/mysqlinit"
)

func init() {
	
	
	// beego.BConfig.WebConfig.Session.SessionProvider = "redis"
	// beego.BConfig.WebConfig.Session.SessionProviderConfig = "redis:6379"
}



func main() {
	if beego.BConfig.RunMode == "dev" {
	
		beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"*"},
			AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
			ExposeHeaders:    []string{"*"},
			AllowCredentials: true,
		}))
		}
	
		beego.InsertFilter("/test/*", beego.BeforeExec, controllers.Filter_user)	
	beego.InsertFilter("/admin/*", beego.BeforeExec, controllers.Filter_user)
	beego.InsertFilter("/api/*", beego.BeforeExec, controllers.Filter_user)

	//beego.InsertFilter("*", beego.BeforeRouter, controllers.Allow_origins)

	// defer task.StopTask()

	f := &logs.PatternLogFormatter{
        Pattern:    "%F:%n|%w%t>> %m",
        WhenFormat: "2006-01-02",
    }
    logs.RegisterFormatter("pattern", f)

    _ = logs.SetGlobalFormatter("pattern")
	beego.AddFuncMap("i18n", i18n.Tr)
	job.InitTask()
	task.StartTask()
	defer task.StopTask()
	beego.Run()
}
