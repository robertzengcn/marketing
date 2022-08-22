package routers

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"marketing/controllers"
	"strings"
)

func init() {
	langs, err := beego.AppConfig.String("langs::types") // 1
	if err != nil {                                      // 2
		fmt.Println("Failed to load languages from the app.conf")
		// return
	}
	// fmt.Println("11111")
	langsArr := strings.Split(langs, "|")
	// fmt.Printf("%v", langsArr)
	for _, lang := range langsArr {
		// fmt.Print(lang)
		if err := i18n.SetMessage(lang, "conf/"+lang+".ini"); err != nil { // 5
			fmt.Println("Failed to set message file for l10n")
			return
		}
	}

	beego.Router("/", &controllers.MainController{})
	beego.Router("/campaign/create", &controllers.CampaignController{}, "post:CreateCampaign")
	beego.Router("/task/create", &controllers.TaskController{}, "post:CreateTask")
	beego.Router("/login/accountlogin", &controllers.AccountController{}, "post:Validaccount")
	beego.Router("/task/updatetask", &controllers.TaskController{}, "post:UpdateTaskstatus")
	//list campagin
	beego.Router("/campaign", &controllers.CampaignController{}, "get:ListCampaign")
	// beego.Router("/addSite", &controllers.CampaignController{}, "post:Createsite")
	beego.Router("/emailtpl/create", &controllers.EmailtplController{}, "post:CreateEmailtpl")
	beego.Router("/emailservice/add", &controllers.EmailserviceController{}, "post:Addemailservice")
	beego.Router("/emailservice/testsendemail", &controllers.EmailserviceController{}, "post:Testsendemail")
	// beego.Router("/welcome", &controllers.CampaignController{}, "get:Welcome")
	beego.Router("/healthcheck", &controllers.MainController{}, "get:Healthcheck")
	ns := beego.NewNamespace("/test",
		beego.NSInclude(
			&controllers.TestController{},
		),
		beego.NSRouter("/savesearchreq", &controllers.TestController{}, "get:Savesearchrequest"),
	)
	beego.AddNamespace(ns)
}
